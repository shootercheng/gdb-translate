package translate

import (
	"testing"

	"github.com/shootercheng/gdb-translate/internal/translate"
)

func TestParseHtmlFileToNodes(t *testing.T) {
	doc, childNodes, err := translate.ParseHtmlFileToNodes("index.html")
	if err != nil {
		t.Fatal(err)
	}
	println(doc)
	nodeLen := len(childNodes)
	println(nodeLen)
	if nodeLen <= 0 {
		t.Fatalf("child nodes length:%d", nodeLen)
	}
}

func TestConvertRequestParam(t *testing.T) {
	_, childNodes, err := translate.ParseHtmlFileToNodes("index.html")
	if err != nil {
		t.Fatal(err)
	}
	batchParam, originMap := translate.ConvertRequestParam(childNodes, 100)
	println(len(batchParam))
	println(len(originMap))
}

func TestRequestLlmBatch(t *testing.T) {
	_, childNodes, err := translate.ParseHtmlFileToNodes("index.html")
	if err != nil {
		t.Fatal(err)
	}
	batchParam, originMap := translate.ConvertRequestParam(childNodes, 100)
	tranlateRes, usage, err := translate.RequestLlmBatch(batchParam, originMap)
	if err != nil {
		t.Fatalf("request llm batch error")
	}
	println(tranlateRes)
	println(usage)
}

func TestTranslateHtmlFile(t *testing.T) {
	translateFileParam := translate.TranslateParam{
		InputPath:        "./",
		FileName:         "index.html",
		RequestBatchSize: 50,
		OutputPath:       "./output",
	}
	translate.TranslateHtmlFile(translateFileParam)
}
