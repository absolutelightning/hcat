package tfunc

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/hashicorp/hcat"
	"github.com/hashicorp/hcat/dep"
)

func testHealthServiceQueryID(service string) string {
	return fmt.Sprintf("health.service(%s|passing)", service)
}

func testKVListQueryID(prefix string) string {
	return fmt.Sprintf("kv.list(%s)", prefix)
}

func TestAllForDups(t *testing.T) {
	all := make(template.FuncMap)
	allfuncs := []func() template.FuncMap{
		ConsulFilters, Env, Control, Helpers, Math}
	for _, f := range allfuncs {
		for k, v := range f() {
			if _, ok := all[k]; ok {
				t.Fatal("duplicate entry")
			}
			all[k] = v
		}
	}
}

// Wrap the new template to use our template library
func NewTemplate(ti hcat.TemplateInput) *hcat.Template {
	switch ti.FuncMapMerge {
	case nil:
		ti.FuncMapMerge = All()
	default:
		for k, v := range All() {
			ti.FuncMapMerge[k] = v
		}
	}
	return hcat.NewTemplate(ti)
}

// fake/stub Watcherer for tests
type fakeWatcher struct {
	*hcat.Store
}

func (fakeWatcher) Buffer(hcat.Notifier) bool     { return false }
func (f fakeWatcher) Complete(hcat.Notifier) bool { return true }
func (f fakeWatcher) Recaller(hcat.Notifier) hcat.Recaller {
	return func(d dep.Dependency) (value interface{}, found bool) {
		return f.Store.Recall(d.ID())
	}
}
