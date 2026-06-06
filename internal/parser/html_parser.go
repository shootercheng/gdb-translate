package parser

import (
	"strings"

	"golang.org/x/net/html"
)

func CollectTextNodes(n *html.Node, nodes *[]*html.Node, skipTags map[string]bool) {
	if n.Type == html.ElementNode && skipTags[n.Data] {
		return
	}
	if n.Type == html.TextNode {
		if strings.TrimSpace(n.Data) != "" {
			*nodes = append(*nodes, n)
		}
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		CollectTextNodes(c, nodes, skipTags)
	}
}
