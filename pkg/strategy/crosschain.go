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

type crossChain struct{}

type crossChainArgs struct {
	SourceChain string  `json:"source_chain"`
	Token       string  `json:"token"`
	Amount      float64 `json:"amount"`
	Receiver    string  `json:"receiver"`
	TargetChain string  `json:"target_chain"`
	Summary     string  `json:"summary"`
}

func (c crossChain) Prompt() string {
	return `I want you act as an experienced cryptocurrency researcher, here are my cross chain expectations. 
			Analyze the core elements in cross chain asset operation defined in functions.`
}

func (c crossChain) Functions() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{{
		Name: "cross_chain_analyze",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"source_chain": {
					Type:        jsonschema.String,
					Description: "The source blockchain name, e.g. Ethereum Mainnet",
				},
				"token": {
					Type:        jsonschema.String,
					Description: "The cross chain token, e.g. USDC",
				},
				"amount": {
					Type:        jsonschema.Number,
					Description: "The cross asset amount, e.g. 100 ether",
				},
				"receiver": {
					Type:        jsonschema.String,
					Description: "The receiver blockchain address, e.g. 0xd8da6bf26964af9d7eed9e03e53415d37aa96045",
				},
				"target_chain": {
					Type:        jsonschema.String,
					Description: "The target blockchain name, e.g. Avalance",
				},
				"summary": {
					Type:        jsonschema.String,
					Description: "The summary of cross chain operation",
				},
			},
			Required: []string{"source_chain", "token", "amount", "target_chain", "summary"},
		},
	}}
}

func (c crossChain) Render(resp *model.DemandResponse, name, args string) error {
	if name == "cross_chain_analyze" {
		in := crossChainArgs{}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return err
		}
		resp.Summary = in.Summary
		resp.Category = "crossChain"
		sourceChainId, err := data.GetChainIdByName(in.SourceChain)
		if err != nil {
			return err
		}
		targetChainId, err := data.GetChainIdByName(in.TargetChain)
		if err != nil {
			return err
		}
		resp.Detail = model.DetailResp{
			Reply: fmt.Sprintf("Ok I will transfer %s %s to %s from %s to %s", utils.Float2String(in.Amount), in.Token, in.Receiver, in.SourceChain, in.TargetChain),
			OPs: []interface{}{
				model.CrossChainResponse{
					Type:            model.CrossChainTransfer,
					SourceChainId:   sourceChainId,
					SourceChainName: in.SourceChain,
					Token:           in.Token,
					Amount:          utils.Float2String(in.Amount),
					Receiver:        in.Receiver,
					TargetChainId:   targetChainId,
					TargetChainName: in.TargetChain,
				},
			},
		}
		return nil
	}
	return ErrFunctionNotDefined
}

type crossChainAbstractionArgs struct {
	Token                          string  `json:"token"`
	SourceChain                    string  `json:"source_chain"`
	SourceChainTokenBalance        float64 `json:"source_chain_token_balance"`
	SourceChainTokenBalanceDecimal decimal.Decimal
	TargetChain                    string  `json:"target_chain"`
	TargetChainTokenBalance        float64 `json:"target_chain_token_balance"`
	TargetChainTokenBalanceDecimal decimal.Decimal
	TransferAmount                 float64 `json:"transfer_amount"`
	TransferAmountDecimal          decimal.Decimal
	Receiver                       string `json:"receiver"`
	Summary                        string `json:"summary"`
}

type chainAbstraction struct{}

func (c chainAbstraction) Prompt() string {
	return `I want you to become an experienced cryptocurrency researcher that allows the use of cross-chain protocols.
Analyze the core elements of cross-chain asset operations as defined in the function. Requirements may be missing parameters, if they are missing please don't use defaults, don't autofill, and just don't give values for the missing parameters. Here is a description of my requirements for the transaction: `
}

func (c chainAbstraction) Functions() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{{
		Name: "cross_chain_abstraction",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"source_chain": {
					Type:        jsonschema.String,
					Description: "The source blockchain name, e.g. Ethereum Mainnet",
				},
				"source_chain_token_balance": {
					Type:        jsonschema.Number,
					Description: "The balance of token on source blockchain, e.g. 1 ether",
				},
				"target_chain_token_balance": {
					Type:        jsonschema.Number,
					Description: "The balance of token on target blockchain, e.g. 1 ether",
				},
				"token": {
					Type:        jsonschema.String,
					Description: "The cross chain token, e.g. USDC",
				},
				"transfer_amount": {
					Type:        jsonschema.Number,
					Description: "The cross asset amount, e.g. 100 ether",
				},
				"receiver": {
					Type:        jsonschema.String,
					Description: "The receiver blockchain address, e.g. 0xd8da6bf26964af9d7eed9e03e53415d37aa96045",
				},
				"target_chain": {
					Type:        jsonschema.String,
					Description: "The target blockchain name, e.g. Avalance",
				},
				"summary": {
					Type:        jsonschema.String,
					Description: "The summary of cross chain operation",
				},
			},
			Required: []string{},
		},
	}}
}

func (c chainAbstraction) Render(resp *model.DemandResponse, name, args string) error {
	if name == "cross_chain_abstraction" {
		in := crossChainAbstractionArgs{}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return err
		}
		if in.TargetChainTokenBalanceDecimal.IsZero() {
			in.TargetChainTokenBalanceDecimal = decimal.NewFromFloat(in.TargetChainTokenBalance)
		}
		if in.SourceChainTokenBalanceDecimal.IsZero() {
			in.SourceChainTokenBalanceDecimal = decimal.NewFromFloat(in.SourceChainTokenBalance)
		}
		if in.TransferAmountDecimal.IsZero() {
			in.TransferAmountDecimal = decimal.NewFromFloat(in.TransferAmount)
		}
		resp.Summary = in.Summary
		resp.Category = "crossChainAbstraction"
		if reply, ok := in.isEmpty(); ok {
			resp.Detail = model.DetailResp{
				Reply: reply,
				OPs:   nil,
			}
			return nil
		}
		// case 1: target chain enough
		if in.TargetChainTokenBalanceDecimal.Cmp(in.TransferAmountDecimal) > 0 {
			targetChainId, err := data.GetChainIdByName(in.TargetChain)
			if err != nil {
				return err
			}
			resp.Detail = model.DetailResp{
				Reply: fmt.Sprintf("Ok I will transfer %s %s to %s on %s", in.TransferAmountDecimal.String(), in.Token, in.Receiver, in.TargetChain),
				OPs: []interface{}{
					model.CrossChainResponse{
						Type:            model.ChainInternalTransfer,
						SourceChainId:   targetChainId,
						SourceChainName: in.TargetChain,
						Token:           in.Token,
						Amount:          in.TransferAmountDecimal.String(),
						Receiver:        in.Receiver,
						TargetChainName: in.TargetChain,
						TargetChainId:   targetChainId,
					},
				},
			}
			return nil
		}
		// case 2: source chain + target chain
		if in.SourceChainTokenBalanceDecimal.Add(in.TargetChainTokenBalanceDecimal).Cmp(in.TransferAmountDecimal) > 0 {
			sourceChainId, err := data.GetChainIdByName(in.SourceChain)
			if err != nil {
				log.Errorf("get chain id error: %v", err)
				resp.Detail = model.DetailResp{
					Reply: "cross chain query failed",
					OPs:   nil,
				}
				return nil
			}
			targetChainId, err := data.GetChainIdByName(in.TargetChain)
			if err != nil {
				log.Errorf("get chain id error: %v", err)
				resp.Detail = model.DetailResp{
					Reply: "cross chain query failed",
					OPs:   nil,
				}
				return nil
			}
			ok, ret := pkg.Base.CheckCross(sourceChainId, targetChainId, strings.ToUpper(in.Token))
			if !ok {
				log.Warnf("token:%s cannot cross chain from %s to %s", in.Token, in.SourceChain, in.TargetChain)
				resp.Detail = model.DetailResp{
					Reply: "cross chain failed",
					OPs:   nil,
				}
				return nil
			}
			crossChainBalance := in.TransferAmountDecimal.Sub(in.TargetChainTokenBalanceDecimal)
			sourceChainId, err = data.GetChainIdByName(in.SourceChain)
			if err != nil {
				return err
			}
			targetChainId, err = data.GetChainIdByName(in.TargetChain)
			if err != nil {
				return err
			}
			internalTransfer := &model.CrossChainResponse{
				Type:            model.ChainInternalTransfer,
				SourceChainName: in.TargetChain,
				SourceChainId:   targetChainId,
				Token:           in.Token,
				Amount:          in.TargetChainTokenBalanceDecimal.String(),
				Receiver:        in.Receiver,
				TargetChainName: in.TargetChain,
				TargetChainId:   targetChainId,
			}
			crossTransfer := &model.CrossChainResponse{
				RawResponse:     ret,
				Type:            model.CrossChainTransfer,
				SourceChainId:   sourceChainId,
				SourceChainName: in.SourceChain,
				Token:           in.Token,
				Amount:          crossChainBalance.String(),
				Receiver:        in.Receiver,
				TargetChainId:   targetChainId,
				TargetChainName: in.TargetChain,
			}
			if in.TargetChainTokenBalanceDecimal.IsZero() {
				resp.Detail = model.DetailResp{
					Reply: fmt.Sprintf("Ok I will transfer %s %s to %s from %s to %s",
						crossChainBalance.String(), in.Token, in.Receiver, in.SourceChain, in.TargetChain),
					OPs: []interface{}{crossTransfer},
				}
			} else {
				resp.Detail = model.DetailResp{
					Reply: fmt.Sprintf("Ok I will transfer %s %s to %s from %s to %s, and transfer %s %s to %s on %s",
						crossChainBalance.String(), in.Token, in.Receiver, in.SourceChain, in.TargetChain, in.TargetChainTokenBalanceDecimal.String(), in.Token, in.Receiver, in.TargetChain),
					OPs: []interface{}{internalTransfer, crossTransfer},
				}
			}
			return nil
		}
		// case 3: not enough
		resp.Detail = model.DetailResp{
			Reply: "Insufficient Balance",
			OPs:   nil,
		}
		return nil
	}
	return ErrFunctionNotDefined
}

func (a *crossChainAbstractionArgs) isEmpty() (reply string, ok bool) {
	ok = true
	if a.Token == "" {
		reply = "missing token"
		return
	}
	if a.SourceChain == "" {
		reply = "missing source chain"
		return
	}
	if a.Receiver == "" {
		reply = "missing receiver"
		return
	}
	if a.TargetChain == "" {
		reply = "missing target token"
		return
	}
	if a.TransferAmountDecimal.IsZero() {
		reply = "missing transfer amount"
		return
	}
	ok = false
	return
}
