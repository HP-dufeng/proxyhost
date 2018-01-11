package proxyhost

import (
	"context"

	"github.com/chromedp/chromedp/runner"
)

//Browser open or close browser
type Browser interface {
	Open(url string) error
	Close() error
}

type chrome struct {
	runner *runner.Runner
	ctxt   context.Context
	proxy  string
}

func (c *chrome) Open(url string) error {
	c.ctxt = context.Background()

	// create chrome instance
	var err error
	c.runner, err = runner.Run(c.ctxt, runner.Proxy(c.proxy), runner.StartURL(url))
	if err != nil {
		return err
	}

	return nil
}

func (c *chrome) Close() error {
	// shutdown chrome
	err := c.runner.Shutdown(c.ctxt)
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
