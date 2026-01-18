package utils

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

type xmlNode struct {
	name     string
	attrs    map[string]string
	children []*xmlNode
	text     strings.Builder
}

func ParseXMLToMap(data []byte) (map[string]any, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	var stack []*xmlNode
	var root *xmlNode

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			node := &xmlNode{name: tok.Name.Local}
			if len(tok.Attr) > 0 {
				node.attrs = make(map[string]string, len(tok.Attr))
				for _, attr := range tok.Attr {
					node.attrs[attr.Name.Local] = attr.Value
				}
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, node)
			}
			stack = append(stack, node)
			if root == nil {
				root = node
			}
		case xml.CharData:
			if len(stack) == 0 {
				continue
			}
			stack[len(stack)-1].text.Write(tok)
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}

	if root == nil {
		return map[string]any{}, nil
	}

	return map[string]any{root.name: root.toValue()}, nil
}

func (n *xmlNode) toValue() any {
	text := strings.TrimSpace(n.text.String())
	if len(n.children) == 0 {
		if len(n.attrs) == 0 {
			return text
		}
		out := make(map[string]any, len(n.attrs)+1)
		for key, value := range n.attrs {
			out["@"+key] = value
		}
		if text != "" {
			out["_text"] = text
		}
		return out
	}

	out := make(map[string]any, len(n.children)+len(n.attrs)+1)
	for key, value := range n.attrs {
		out["@"+key] = value
	}
	if text != "" {
		out["_text"] = text
	}

	for _, child := range n.children {
		val := child.toValue()
		if existing, ok := out[child.name]; ok {
			switch v := existing.(type) {
			case []any:
				out[child.name] = append(v, val)
			default:
				out[child.name] = []any{v, val}
			}
			continue
		}
		out[child.name] = val
	}

	return out
}
