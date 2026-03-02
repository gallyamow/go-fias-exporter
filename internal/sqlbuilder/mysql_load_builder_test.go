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
}
