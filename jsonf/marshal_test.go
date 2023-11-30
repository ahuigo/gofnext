package jsonf

import (
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	type Stu struct {
		name string
	}
	type MyStruct struct {
		Name       string
		nums       []int
		Nil        *int
		Remark     [1]string `json:"remark"`
		pointerMap map[string]*Stu
	}

	myStruct := MyStruct{
		Name:       "public",
		nums:       []int{1, 2},
		Remark:     [1]string{"remark"},
		pointerMap: map[string]*Stu{"stu1": {name: "private"}},
	}

	data, err := Marshal(myStruct)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(string(data))
}
