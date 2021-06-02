# niffler-go

## 本SDK 日志收集工具在结构上分为两层， 第一层为业务模块，第二层为事件名称

## 本SDK 仅提供日志收集操作 ，具体用法参见以下 demo code :

### debug 模式调试代码，此模式下，日志信息直接以post形式发送至指定服务器. (_注意此模式不要应用于生产环境_)
````go
func TestCase(t *testing.T) {
    properties := map[string]interface{}{
		"Name":          "name value",
		"IsTrueOrFalse": false,
		"Age":           int64(234),
		"aa":            34.0034,
		"array":         []string{"1", "323", "545"},
		"timeTT":        time.Now(),
        "test_value":    "value not null",
	}

    // 通过debug模式调试代码
    niffer,err := niffler.InitDebuggerConsumer("test", "http://127.0.0.1:8071/api/v1/order/abc")
    if err != nil {
    	t.Log(err.Error())
    }
	err = niffer.AddUserEvent("distinctId", "test_event", properties)
	if err != nil {
		t.Log(err.Error())
	}

    // 打印数据到文件内
    nifferLog,err1 := niffler.InitConcurrentLoggingConsumer("test", "./event_log",false)
    if err1 != nil {
    	t.Log(err1.Error())
    }
	defer nifferLog.Close() // 切记一定要记得关闭、落盘
    err = nifferLog.AddUserEvent("distinctId", "test_event_log", properties)
	if err != nil {
		t.Log(err.Error())
	}
    
    // 在控制台打印
    nifflerConsole,err2 := niffler.InitConsoleConsumer("test")
    if err2 != nil {
    	t.Log(err2.Error())
    }
    consoleErr := nifflerConsole.AddUserEvent("distinctId", "test_event_log", properties)
	if consoleErr != nil {
		t.Log(consoleErr.Error())
	}
    
    // 打印日志模式 ， 当logType为空字符串时，则不进入es
    niffer.Log("info",properties)




}

````



### 发送日志数据，程序会自动校验不可使用的关键字及字段名称、长度、类型
`关键字: distinct_id、time、properties、events、event、user_id、date、datetime`
`String 类型,长度不应超过 8192 `
`支持的数据类型: string,int64,float64,bool,time.Time,[]string `
`所有key值命名需遵循驼峰命名法, 如: test_event_log `


### 线上环境中产生的日志，需要单独通过 filebeat 软件监控并增量的传输至kafka ，再由数据团队解析使用， 安装及配置filebeat方式如下:
```
# 1、首先下载 filebeat 软件, 个人偏好下载放置于  /usr/local 内， 方便好找
wget -P /usr/local https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-7.11.1-linux-x86_64.tar.gz

# 2、解压
tar -xvf filebeat-7.11.1-linux-x86_64.tar.gz

# 3、进入文件夹 ， 修改 filebeat.yml 文件内容为: 
#=========================== Filebeat inputs =============================
filebeat.inputs:
- type: log

  # Change to true to enable this input configuration.
  enabled: true

  # 配置监控的日志地址， 请自行修改监控的日志地址
  paths:
    - /usr/local/log-data/event_log.*
    #- c:\programdata\elasticsearch\logs\*

  close_renamed: true
  clean_removed: true
  close_removed: true

# -------------------------------- Kafka Output -------------------------------
output.kafka:
  # Boolean flag to enable or disable the output module.
  enabled: true

  # 设置 kafka的访问地址， 集群可以设置多个。
  hosts: ["XXX.XXX:9092"]

  # The Kafka topic used for produced events. 
  # 设置kafka 的主题, 需要自行指定
  topic: filebeats-topic-test


#================================ Processors =====================================
# 可以在这里删除一些不需要的字段，
processors:
  - add_host_metadata: ~
  - add_cloud_metadata: ~
  - drop_fields:
     fields: ["beat.hostname", "beat.name", "beat.version", "beat","host","agent","input","ecs","@metadata"]


# 4、创建logs文件夹，方便存放 filebaat 日志文件
mkdir logs 

# 5、启动执行
nohup ./filebeat -e -c filebeat.yml >> logs/output.log 2>&1 &
```

