package tdc

import (
	"bytes"
	"testing"
)

type FakeSolver struct {
}

func (s FakeSolver) ConfigurationRefByName(name string) ([]byte, error) {
	return []byte(`server="100"`), nil
}
func TestDynamicToml_Load(t *testing.T) {
	data, err := NewDynamicToml(FakeSolver{}).Parse([]byte(`
	# refs serverConf
	a=1
	b=2
	`))
	if err != nil {
		t.Fatal(err)
	}
	expected := `# refs serverConf
server="100"
a=1
b=2`
	if !bytes.Equal(bytes.TrimSpace(data), []byte(expected)) {
		t.Fatal(string(data), expected)
	}
}
