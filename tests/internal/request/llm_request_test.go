package request_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/shootercheng/gdb-translate/internal/request"
)

func TestRequestLlm(t *testing.T) {
	filePath := "/home/scd/code/c-lang/gdb-17.2/gdb/doc/gdb/Summary.html"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read file error: %s", err.Error())
	}
	strContent := string(fileContent)
	reqBody := request.ChatRequest{
		Model: "doubao-1-5-lite-32k-250115",
		Messages: []request.Message{
			{
				Role:    "system",
				Content: "您是一个专业的翻译引擎。请将用户输入的英文html翻译成中文，不要翻译html标签。保持原来的格式输出html",
			},
			{
				Role:    "user",
				Content: strContent,
			},
		},
		Stream: false,
	}
	result, err := request.RequestLlm(&reqBody)
	if err != nil {
		t.Fatalf("request llm error %s", err.Error())
	}
	fmt.Printf("request result %v", result)
	if len(result.Choices) > 0 {
		res := result.Choices[0].Message.Content
		err = os.WriteFile("test.html", []byte(res), 0644)
		if err != nil {
			fmt.Printf("write file error:%s\n", err.Error())
		}
	} else {
		t.Fatal("not found translate result")
	}
}
