package config

type TableConfig struct {
	Name       string
	PrimaryKey string
}

type NullabilityRule struct {
	Table  string
	Column string
}

type SchemaConfig struct {
	TableMappings map[string]string
	PrimaryKeys   map[string]string
	Nullability   []NullabilityRule
}

// DefaultSchemaConfig store some hardcoded things
func DefaultSchemaConfig() *SchemaConfig {
	return &SchemaConfig{
		TableMappings: map[string]string{
			"rooms_params":      "param",
			"carplaces_params":  "param",
			"addr_obj_params":   "param",
			"apartments_params": "param",
			"houses_params":     "param",
			"steads_params":     "param",
		},
		PrimaryKeys: map[string]string{
			"change_history": "changeid",
			"reestr_objects": "objectid",
			"object_levels":  "level",
		},
		Nullability: []NullabilityRule{
			{Table: "normative_docs", Column: "name"},
			{Table: "steads", Column: "number"},
		},
	}
}

func (c *SchemaConfig) GetTableName(tableName string) string {
	if mapped, exists := c.TableMappings[tableName]; exists {
		return mapped
	}
	return tableName
}

func (c *SchemaConfig) GetPrimaryKey(tableName string) string {
	if pk, exists := c.PrimaryKeys[tableName]; exists {
		return pk
	}
	return "id"
}

func (c *SchemaConfig) ShouldBeNull(table, column string) bool {
	for _, rule := range c.Nullability {
		if rule.Table == table && rule.Column == column {
			return true
		}
	}
	return false
}
