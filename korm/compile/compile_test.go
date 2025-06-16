package compile

import (
	"testing"
)

func TestComplie(t *testing.T) {
	err := Compile("../../examples", "./", []string{"../../examples/models/user.go", "../../examples/models/user_ext.go", "../../examples/models/user_tests.go", "../../examples/models/user_ext_join.go"})
	if err != nil {
		t.Fatal(err)
	}
}
