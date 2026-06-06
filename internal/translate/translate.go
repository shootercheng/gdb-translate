package translate

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/shootercheng/gdb-translate/internal/parser"
	"github.com/shootercheng/gdb-translate/internal/request"
	"golang.org/x/net/html"
)

type TranslateParam struct {
	InputPath        string
	FileName         string
	RequestBatchSize int
	OutputPath       string
}

func ParseHtmlFileToNodes(htmlPath string) (*html.Node, []*html.Node, error) {
	fileContent, err := os.ReadFile(htmlPath)
	if err != nil {
		fmt.Printf("read file error: %s \n", err.Error())
		return nil, nil, err
	}
	htmlContent := string(fileContent)
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Printf("parse html file error: %s \n", err.Error())
		return nil, nil, err
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
	return doc, textNodes, nil
}

func ConvertRequestParam(textNodes []*html.Node, requestBatchSize int) ([]map[int]string, map[int]string) {
	batchData := make([]map[int]string, 0)
	idx := 0
	batchCount := 0
	batchMap := make(map[int]string)
	allOriginMap := make(map[int]string, len(batchData))
	for _, node := range textNodes {
		if trimmed := strings.TrimSpace(node.Data); trimmed != "" {
			batchMap[idx] = trimmed
			allOriginMap[idx] = trimmed
			idx++
			batchCount++
			if batchCount >= requestBatchSize {
				batchData = append(batchData, batchMap)
				batchMap = make(map[int]string)
				batchCount = 0
			}
		}
	}
	if batchCount > 0 {
		batchData = append(batchData, batchMap)
	}
	return batchData, allOriginMap
}

func RequestLlmBatch(batchData []map[int]string, allOriginMap map[int]string) (map[int]string, *request.Usage, error) {
	allTranslatedData := make(map[int]string, len(allOriginMap))
	allUsage := request.Usage{}
	for _, item := range batchData {
		jsonData, err := json.Marshal(item)
		if err != nil {
			fmt.Println("JSON marshal error:", err)
			return nil, nil, err
		}
		reTry := 5
		for i := range reTry {
			fmt.Printf("the %d request to the large model interface\n", i+1)
			reqBody := request.ChatRequest{
				Model: os.Getenv("LLM_MODEL_ID"),
				Messages: []request.Message{
					{
						Role:    "system",
						Content: "您是一个专业的翻译引擎。请将用户输入的json翻译成中文, 输入示例：{\"1\": \"test\"}, 输出示例：{\"1\": \"测试\"}, 保证所有的json英文值都翻译完成，符号不需要翻译。请严格按照json格式输出，程序需要解析输出的json",
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
				fmt.Printf("request llm error: %s\n", err.Error())
				continue
			}
			fmt.Printf("request result %v \n", result)
			if len(result.Choices) == 0 {
				fmt.Println("request llm result empty")
				continue
			}
			res := result.Choices[0].Message.Content
			fmt.Printf("translate res %s \n", res)
			var translatedTexts map[int]string
			err = json.Unmarshal([]byte(res), &translatedTexts)
			if err != nil {
				fmt.Printf("json unmarshal error: %s\n", err.Error())
				continue
			}
			maps.Copy(allTranslatedData, translatedTexts)
			allUsage.CompletionTokens += result.Usage.CompletionTokens
			allUsage.PromptTokens += result.Usage.CompletionTokens
			allUsage.TotalTokens += result.Usage.TotalTokens
			break
		}
	}
	return allTranslatedData, &allUsage, nil
}

func TranslateHtmlFile(translateFileParam TranslateParam) {
	inputPath := translateFileParam.InputPath
	fileName := translateFileParam.FileName
	fmt.Printf("traneslate %s \n", fileName)
	htmlPath := inputPath + string(os.PathSeparator) + fileName
	doc, textNodes, err := ParseHtmlFileToNodes(htmlPath)
	if err != nil {
		fmt.Printf("parse html file error:%s\n", err.Error())
		return
	}
	batchData, allOriginMap := ConvertRequestParam(textNodes, translateFileParam.RequestBatchSize)
	allTranslatedData, allUsage, err := RequestLlmBatch(batchData, allOriginMap)
	if err != nil {
		fmt.Printf("request llm error:%s", err.Error())
	}
	idx := 0
	for _, node := range textNodes {
		if trimmed := strings.TrimSpace(node.Data); trimmed != "" {
			val, ok := allTranslatedData[idx]
			if !ok {
				val = allOriginMap[idx]
				fmt.Printf("translate missing: %s\n", allOriginMap[idx])
			}
			node.Data = strings.Replace(node.Data, trimmed, val, 1)
			idx++
		}
	}
	outputHtmlPath := translateFileParam.OutputPath + string(os.PathSeparator) + fileName
	outFile, err := os.Create(outputHtmlPath)
	if err != nil {
		fmt.Printf("create file error: %s\n", err.Error())
		return
	}
	defer outFile.Close()

	if err := html.Render(outFile, doc); err != nil {
		fmt.Printf("render html file error: %s\n", err.Error())
		return
	}
	outputTokePath := outputHtmlPath + ".json"
	jsonData, err := json.Marshal(allUsage)
	if err != nil {
		fmt.Printf("json Marshal error:%s\n", err.Error())
		return
	}
	err = os.WriteFile(outputTokePath, jsonData, 0644)
	if err != nil {
		fmt.Printf("write json file error:%s\n", err.Error())
	}
}
