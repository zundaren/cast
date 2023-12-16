package cast_test

import (
	"fmt"
	"testing"

	"github.com/lvan100/cast"
)

func TestToInt(t *testing.T) {
	v, err := cast.To[int](3.0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}

func TestToString(t *testing.T) {
	v, err := cast.To[string](3.0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}

func TestToBool(t *testing.T) {
	v, err := cast.To[bool](false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}

func TestToStruct(t *testing.T) {
	type Temp struct {
		A string `json:"a"`
	}
	i := map[string]interface{}{
		"a": "123",
	}
	v, err := cast.To[Temp](i)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}
