package constants

import "fmt"

const (
	SDK_LIB 				= "Golang sdk"
	// sdk 版本
	SDK_VERSION 			= "1.0.6"
	// 字符串最大长度
	STRING_VALUE_LEN_MAX 	= 8192
	// key最大长度
	KEY_VALUE_LEN_MAX 		= 255
	// 禁用的关键字
	KEYWORD_PATTERN 		= "^(^distinct_id$|^original_id$|^id$|^first_id$|^second_id$|^time$|^properties$|^users$|^events$|^event$|^user_id$|^date$|^datetime$)$"
	// key的正则过滤，防止奇怪的字符等等
	VALUE_PATTERN 			= "^[a-zA-Z_$][a-zA-Z\\d_$]{0,99}$"
	// 管道的长度
	CHANNEL_SIZE 			= 1000
)


// sensor枚举类型
type SensorType string
const (
	TRACK             SensorType = "sensor_track"             // 记录一个没有任何属性的事件
	ITEM_SET          SensorType = "sensor_item_set"          // 设置 item
	ITEM_DELETE       SensorType = "sensor_item_delete"       // 删除 item
	TRACK_SIGNUP      SensorType = "sensor_track_signup"      // 若属性包含 $time 字段，它会覆盖事件的默认时间属性，该字段只接受Date类型； 若属性包含 $project 字段，则它会指定事件导入的项目；
	PROFILE_SET       SensorType = "sensor_profile_set"       // 设置用户的属性。这个接口只能设置单个key对应的内容，同样，如果已经存在，则覆盖，否则，新创建
	PROFILE_SET_ONCE  SensorType = "sensor_profile_set_once"  // 首次设置用户的属性。这个接口只能设置单个key对应的内容。 与profileSet接口不同的是，如果key的内容之前已经存在，则不处理，否则，重新创建
	PROFILE_INCREMENT SensorType = "sensor_profile_increment" // 为用户的数值类型的属性累加一个数值，若该属性不存在，则创建它并设置默认值为0
	PROFILE_APPEND    SensorType = "sensor_profile_append"    // 为用户的数组类型的属性追加一个字符串
	PROFILE_UNSET     SensorType = "sensor_profile_unset"     // 删除用户某一个属性
	PROFILE_DELETE    SensorType = "sensor_profile_delete"    // 删除用户所有属性
)

func GetInstance(sensorType string) *SensorType {
	switch sensorType {
	case fmt.Sprintf("%v", TRACK): return convert(TRACK)
	case fmt.Sprintf("%v", ITEM_SET): return convert(ITEM_SET)
	case fmt.Sprintf("%v", ITEM_DELETE): return convert(ITEM_DELETE)
	case fmt.Sprintf("%v", TRACK_SIGNUP): return convert(TRACK_SIGNUP)
	case fmt.Sprintf("%v", PROFILE_SET): return convert(PROFILE_SET)
	case fmt.Sprintf("%v", PROFILE_SET_ONCE): return convert(PROFILE_SET_ONCE)
	case fmt.Sprintf("%v", PROFILE_INCREMENT): return convert(PROFILE_INCREMENT)
	case fmt.Sprintf("%v", PROFILE_APPEND): return convert(PROFILE_APPEND)
	case fmt.Sprintf("%v", PROFILE_UNSET): return convert(PROFILE_UNSET)
	case fmt.Sprintf("%v", PROFILE_DELETE): return convert(PROFILE_DELETE)
	default: return nil
	}
}
func convert(sen SensorType) *SensorType { return &sen }

// 神策相关数据
type SensorModel struct {
	ItemType string
	ItemId   string
	// 用户 ID 是否是登录 ID，false 表示该 ID 是一个匿名 ID
	IsLoginId bool
}