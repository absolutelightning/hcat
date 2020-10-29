package tfunc

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hashicorp/hcat"
)

func TestParseExecute(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		ti   hcat.TemplateInput
		i    hcat.Recaller
		e    string
		err  bool
	}{
		{
			"parseBool",
			hcat.TemplateInput{
				Contents: `{{ "true" | parseBool }}`,
			},
			hcat.NewStore(),
			"true",
			false,
		},
		{
			"parseFloat",
			hcat.TemplateInput{
				Contents: `{{ "1.2" | parseFloat }}`,
			},
			hcat.NewStore(),
			"1.2",
			false,
		},
		{
			"parseInt",
			hcat.TemplateInput{
				Contents: `{{ "-1" | parseInt }}`,
			},
			hcat.NewStore(),
			"-1",
			false,
		},
		{
			"parseJSON",
			hcat.TemplateInput{
				Contents: `{{ "{\"foo\": \"bar\"}" | parseJSON }}`,
			},
			hcat.NewStore(),
			"map[foo:bar]",
			false,
		},
		{
			"parseUint",
			hcat.TemplateInput{
				Contents: `{{ "1" | parseUint }}`,
			},
			hcat.NewStore(),
			"1",
			false,
		},
		{
			"parseYAML",
			hcat.TemplateInput{
				Contents: `{{ "foo: bar" | parseYAML }}`,
			},
			hcat.NewStore(),
			"map[foo:bar]",
			false,
		},
		{
			"parseYAMLv2",
			hcat.TemplateInput{
				Contents: `{{ "foo: bar\nbaz: \"foo\"" | parseYAML }}`,
			},
			hcat.NewStore(),
			"map[baz:foo foo:bar]",
			false,
		},
		{
			"parseYAMLnested",
			hcat.TemplateInput{
				Contents: `{{ "foo:\n  bar: \"baz\"\n  baz: 7" | parseYAML }}`,
			},
			hcat.NewStore(),
			"map[foo:map[bar:baz baz:7]]",
			false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d_%s", i, tc.name), func(t *testing.T) {
			tpl := NewTemplate(tc.ti)

			a, err := tpl.Execute(tc.i)
			if (err != nil) != tc.err {
				t.Fatal(err)
			}
			if a != nil && !bytes.Equal([]byte(tc.e), a.Output) {
				t.Errorf("\nexp: %#v\nact: %#v", tc.e, string(a.Output))
			}
		})
	}
}