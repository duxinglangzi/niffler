package niffler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duxinglangzi/niffler/constants"
	"github.com/duxinglangzi/niffler/consumers"
	"github.com/duxinglangzi/niffler/util"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Niffler struct {
	Consumer consumers.Consumer
	// 公共属性, 用户在事件内自定义上传属性值时，会覆盖公共属性
	superProperties map[string]interface{}
	// 当前的项目名称
	ProjectName string
	// 正则
	keywordPattern *regexp.Regexp
	valuePattern   *regexp.Regexp
}

// 初始化 debug模式发送数据， 此模式为实时在线发送数据，仅在测试环节使用，上线前请切换至  InitConcurrentLoggingConsumer 方式
func InitDebuggerConsumer(projectName, serverUrl string) (*Niffler, error) {
	debuggerConsumer, err := consumers.InitDebugger(serverUrl)
	if err != nil {
		return nil, errors.New("init debug consumer error : " + err.Error())
	}
	return InitNiffler(projectName, debuggerConsumer)
}

// 初始化多进程方式写入文件
func InitConcurrentLoggingConsumer(projectName, fileName string, day bool) (*Niffler, error) {
	concurrentLog, err := consumers.InitConcurrentLoggingConsumer(fileName, day)
	if err != nil {
		return nil, errors.New("init concurrent log consumer error : " + err.Error())
	}
	return InitNiffler(projectName, concurrentLog)
}

// 初始化 console 打印日志
func InitConsoleConsumer(projectName string) (*Niffler, error) {
	consoleConsumer, err := consumers.InitConsole()
	if err != nil {
		return nil, errors.New("init console consumer error : " + err.Error())
	}
	return InitNiffler(projectName, consoleConsumer)
}

func InitNiffler(projectName string, consumer consumers.Consumer) (*Niffler, error) {
	if projectName == "" {
		return nil, errors.New("project name is null ")
	}
	marshal, _ := json.Marshal(map[string]string{
		"niffler_sdk_name":    constants.SDK_LIB,
		"niffler_sdk_version": constants.SDK_VERSION,
	})
	fmt.Fprintln(os.Stdout, string(marshal))
	return &Niffler{
		ProjectName:     projectName,
		Consumer:        consumer,
		superProperties: map[string]interface{}{},
		keywordPattern:  regexp.MustCompile(constants.KEYWORD_PATTERN),
		valuePattern:    regexp.MustCompile(constants.VALUE_PATTERN),
	}, nil
}

func (n *Niffler) Flush() {
	n.Consumer.Flush()
}

func (n *Niffler) Close() {
	n.Consumer.Close()
}

// 启动后可增加全局公共属性
func (n *Niffler) RegisterSuperProperties(superMap map[string]interface{}) {
	for key, val := range superMap {
		n.superProperties[key] = val
	}
}

// 清理公共属性
func (n *Niffler) ClearSuperProperties() {
	n.superProperties = make(map[string]interface{})
}

//  标记该事件信息为神策的数据，此事件数据不会进入es
//  distinctId 用户 ID
//  eventName  事件名称
//  properties 事件的属性
//  sensorModel 神策事件参数, 当调用 ITEM_SET、ITEM_DELETE、TRACK_SIGNUP时需要此参数对象
//  sensorProjectName 神策项目名称(必传)
//
//  return error:  eventName 或 properties 不符合命名规范和类型规范时抛出该异常
func (n *Niffler) AddSensorEvent(distinctId, sensorProjectName string, sensorType constants.SensorType, eventName string, sensorModel *constants.SensorModel, properties map[string]interface{}) error {
	if sensorProjectName == "" {
		return errors.New(" The 'sensorProjectName' is null.")
	}
	if fmt.Sprintf("%v", sensorType) == "" {
		sensorType = constants.TRACK
	}
	// return n.AddEvent(distinctId, projectName, fmt.Sprintf("%v", sensorType), eventName, sensorModel, properties)
	err := n.assertProperties(properties)
	if err != nil {
		return err
	}
	event := make(map[string]interface{})
	var eventProject string
	switch sensorType {
	case constants.TRACK:
		err = n.assertKey("Distinct Id", distinctId)
		if err != nil {
			return err
		}
		err = n.assertKeyWithRegex("Event Name", eventName)
		if err != nil {
			return err
		}
	case constants.ITEM_DELETE, constants.ITEM_SET:
		var itemType, itemId = "", ""
		if sensorModel != nil {
			itemType = sensorModel.ItemType
			itemId = sensorModel.ItemId
		}
		err = n.assertKeyWithRegex("Item Type", itemType)
		if err != nil {
			return err
		}
		err = n.assertKey("Item Id", itemId)
		if err != nil {
			return err
		}
		eventName = ""
		event["item_type"] = itemType
		event["item_id"] = itemId
		if val, ok := properties["$project"]; ok {
			eventProject = val.(string)
			delete(properties, "$project")
		}
		distinctId = ""
	
	case constants.TRACK_SIGNUP:
		err = n.assertKey("Distinct Id", distinctId)
		if err != nil {
			return err
		}
		originDistinctId := ""
		if sensorModel != nil {
			originDistinctId = sensorModel.OriginDistinctId
		}
		err = n.assertKey("Original Distinct Id", originDistinctId)
		if err != nil {
			return err
		}
		event["original_id"] = originDistinctId
		eventName = "$SignUp"
	case constants.PROFILE_SET, constants.PROFILE_SET_ONCE, constants.PROFILE_INCREMENT, constants.PROFILE_APPEND:
		err = n.assertKey("Distinct Id", distinctId)
		if err != nil {
			return err
		}
		eventName = ""
	case constants.PROFILE_UNSET:
		err = n.assertKey("Distinct Id", distinctId)
		if err != nil {
			return err
		}
		if len(properties) == 0 {
			return errors.New(" The properties is null.")
		}
		for k, v := range properties {
			if k != "$project" {
				val, ok := v.(bool)
				if ok && val {
					continue
				}
				return errors.New(" The property value of " + k + " should be true.")
			}
		}
		eventName = ""
	case constants.PROFILE_DELETE:
		eventName = ""
		properties = make(map[string]interface{})
	}
	
	//  判断sensor是否需要登录
	if sensorModel != nil && sensorModel.IsLoginId {
		event["$is_login_id"] = true
	}
	
	// event properties
	var eventProperties map[string]interface{}
	if n.superProperties != nil {
		eventProperties = util.DeepCopy(n.superProperties)
	} else {
		eventProperties = make(map[string]interface{})
	}
	
	if properties != nil {
		eventProperties = util.MergeCopy(properties, eventProperties)
	}
	// Event time
	eventTime := n.extractEventTime(eventProperties)
	if distinctId != "" {
		event["distinct_id"] = distinctId
	}
	event["type"] = fmt.Sprintf("%v", sensorType)
	if eventName != "" {
		event["event"] = eventName
	}
	if eventProject != "" {
		event["project"] = eventProject
	}
	event["sensor_project_name"] = sensorProjectName
	event["time"] = eventTime
	event["lib"] = n.getLibProperties()
	event["properties"] = eventProperties
	event["project_name"] = n.ProjectName
	return n.Consumer.Send(event)
	
}

//  distinctId 用户 ID
//  eventName  事件名称
//  properties 事件的属性
//  return error:  eventName 或 properties 不符合命名规范和类型规范时抛出该异常
func (n *Niffler) AddUserEvent(distinctId, eventName string, properties map[string]interface{}) error {
	return n.AddEvent(distinctId, "user", eventName, properties)
}

// 商品事件
func (n *Niffler) AddGoodsEvent(distinctId, eventName string, properties map[string]interface{}) error {
	return n.AddEvent(distinctId, "goods", eventName, properties)
}

// 订单事件
func (n *Niffler) AddOrderEvent(distinctId, eventName string, properties map[string]interface{}) error {
	return n.AddEvent(distinctId, "order", eventName, properties)
}

// 购物车事件
func (n *Niffler) AddCartEvent(distinctId, eventName string, properties map[string]interface{}) error {
	return n.AddEvent(distinctId, "cart", eventName, properties)
}

// 记录一个拥有一个或多个属性的事件。属性取值可接受类型为 string,int64,float64,bool,time.Time,[]string ;
// 若属性包含 $time 字段，则它会覆盖事件的默认时间属性，该字段只接受 time.Time 类型.
//
//  distinctId 用户 ID
//  eventType  事件类型(如: 用户、商品、订单、购物车 等等)
//  eventName  事件名称
//  properties 事件的属性
//  returns error  distinctId 或 eventName、properties 不符合命名规范和类型规范时抛出该异常
func (n *Niffler) AddEvent(distinctId, eventType, eventName string, properties map[string]interface{}) error {
	err := n.assertKey("Distinct Id", distinctId)
	if err != nil {
		return err
	}
	err = n.assertProperties(properties)
	if err != nil {
		return err
	}
	err = n.assertKeyWithRegex("Event Name", eventName)
	if err != nil {
		return err
	}
	
	// event properties
	var eventProperties map[string]interface{}
	if n.superProperties != nil {
		eventProperties = util.DeepCopy(n.superProperties)
	} else {
		eventProperties = make(map[string]interface{})
	}
	
	if properties != nil {
		eventProperties = util.MergeCopy(properties, eventProperties)
	}
	// Event time
	eventTime := n.extractEventTime(eventProperties)
	event := make(map[string]interface{})
	event["distinct_id"] = distinctId
	event["type"] = eventType
	event["event"] = eventName
	event["time"] = eventTime
	event["lib"] = n.getLibProperties()
	event["event_id"] = strings.ReplaceAll(util.NewUUID(), "-", "")
	event["properties"] = eventProperties
	event["project_name"] = n.ProjectName
	return n.Consumer.Send(event)
}

// 日志打印， 当不想让日志信息进入到es时，传入空字符串即可
func (n *Niffler) Log(logType string, properties map[string]interface{}) error {
	if properties == nil || len(properties) < 1 {
		return errors.New("properties is null ")
	}
	logMap := make(map[string]interface{})
	if logType != "" {
		logMap["type"] = logType
	}
	logMap["time"] = util.NowMilliseconds()
	logMap["log_id"] = strings.ReplaceAll(util.NewUUID(), "-", "")
	logMap["properties"] = util.DeepCopy(properties)
	logMap["project_name"] = n.ProjectName
	return n.Consumer.Send(logMap)
}

// if event time is null , return current milliseconds
func (n *Niffler) extractEventTime(m map[string]interface{}) int64 {
	if t, contain := m["$time"]; contain {
		if v, ok := t.(int64); ok {
			delete(m, "$time")
			return v
		} else {
			return util.NowMilliseconds()
		}
	}
	return util.NowMilliseconds()
}

func (n *Niffler) getLibProperties() map[string]string {
	resultMap := make(map[string]string)
	resultMap["$lib"] = constants.SDK_LIB
	resultMap["$lib_version"] = constants.SDK_VERSION
	if pc, file, line, ok := runtime.Caller(3); ok {
		fun := runtime.FuncForPC(pc)
		resultMap["$lib_detail"] = fmt.Sprintf("%s##%s##%d", file, fun.Name(), line)
	}
	return resultMap
}

func (n *Niffler) assertProperties(properties map[string]interface{}) error {
	if properties == nil {
		return nil
	}
	for k, v := range properties {
		err := n.assertKeyWithRegex("property", k)
		if err != nil {
			return err
		}
		// 单独检查 $is_login_id
		if k == "$is_login_id" {
			switch v.(type) {
			case bool:
				continue
			default:
				return errors.New(" The property value of '$is_login_id' should be bool.")
			}
		}
		
		// 检查value 值，只支持部分类型
		switch v.(type) {
		case bool:
		case int64:
		case float64:
		case string:
			if len(v.(string)) > constants.STRING_VALUE_LEN_MAX {
				return errors.New(" The maximum length of the string is " + strconv.Itoa(constants.STRING_VALUE_LEN_MAX))
			}
		case []string:
			for _, e := range v.([]string) {
				if len(e) > constants.STRING_VALUE_LEN_MAX {
					return errors.New(" The maximum length of the string is " + strconv.Itoa(constants.STRING_VALUE_LEN_MAX))
				}
			}
		case time.Time:
			properties[k] = v.(time.Time).Format("2006-01-02 15:04:05.999")
		default:
			return errors.New("The property '" + k + "' should be a basic type:  string,int64,float64,bool,time.Time,[]string")
		}
	}
	return nil
}

func (n *Niffler) assertKey(typeName, key string) error {
	if key == "" {
		return errors.New("The " + typeName + " is empty.")
	}
	if len(key) > constants.KEY_VALUE_LEN_MAX {
		return errors.New("The " + typeName + " is too long, max length is " + strconv.Itoa(constants.KEY_VALUE_LEN_MAX))
	}
	return nil
}

func (n *Niffler) assertKeyWithRegex(typeName, key string) error {
	err := n.assertKey(typeName, key)
	if err != nil {
		return err
	}
	if !n.checkPattern([]byte(key)) {
		return errors.New("The " + typeName + " '" + key + "' is invalid.")
	}
	return err
}

func (n *Niffler) checkPattern(name []byte) bool {
	return !n.keywordPattern.Match(name) && n.valuePattern.Match(name)
}
