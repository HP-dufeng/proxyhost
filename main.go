// Step 0: read host config
// Step 1: create proxy server
// Step 2: change incoming request host
// Step 3: start chrome with proxy

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/chromedp/chromedp"
	"github.com/elazarl/goproxy"
)

var hostPath = "./hosts"

func main() {
	// initialize our channels
	signals := make(chan os.Signal)
	done := make(chan bool)

	// hook them up to the signals lib
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// if a signal is caught by this go routine
	// it will write to done
	go catchSig(signals, done)

	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	err = c.Run(ctxt, click())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press ctrl-c to terminate...")
	// the program blocks until someone writes to done
	<-done
	fmt.Println("Done!")

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
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

func click() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://golang.org/pkg/time/`),
		// chromedp.WaitVisible(`#footer`),
		// chromedp.Click(`#pkg-overview`, chromedp.NodeVisible),
		// chromedp.Sleep(10 * time.Second),
	}
}

func proxy() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			// r.Header.Set("X-GoProxy","yxorPoG-X")
			// address, err := url.Parse(r.URL.Path)
			// if(err != nil) {
			// 	fmt.Print(err)
			// }
			log.Println("请求的url:", r.URL.Hostname())
			ip, ok := findIp(r.URL.Hostname())
			if ok {
				r.Host = ip
				fmt.Println("Proxy on: ", ip)
			}
			return r, nil
		})

	log.Fatal(http.ListenAndServe(":9999", proxy))
}

func findIp(host string) (string, bool) {
	hosts, err := readHosts()
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	val, ok := hosts[host]
	return val, ok
}

func readHosts() (map[string]string, error) {
	hosts := make(map[string]string)
	file, err := os.Open(hostPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) <= 0 {
			continue
		}
		if i := strings.IndexAny("#", line); i >= 0 {
			continue
		}
		f := strings.Fields(line)
		hosts[f[1]] = f[0]
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}
