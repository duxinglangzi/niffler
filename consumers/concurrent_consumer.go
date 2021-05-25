package consumers

import (
	"encoding/json"
	"fmt"
	"github.com/duxinglangzi/niffler/constants"
	"os"
	"sync"
	"syscall"
	"time"
)

type ConcurrentLoggingConsumer struct {
	// 文件名称
	FileName 	string
	// 是否按天分隔， 否则按小时分隔
	Day 		bool
	cw  		*ConcurrentWriter
}

func InitConcurrentLoggingConsumer(fileName string, Day bool) (*ConcurrentLoggingConsumer, error) {
	cw, err := InitConcurrentWriter(fileName, Day)
	if err != nil {
		return nil, err
	}
	clc := &ConcurrentLoggingConsumer{FileName: fileName, Day: Day, cw: cw}
	return clc, nil
}

func (c *ConcurrentLoggingConsumer) Send(data map[string]interface{}) error {
	return c.cw.Write(data)
}

func (c *ConcurrentLoggingConsumer) Flush() error {
	c.cw.Flush()
	return nil
}

func (c *ConcurrentLoggingConsumer) Close() error {
	c.cw.Close()
	return nil
}

type ConcurrentWriter struct {
	rec     chan string
	filName string
	file    *os.File
	
	day  int
	hour int
	
	dayRotate bool
	
	wg sync.WaitGroup
}

func (w *ConcurrentWriter) Write(data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.rec <- string(jsonData)
	return nil
}

func (w *ConcurrentWriter) Flush() {
	w.file.Sync()
}

func (w *ConcurrentWriter) Close() {
	close(w.rec)
	w.wg.Wait()
}

func (w *ConcurrentWriter) intRotate() error {
	fileName := ""
	
	if w.file != nil {
		w.file.Close()
	}
	
	now := time.Now()
	today := now.Format("2006-01-02")
	w.day = time.Now().Day()
	
	if w.dayRotate {
		fileName = fmt.Sprintf("%s.%s", w.filName, today)
	} else {
		hour := now.Hour()
		w.hour = hour
		fileName = fmt.Sprintf("%s.%s.%02d", w.filName, today, hour)
	}
	
	fd, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
		return err
	}
	w.file = fd
	
	return nil
}

func InitConcurrentWriter(filName string, day bool) (*ConcurrentWriter, error) {
	w := &ConcurrentWriter{
		filName:   filName,
		day:       time.Now().Day(),
		hour:      time.Now().Hour(),
		dayRotate: day,
		rec:       make(chan string, constants.CHANNEL_SIZE),
	}
	
	if err := w.intRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "ConcurrentWriter(%q): %s\n", w.filName, err)
		return nil, err
	}
	
	w.wg.Add(1)
	
	go func() {
		defer func() {
			if w.file != nil {
				w.file.Sync()
				w.file.Close()
			}
			w.wg.Done()
		}()
		
		for {
			select {
			case rec, ok := <-w.rec:
				if !ok {
					return
				}
				
				now := time.Now()
				
				if (!w.dayRotate && now.Hour() != w.hour) || (now.Day() != w.day) {
					if err := w.intRotate(); err != nil {
						fmt.Fprintf(os.Stderr, "ConcurrentWriter(%q): %s\n", w.filName, err)
						return
					}
				}
				
				syscall.Flock(int(w.file.Fd()), syscall.LOCK_EX)
				_, err := fmt.Fprintln(w.file, rec)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ConcurrentWriter(%q): %s\n", w.filName, err)
					syscall.Flock(int(w.file.Fd()), syscall.LOCK_UN)
					return
				}
				syscall.Flock(int(w.file.Fd()), syscall.LOCK_UN)
			}
		}
	}()
	
	return w, nil
}
