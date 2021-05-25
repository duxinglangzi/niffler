package consumers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type DebugConsumer struct {
	ServerUrl string
}

func InitDebugger(serverUrl string, ) (*DebugConsumer, error) {
	urls, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}
	return &DebugConsumer{ServerUrl: urls.String()}, nil
}

func (d *DebugConsumer) Send(data map[string]interface{}) error {
	sendData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	statusCode, body, err := d.DebugPostRequest(d.ServerUrl, string(sendData))
	if statusCode == 200 {
		fmt.Fprintf(os.Stdout, "Valid message: %s\n", body)
		return nil
	} else {
		fmt.Fprintf(os.Stderr, "Invalid message: %s\n", body)
		fmt.Fprintf(os.Stderr, "Reponse_code: %d\n", statusCode)
		fmt.Fprintf(os.Stderr, "Reponse_content: %s\n", body)
	}
	if statusCode >= 300 {
		return errors.New(" Bad http status. ")
	}
	return err
}
func (d *DebugConsumer) Flush() error {
	// do nothing
	return nil
}
func (d *DebugConsumer) Close() error {
	// do nothing
	return nil
}

func (d *DebugConsumer) DebugPostRequest(url, args string) (StatusCode int, ResponseBody string, Err error) {
	var resp *http.Response
	data := bytes.NewBufferString(args)
	req, _ := http.NewRequest("POST", url, data)
	client := &http.Client{Timeout: time.Second * 6}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return 0, "", err
	}
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return resp.StatusCode, string(body), nil
	}
	return 0, "", err
}
