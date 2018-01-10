package proxyhost

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

//Browser open or close browser
type Browser interface {
	Open(url string) error
	Close() error
}

type chrome struct {
	cdp   *chromedp.CDP
	ctxt  context.Context
	proxy string
}

func (c *chrome) Open(url string) error {
	// create context
	// ctxt, cancel := context.WithCancel(context.Background())
	// defer cancel()
	c.ctxt = context.Background()

	// create chrome instance
	var err error
	c.cdp, err = chromedp.New(c.ctxt, chromedp.WithRunnerOptions(runner.Proxy(c.proxy)))
	if err != nil {
		return err
	}

	// run task list
	err = c.cdp.Run(c.ctxt, chromedp.Tasks{
		chromedp.Navigate(url),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *chrome) Close() error {
	// shutdown chrome
	err := c.cdp.Shutdown(c.ctxt)
	if err != nil {
		return err
	}

	// wait for chrome to finish
	err = c.cdp.Wait()
	if err != nil {
		return err
	}

	return nil
}

//var ChromeBrowser Browser = (*chrome)(nil)

//Chrome init a chrome instance with proxyURL
func Chrome(proxy string) Browser {
	return &chrome{
		proxy: proxy,
	}
}
