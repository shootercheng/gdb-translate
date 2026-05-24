package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// 定义请求结构体
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	Stream      bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 定义响应结构体 (根据提供的 JSON 结果简化)
type ChatResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role             string      `json:"role"`
			Content          string      `json:"content"`
			ToolCalls        interface{} `json:"tool_calls"`
			ReasoningContent interface{} `json:"reasoning_content"`
		} `json:"message"`
		FinishReason string      `json:"finish_reason"`
		Logprobs     interface{} `json:"logprobs"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
}

func RequestLlm(chatRequest *ChatRequest) (*ChatResponse, error) {
	// 将结构体序列化为 JSON
	jsonData, err := json.Marshal(chatRequest)
	if err != nil {
		fmt.Println("JSON 序列化失败:", err)
		return nil, err
	}
	fmt.Printf("request json %s \n", jsonData)

	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	url := os.Getenv("LLM_REQUEST_URL")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LLM_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return nil, err
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("请求失败，状态码: %d, 响应内容: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("request error staus code : %d body : %s", resp.StatusCode, string(body))
	}

	var result ChatResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("解析响应 JSON 失败:", err)
		fmt.Println("原始响应:", string(body))
		return nil, err
	}
	return &result, nil
}
