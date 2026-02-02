package sqlbuilder

import (
	"fmt"
	"maps"
	"path/filepath"
	"regexp"
	"strings"
)

// these values have some differences in building
const (
	TableRoomsParams      = "rooms_params"
	TableCarplacesParams  = "carplaces_params"
	TableAddrObjParams    = "addr_obj_params"
	TableApartmentsParams = "apartments_params"
	TableHousesParams     = "houses_params"
	TableSteadsParams     = "steads_params"
	TableChangeHistory    = "change_history"
	TableReestrObjects    = "reestr_objects"
	TableObjectLevels     = "object_levels"
	TableNormativeDocs    = "normative_docs"
)

// ResolveTableName resolves table name from filename.
// Examples:
// - AS_ADDR_OBJ_20250626_bc6f64d9-fb28-40d6-8a99-57e44b920d07.XML => addr_obj
// - AS_CHANGE_HISTORY_20250626_d1a57485-156c-4463-8a23-2328fb0f6f9d => change_history
func ResolveTableName(filename string) (string, error) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))

	re := regexp.MustCompile(`(?i)^AS_([A-Z_]+)_\d*`)
	m := re.FindStringSubmatch(base)
	if len(m) != 2 {
		return "", fmt.Errorf("failed to resolve table name by filename %s", filename)
	}

	tableName := strings.ToLower(m[1])

	switch tableName {
	// (hardcoded)
	case TableRoomsParams, TableCarplacesParams, TableAddrObjParams, TableApartmentsParams, TableHousesParams, TableSteadsParams:
		return "param", nil
	default:
		return tableName, nil
	}
}

func ResolveAttrs(row map[string]string) []string {
	var res []string
	for k := range maps.Keys(row) {
		res = append(res, k)
	}
	return res
}

// resolveColumnName converts an attribute name into a safe SQL column identifier
// ITEM_ID => item_id
// CHANGEID => changeid
func resolveColumnName(attrName string) string {
	return strings.ToLower(attrName)
}

// resolvePrimaryKey resolves primary key by table name.
func resolvePrimaryKey(tableName string) string {
	// (hardcoded)
	switch tableName {
	case TableChangeHistory:
		return "changeid"
	case TableReestrObjects:
		return "objectid"
	case TableObjectLevels:
		return "level"
	}
	return "id"
}

func resolveNullability(tableName string, columnName string, attr attribute) string {
	// (hardcoded)
	if tableName == TableNormativeDocs && columnName == "name" {
		// example: /63_sql/AS_NORMATIVE_DOCS_20260127_29897f0f-87b4-43b9-bea9-54152f80d42f.sql:493710: ERROR:  null value in column "name" of relation "normative_docs" violates not-null constraint
		return ""
	}

	if attr.Use == "required" {
		return "NOT NULL"
	}

	return ""
}

func buildFullTableName(dbSchema string, tableName string) string {
	if dbSchema != "" {
		return fmt.Sprintf("%s.%s", dbSchema, tableName)
	}
	return tableName
}

func escapeColumnName(columnName string) string {
	// (hardcoded)
	if columnName == "desc" {
		return fmt.Sprintf(`"%s"`, columnName)
	}
	return columnName
}

func escapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
