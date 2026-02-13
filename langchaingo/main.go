package main

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
)

func main() {
	ctx := context.Background()
	llm, _ := ollama.New(ollama.WithModel("deepseek-r1:7b"))

	// 1. 创建一个简单的对话内存
	mem := memory.NewConversationBuffer()

	// 2. 创建一个对话链
	executor := chains.NewConversation(llm, mem)

	// 第一轮对话
	res1, _ := chains.Predict(ctx, executor, map[string]any{"input": "你好，我是程序员小明。"})
	fmt.Println("AI:", res1)

	// 第二轮对话（AI 会记得你叫小明）
	res2, _ := chains.Predict(ctx, executor, map[string]any{"input": "请问我叫什么名字？"})
	fmt.Println("AI:", res2)
}
