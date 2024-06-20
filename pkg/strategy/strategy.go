package strategy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/smarterwallet/demand-abstraction-serv/model"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type IStrategy interface {
	Prompt() string
	Functions() []openai.FunctionDefinition
	Render(resp *model.DemandResponse, name, args string) error
}

var (
	ErrFunctionNotDefined           = errors.New("function not defined")
	_                     IStrategy = &transfer{}
	_                     IStrategy = &trade2Earn{}
	_                     IStrategy = &crossChain{}
	_                     IStrategy = &chainAbstraction{}
	_                     IStrategy = &selectStrategy{}
	strategy                        = map[string]IStrategy{
		"trade2Earn":            trade2Earn{},
		"transfer":              transfer{},
		"crossChain":            crossChain{},
		"crossChainAbstraction": chainAbstraction{},
		"selectStrategy":        selectStrategy{},
	}
)

func MatchStrategy(category string, ctx *model.CtxRequest) (IStrategy, error) {
	if category == "transfer" {
		return transfer{balance: ctx}, nil
	}
	st, ok := strategy[category]
	if !ok {
		return nil, errors.New("strategy not support")
	}
	return st, nil
}

func percentStr2Decimal(percentageStr string) string {
	percentageStr = strings.TrimRight(percentageStr, "%")
	percentage, err := strconv.ParseFloat(percentageStr, 64)
	if err != nil {
		return ""
	}
	decimal := percentage / 100
	return fmt.Sprintf("%.2f", decimal)
}

type selectStrategy struct{}

func (s selectStrategy) Prompt() string {
	return `As an experienced cryptocurrency investor, I'd like you to analyze user's operations. 
	Based on these, choose the most suitable strategy from the two available options: trade2Earn or transfer.
	If blockchain chains are detected in user's demand, such as token transfers among Ethereum, Goerli, Fuji etc, the strategy is transfer.	
	trade2Earn focuses on user's financial investments such as High/Low Return expectations.`
}

func (s selectStrategy) Functions() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{{
		Name: "select_strategy",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"strategy": {
					Type:        jsonschema.String,
					Description: "The matched strategy, e.g. transfer, trade2Earn",
				},
			},
			Required: []string{"strategy"},
		},
	}}
}

func (s selectStrategy) Render(resp *model.DemandResponse, name, args string) error {
	return nil
}
