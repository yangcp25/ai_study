package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"softgen/internal/pkg"
	"strings"
)

func CallDeepSeekStream(ctx context.Context, prompt, apiKey string, model uint) (string, error) {
	if strings.TrimSpace(apiKey) == "" {
		return "", errors.New("deepseek api key empty")
	}

	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	payload := map[string]interface{}{
		"model": pkg.ChatModel(model).String(),
		"messages": []Message{
			{Role: "user", Content: prompt},
		},
		"stream":            true,
		"max_tokens":        8192 * 3, // 必须拉满
		"temperature":       0.7,      // 略微提高，增加“话痨”程度
		"presence_penalty":  0.5,      // 强制开拓新内容，避免重复
		"frequency_penalty": 0.3,      // 减少词汇重复
		"top_p":             0.95,     // 保证文笔流畅
	}

	bs, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.deepseek.com/v1/chat/completions",
		bytes.NewReader(bs),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("deepseek status=%d body=%s", resp.StatusCode, body)
	}

	reader := bufio.NewReader(resp.Body)
	var buf strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		line = strings.TrimSpace(line)

		// DeepSeek / OpenAI stream 结束标志
		if line == "data: [DONE]" {
			break
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 {
			buf.WriteString(chunk.Choices[0].Delta.Content)
		}
	}

	return buf.String(), nil
}
