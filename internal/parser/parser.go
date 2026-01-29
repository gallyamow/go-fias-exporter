package model

import (
	"encoding/xml"
	"fmt"
	"io"
)

func HandleItems(file io.Reader, handler func()) error {
	decoder := xml.NewDecoder(file)

	for {
		tok, err := decoder.Token()
		if err != nil {
			break // EOF
		}

		fmt.Println(tok)
		//handler(tok)

		//switch t := tok.(type) {
		//case xml.StartElement:
		//	switch t := tok.(type) {
		//	case xml.StartElement:
		//		if t.Name.Local == "id" {
		//			var v string
		//			decoder.DecodeElement(&v, &t)
		//			fmt.Println("id:", v)
		//		}
		//	}
		//}
	}

	return nil
}
