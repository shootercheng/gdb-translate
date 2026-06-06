package parser_tes

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/shootercheng/gdb-translate/internal/parser"
	"golang.org/x/net/html"
)

func TestCollectTextNodes(t *testing.T) {
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
	for i, node := range textNodes {
		fmt.Printf("%d: %q\n", i, node.Data)
	}
}
