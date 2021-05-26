package examples

import (
	"github.com/duxinglangzi/niffler"
	"testing"
	"time"
)

func TestConsoleConsumer(t *testing.T) {
	niffer := niffler.InitConsoleConsumer("test")
	superMap := make(map[string]interface{})
	superMap["aaaaa"] = int64(2345)
	superMap["bbbbbb"] = float64(2345.23123)
	niffer.RegisterSuperProperties(superMap)
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
	// 日志打印
	niffer.Log("info",properties)
}
