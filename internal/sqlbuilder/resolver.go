package sqlbuilder

import (
	"fmt"
	"maps"
	"path/filepath"
	"regexp"
	"strings"
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
	case "rooms_params", "carplaces_params", "addr_obj_params", "apartments_params", "houses_params", "steads_params":
		return "param", nil
	default:
		return tableName, nil
	}
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
	case "change_history":
		return "changeid"
	case "reestr_objects":
		return "objectid"
	case "object_levels":
		return "level"
	}
	return "id"
}

func ResolveAttrs(row map[string]string) []string {
	var res []string
	for k := range maps.Keys(row) {
		res = append(res, k)
	}
	return res
}

func buildFullTableName(dbSchema string, tableName string) string {
	if dbSchema != "" {
		return fmt.Sprintf("%s.%s", dbSchema, tableName)
	}
	return tableName
}

func escapeColumnName(c string) string {
	// (hardcoded)
	if c == "desc" {
		return fmt.Sprintf(`"%s"`, c)
	}
	return c
}

func escapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
