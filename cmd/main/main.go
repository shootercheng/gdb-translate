package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/shootercheng/gdb-translate/internal/request"
)

func main() {
	requireEnvList := []string{"LLM_REQUEST_URL", "LLM_API_KEY", "LLM_MODEL_ID"}
	for _, envKey := range requireEnvList {
		envVal := os.Getenv(envKey)
		if envVal == "" {
			fmt.Printf("env key %s not config\n", envKey)
			return
		}
	}
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("get pwd error: %s \n", err.Error())
		return
	}
	outputPath := pwd + string(os.PathSeparator) + "output"
	_, err = os.Stat(outputPath)
	translatedMap := make(map[string]int)
	if os.IsNotExist(err) {
		err = os.MkdirAll(outputPath, 0755)
		if err != nil {
			fmt.Printf("create file dir error: %s", err.Error())
			return
		}
	} else {
		entities, err := os.ReadDir(outputPath)
		if err != nil {
			fmt.Printf("read output dir error: %s \n", err.Error())
			return
		}
		for _, entry := range entities {
			translatedMap[entry.Name()] = 1
		}
	}

	inputPath := "/home/scd/code/c-lang/gdb-17.2/gdb/doc/gdb"
	entities, err := os.ReadDir(inputPath)
	if err != nil {
		fmt.Printf("read dir error: %s \n", err.Error())
		return
	}
	todoFileNames := make([]string, 0, len(entities))
	for _, entry := range entities {
		fileName := entry.Name()
		_, ok := translatedMap[fileName]
		if ok {
			fmt.Printf("file name :%s translated\n", fileName)
			continue
		}
		if strings.HasSuffix(fileName, ".html") {
			todoFileNames = append(todoFileNames, fileName)
		} else {
			fmt.Printf("can't translate %s\n", fileName)
		}
	}

	fmt.Println("to do file names len ", len(todoFileNames))

	var wg sync.WaitGroup
	maxConcurrency := 10
	sem := make(chan int, maxConcurrency)
	for index, fileName := range todoFileNames {
		wg.Go(func() {

			sem <- index
			defer func() {
				<-sem
			}()

			fmt.Printf("traneslate %s \n", fileName)
			htmlPath := inputPath + string(os.PathSeparator) + fileName
			fileContent, err := os.ReadFile(htmlPath)
			if err != nil {
				fmt.Printf("read file error: %s", err.Error())
				return
			}
			strContent := string(fileContent)
			reqBody := request.ChatRequest{
				Model: os.Getenv("LLM_MODEL_ID"),
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
			response, err := request.RequestLlm(&reqBody)
			if err != nil {
				fmt.Printf("reuqest llm error : %s \n", err.Error())
				return
			}
			outputHtmlPath := outputPath + string(os.PathSeparator) + fileName
			if len(response.Choices) > 0 {
				res := response.Choices[0].Message.Content
				err = os.WriteFile(outputHtmlPath, []byte(res), 0644)
				if err != nil {
					fmt.Printf("write html file error:%s\n", err.Error())
				}
			}
			outputTokePath := outputHtmlPath + ".json"
			jsonData, err := json.Marshal(response.Usage)
			if err != nil {
				fmt.Println("json Marshal error:", err)
				return
			}
			err = os.WriteFile(outputTokePath, jsonData, 0644)
			if err != nil {
				fmt.Printf("write json file error:%s\n", err.Error())
			}
		})
	}

	wg.Wait()
}
