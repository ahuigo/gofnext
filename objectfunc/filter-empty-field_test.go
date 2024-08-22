package objectfunc

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFilterEmptyFieldOfMap(t *testing.T) {
	objMap := map[string][]byte{
		"k1": []byte("v1"),
		"k2": []uint8{},
		"k3": nil,
	}
	out, _ := json.Marshal(objMap)
	fmt.Println(string(out)) //output: {"k1":"djE=","k2":"djI="}

	objString := FilterObjectEmptyField(objMap)
	out, _ = json.Marshal(objString)
	fmt.Println(string(out)) //output: {"k1":"v1","k2":"v2"}

	expectedOut := `{"k1":"djE="}`
	if string(out) != expectedOut {
		t.Fatalf("expected out:%v, unexpected out: %v", expectedOut, string(out))
	}
}

func TestFilterEmptyFieldOfStruct(t *testing.T) {
	type HistoryEvent struct {
		EventId *int64 `json:"eventId"`
		TaskId  *int64 `json:"taskId"`
	}
	i := int64(1)
	obj := HistoryEvent{
		EventId: &i,
	}
	if out, err := json.Marshal(FilterObjectEmptyField(obj)); err != nil {
		t.Fatal(err)
	} else {
		expectedOut := `{"eventId":1}`
		if string(out) != expectedOut {
			t.Fatalf("expected out:%v, unexpected out: %v", expectedOut, string(out))
		}
	}
}
