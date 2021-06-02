package examples

import (
	"github.com/duxinglangzi/niffler"
	"testing"
	"time"
)

func TestDebugConsumer(t *testing.T) {
	niffer ,err := niffler.InitDebuggerConsumer("test", "http://127.0.0.1:8071/api/v1/order/abc")
	if err != nil {
		t.Log(err.Error())
	}
	properties := map[string]interface{}{
		"Name":          "name value",
		"IsTrueOrFalse": false,
		"Age":           int64(234),
		"aa":            34.0034,
		"array":         []string{"1", "323", "545"},
		"timeTT":        time.Now(),
	}
	err = niffer.AddUserEvent("distinctId", "test_event", properties)
	if err != nil {
		t.Log(err.Error())
	}
}
