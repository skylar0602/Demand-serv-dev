package strategy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/smarterwallet/demand-abstraction-serv/data"

	"github.com/smarterwallet/demand-abstraction-serv/pkg"

	log "github.com/cihub/seelog"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/smarterwallet/demand-abstraction-serv/model"
	"github.com/smarterwallet/demand-abstraction-serv/utils"
)

type transferArgs struct {
	SourceChain string  `json:"source_chain"`
	Token       string  `json:"token"`
	Amount      float64 `json:"amount"`
	AmtDecimal  decimal.Decimal
	Receiver    string `json:"receiver"`
	TargetChain string `json:"target_chain"`
	IsUsd       bool   `json:"is_usd"`
}

type transfer struct {
	balance *model.CtxRequest
}

func (t transfer) Prompt() string {
	return fmt.Sprintf(`As a seasoned cryptocurrency researcher, your task is to analyze cross-chain transfer demands. transfer token from chain:%s.
Focus on identifying key elements in the transactions, particularly noting the source and target chains involved.
If the user explicitly mentions that the transfer is in US dollars but not stable coins, set is_usd true`, t.balance.BaseChain)
}

func (t transfer) Functions() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{
		{
			Name: "get_trade_strategy",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"source_chain": {
						Type:        jsonschema.String,
						Description: "The source blockchain name, e.g. Ethereum Mainnet",
					},
					"token": {
						Type:        jsonschema.String,
						Description: "The transfer token, e.g. USDC",
					},
					"amount": {
						Type:        jsonschema.Number,
						Description: "The transfer amount, e.g. 100",
					},
					"receiver": {
						Type:        jsonschema.String,
						Description: "The receiver address, e.g. 0xd8da6bf26964af9d7eed9e03e53415d37aa96045",
					},
					"target_chain": {
						Type:        jsonschema.String,
						Description: "The target blockchain name, e.g. Ethereum Mainnet",
					},
					"is_usd": {
						Type:        jsonschema.Boolean,
						Description: "Whether the transfer is in US dollars or an equivalent (e.g., through a USD-pegged stablecoin like USDT)",
					},
				},
				Required: []string{"source_chain", "token", "amount", "receiver", "target_chain", "is_usd"},
			},
		},
	}
}

func (t transfer) Render(resp *model.DemandResponse, name, args string) error {
	if name != "get_trade_strategy" {
		return ErrFunctionNotDefined
	}
	in := transferArgs{}
	if err := json.Unmarshal([]byte(args), &in); err != nil {
		return err
	}
	in.AmtDecimal = decimal.NewFromFloat(in.Amount)
	// 1. is usd
	if in.IsUsd {
		// todo choose stable usd coin
		in.Token = "USDC"
	}
	in.SourceChain = strings.ToLower(in.SourceChain)
	in.TargetChain = strings.ToLower(in.TargetChain)
	in.Token = strings.ToUpper(in.Token)
	if in.SourceChain != t.balance.BaseChain {
		log.Warnf("unexpected source chain: %s", in.SourceChain)
		in.SourceChain = t.balance.BaseChain
	}
	// 1. internal
	if in.SourceChain == in.TargetChain || in.TargetChain == "" {
		t.internalTransfer(in, resp)
		return nil
	}
	// 2. cross chain: swap on source chain first
	potentialSwapPairs := make(map[model.Reserve]struct{})
	tokens := t.balance.GetTokens(in.SourceChain)
	for _, token := range tokens {
		if token.Symbol == in.Token {
			continue
		}
		potentialSwapPairs[token] = struct{}{}
	}
	currTokenBalance := t.balance.GetTokenBalance(in.SourceChain, in.Token)
	targetTokenBalance := t.balance.GetTokenBalance(in.TargetChain, in.Token)
	var (
		swapOp model.SwapResponse
		ok     bool
	)
	if currTokenBalance.Add(targetTokenBalance).Cmp(in.AmtDecimal) < 0 {
		// need to swap
		swapOp, ok = t.potentialSwap(potentialSwapPairs, in.SourceChain, in.Token, in.AmtDecimal.Sub(targetTokenBalance))
		if !ok {
			resp.Detail = model.DetailResp{
				Reply: "swap not support",
				OPs:   nil,
			}
			return nil
		}
	}
	crossArgs := crossChainAbstractionArgs{
		Token:                          in.Token,
		SourceChain:                    in.SourceChain,
		SourceChainTokenBalanceDecimal: t.balance.GetTokenBalance(in.SourceChain, in.Token).Add(in.AmtDecimal),
		TargetChain:                    in.TargetChain,
		TargetChainTokenBalanceDecimal: targetTokenBalance,
		TransferAmountDecimal:          in.AmtDecimal,
		Receiver:                       in.Receiver,
		Summary:                        "",
	}
	msg, _ := json.Marshal(crossArgs)
	_ = chainAbstraction{}.Render(resp, "cross_chain_abstraction", string(msg))
	if swapOp.Dex != "" {
		newOps := []interface{}{swapOp}
		newOps = append(newOps, resp.Detail.OPs...)
		resp.Detail.OPs = newOps
	}
	return nil
}

func (t transfer) internalTransfer(in transferArgs, resp *model.DemandResponse) {
	tokenBalance := t.balance.GetTokenBalance(in.SourceChain, in.Token)
	// 2.1 no need to swap
	if tokenBalance.Cmp(in.AmtDecimal) > 0 {
		sourceChainId, err := data.GetChainIdByName(in.SourceChain)
		if err != nil {
			log.Errorf("get chain id error: %v", err)
			return
		}
		targetChainId, err := data.GetChainIdByName(in.TargetChain)
		if err != nil {
			log.Errorf("get chain id error: %v", err)
			return
		}
		resp.Detail = model.DetailResp{
			Reply: fmt.Sprintf("Ok I will transfer %s %s to %s on %s", utils.Float2String(in.Amount), in.Token, in.Receiver, in.SourceChain),
			OPs: []interface{}{
				model.CrossChainResponse{
					Type:            model.ChainInternalTransfer,
					SourceChainId:   targetChainId,
					SourceChainName: in.TargetChain,
					Token:           in.Token,
					Amount:          utils.Float2String(in.Amount),
					Receiver:        in.Receiver,
					TargetChainName: in.SourceChain,
					TargetChainId:   sourceChainId,
				},
			},
		}
		return
	}
	// 2.2 need to swap
	potentialSwapPairs := make(map[model.Reserve]struct{})
	tokens := t.balance.GetTokens(in.SourceChain)
	for _, token := range tokens {
		if token.Symbol == in.Token {
			continue
		}
		potentialSwapPairs[token] = struct{}{}
	}
	swapOp, ok := t.potentialSwap(potentialSwapPairs, in.SourceChain, in.Token, in.AmtDecimal)
	if !ok {
		resp.Detail = model.DetailResp{
			Reply: "swap not support",
			OPs:   nil,
		}
		return
	}
	sourceChainId, err := data.GetChainIdByName(in.SourceChain)
	if err != nil {
		log.Errorf("get chain id error: %v", err)
		return
	}
	targetChainId, err := data.GetChainIdByName(in.TargetChain)
	if err != nil {
		log.Errorf("get chain id error: %v", err)
		return
	}
	resp.Detail = model.DetailResp{
		Reply: fmt.Sprintf("Ok I will transfer %s %s to %s on %s", utils.Float2String(in.Amount), in.Token, in.Receiver, in.SourceChain),
		OPs: []interface{}{swapOp, model.CrossChainResponse{
			Type:            model.ChainInternalTransfer,
			SourceChainId:   targetChainId,
			SourceChainName: in.TargetChain,
			Token:           in.Token,
			Amount:          utils.Float2String(in.Amount),
			Receiver:        in.Receiver,
			TargetChainName: in.SourceChain,
			TargetChainId:   sourceChainId,
		},
		},
	}
}

func (t transfer) potentialSwap(pairs map[model.Reserve]struct{}, chain, outToken string, minOut decimal.Decimal) (model.SwapResponse, bool) {
	id, err := data.GetChainIdByName(chain)
	if err != nil {
		log.Errorf("get chain id error: %v", err)
		return model.SwapResponse{}, false
	}
	currBalance := t.balance.GetTokenBalance(chain, outToken)
	swapOutAmt := minOut.Sub(currBalance)
	var (
		bestSwapToken = ""
		swapIn        = ""
		swapOut       = swapOutAmt.String()
		dex           = "uniswap"
		rawResp       json.RawMessage
	)
	for reserve := range pairs {
		if reserve.Balance <= 0 {
			continue
		}
		out, _ := swapOutAmt.Float64()
		req := model.SwapReq{
			ChainId:         id,
			TokenInAddress:  reserve.Address,
			TokenOutAddress: t.balance.GetTokenAddress(chain, outToken),
			AmountOut:       out,
		}
		minIn, body, err := pkg.Base.CheckSwap(req)
		if err != nil {
			continue
		}
		bestSwapToken = reserve.Symbol
		swapIn = minIn
		rawResp = body
		break
	}
	if bestSwapToken == "" {
		return model.SwapResponse{}, false
	}
	return model.SwapResponse{
		Type:        "swap",
		ChainId:     id,
		ChainName:   chain,
		RawResponse: rawResp,
		SourceToken: bestSwapToken,
		TargetToken: outToken,
		SwapIn:      swapIn,
		SwapOut:     swapOut,
		Dex:         dex,
	}, true
}
