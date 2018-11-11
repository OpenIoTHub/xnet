package xhttp

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func HttpHeadlessGet(url, proxy, sel string, timeout time.Duration) (string, error) {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		return "", err
	}

	// run task list
	var buf []byte
	err = c.Run(ctxt, headlessGet(url, sel, &buf))
	if err != nil {
		return "", err
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		return "", err
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		return "", err
	}

}

func headlessGet(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second),
		chromedp.WaitVisible(sel, chromedp.ByID),
	}
}