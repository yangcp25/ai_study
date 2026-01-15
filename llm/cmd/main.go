package main

import (
	"context"
	"encoding/json"
	"fmt"
	"llm/internal/llm"
	"os"
)

type ROIAnswer struct {
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

func main() {
	client := llm.NewClient(os.Getenv("DEEPSEEK_API_KEY"))

	resp, err := client.Chat(context.Background(), llm.ChatReq{
		Prompt: `
你是一个广告投放策略助手。

请严格以 JSON 格式回答，不要输出任何多余内容。

JSON 格式如下：
{
  "definition": string,
  "example": string
}

问题：什么是 ROI？
`,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Content)

	var ans ROIAnswer
	if err := json.Unmarshal([]byte(resp.Content), &ans); err != nil {
		panic(err)
	}

	fmt.Println("定义:", ans.Definition)
	fmt.Println("例子:", ans.Example)
}
