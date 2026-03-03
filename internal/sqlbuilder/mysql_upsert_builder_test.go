package sqlbuilder

import (
	"testing"
)

func TestMySQLInsertBuilder_Build(t *testing.T) {
	t.Run("basic insert", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":          "1",
				"name":        "test",
				"description": "test description",
			},
			{
				"id":          "2",
				"name":        "test2",
				"description": "test description 2",
			},
		}

		builder := NewMySQLUpsertBuilder("", "test_table", []string{"id", "name", "description"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `INSERT INTO test_table (id,name,description) VALUES ('1','test','test description'),('2','test2','test description 2') ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description);`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})

	t.Run("with schema", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":   "1",
				"name": "test",
			},
		}

		builder := NewMySQLUpsertBuilder("test_schema", "test_table", []string{"id", "name"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `INSERT INTO test_schema.test_table (id,name) VALUES ('1','test') ON DUPLICATE KEY UPDATE name=VALUES(name);`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})

	t.Run("single row", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":   "1",
				"name": "test",
			},
		}

		builder := NewMySQLUpsertBuilder("", "test_table", []string{"id", "name"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `INSERT INTO test_table (id,name) VALUES ('1','test') ON DUPLICATE KEY UPDATE name=VALUES(name);`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})

	t.Run("with special characters", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":   "1",
				"name": "test'with'quotes",
				"desc": "test\nwith\nnewlines",
			},
		}

		builder := NewMySQLUpsertBuilder("", "test_table", []string{"id", "name", "desc"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := "INSERT INTO test_table (id,name,`desc`) VALUES ('1','test''with''quotes','test\nwith\nnewlines') ON DUPLICATE KEY UPDATE name=VALUES(name),`desc`=VALUES(`desc`);"

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})

	t.Run("empty rows", func(t *testing.T) {
		rows := []map[string]string{}

		builder := NewMySQLUpsertBuilder("", "test_table", []string{"id", "name"})
		_, err := builder.Build(rows)
		if err == nil {
			t.Fatal("Build() expected error for empty rows")
		}

		expectedError := "no rows to build"
		if err.Error() != expectedError {
			t.Errorf("Build() error = %v, expected %v", err.Error(), expectedError)
		}
	})

	t.Run("different primary key", func(t *testing.T) {
		rows := []map[string]string{
			{
				"objectid": "1",
				"name":     "test",
			},
		}

		builder := NewMySQLUpsertBuilder("", "reestr_objects", []string{"objectid", "name"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `INSERT INTO reestr_objects (objectid,name) VALUES ('1','test') ON DUPLICATE KEY UPDATE name=VALUES(name);`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})
}

func TestMySQLInsertBuilder_buildValues(t *testing.T) {
	t.Run("multiple rows", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":   "1",
				"name": "test1",
			},
			{
				"id":   "2",
				"name": "test2",
			},
		}

		builder := NewMySQLUpsertBuilder("", "test_table", []string{"id", "name"})
		result, err := builder.buildValues(rows)
		if err != nil {
			t.Fatalf("buildValues() error = %v", err)
		}

		expected := "('1','test1'),('2','test2')"
		if result != expected {
			t.Errorf("buildValues() result = %s, expected %s", result, expected)
		}
	})

	t.Run("with null values", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":   "1",
				"name": "",
				"desc": "",
			},
		}

		builder := NewMySQLUpsertBuilder("", "test_table", []string{"id", "name", "desc"})
		result, err := builder.buildValues(rows)
		if err != nil {
			t.Fatalf("buildValues() error = %v", err)
		}

		expected := "('1','','')"
		if result != expected {
			t.Errorf("buildValues() result = %s, expected %s", result, expected)
		}
	})
}

func TestMySQLInsertBuilder_buildColumns(t *testing.T) {
	t.Run("basic columns", func(t *testing.T) {
		builder := NewMySQLUpsertBuilder("", "test_table", []string{"id", "name", "description"})
		result := builder.buildColumns()

		expected := "id,name,description"
		if result != expected {
			t.Errorf("buildColumns() result = %s, expected %s", result, expected)
		}
	})

	t.Run("with reserved word", func(t *testing.T) {
		builder := NewMySQLUpsertBuilder("", "test_table", []string{"id", "desc"})
		result := builder.buildColumns()

		expected := "id,`desc`"
		if result != expected {
			t.Errorf("buildColumns() result = %s, expected %s", result, expected)
		}
	})
}

func TestMySQLInsertBuilder_buildOnDuplicateKeyUpdate(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		builder := NewMySQLUpsertBuilder("", "test_table", []string{"id", "name", "description"})
		result := builder.buildOnConflict()

		expected := "ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description)"
		if result != expected {
			t.Errorf("buildOnConflict() result = %s, expected %s", result, expected)
		}
	})

	t.Run("single column", func(t *testing.T) {
		builder := NewMySQLUpsertBuilder("", "test_table", []string{"id", "name"})
		result := builder.buildOnConflict()

		expected := "ON DUPLICATE KEY UPDATE name=VALUES(name)"
		if result != expected {
			t.Errorf("buildOnConflict() result = %s, expected %s", result, expected)
		}
	})

	t.Run("different primary key", func(t *testing.T) {
		builder := NewMySQLUpsertBuilder("", "reestr_objects", []string{"objectid", "name"})
		result := builder.buildOnConflict()

		expected := "ON DUPLICATE KEY UPDATE name=VALUES(name)"
		if result != expected {
			t.Errorf("buildOnConflict() result = %s, expected %s", result, expected)
		}
	})
}
