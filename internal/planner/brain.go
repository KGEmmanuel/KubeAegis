package planner

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	// CORRECTION ICI : Ajout de "-ext" dans les chemins
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
)

// Brain is the cognitive engine wrapper
type Brain struct {
	chatModel model.ChatModel
}

// NewBrain initializes the connection to the LLM.
// If OPENAI_API_KEY is set, it uses OpenAI (or DeepSeek).
// Otherwise, it defaults to Local Ollama (Model: llama3).
func NewBrain(ctx context.Context) (*Brain, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	var chat model.ChatModel
	var err error

	if apiKey != "" {
		// Use Cloud LLM (OpenAI / DeepSeek)
		chat, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
			APIKey: apiKey,
			Model:  "gpt-4o",
		})
	} else {
		// Use Local LLM (Ollama) - FREE
		fmt.Println("🧠 No API Key found. Using Local Ollama (llama3)...")
		chat, err = ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
			BaseURL: "http://localhost:11434",
			Model:   "llama3",
		})
	}

	if err != nil {
		return nil, err
	}
	return &Brain{chatModel: chat}, nil
}

// GenerateManifest asks the LLM to convert intent to K8s YAML
func (b *Brain) GenerateManifest(ctx context.Context, intent string) (string, error) {
	systemPrompt := `You are a Kubernetes Expert. 
    Your ONLY output must be valid Kubernetes YAML. 
    Do not add markdown formatting, explanations, or quotes.
    Just the YAML.`

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(fmt.Sprintf("Task: %s", intent)),
	}

	resp, err := b.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}
