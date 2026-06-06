package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/shootercheng/gdb-translate/internal/translate"
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
	outputPath := pwd + string(os.PathSeparator) + "output-v2"
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

	inputPath := pwd + string(os.PathSeparator) + "input"
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
	maxConcurrency := 5
	sem := make(chan int, maxConcurrency)
	requestBatchSize := 100
	for index, fileName := range todoFileNames {
		wg.Go(func() {

			sem <- index
			defer func() {
				<-sem
			}()

			translateParam := translate.TranslateParam{
				InputPath:        inputPath,
				FileName:         fileName,
				RequestBatchSize: requestBatchSize,
				OutputPath:       outputPath,
			}
			translate.TranslateHtmlFile(translateParam)
		})
	}

	wg.Wait()
}
