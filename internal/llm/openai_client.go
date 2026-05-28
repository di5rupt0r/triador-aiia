package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

// AnalysisOutput é a struct que define o schema de saída via reflection.
// Exportada para que o jsonschema reflector consiga inspecioná-la.
type AnalysisOutput struct {
	CandidateName   string   `json:"candidate_name"   jsonschema_description:"Full name extracted from the resume"`
	TechnicalSkills []string `json:"technical_skills" jsonschema_description:"List of technical skills identified"`
	YearsExperience float64  `json:"years_experience" jsonschema_description:"Total years of professional experience"`
	FitScore        int      `json:"fit_score"        jsonschema_description:"Fit score from 0 to 100"`
	Summary         string   `json:"summary"          jsonschema_description:"2 to 3 sentences justifying the fit score"`
}

func generateSchema[T any]() map[string]any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)

	data, _ := json.Marshal(schema)
	var result map[string]any
	json.Unmarshal(data, &result)
	return result
}

var analysisSchema = generateSchema[AnalysisOutput]()

type OpenAIClient struct {
	client openai.Client
	model  string
}

func NewOpenAIClient(apiKey, baseURL, model string) *OpenAIClient {
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	return &OpenAIClient{
		client: openai.NewClient(opts...),
		model:  model,
	}
}

func (c *OpenAIClient) Complete(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(c.model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userMessage),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "resume_analysis",
					Schema: analysisSchema,
					Strict: openai.Bool(true),
				},
			},
		},
		Temperature: openai.Float(0.2),
	})
	if err != nil {
		return "", fmt.Errorf("openai completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	return resp.Choices[0].Message.Content, nil
}
