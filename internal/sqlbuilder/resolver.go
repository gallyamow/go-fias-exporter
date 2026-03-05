package sqlbuilder

import (
	"encoding/xml"
	"fmt"
	"maps"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gallyamow/go-fias-exporter/internal/config"
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
	schemaConfig := config.DefaultSchemaConfig()
	return schemaConfig.GetTableName(tableName), nil
}

func ResolveAttrs(row map[string]string) []string {
	var res []string
	for k := range maps.Keys(row) {
		res = append(res, k)
	}
	return res
}

func ResolveSchemaAttrs(tableName string, data []byte) ([]attribute, string, error) {
	var schema schema
	if err := xml.Unmarshal(data, &schema); err != nil {
		return nil, "", err
	}

	if len(schema.Element) == 0 {
		return nil, "", fmt.Errorf("invalid schema attrs for '%s'", tableName)
	}

	var attrs []attribute
	var descr string

	// (hardcoded)
	if tableName == "normative_docs_kinds" || tableName == "normative_docs_types" {
		attrs = schema.Element[1].ComplexType.Attributes
		descr = schema.Element[0].ComplexType.Sequence.Elements[0].Annotation.Documentation
	} else {
		attrs = schema.Element[0].ComplexType.Sequence.Elements[0].ComplexType.Attributes
		descr = schema.Element[0].Annotation.Documentation
	}

	if len(attrs) == 0 {
		return nil, "", fmt.Errorf("empty attrs for '%s'", tableName)
	}

	return attrs, descr, nil
}

// resolveColumnName converts an attribute name into a safe SQL column identifier
// ITEM_ID => item_id
// CHANGEID => changeid
func resolveColumnName(attrName string) string {
	return strings.ToLower(attrName)
}

// resolvePrimaryKey resolves primary key by table name.
func resolvePrimaryKey(tableName string) string {
	schemaConfig := config.DefaultSchemaConfig()
	return schemaConfig.GetPrimaryKey(tableName)
}

func resolveNullability(tableName string, columnName string, attr attribute) string {
	schemaConfig := config.DefaultSchemaConfig()
	if schemaConfig.ShouldBeNull(tableName, columnName) {
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

func escapeColumnNamePostgreSQL(columnName string) string {
	// (hardcoded)
	if columnName == "desc" {
		return fmt.Sprintf(`"%s"`, columnName)
	}
	return columnName
}

func escapeColumnNameMySQL(columnName string) string {
	// (hardcoded)
	if columnName == "desc" {
		return fmt.Sprintf("`%s`", columnName)
	}
	return columnName
}

func escapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func convertBooleanToMySQL(value string) string {
	switch strings.ToLower(value) {
	case "true":
		return "1"
	case "false":
		return "0"
	default:
		return value
	}
}
