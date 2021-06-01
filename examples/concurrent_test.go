package examples

import (
	"github.com/duxinglangzi/niffler"
	"testing"
	"time"
)

func TestConcurrentLogConsumer(t *testing.T) {
	niffer := niffler.InitConcurrentLoggingConsumer("test", "./event_log", false)
	defer niffer.Close() // 切记一定要记得关闭、落盘
	properties := map[string]interface{}{
		"Name":          "name value",
		"IsTrueOrFalse": false,
		"Age":           int64(234),
		"aa":            34.0034,
		"array":         []string{"1", "323", "545"},
		"timeTT":        time.Now(),
	}
	err := niffer.AddCartEvent("distinctId", "test_event", properties)
	if err != nil {
		t.Log(err.Error())
	}
}
