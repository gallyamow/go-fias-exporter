package model

import (
	"strings"
	"testing"
)

func TestHandleItems(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="utf-8"?><ITEMS>
<ITEM ID="8" PARENTID="6817" CHILDID="6951" CHANGEID="22008" /><ITEM ID="9" PARENTID="6869" CHILDID="6068" CHANGEID="21760" />
<ITEM ID="10" PARENTID="7062" CHILDID="7147" CHANGEID="22416" />`

		rdr := strings.NewReader(xmlData)
		HandleItems(rdr, func() {
			// todo
		})
	})
}
