package itemiterator

import (
	"io"
	"strings"
	"testing"
)

func TestItemIterator_Next(t *testing.T) {
	xmlData := `
	<root>
		<item ID="1"/>
		<item ID="2"/>
		<item ID="3"/>
	</root>`

	it := New(strings.NewReader(xmlData))

	tests := []struct {
		name    string
		n       int
		wantLen int
		wantEOF bool
		wantIDs []string
	}{
		{
			name:    "single item",
			n:       1,
			wantLen: 1,
			wantIDs: []string{"1"},
		},
		{
			name:    "multiple items",
			n:       2,
			wantLen: 2,
			wantIDs: []string{"2", "3"},
		},
		{
			name:    "no items left",
			n:       1,
			wantLen: 0,
			wantEOF: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := it.Next(tt.n)

			if tt.wantEOF {
				if err != io.EOF {
					t.Fatalf("expected EOF, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(items) != tt.wantLen {
				t.Fatalf("expected %d items, got %d", tt.wantLen, len(items))
			}

			for i, id := range tt.wantIDs {
				if items[i]["id"] != id {
					t.Fatalf("expected id=%s, got %v", id, items[i])
				}
			}
		})
	}
}

func TestItemIterator_SkipsRootElement(t *testing.T) {
	xmlData := `
	<root>
		<item ID="42"/>
	</root>`

	it := New(strings.NewReader(xmlData))

	items, err := it.Next(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if items[0]["id"] != "42" {
		t.Fatalf("expected id=42, got %v", items[0])
	}
}

func TestItemIterator_AttributeNamesLowercased(t *testing.T) {
	xmlData := `
	<root>
		<item ITEM_ID="123" ChangeID="456"/>
	</root>`

	it := New(strings.NewReader(xmlData))

	items, err := it.Next(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	item := items[0]

	if _, ok := item["item_id"]; !ok {
		t.Fatalf("expected attribute item_id, got %v", item)
	}

	if _, ok := item["changeid"]; !ok {
		t.Fatalf("expected attribute changeid, got %v", item)
	}
}

func TestItemIterator_EmptyDocument(t *testing.T) {
	xmlData := `<root></root>`

	it := New(strings.NewReader(xmlData))

	items, err := it.Next(1)

	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}

	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}
