package constants

const (
	SDK_LIB 				= "Golang sdk"
	// sdk 版本
	SDK_VERSION 			= "1.0.5"
	// 字符串最大长度
	STRING_VALUE_LEN_MAX 	= 8192
	// key最大长度
	KEY_VALUE_LEN_MAX 		= 255
	// 禁用的关键字
	KEYWORD_PATTERN 		= "^(^distinct_id$|^time$|^properties$|^users$|^events$|^event$|^user_id$|^date$)$"
	// key的正则过滤，防止奇怪的字符等等
	VALUE_PATTERN 			= "^[a-zA-Z_$][a-zA-Z\\d_$]{0,99}$"
)

const (
	// 管道的长度
	CHANNEL_SIZE = 1000
)
