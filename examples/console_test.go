package examples

import (
	"github.com/duxinglangzi/niffler"
	"github.com/duxinglangzi/niffler/constants"
	"testing"
	"time"
)

func TestConsoleConsumer(t *testing.T) {
	niffer,err := niffler.InitConsoleConsumer("test")
	if err != nil {
		t.Log(err.Error())
	}
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
		"keyvalue":     "values",
		// "original_id":	"asdfasf",
	}
	// err = niffer.AddCartEvent("distinctId", "test_event", properties)
	// if err != nil {
	// 	t.Log(err.Error())
	// }
	// 日志打印
	// niffer.Log("info",properties)
	
	sensorModel := &constants.SensorModel{
		ItemId:   "asfhkasjdhf",
		ItemType: "asdkufasdf",
	}
	err = niffer.AddSensorEvent("dis_id","production", constants.PROFILE_SET, "test_event_ss", sensorModel, properties)
	if err != nil {
		t.Log(err.Error())
	}
	
}
