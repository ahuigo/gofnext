package examples

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ahuigo/gofnext/object"
)

func TestConvertMapBytes(t *testing.T) {
	objBytes := map[string][]byte{
		"k1": []byte("v1"),
		"k2": []byte("v2"),
	}
	out, _ := json.Marshal(objBytes)
	fmt.Println(string(out)) //output: {"k1":"djE=","k2":"djI="}

	objString := object.ConvertObjectByte2String(objBytes)
	out, _ = json.Marshal(objString)
	fmt.Println(string(out)) //output: {"k1":"v1","k2":"v2"}

	expectedOut := `{"k1":"v1","k2":"v2"}`
	if string(out) != expectedOut {
		t.Fatalf("expected out:%v, unexpected out: %v", expectedOut, string(out))
	}
}

func TestConvertOmitEmpty(t *testing.T) {
	type HistoryEvent struct {
		EventId *int64 `json:"eventId,omitempty"`
		TaskId  *int64 `json:"taskId,omitempty"`
	}
	i := int64(1)
	obj := HistoryEvent{
		EventId: &i,
	}
	if out, err := json.Marshal(object.ConvertObjectByte2String(obj)); err != nil {
		t.Fatal(err)
	} else {
		expectedOut := `{"eventId":1}`
		if string(out) != expectedOut {
			t.Fatalf("expected out:%v, unexpected out: %v", expectedOut, string(out))
		}
	}
}
