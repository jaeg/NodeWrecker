package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var shouldStop = false

var threads = flag.Int("threads", 4, "Number of threads to run")
var msSleep = flag.Int64("sleep", 1, "milliseconds to sleep")
var shouldEscalate = flag.Bool("escalate", false, "Keep creating threads")
var escalateRate = flag.Int64("escalate-rate", 1000, "milliseconds between creating new threads")
var stringLength = flag.Int("string-length", 1000, "length of randomly generated string")
var abuseMemory = flag.Bool("abuse-memory", false, "if true nodewrecker will store all generated values in memory")

var threadCount = 0
var memory sync.Map

func main() {
	flag.Parse()
	fmt.Println("Let's wreck it")

	//Capture sigterm
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Sign to stop")
		shouldStop = true
	}()

	var wg sync.WaitGroup
	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go cpuThread(&wg)
	}

	if *shouldEscalate {
		wg.Add(1)
		go escalate(&wg)
	}

	wg.Wait()
	fmt.Println("Stopping")
}

func escalate(wg *sync.WaitGroup) {
	for !shouldStop {
		wg.Add(1)
		go cpuThread(wg)
		time.Sleep(time.Duration(*escalateRate) * time.Millisecond)
	}
	wg.Done()
}

func generateRandomString(length int) (out string) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i := 0; i < length; i++ {
		out += string(chars[rand.Intn(len(chars))])
	}

	return
}

func cpuThread(wg *sync.WaitGroup) {
	id := threadCount
	threadCount++
	fmt.Println("Thread ", id, " has started")
	for {
		if shouldStop {
			fmt.Println("Got signal to stop.")
			break
		}

		a := generateRandomString(*stringLength)
		b := encodeBase64([]byte(a))
		c := decodeBase64(b)
		if *abuseMemory {
			memory.Store(a, b)
		}
		fmt.Println("Thread ", id, " : ", string(c))
		time.Sleep(time.Duration(*msSleep) * time.Millisecond)
	}

	wg.Done()
}

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
