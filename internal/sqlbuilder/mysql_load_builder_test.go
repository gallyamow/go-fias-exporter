package sqlbuilder

import (
	"testing"
)

func TestMySQLLoadDataBuilder_Build(t *testing.T) {
	t.Run("basic load data", func(t *testing.T) {
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

		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "name", "description"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `LOAD DATA LOCAL INFILE 'stdin' INTO TABLE test_table FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' (id,name,description);
1,test,test description
2,test2,test description 2
`

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

		builder := NewMySQLLoadDataBuilder("test_schema", "test_table", []string{"id", "name"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `LOAD DATA LOCAL INFILE 'stdin' INTO TABLE test_schema.test_table FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' (id,name);
1,test
`

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

		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "name"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `LOAD DATA LOCAL INFILE 'stdin' INTO TABLE test_table FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' (id,name);
1,test
`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})

	t.Run("with special characters", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":   "1",
				"name": "test,with,commas",
				"desc": "test\"with\"quotes",
			},
		}

		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "name", "desc"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `LOAD DATA LOCAL INFILE 'stdin' INTO TABLE test_table FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' (id,name,"desc");
1,"test,with,commas","test""with""quotes"
`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})

	t.Run("with newlines", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":   "1",
				"name": "test\nwith\nnewlines",
				"desc": "description",
			},
		}

		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "name", "desc"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `LOAD DATA LOCAL INFILE 'stdin' INTO TABLE test_table FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' (id,name,"desc");
1,"test
with
newlines",description
`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})

	t.Run("empty rows", func(t *testing.T) {
		rows := []map[string]string{}

		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "name"})
		_, err := builder.Build(rows)
		if err == nil {
			t.Fatal("Build() expected error for empty rows")
		}

		expectedError := "no rows to build"
		if err.Error() != expectedError {
			t.Errorf("Build() error = %v, expected %v", err.Error(), expectedError)
		}
	})

	t.Run("null values", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":   "1",
				"name": "",
				"desc": "",
			},
		}

		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "name", "desc"})
		result, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		expected := `LOAD DATA LOCAL INFILE 'stdin' INTO TABLE test_table FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' (id,name,"desc");
1,,
`

		if result != expected {
			t.Errorf("Build() result =\n%s\n\nexpected =\n%s", result, expected)
		}
	})
}

func TestMySQLLoadDataBuilder_buildValues(t *testing.T) {
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

		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "name"})
		result, err := builder.buildValues(rows)
		if err != nil {
			t.Fatalf("buildValues() error = %v", err)
		}

		expected := "1,test1\n2,test2\n"
		if result != expected {
			t.Errorf("buildValues() result = %q, expected %q", result, expected)
		}
	})

	t.Run("with commas and quotes", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":   "1",
				"name": "test,with,commas",
				"desc": "test\"with\"quotes",
			},
		}

		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "name", "desc"})
		result, err := builder.buildValues(rows)
		if err != nil {
			t.Fatalf("buildValues() error = %v", err)
		}

		expected := "1,\"test,with,commas\",\"test\"\"with\"\"quotes\"\n"
		if result != expected {
			t.Errorf("buildValues() result = %q, expected %q", result, expected)
		}
	})

	t.Run("with newlines in data", func(t *testing.T) {
		rows := []map[string]string{
			{
				"id":   "1",
				"name": "test\nwith\nnewlines",
			},
		}

		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "name"})
		result, err := builder.buildValues(rows)
		if err != nil {
			t.Fatalf("buildValues() error = %v", err)
		}

		expected := "1,\"test\nwith\nnewlines\"\n"
		if result != expected {
			t.Errorf("buildValues() result = %q, expected %q", result, expected)
		}
	})
}

func TestMySQLLoadDataBuilder_buildColumns(t *testing.T) {
	t.Run("basic columns", func(t *testing.T) {
		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "name", "description"})
		result := builder.buildColumns()

		expected := "id,name,description"
		if result != expected {
			t.Errorf("buildColumns() result = %s, expected %s", result, expected)
		}
	})

	t.Run("with reserved word", func(t *testing.T) {
		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id", "desc"})
		result := builder.buildColumns()

		expected := "id,\"desc\""
		if result != expected {
			t.Errorf("buildColumns() result = %s, expected %s", result, expected)
		}
	})

	t.Run("single column", func(t *testing.T) {
		builder := NewMySQLLoadDataBuilder("", "test_table", []string{"id"})
		result := builder.buildColumns()

		expected := "id"
		if result != expected {
			t.Errorf("buildColumns() result = %s, expected %s", result, expected)
		}
	})
}
