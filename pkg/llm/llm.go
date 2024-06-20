package llm

import (
	"context"

	log "github.com/cihub/seelog"

	"github.com/pkg/errors"

	"github.com/sashabaranov/go-openai"
	"github.com/smarterwallet/demand-abstraction-serv/config"
)

type Llm interface {
	Chat(ctx context.Context, prompt, content string, functions []openai.FunctionDefinition) (string, string, error)
}

type OpenAI struct {
	client *openai.Client
	model  string
}

func NewOpenAI(cfg *config.AiConfig) *OpenAI {
	if cfg.APIKey == "" || cfg.Model == "" {
		panic("missing apikey or model")
	}
	openaiCfg := openai.DefaultConfig(cfg.APIKey)
	if cfg.Endpoint != "" {
		openaiCfg.BaseURL = cfg.Endpoint
	}
	return &OpenAI{client: openai.NewClientWithConfig(openaiCfg), model: cfg.Model}
}

func (a *OpenAI) Chat(ctx context.Context, prompt, content string, functions []openai.FunctionDefinition) (string, string, error) {
	tools := make([]openai.Tool, 0, len(functions))
	for _, function := range functions {
		tools = append(tools, openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: function,
		})
	}
	request := openai.ChatCompletionRequest{
		Model: a.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: prompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
		Tools: tools,
	}
	resp, err := a.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", "", err
	}
	if len(resp.Choices) == 0 || len(resp.Choices[0].Message.ToolCalls) == 0 {
		log.Errorf("empty response resp=%+v", resp)
		return "", "", errors.New("empty response")
	}
	call := resp.Choices[0].Message.ToolCalls[0].Function
	return call.Name, call.Arguments, nil
}

type MockOpenAI struct {
}

func NewMockOpenAI() *OpenAI {
	return &OpenAI{}
}

func (a *MockOpenAI) Chat(ctx context.Context, prompt string, functions []openai.FunctionDefinition) (string, string, error) {
	return "mock name", "{}", nil
}
