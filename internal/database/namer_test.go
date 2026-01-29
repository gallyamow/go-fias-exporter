package database

import (
	"testing"
)

func TestResolveTableName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "AS_ADDR_OBJ_20260126_f4bfcf47-40cb-4dd1-922d-1c789141498a.XML", want: "addr_obj"},
		{input: "AS_ADDR_OBJ_DIVISION_20260127_36d1e18d-6acf-4755-a7b0-49d9a30a5dae.XML", want: "addr_obj_division"},
		{input: "AS_ADDR_OBJ_PARAMS_20260127_87bc98cf-a9d1-4b1f-b689-3eec3c903266.XML", want: "addr_obj_params"},
		{input: "AS_ADM_HIERARCHY_20260126_483d1f66-0960-4f86-a0ba-eddc3c414b16.XML", want: "adm_hierarchy"},
		{input: "AS_APARTMENTS_20260126_fda601f3-15ff-4b31-947a-5d9cc318d0f4.XML", want: "apartments"},
		{input: "AS_APARTMENTS_PARAMS_20260127_88a6a8b1-f43d-49c4-b526-af5ff7827af1.XML", want: "apartments_params"},
		{input: "AS_CARPLACES_20260126_6b576c86-b58c-4310-bf21-c0d642f56bc4.XML", want: "carplaces"},
		{input: "AS_CARPLACES_PARAMS_20260127_c82166cb-480e-4015-875c-eca28332999e.XML", want: "carplaces_params"},
		{input: "AS_CHANGE_HISTORY_20260126_84625c36-7ade-4cc3-aee2-f16721714c5c.XML", want: "change_history"},
		{input: "AS_HOUSES_20260126_7ecfbbc1-ac67-488c-aea7-a7726bdfa7f0.XML", want: "houses"},
		{input: "AS_HOUSES_PARAMS_20260127_96c2d148-a773-4d44-b069-037142bafc01.XML", want: "houses_params"},
		{input: "AS_MUN_HIERARCHY_20260126_78e57048-0da1-4f48-acf5-7f7bde666e53.XML", want: "mun_hierarchy"},
		{input: "AS_NORMATIVE_DOCS_20260127_20486cd3-eaeb-417d-934c-f4b26aa16975.XML", want: "normative_docs"},
		{input: "AS_REESTR_OBJECTS_20260126_a5050e20-ecfc-4770-93c9-b831be3627a9.XML", want: "reestr_objects"},
		{input: "AS_ROOMS_20260126_21d02043-f06d-4da6-84a1-26d219cbf2be.XML", want: "rooms"},
		{input: "AS_ROOMS_PARAMS_20260127_86bd2561-24f3-4a66-b609-68699b2ef9fa.XML", want: "rooms_params"},
		{input: "AS_STEADS_20260126_de4b81b3-1ddd-4126-86e4-f7011bb492f3.XML", want: "steads"},
		{input: "AS_STEADS_PARAMS_20260127_28cd340b-98e1-4a17-bbe3-3d0957e999f4.XML", want: "steads_params"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ResolveTableName(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
		})
	}
}
