package request_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/shootercheng/gdb-translate/internal/parser"
	"github.com/shootercheng/gdb-translate/internal/request"
	"golang.org/x/net/html"
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

func TestCollectTextNodesWithBatchLLmRequest(t *testing.T) {
	filePath := "/home/scd/code/c-lang/gdb-17.2/gdb/doc/gdb/Summary.html"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read file error: %s", err.Error())
	}
	htmlContent := string(fileContent)
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatal(err)
	}

	var textNodes []*html.Node
	skip := map[string]bool{
		"script":   true,
		"style":    true,
		"noscript": true,
		"meta":     true,
		"link":     true,
		"small":    true,
		"code":     true,
	}
	parser.CollectTextNodes(doc, &textNodes, skip)
	texts := make(map[int]string)
	idx := 0
	for _, node := range textNodes {
		if trimmed := strings.TrimSpace(node.Data); trimmed != "" {
			texts[idx] = trimmed
			idx++
		}
	}

	jsonData, err := json.Marshal(texts)
	if err != nil {
		fmt.Println("JSON 序列化失败:", err)
		t.Fatal(err)
	}
	reqBody := request.ChatRequest{
		Model: "doubao-1-5-lite-32k-250115",
		Messages: []request.Message{
			{
				Role:    "system",
				Content: "您是一个专业的翻译引擎。请将用户输入的json翻译成中文, 输入示例：{\"1\": \"test\"}, 输出示例：{\"1\": \"测试\"}, 保证所有的json英文值都翻译完成，符号不需要翻译",
			},
			{
				Role:    "user",
				Content: string(jsonData),
			},
		},
		Stream: false,
	}
	result, err := request.RequestLlm(&reqBody)
	if err != nil {
		t.Fatalf("request llm error %s", err.Error())
	}
	fmt.Printf("request result %v \n", result)
	if len(result.Choices) == 0 {
		t.Fatalf("request llm error %s", "result empty")
	}
	res := result.Choices[0].Message.Content
	fmt.Printf("translate res %s \n", res)
	var translatedTexts map[int]string
	err = json.Unmarshal([]byte(res), &translatedTexts)
	if err != nil {
		t.Fatalf("deseralize json error %s", err)
	}
	// 写回，需要保持对应关系（跳过空白节点）
	idx = 0
	for _, node := range textNodes {
		if trimmed := strings.TrimSpace(node.Data); trimmed != "" {
			node.Data = strings.Replace(node.Data, trimmed, translatedTexts[idx], 1)
			idx++
		}
	}
	outFile, err := os.Create("test-node.html")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	if err := html.Render(outFile, doc); err != nil {
		log.Fatal(err)
	}
}
