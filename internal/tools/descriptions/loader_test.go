package descriptions

import (
	"strings"
	"testing"
)

func TestParse_ExtractsToolAndParams(t *testing.T) {
	input := []byte(`
## my_tool

This is the tool description.
It can span multiple lines.

### param_one
First parameter description.

### param_two
Second parameter description.
`)
	result := parse(input)
	doc, ok := result["my_tool"]
	if !ok {
		t.Fatal("my_tool not found in catalog")
	}
	if !strings.Contains(doc.Description, "tool description") {
		t.Errorf("description: got %q", doc.Description)
	}
	if doc.Params["param_one"] == "" {
		t.Error("param_one missing")
	}
	if doc.Params["param_two"] == "" {
		t.Error("param_two missing")
	}
}

func TestLoad_AllCategoriesPresent(t *testing.T) {
	d := Load()
	if len(d) == 0 {
		t.Fatal("catalog is empty — embedded .md files not found or not parseable")
	}
}

func TestParse_MultipleToolsInOneFile(t *testing.T) {
	input := []byte(`
## tool_a

Description A.

### foo
Foo param.

## tool_b

Description B.

### bar
Bar param.
`)
	result := parse(input)
	if _, ok := result["tool_a"]; !ok {
		t.Error("tool_a missing")
	}
	if _, ok := result["tool_b"]; !ok {
		t.Error("tool_b missing")
	}
	if result["tool_a"].Params["foo"] == "" {
		t.Error("tool_a.foo missing")
	}
	if result["tool_b"].Params["bar"] == "" {
		t.Error("tool_b.bar missing")
	}
}
