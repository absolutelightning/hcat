package tfunc

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hashicorp/hcat"
)

func TestLoopExecute(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		ti   hcat.TemplateInput
		i    hcat.Recaller
		e    string
		err  bool
	}{
		{
			"helper_loop",
			hcat.TemplateInput{
				Contents: `{{ range loop 3 }}1{{ end }}`,
			},
			hcat.NewStore(),
			"111",
			false,
		},
		{
			"helper_loop__i",
			hcat.TemplateInput{
				Contents: `{{ range $i := loop 3 }}{{ $i }}{{ end }}`,
			},
			hcat.NewStore(),
			"012",
			false,
		},
		{
			"helper_loop_start",
			hcat.TemplateInput{
				Contents: `{{ range loop 1 3 }}1{{ end }}`,
			},
			hcat.NewStore(),
			"11",
			false,
		},
		{
			"helper_loop_text",
			hcat.TemplateInput{
				Contents: `{{ range loop 1 "3" }}1{{ end }}`,
			},
			hcat.NewStore(),
			"11",
			false,
		},
		{
			"helper_loop_parseInt",
			hcat.TemplateInput{
				Contents: `{{ $i := print "3" | parseInt }}{{ range loop 1 $i }}1{{ end }}`,
			},
			hcat.NewStore(),
			"11",
			false,
		},
		{
			// GH-1143
			"helper_loop_var",
			hcat.TemplateInput{
				Contents: `{{$n := 3 }}` +
					`{{ range $i := loop $n }}{{ $i }}{{ end }}`,
			},
			hcat.NewStore(),
			"012",
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
