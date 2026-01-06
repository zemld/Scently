package fetching

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
	"github.com/zemld/config-manager/pkg/cm"
)

const (
	systemPrompt = `You are an expert perfume recommender. 
             Your task is to suggest alternative perfumes 
             based on the user's favorite one. 
             Do not include explanations, text, or markdown — only JSON.`
	userPrompt = `User's favorite perfume:
            Brand: %s
            Name: %s
            Sex: %s

            Return exactly 4 other perfumes that the user might also like.
            Each item must include:
            - brand
            - name
            - sex (must be one of: "male", "female", "unisex")

            Sex constraint based on user's Sex:
            - if user's Sex is "unisex": returned items' sex MUST be only "unisex"
            - if user's Sex is "male": returned items' sex MUST be one of: "male", "unisex"
            - if user's Sex is "female": returned items' sex MUST be one of: "female", "unisex"
            Allowed for this request: %v

            Respond strictly in this JSON format:
            [
            {{"brand": "string", "name": "string", "sex": "string"}},
            {{"brand": "string", "name": "string", "sex": "string"}},
            {{"brand": "string", "name": "string", "sex": "string"}},
            {{"brand": "string", "name": "string", "sex": "string"}}
            ]`
)

type requestBody struct {
	ModelUri          string            `json:"modelUri"`
	Messages          []message         `json:"messages"`
	CompletionOptions completionOptions `json:"completionOptions"`
	JsonSchema        jsonSchema        `json:"jsonSchema"`
}

type message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type completionOptions struct {
	MaxTokens   int     `json:"maxTokens"`
	Temperature float64 `json:"temperature"`
	Stream      bool    `json:"stream"`
}

type jsonSchema struct {
	Schema schema `json:"schema"`
}

type schema struct {
	Type_ string `json:"type"`
	Items items  `json:"items"`
}

type items struct {
	Type_                string     `json:"type"`
	Properties           properties `json:"properties"`
	Required             []string   `json:"required"`
	AdditionalProperties bool       `json:"additionalProperties"`
}

type properties struct {
	Brand valueType `json:"brand"`
	Name  valueType `json:"name"`
	Sex   valueType `json:"sex"`
}

type valueType struct {
	Type_ string   `json:"type"`
	Enum  []string `json:"enum,omitempty"`
}

type aiCompletionResponse struct {
	Result struct {
		Alternatives []struct {
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
		} `json:"alternatives"`
	} `json:"result"`
}

type AI struct {
	url       string
	folderId  string
	modelName string
	apiKey    string
	client    *http.Client
	cm        cm.ConfigManager
}

func NewAI(url string, folderId string, modelName string, apiKey string, cm cm.ConfigManager) *AI {
	return &AI{
		url:       url,
		folderId:  folderId,
		modelName: modelName,
		apiKey:    apiKey,
		client:    http.DefaultClient,
		cm:        cm,
	}
}

func (f *AI) Fetch(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool) {
	if len(params) == 0 {
		return nil, false
	}

	ctx, cancel := context.WithTimeout(ctx, f.cm.GetDurationWithDefault("ai_fetcher_timeout", 20*time.Second))
	defer cancel()

	r, err := f.createRequest(ctx, params[0])
	if err != nil {
		return nil, false
	}

	response, err := f.client.Do(r)
	if err != nil {
		return nil, false
	}
	defer response.Body.Close()

	perfumes, err := f.tryParseResponse(response)
	if err != nil {
		return nil, false
	}
	return perfumes, true
}

func (f *AI) createRequest(ctx context.Context, perfume parameters.RequestPerfume) (*http.Request, error) {
	body := requestBody{
		ModelUri: f.createModelUri(f.folderId, f.modelName),
		Messages: []message{
			{Role: "system", Text: systemPrompt},
			{Role: "user", Text: fmt.Sprintf(userPrompt, perfume.Brand, perfume.Name, perfume.Sex, f.getAllowedSexes(perfume.Sex))},
		},
		CompletionOptions: completionOptions{
			MaxTokens:   500,
			Temperature: 0.4,
			Stream:      false,
		},
		JsonSchema: jsonSchema{
			Schema: schema{
				Type_: "array",
				Items: items{
					Type_: "object",
					Properties: properties{
						Brand: valueType{Type_: "string"},
						Name:  valueType{Type_: "string"},
						Sex:   valueType{Type_: "string", Enum: f.getAllowedSexes(perfume.Sex)},
					},
					Required:             []string{"brand", "name", "sex"},
					AdditionalProperties: false,
				},
			},
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequestWithContext(ctx, "POST", f.url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", f.apiKey))
	return r, nil
}

func (f *AI) createModelUri(folderId string, modelName string) string {
	return fmt.Sprintf("gpt://%s/%s", folderId, modelName)
}

func (f *AI) getAllowedSexes(sex string) []string {
	allowed := []string{parameters.SexUnisex}
	if sex == parameters.SexMale {
		allowed = append(allowed, parameters.SexMale)
	}
	if sex == parameters.SexFemale {
		allowed = append(allowed, parameters.SexFemale)
	}
	return allowed
}

func (f *AI) tryParseResponse(response *http.Response) ([]models.Perfume, error) {
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad request to llm service: %d", response.StatusCode)
	}

	var completionResponse aiCompletionResponse
	if err := json.NewDecoder(response.Body).Decode(&completionResponse); err != nil {
		return nil, err
	}
	if len(completionResponse.Result.Alternatives) == 0 {
		return nil, fmt.Errorf("no alternatives in response")
	}

	var perfumes []models.Perfume
	if err := json.Unmarshal([]byte(completionResponse.Result.Alternatives[0].Message.Text), &perfumes); err != nil {
		return nil, err
	}
	return perfumes, nil
}
