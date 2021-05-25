package examples

import (
	"github.com/duxinglangzi/niffler"
	"testing"
	"time"
)

func TestConsoleConsumer(t *testing.T) {
	niffer := niffler.InitConsoleConsumer("test")
	properties := map[string]interface{}{
		"Name":          "name value",
		"IsTrueOrFalse": false,
		"Age":           int64(234),
		"aa":            34.0034,
		"array":         []string{"1", "323", "545"},
		"timeTT":        time.Now(),
	}
	err := niffer.AddCartEvent("distinctId", "test-event", properties)
	if err != nil {
		t.Log(err.Error())
	}
}
