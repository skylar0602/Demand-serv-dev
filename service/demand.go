package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/smarterwallet/demand-abstraction-serv/pkg"

	log "github.com/cihub/seelog"

	"github.com/google/uuid"

	"github.com/smarterwallet/demand-abstraction-serv/data"

	"github.com/pkg/errors"
	"github.com/smarterwallet/demand-abstraction-serv/config"
	"github.com/smarterwallet/demand-abstraction-serv/model"
	"github.com/smarterwallet/demand-abstraction-serv/pkg/llm"
	"github.com/smarterwallet/demand-abstraction-serv/pkg/strategy"
)

type DemandService struct {
	cfg                *config.Config
	llm                llm.Llm
	cache              *data.Cache
	localConversations map[string]int64
	mu                 sync.Mutex
	tokens             sync.Map
}

type TokenInfo struct {
	Symbol  string
	Address string
	Decimal int
}

func NewDemandService(cfg *config.Config) *DemandService {
	var llmInstance llm.Llm
	if cfg.AiConfig.APIKey == "" {
		llmInstance = llm.NewMockOpenAI()
	} else {
		llmInstance = llm.NewOpenAI(cfg.AiConfig)
	}
	cache, err := data.NewCache(cfg.Redis)
	if err != nil {
		log.Errorf("init cache error: %v", err)
		return nil
	}
	ds := &DemandService{cfg: cfg, llm: llmInstance, cache: cache, localConversations: make(map[string]int64), mu: sync.Mutex{}, tokens: sync.Map{}}
	if err := ds.loadTokens(); err != nil {
		log.Errorf("init tokens error: %v", err)
		return nil
	}
	go ds.invalidConversation()
	return ds
}

func (s *DemandService) invalidConversation() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		s.loadTokens()
		now := time.Now().Unix()
		for cid, ttl := range s.localConversations {
			if now > ttl {
				s.mu.Lock()
				log.Errorf("invalidConversation cid=%s\n", cid)
				delete(s.localConversations, cid)
				s.cache.Invalid(context.Background(), cid)
				s.mu.Unlock()
			}
		}
	}
}

func (s *DemandService) loadTokens() error {
	res, err := pkg.Base.LoadTokens()
	if err != nil {
		log.Errorf("load tokens error: %v", err)
		return err
	}
	for _, chain := range res.Result.Chain {
		data.ChainIDMap[strings.ToLower(chain.Name)] = chain.ID
		tokens := make(map[string]TokenInfo)
		for _, token := range chain.Tokens {
			tokens[token.Name] = TokenInfo{
				Symbol:  token.Name,
				Address: token.Address,
				Decimal: token.Decimal,
			}
		}
		s.tokens.Store(strings.ToLower(chain.Name), tokens)
	}
	return nil
}

func (s *DemandService) NewChat() string {
	cid := uuid.NewString()
	s.mu.Lock()
	s.localConversations[cid] = time.Now().Add(10 * time.Minute).Unix()
	s.mu.Unlock()
	return cid
}

func (s *DemandService) analyzeStrategy(ctx context.Context, demand string, demandCtx *model.CtxRequest) strategy.IStrategy {
	selectStrategy, err := strategy.MatchStrategy("selectStrategy", demandCtx)
	if err != nil {
		return nil
	}
	name, args, err := s.llm.Chat(ctx, selectStrategy.Prompt(), demand, selectStrategy.Functions())
	if err != nil {
		log.Errorf("analyzeStrategy err=%v\n", err)
		return nil
	}
	type selectStrategyArgs struct {
		Strategy string `json:"strategy"`
	}
	if name == "select_strategy" {
		in := selectStrategyArgs{}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return nil
		}
		log.Infof("selectStrategy=%s\n", in.Strategy)
		st, err := strategy.MatchStrategy(in.Strategy, demandCtx)
		if err != nil {
			return nil
		}
		return st
	}
	return nil
}

func (s *DemandService) InitCtx(ctx context.Context, cid string, req *model.CtxRequest) error {
	req.Format()
	if req.Address == "" {
		return errors.New("address not found")
	}
	_, ok := s.tokens.Load(strings.ToLower(req.BaseChain))
	if !ok {
		return errors.New("chain not found")
	}
	userBalances := make(map[string][]model.Reserve)
	for c, reserves := range req.Balances {
		if _, ok := s.tokens.Load(c); !ok {
			return errors.New("chain not found")
		}
		for i, reserve := range reserves {
			mm, ok := s.tokens.Load(c)
			tt, ok := mm.(map[string]TokenInfo)[reserve.Symbol]
			if !ok {
				return errors.New("token not found symbol")
			}
			reserves[i].Address = tt.Address
		}
		userBalances[c] = reserves
	}
	req.Balances = userBalances
	if err := s.cache.SetCtx(ctx, cid, req, time.Hour); err != nil {
		return err
	}
	return nil
}

func (s *DemandService) ChatDemand(ctx context.Context, cid, demand string) (*model.DemandResponse, error) {
	demandCtx := s.prepareCtx(ctx, cid)
	st := s.analyzeStrategy(ctx, demand, demandCtx)
	if st == nil {
		return nil, errors.New("strategy not found")
	}
	history, err := s.getHistory(ctx, cid)
	if err != nil {
		log.Errorf("getHistory err=%s\n", err)
		return nil, err
	}
	injectedPrompt := fmt.Sprintf("%s, Here is the chat history:%s", st.Prompt(), history)
	name, args, err := s.llm.Chat(ctx, injectedPrompt, demand, st.Functions())
	if err != nil {
		return nil, errors.Wrap(err, "ChatDemand")
	}
	resp := &model.DemandResponse{}
	if err := st.Render(resp, name, args); err != nil {
		return nil, errors.Wrap(err, "ChatDemand")
	}
	if err := s.appendToHistory(ctx, cid, demand, resp.Detail.Reply); err != nil {
		log.Errorf("appendToHistory err=%s\n", err)
		return nil, err
	}
	return resp, nil
}

func (s *DemandService) prepareCtx(ctx context.Context, cid string) *model.CtxRequest {
	demandCtx := &model.CtxRequest{}
	res, err := s.cache.GetCtx(ctx, cid)
	if err == nil {
		if err = json.Unmarshal([]byte(res), demandCtx); err != nil {
			log.Errorf("prepareCtx err=%s\n", err)
		}
	}
	return demandCtx
}

func (s *DemandService) getHistory(ctx context.Context, cid string) (string, error) {
	dialogues, err := s.cache.ChatHistory(ctx, cid)
	if err != nil {
		return "", err
	}
	history, err := json.Marshal(dialogues)
	if err != nil {
		return "", err
	}
	log.Infof("cid:%s history: %s\n", cid, string(history))
	return string(history), nil
}

func (s *DemandService) appendToHistory(ctx context.Context, cid, demand, reply string) error {
	if err := s.cache.AppendChat(ctx, cid, demand, model.DialogueRoleUser); err != nil {
		return err
	}
	if err := s.cache.AppendChat(ctx, cid, reply, model.DialogueRoleAI); err != nil {
		return err
	}
	return nil
}
