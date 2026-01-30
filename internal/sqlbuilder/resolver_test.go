package sqlbuilder

import "testing"

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

func TestResolveColumnName(t *testing.T) {
	tests := []struct {
		name     string
		attrName string
		want     string
	}{
		{
			name:     "uppercase",
			attrName: "CHANGEID",
			want:     "changeid",
		},
		{
			name:     "lowercase",
			attrName: "changeid",
			want:     "changeid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveColumnName(tt.attrName)

			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestResolvePrimaryKey(t *testing.T) {
	tests := []struct {
		name  string
		table string
		row   map[string]string
		want  string
	}{
		{
			name:  "has ID",
			table: "some_table",
			row:   map[string]string{"ID": "1", "NAME": "test"},
			want:  "id",
		},
		{
			name:  "changeid",
			table: "change_history",
			row:   map[string]string{"CHANGEID": "1", "NAME": "test"},
			want:  "changeid",
		},
		{
			name:  "objectid",
			table: "reestr_objects",
			row:   map[string]string{"OBJECTID": "1", "NAME": "test"},
			want:  "objectid",
		},
		{
			name:  "level",
			table: "object_levels",
			row:   map[string]string{"LEVEL": "1", "NAME": "test"},
			want:  "level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolvePrimaryKey(tt.table, tt.row)

			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}
