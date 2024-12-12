package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var wait sync.WaitGroup
var sig = make(chan bool, 1)
var log chan []byte
var Writer io.Writer = os.Stdout

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
		case <-sig:
			return
		case logBytes, ok := <-log:
			if !ok {
				return
			}

			if _, err := Writer.Write(logBytes); err != nil {
				fmt.Printf("record log error: %s\n", err)
			}
			if _, err := Writer.Write([]byte("\n")); err != nil {
				fmt.Printf("record log error: %s\n", err)
			}
		}
	}
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
}
