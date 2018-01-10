// Step 0: read host config
// Step 1: create proxy server
// Step 2: resolve incoming request host
// Step 3: start chrome with proxy

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fengdu/proxyhost"
)

func main() {
	url := flag.String("url", "about:blank", "start url path")
	flag.Parse()
	port := strconv.Itoa(proxyhost.RandomPort())

	go proxyhost.StartServer(port)

	// initialize our channels
	signals := make(chan os.Signal)
	done := make(chan bool)

	// hook them up to the signals lib
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// if a signal is caught by this go routine
	// it will write to done
	go catchSig(signals, done)

	// create chrome instance
	browser := proxyhost.Chrome("http://127.0.0.1:" + port)
	browser.Open(*url)
	defer browser.Close()

	fmt.Println("Press ctrl-c to terminate...")
	// the program blocks until someone writes to done
	<-done
	fmt.Println("Done!")
}

func catchSig(ch chan os.Signal, done chan bool) {
	// block on waiting for a signal
	sig := <-ch
	// print it when it's received
	fmt.Println("\nsig received:", sig)

	// we can set up handlers for all types of
	// sigs here
	switch sig {
	case syscall.SIGINT:
		fmt.Println("handling a SIGINT now!")
	case syscall.SIGTERM:
		fmt.Println("handling a SIGTERM in an entirely different way!")
	default:
		fmt.Println("unexpected signal received")
	}

	// terminate
	done <- true
}
