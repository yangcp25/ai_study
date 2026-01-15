package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	apiKey string
	http   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Chat(ctx context.Context, req ChatReq) (ChatResp, error) {
	payload := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "user", "content": req.Prompt},
		},
	}

	body, _ := json.Marshal(payload)

	httpReq, _ := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.deepseek.com/chat/completions",
		bytes.NewBuffer(body),
	)

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return ChatResp{}, err
	}
	defer resp.Body.Close()

	var raw struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return ChatResp{}, err
	}

	return ChatResp{Content: raw.Choices[0].Message.Content}, nil
}
