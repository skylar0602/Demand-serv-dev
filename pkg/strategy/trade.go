package strategy

import (
	"encoding/json"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/smarterwallet/demand-abstraction-serv/model"
)

type trade2Earn struct{}

type trade2EarnArgs struct {
	Maximum string `json:"maximum"`
	Minimum string `json:"minimum"`
	Summary string `json:"summary"`
}

func (t trade2Earn) Prompt() string {
	return `I want you act as an experienced cryptocurrency investor, here are my investment expectations. 
			Give me the advice investment rate of return range. Mention that if the input is irrelevant about cryptocurrency invest, reply with zero.`
}

func (t trade2Earn) Functions() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{{
		Name: "get_trade_to_earn_strategy",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"minimum": {
					Type:        jsonschema.String,
					Description: "The minimum rate of return, e.g. 6%",
				},
				"maximum": {
					Type:        jsonschema.String,
					Description: "The maximum rate of return, e.g. 10%",
				},
				"summary": {
					Type:        jsonschema.String,
					Description: "Summary of investment",
				},
			},
			Required: []string{"minimum", "maximum", "summary"},
		},
	}}
}

func (t trade2Earn) Render(resp *model.DemandResponse, name, args string) error {
	if name == "get_trade_to_earn_strategy" {
		in := trade2EarnArgs{}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return err
		}
		resp.Summary = in.Summary
		resp.Category = "trade2Earn"
		resp.Detail = model.DetailResp{
			Reply: "",
			OPs: []interface{}{
				model.TradeStrategyResponse{
					BotName:   "Recommended strategy One-time decentralized automated trading botStrategy",
					Strategy:  "Simple spot grid",
					MinReturn: percentStr2Decimal(in.Minimum),
					MaxReturn: percentStr2Decimal(in.Maximum),
					Operations: []model.Operation{
						{
							Seq:  1,
							Type: "swap",
							Param: model.OperationParam{
								From:             "USWT",
								To:               "SWT",
								GasFee:           "0.5",
								FeeUint:          "SWT",
								ConditionsSymbol: nil,
								Conditions:       nil,
							},
						},
						{
							Seq:  2,
							Type: "sell",
							Param: model.OperationParam{
								From:             "",
								To:               "USWT",
								GasFee:           "0.5",
								FeeUint:          "SWT",
								ConditionsSymbol: []string{"AND"},
								Conditions: []model.Condition{
									{
										TokenName:  "USWT",
										Trend:      "rise",
										Percentage: percentStr2Decimal(in.Maximum),
									},
								},
							},
						},
					},
				},
			},
		}
		return nil
	}
	return ErrFunctionNotDefined
}
