package sqlbuilder

import "testing"

func TestUpsertBuilder_Build(t *testing.T) {
	t.Run("with primary key", func(t *testing.T) {
		builder := NewUpsertBuilder(
			"tmp",
			"test_table",
			"id",
			[]string{"ID", "NAME"},
		)

		rows := []map[string]string{
			{
				"ID":   "1",
				"NAME": "Alice",
			},
			{
				"ID":   "2",
				"NAME": "Bob",
			},
		}

		got, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `INSERT INTO tmp.test_table (id,name) VALUES ('1','Alice'),('2','Bob') ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name`
		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("no primary key", func(t *testing.T) {
		builder := NewUpsertBuilder(
			"tmp",
			"test_table",
			"",
			[]string{"FIELD", "NAME"},
		)

		rows := []map[string]string{
			{
				"FIELD": "1",
				"NAME":  "Alice",
			},
			{
				"FIELD": "2",
				"NAME":  "Bob",
			},
		}

		got, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `INSERT INTO tmp.test_table (field,name) VALUES ('1','Alice'),('2','Bob')`
		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("escape values", func(t *testing.T) {
		builder := NewUpsertBuilder(
			"tmp",
			"test_table",
			"id",
			[]string{"ID", "NAME"},
		)

		rows := []map[string]string{
			{
				"ID":   "1",
				"NAME": "Alice's brother",
			},
			{
				"ID":   "2",
				"NAME": "Bob's car",
			},
		}

		got, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `INSERT INTO tmp.test_table (id,name) VALUES ('1','Alice''s brother'),('2','Bob''s car') ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name`
		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("no schema", func(t *testing.T) {
		builder := NewUpsertBuilder(
			"",
			"test_table",
			"id",
			[]string{"ID", "NAME"},
		)

		rows := []map[string]string{
			{
				"ID":   "1",
				"NAME": "Alice",
			},
			{
				"ID":   "2",
				"NAME": "Bob",
			},
		}

		got, err := builder.Build(rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `INSERT INTO test_table (id,name) VALUES ('1','Alice'),('2','Bob') ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name`
		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})
}

func TestUpsertBuilder_PreservesColumnOrder(t *testing.T) {
	builder := NewUpsertBuilder(
		"tmp",
		"test_table",
		"id",
		[]string{"NAME", "ID"},
	)

	rows := []map[string]string{
		{
			"ID":   "1",
			"NAME": "Alice",
		},
		{
			"NAME": "Bob",
			"ID":   "2",
		},
	}

	got, err := builder.Build(rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `INSERT INTO tmp.test_table (name,id) VALUES ('Alice','1'),('Bob','2') ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name`
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
