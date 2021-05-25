package consumers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ConsoleConsumer struct {
	standardOutput io.Writer
}

func InitConsole() (*ConsoleConsumer, error) {
	return &ConsoleConsumer{standardOutput: os.Stdout}, nil
}

func (c *ConsoleConsumer) Send(data map[string]interface{}) error {
	sendData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(c.standardOutput, string(sendData))
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsoleConsumer) Flush() error {
	// do nothing
	return nil
}
func (c *ConsoleConsumer) Close() error {
	// do nothing
	return nil
}
