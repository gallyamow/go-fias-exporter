package sqlbuilder

import "testing"

func TestUpsertBuilder_Build(t *testing.T) {
	builder := New(
		"test_table",
		[]string{"id"},
		[]string{"id", "name"},
	)

	rows := []map[string]string{
		{
			"id":   "1",
			"name": "Alice",
		},
		{
			"id":   "2",
			"name": "Bob",
		},
	}

	got, err := builder.Build(rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `INSERT INTO test_table (id,name) VALUES ('1','Alice'),('2','Bob') ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name`
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestUpsertBuilder_PreservesColumnOrder(t *testing.T) {
	builder := New(
		"test_table",
		[]string{"id"},
		[]string{"name", "id"},
	)

	rows := []map[string]string{
		{
			"id":   "1",
			"name": "Alice",
		},
		{
			"name": "Bob",
			"id":   "2",
		},
	}

	got, err := builder.Build(rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `INSERT INTO test_table (name,id) VALUES ('Alice','1'),('Bob','2') ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name`
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestResolveTableName(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		want      string
		wantError bool
	}{
		{
			name:     "addr obj",
			filename: "AS_ADDR_OBJ_20250626_bc6f64d9-fb28-40d6-8a99-57e44b920d07.XML",
			want:     "addr_obj",
		},
		{
			name:     "addr obj division",
			filename: "AS_ADDR_OBJ_DIVISION_20260127_36d1e18d-6acf-4755-a7b0-49d9a30a5dae.XML",
			want:     "addr_obj_division",
		},
		{
			name:     "change history",
			filename: "AS_CHANGE_HISTORY_20250626_d1a57485-156c-4463-8a23-2328fb0f6f9d",
			want:     "change_history",
		},
		{
			name:      "invalid filename",
			filename:  "random_file.xml",
			wantError: true,
		},
		{
			name:     "lowercase extension",
			filename: "AS_TEST_TABLE_20240101_xxx.xml",
			want:     "test_table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveTableName(tt.filename)

			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}
