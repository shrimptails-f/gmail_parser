package openai

import (
	cd "business/internal/common/domain"
	"context"
	"encoding/json"
	"log"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type client struct {
	sdk openai.Client
}

func New(apiKey string) ClientInterFace {
	return &client{
		sdk: openai.NewClient(
			option.WithAPIKey(apiKey),
		),
	}
}

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	return reflector.Reflect(v)
}

// type AnalysisResults struct {
// 	Items []cd.AnalysisResult `json:"results" jsonschema_description:"分析結果の配列"`
// }

func (c *client) Chat(ctx context.Context, prompt string) ([]cd.AnalysisResult, error) {
	// schema := GenerateSchema[AnalysisResults]()

	// schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
	// 	Name:        "email_analysis_result",
	// 	Description: openai.String("メールの構造化分析結果"),
	// 	Schema:      schema,
	// 	Strict:      openai.Bool(true),
	// }

	resp, err := c.sdk.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4_1Mini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		// 低コストバージョンを使いたいので指定できない。
		// ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
		// 	OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
		// 		JSONSchema: schemaParam,
		// 	},
		// },
	})
	if err != nil {
		return nil, err
	}
	raw := resp.Choices[0].Message.Content

	var results []cd.AnalysisResult
	if err := json.Unmarshal([]byte(raw), &results); err != nil {
		log.Printf("構造エラー: JSON→構造体変換失敗:\n%s\nエラー: %v", raw, err)
		return nil, err
	}

	return results, nil
}
