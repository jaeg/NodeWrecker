package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var shouldStop = false
var ended = false

var threads = flag.Int("threads", 4, "Number of threads to run")
var msSleep = flag.Int64("sleep", 1, "milliseconds to sleep")
var shouldEscalate = flag.Bool("escalate", false, "Keep creating threads")
var escalateRate = flag.Int64("escalate-rate", 1000, "milliseconds between creating new threads")
var stringLength = flag.Int("string-length", 1000, "length of randomly generated string")
var abuseMemory = flag.Bool("abuse-memory", false, "if true nodewrecker will store all generated values in memory")
var chaos = flag.Bool("chaos", false, "When true stress testing starts and stops randomly.")
var minDuration = flag.Int64("min-duration", 10, "minimum seconds a test lasts")
var maxDuration = flag.Int64("max-duration", 60, "max seconds a test lasts")
var maxDelay = flag.Int64("max-delay", 10, "max seconds between tests")
var minDelay = flag.Int64("min-delay", 1, "minimum seconds between tests")
var verbose = flag.Bool("verbose", false, "output everything from threads")
var output = flag.Bool("output", false, "Output content generated from threads")
var outputDir = flag.String("output-dir", "./", "directory to put output")

var threadCount = 0
var memory sync.Map

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	fmt.Println("Let's wreck it")

	//Capture sigterm
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Sign to stop")
		shouldStop = true
		ended = true
	}()

	if *chaos {
		shouldStop = true
		go makeChaos()
	}

	var wg sync.WaitGroup

	for !ended {
		if !shouldStop {
			for i := 0; i < *threads; i++ {
				wg.Add(1)
				go cpuThread(&wg)
			}

			if *shouldEscalate {
				wg.Add(1)
				go escalate(&wg)
			}

			wg.Wait()
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Println("Stopping")
}

func makeChaos() {
	for !ended {
		shouldStop = true
		threadCount = 0
		memory = sync.Map{}
		max := *maxDelay - *minDelay
		delay := *minDelay + rand.Int63n(max)
		fmt.Println("Chaos is sleeping for ", delay, " seconds..")
		time.Sleep(time.Duration(delay) * time.Second)

		shouldStop = false
		max = *maxDuration - *minDuration
		delay = *minDuration + rand.Int63n(max)
		fmt.Println("Chaos is awake for ", delay, " seconds!")
		time.Sleep(time.Duration(delay) * time.Second)
	}
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
	var f *os.File
	if *output {
		var err error
		f, err = os.Create(*outputDir + "/" + strconv.Itoa(id) + "-" + time.Now().String() + ".txt")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	for {
		if shouldStop {
			fmt.Println("Thread ", id, ": Got signal to stop.")
			break
		}

		a := generateRandomString(*stringLength)
		b := encodeBase64([]byte(a))
		c := decodeBase64(b)
		if *abuseMemory {
			memory.Store(a, b)
		}
		if *verbose {
			fmt.Println("Thread ", id, " : ", string(c))
		}
		if *output {
			f.WriteString(string(c))
		}

		time.Sleep(time.Duration(*msSleep) * time.Millisecond)
	}

	if *output {
		f.Close()
	}

	wg.Done()
}

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
