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
		Name   string
		nums   []int
		Remark [1]string `json:"remark"`
		stus    map[int]*Stu
	}

	myStruct := MyStruct{
		Name:   "public",
		nums:   []int{1, 2},
		Remark: [1]string{"remark"},
		stus:    map[int]*Stu{1: {name: "private"}},
	}

	data, err := Marshal(myStruct)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(string(data))
}
