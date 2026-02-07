package sqlbuilder

import "testing"

func TestCopyBuilder_Build(t *testing.T) {
	t.Run("with dbSchema", func(t *testing.T) {
		builder := NewCopyBuilder(
			"tmp",
			"test_table",
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

		want := `COPY tmp.test_table (id,name) FROM STDIN WITH (FORMAT csv);
1,Alice
2,Bob
\.`
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("no dbSchema", func(t *testing.T) {
		builder := NewCopyBuilder(
			"",
			"test_table",
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

		want := `COPY test_table (id,name) FROM STDIN WITH (FORMAT csv);
1,Alice
2,Bob
\.`
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

func TestCopyBuilder_PreservesColumnOrder(t *testing.T) {
	builder := NewCopyBuilder(
		"tmp",
		"test_table",
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

	want := `COPY tmp.test_table (name,id) FROM STDIN WITH (FORMAT csv);
Alice,1
Bob,2
\.`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
