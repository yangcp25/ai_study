package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"softgen/configs"
	"softgen/internal/llm"
	"softgen/internal/pkg"
	"strings"
)

func Run(name, typ string, model uint) error {
	switch typ {
	case "manual":
		fmt.Println("生成操作说明书:", name)
		err := GenerateManual(context.Background(), name, "prompts/manual_v1.txt", "template/doc.txt",
			"/Users/ycp/work/code/own/game/dy/docs_editor2pdf/template/md", os.Getenv("DEEPSEEK_API_KEY"), model)
		if err != nil {
			return err
		}
	case "code":
		fmt.Println("生成代码:", name)
		err := generateCode(context.Background(), "prompts/code_v1.txt", "template/code.txt",
			"/Users/ycp/work/code/own/game/dy/docs_editor2pdf/template/md", os.Getenv("DEEPSEEK_API_KEY"), model)
		if err != nil {
			return err
		}
	case "all":
		fmt.Println("生成操作说明书:", name)
		err := GenerateManual(context.Background(), name, "prompts/manual_v1.txt", "template/doc.txt",
			"/Users/ycp/work/code/own/game/dy/docs_editor2pdf/template/md", os.Getenv("DEEPSEEK_API_KEY"), model)
		if err != nil {
			return err
		}

		fmt.Println("生成代码:", name)
		err = generateCode(context.Background(), "prompts/code_v1.txt", "template/code.txt",
			"/Users/ycp/work/code/own/game/dy/docs_editor2pdf/template/md", os.Getenv("DEEPSEEK_API_KEY"), model)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown type: " + typ)
	}
	return nil
}

// ManualAST 定义（可按需扩展）
type ManualAST struct {
	Title        string    `json:"title"`
	SoftwareName string    `json:"software_name"`
	Sections     []Section `json:"sections"`
}

type Section struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 完成你要求的最小闭环
// ctx: 上下文
// softwareName: 软件名称
// promptPath: 本地 prompt 文件路径
// examplePath: 本地 example 模板路径（参考说明书）
// outputPath: 输出 docx 路径，例如 "xx小游戏_说明书.docx"
// deepseekAPIKey: DeepSeek API Key（或传空以 mock）
func GenerateManual(ctx context.Context, softwareName, promptPath, examplePath, outputPath, deepseekAPIKey string, model uint) error {
	// 1. 读 prompt & example
	promptTmpl, err := configs.Config.ReadFile(promptPath)
	if err != nil {
		return fmt.Errorf("read prompt: %w", err)
	}
	example, err := configs.Config.ReadFile(examplePath)
	if err != nil {
		return fmt.Errorf("read example: %w", err)
	}

	// 2. 组合 prompt（你可以改这里的拼接逻辑或使用 text/template）
	prompt := pkg.BuildPrompt(string(promptTmpl), string(example), softwareName)

	// 3. 调用 LLM（DeepSeek），如果 deepseekAPIKey 为空则会返回 mock 示例（方便本地测试）
	md, err := llm.CallDeepSeekStream(ctx, prompt, deepseekAPIKey, model)
	if err != nil {
		return fmt.Errorf("call llm: %w", err)
	}
	if strings.TrimSpace(md) == "" {
		return errors.New("empty markdown from llm")
	}

	// 4. 保存 Markdown
	return pkg.SaveMarkdown(pkg.FileTypeDoc, md, outputPath)
}

// generateCode 生成代码目录及Main文件 (新增逻辑)
func generateCode(ctx context.Context, promptPath, examplePath, outputPath, deepseekAPIKey string, model uint) error {
	// 1. 读 prompt & example
	promptTmpl, err := configs.Config.ReadFile(promptPath)
	if err != nil {
		return fmt.Errorf("read prompt: %w", err)
	}
	example, err := configs.Config.ReadFile(examplePath)
	if err != nil {
		return fmt.Errorf("read example: %w", err)
	}

	// 2. 组合 prompt（你可以改这里的拼接逻辑或使用 text/template）
	prompt := pkg.BuildCodePrompt(string(promptTmpl), string(example))

	// 3. 调用 LLM（DeepSeek），如果 deepseekAPIKey 为空则会返回 mock 示例（方便本地测试）
	md, err := llm.CallDeepSeekStream(ctx, prompt, deepseekAPIKey, model)
	if err != nil {
		return fmt.Errorf("call llm: %w", err)
	}
	if strings.TrimSpace(md) == "" {
		return errors.New("empty markdown from llm")
	}

	// 4. 保存 Markdown
	return pkg.SaveMarkdown(pkg.FileTypeCode, md, outputPath)
}
