package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var wait sync.WaitGroup
var sig = make(chan bool, 1)
var log chan []byte
var Writer io.Writer = os.Stdout
var ticker = time.NewTicker(time.Second * 1)
var useFile = false

func Open(logMax int) {
	if logMax < 1 {
		logMax = 2048
	}

	log = make(chan []byte, logMax)
	wait.Add(1)
	go loop()
}

func loop() {
	defer wait.Done()
	for {
		select {
		case now := <-ticker.C:
			if useFile {
				if w, ok := Writer.(*File); ok {
					if err := w.Check(now); err != nil {
						fmt.Printf("check file error: %s\n", err)
					}
				}
			}
		case <-sig:
			return
		case logBytes, ok := <-log:
			if !ok {
				return
			}

			if _, err := Writer.Write(logBytes); err != nil {
				fmt.Printf("record log error: %s\n", err)
			}
		}
	}
}

func UseFile(path string) {
	file := NewFile(path)
	if file == nil {
		return
	}

	if err := file.Open(); err != nil {
		return
	}

	Writer = file
	useFile = true
}

func Append(logBytes []byte) {
	if logBytes == nil {
		return
	}

	log <- logBytes
}

func Close() {
	sig <- true
	wait.Wait()
	if useFile {
		if w, ok := Writer.(*File); ok {
			w.Close()
		}
	}
}
