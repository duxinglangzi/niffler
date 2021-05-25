package consumers

type Consumer interface {
	Send(data map[string]interface{}) error
	Flush() error
	Close() error
}