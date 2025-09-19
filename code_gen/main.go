package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// OllamaRequest 请求结构体
type OllamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options"`
}

// OllamaResponse 响应结构体
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// CodeGenerator 代码生成器结构体
type CodeGenerator struct {
	ModelName string
	BaseURL   string
	Client    *http.Client
}

// NewCodeGenerator 创建新的代码生成器实例
func NewCodeGenerator(modelName, baseURL string) *CodeGenerator {
	if modelName == "" {
		modelName = "deepseek-coder:6.7b"
	}
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	return &CodeGenerator{
		ModelName: modelName,
		BaseURL:   baseURL,
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GenerateCode 生成代码
func (cg *CodeGenerator) GenerateCode(prompt, language string, temperature float64) (string, error) {
	fullPrompt := fmt.Sprintf(`请用%s编写代码，要求如下：
%s

请只返回可执行的代码，不要包含其他解释。代码应该：
1. 包含必要的注释
2. 遵循最佳实践和Go语言规范
3. 处理错误情况
4. 包含简单的使用示例

代码：
`+"```"+language, language, prompt)

	request := OllamaRequest{
		Model:  cg.ModelName,
		Prompt: fullPrompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": temperature,
			"top_p":       0.9,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", cg.BaseURL)
	resp, err := cg.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	return cg.extractCode(ollamaResp.Response), nil
}

// extractCode 提取代码块
func (cg *CodeGenerator) extractCode(response string) string {
	response = strings.TrimSpace(response)

	// 移除markdown代码块标记
	re := regexp.MustCompile("(?s)```[a-zA-Z]*\\n?(.*)```")
	matches := re.FindStringSubmatch(response)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return response
}

// GenerateStruct 生成Go结构体
func (cg *CodeGenerator) GenerateStruct(structName, description string, fields []string) (string, error) {
	fieldsStr := "基本字段"
	if len(fields) > 0 {
		fieldsStr = strings.Join(fields, "、")
	}

	prompt := fmt.Sprintf("创建一个名为 %s 的Go结构体。功能：%s。包含字段：%s。请添加JSON标签和必要的验证标签。",
		structName, description, fieldsStr)

	return cg.GenerateCode(prompt, "go", 0.2)
}

// GenerateHandler 生成HTTP处理器
func (cg *CodeGenerator) GenerateHandler(handlerName, method, path, description string) (string, error) {
	prompt := fmt.Sprintf("创建一个Go HTTP处理器函数 %s，处理 %s %s 请求。功能：%s。使用标准库net/http或gin框架。",
		handlerName, method, path, description)

	return cg.GenerateCode(prompt, "go", 0.2)
}

// GenerateService 生成服务层代码
func (cg *CodeGenerator) GenerateService(serviceName, description string, methods []string) (string, error) {
	methodsStr := "基本的CRUD方法"
	if len(methods) > 0 {
		methodsStr = strings.Join(methods, "、")
	}

	prompt := fmt.Sprintf("创建一个Go服务层接口和实现 %s。功能：%s。包含方法：%s。请添加错误处理和日志记录。",
		serviceName, description, methodsStr)

	return cg.GenerateCode(prompt, "go", 0.2)
}

// GenerateRepository 生成数据访问层
func (cg *CodeGenerator) GenerateRepository(repoName, description, dbType string) (string, error) {
	prompt := fmt.Sprintf("创建一个Go数据访问层 %s，使用%s数据库。功能：%s。包含基本的CRUD操作，使用database/sql或GORM。",
		repoName, dbType, description)

	return cg.GenerateCode(prompt, "go", 0.2)
}

// GenerateMiddleware 生成中间件
func (cg *CodeGenerator) GenerateMiddleware(middlewareName, description string) (string, error) {
	prompt := fmt.Sprintf("创建一个Go HTTP中间件 %s。功能：%s。兼容标准库和gin框架。",
		middlewareName, description)

	return cg.GenerateCode(prompt, "go", 0.2)
}

// GenerateTest 生成测试代码
func (cg *CodeGenerator) GenerateTest(functionName, description string) (string, error) {
	prompt := fmt.Sprintf("为Go函数 %s 生成完整的单元测试。功能：%s。使用testing包，包含正常情况和边界情况的测试。",
		functionName, description)

	return cg.GenerateCode(prompt, "go", 0.1)
}

// CheckOllamaService 检查Ollama服务状态
func (cg *CodeGenerator) CheckOllamaService() error {
	url := fmt.Sprintf("%s/api/tags", cg.BaseURL)
	resp, err := cg.Client.Get(url)
	if err != nil {
		return fmt.Errorf("连接Ollama服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama服务响应异常: HTTP %d", resp.StatusCode)
	}

	return nil
}

func main() {
	generator := NewCodeGenerator("deepseek-coder:6.7b", "")

	fmt.Println("=== Go版本本地AI代码生成器演示 ===\n")

	// 检查Ollama服务
	if err := generator.CheckOllamaService(); err != nil {
		fmt.Printf("❌ %v\n", err)
		fmt.Println("请检查:")
		fmt.Println("1. 是否已安装Ollama")
		fmt.Println("2. 是否运行了: ollama serve")
		fmt.Println("3. 是否下载了代码生成模型: ollama run deepseek-coder:6.7b")
		return
	}

	fmt.Println("✅ Ollama服务正在运行\n")

	// 示例1：生成结构体
	fmt.Println("1. 生成用户结构体：")
	structCode, err := generator.GenerateStruct(
		"User",
		"用户信息管理",
		[]string{"ID", "Username", "Email", "CreatedAt"},
	)
	if err != nil {
		fmt.Printf("生成结构体失败: %v\n", err)
	} else {
		fmt.Println(structCode)
	}
	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例2：生成HTTP处理器
	fmt.Println("2. 生成用户API处理器：")
	handlerCode, err := generator.GenerateHandler(
		"GetUsersHandler",
		"GET",
		"/api/users",
		"获取用户列表，支持分页和搜索",
	)
	if err != nil {
		fmt.Printf("生成处理器失败: %v\n", err)
	} else {
		fmt.Println(handlerCode)
	}
	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例3：生成服务层
	fmt.Println("3. 生成用户服务：")
	serviceCode, err := generator.GenerateService(
		"UserService",
		"用户业务逻辑处理",
		[]string{"CreateUser", "GetUser", "UpdateUser", "DeleteUser", "ListUsers"},
	)
	if err != nil {
		fmt.Printf("生成服务失败: %v\n", err)
	} else {
		fmt.Println(serviceCode)
	}
	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例4：生成中间件
	fmt.Println("4. 生成JWT认证中间件：")
	middlewareCode, err := generator.GenerateMiddleware(
		"JWTAuthMiddleware",
		"JWT令牌验证，支持从Header或Cookie中获取token",
	)
	if err != nil {
		fmt.Printf("生成中间件失败: %v\n", err)
	} else {
		fmt.Println(middlewareCode)
	}
}
