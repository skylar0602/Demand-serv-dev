package llm

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sashabaranov/go-openai/jsonschema"

	"github.com/sashabaranov/go-openai"

	"golang.org/x/net/context"

	"github.com/smarterwallet/demand-abstraction-serv/config"
)

var (
	cfg = &config.AiConfig{
		Endpoint: "",
		Model:    openai.GPT3Dot5Turbo,
		APIKey:   "",
	}
	ctx    = context.Background()
	prompt = `I want you act as an experienced cryptocurrency investor, here are my investment expectations: %s. 
			Give me the advice investment rate of return range. Mention that if the input is irrelevant about cryptocurrency invest, reply with zero.`
	functions = []openai.FunctionDefinition{{
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
			},
			Required: []string{"minimum", "maximum"},
		},
	}}
)

type input struct {
	Maximum string `json:"maximum"`
	Minimum string `json:"minimum"`
}

func get_trade_to_earn_strategy(input input) {
	fmt.Printf("get_trade_to_earn_strategy range:%s-%s\n", input.Minimum, input.Maximum)
}

func TestChat(t *testing.T) {
	client := NewOpenAI(cfg)
	t.Run("high return", func(t *testing.T) {
		name, args, err := client.Chat(ctx, prompt, "I want high return with MATIC", functions)
		assert.Nil(t, err)
		assert.Equal(t, name, "get_trade_to_earn_strategy")
		in := input{}
		err = json.Unmarshal([]byte(args), &in)
		assert.Nil(t, err)
		get_trade_to_earn_strategy(in)
	})
	t.Run("low return", func(t *testing.T) {
		name, args, err := client.Chat(ctx, prompt, "I want low return with MATIC", functions)
		assert.Nil(t, err)
		assert.Equal(t, name, "get_trade_to_earn_strategy")
		in := input{}
		err = json.Unmarshal([]byte(args), &in)
		assert.Nil(t, err)
		get_trade_to_earn_strategy(in)
	})
	t.Run("irrelevant", func(t *testing.T) {
		name, args, err := client.Chat(ctx, prompt, "what's the weather like today", functions)
		assert.Nil(t, err)
		assert.Equal(t, name, "get_trade_to_earn_strategy")
		in := input{}
		err = json.Unmarshal([]byte(args), &in)
		assert.Nil(t, err)
		get_trade_to_earn_strategy(in)
	})
}
