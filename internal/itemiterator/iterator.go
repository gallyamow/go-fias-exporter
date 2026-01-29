package itemiterator

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type ItemIterator struct {
	file        io.Reader
	decoder     xml.Decoder
	rootVisited bool
}

func New(file io.Reader) *ItemIterator {
	return &ItemIterator{
		file:    file,
		decoder: *xml.NewDecoder(file),
	}
}

func (r *ItemIterator) Next(n int) ([]map[string]string, error) {
	var res []map[string]string

	for {
		tok, err := r.decoder.Token()
		if err != nil {
			if err == io.EOF {
				return res, err
			}
			return nil, fmt.Errorf("read error: %w", err)
		}

		el, ok := tok.(xml.StartElement)
		if !ok {
			// comments, directives, char datas and etc ignored
			continue
		}

		if !r.rootVisited {
			// skip root
			r.rootVisited = true
			continue
		}

		m := make(map[string]string)
		for _, attr := range el.Attr {
			m[resolveColumnName(attr.Name.Local)] = attr.Value
		}

		res = append(res, m)

		if len(res) == n {
			return res, nil
		}
	}
}

// ITEM_ID => item_id
// CHANGEID => changeid
func resolveColumnName(attrName string) string {
	return strings.ToLower(attrName)
}
