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
		// params tables
		{
			name:     "rooms_params",
			filename: "AS_ROOMS_PARAMS_20260129_dc1a4296-0d07-4740-8cf7-3a2f9de8ffa3.XML",
			want:     "param",
		},
		{
			name:     "carplaces_params",
			filename: "AS_CARPLACES_PARAMS_20260129_1fef63bd-e097-4a9c-bf6b-24ea3a03f677.XML",
			want:     "param",
		},
		{
			name:     "addr_obj_params",
			filename: "AS_ADDR_OBJ_PARAMS_20260129_63da61e3-dfcb-49e9-b86d-fab7a11fd814.XML",
			want:     "param",
		},
		{
			name:     "apartments_params",
			filename: "AS_APARTMENTS_PARAMS_20260129_d90d26f3-db84-49e0-aeb6-fa44e2f3811f.XML",
			want:     "param",
		},
		{
			name:     "houses_params",
			filename: "AS_HOUSES_PARAMS_20260129_3182ca06-d1cf-4b37-8663-3cf47847bfae.XML",
			want:     "param",
		},
		{
			name:     "steads_params",
			filename: "AS_STEADS_PARAMS_20260129_42707246-af85-4544-a092-cc258a4ca1fe.XML",
			want:     "param",
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
			got := resolveColumnName(tt.attrName)

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
			got := resolvePrimaryKey(tt.table)

			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}
