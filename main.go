package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Attrs map[string]string

// references:
// stream-parsing example: https://eli.thegreenplace.net/2019/faster-xml-stream-processing-in-go/
// having maps in structs: https://stackoverflow.com/a/34972468/605846
type EADNode struct {
	Name     string
	Attrs    Attrs
	Value    string
	Children []*EADNode
}

type Stack struct {
	S []*EADNode
}

// TODO: could generalize by using "any" for the param and return type
func (s *Stack) Peek() (*EADNode, error) {
	idx := len(s.S) - 1
	if idx < 0 {
		return nil, fmt.Errorf("Stack is empty")
	}
	return s.S[idx], nil
}

func (s *Stack) Push(val *EADNode) {
	s.S = append(s.S, val)
	return
}

func (s *Stack) Pop() (*EADNode, error) {
	idx := len(s.S) - 1
	if idx < 0 {
		return nil, fmt.Errorf("Stack is empty")
	}
	retval := s.S[idx]
	s.S = s.S[:idx]
	return retval, nil
}

func (s *Stack) Len() int {
	return len(s.S)
}

type EADState struct {
	Stack  Stack
	Tree   []*EADNode
	Errors []error
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := xml.NewDecoder(f)

	indent := 0
	for {
		token, err := d.Token()
		if token == nil || err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error decoding token: %s", err)
		}

		// A Token is an interface holding one of the token types:
		//   StartElement, EndElement, CharData, Comment, ProcInst, or Directive.
		// https://stackoverflow.com/a/33139049/605846
		// https://www.socketloop.com/tutorials/golang-read-xml-elements-data-with-xml-chardata-example
		// https://code-maven.com/slides/golang/parse-html-extract-tags-and-attributes
		switch el := token.(type) {
		case xml.StartElement:
			fmt.Printf("%sStartElement --> %s\n", strings.Repeat(" ", indent), el.Name.Local)

			// create new node
			var en *EADNode
			en = new(EADNode)
			en.Name = el.Name.Local
			// TODO: add size to make? len(el.Attr)
			en.Attrs = make(Attrs)

			indent += 4
			for _, attr := range el.Attr {
				en.Attrs[attr.Name.Local] = attr.Value
				fmt.Printf("%s@%s = %s\n", strings.Repeat(" ", indent), attr.Name.Local, attr.Value)
			}

		case xml.CharData:
			str := strings.TrimSpace(string([]byte(el)))
			if len(str) != 0 {
				fmt.Printf("%sCharData --> %s\n", strings.Repeat(" ", indent), str)
			}
		case xml.EndElement:
			indent -= 4
			fmt.Printf("%sEndElement --> %s\n", strings.Repeat(" ", indent), el.Name.Local)
		}

	}

}

// Golang stack implementation
// https://yourbasic.org/golang/implement-stack/
// https://go.dev/play/p/uiYfmQHR1b9
// https://go.dev/play/p/VkWkOFadSYh

// 	package main

// import (
// 	"encoding/json"
// 	"encoding/xml"
// 	"fmt"
// 	"io/ioutil"
// )

// type EAD struct {
// 	Head     string  `xml:"head" json:"head"`
// 	Contents []Mixed `xml:",any" json:"contents"`
// }

// type Mixed struct {
// 	Type  string
// 	Value interface{}
// }

// func main() {
// 	bytes, err := ioutil.ReadFile("example.xml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	var doc EAD
// 	if err := xml.Unmarshal([]byte(bytes), &doc); err != nil {
// 		panic(err)
// 	}

// 	jdoc, err := json.MarshalIndent(doc, "", "  ")
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(string(jdoc))
// }

// func (m *Mixed) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	switch start.Name.Local {
// 	case "head", "p", "list":
// 		var e string
// 		if err := d.DecodeElement(&e, &start); err != nil {
// 			return err
// 		}
// 		m.Value = e
// 		m.Type = start.Name.Local
// 	default:
// 		return fmt.Errorf("unknown element: %s", start)
// 	}
// 	return nil
// }
