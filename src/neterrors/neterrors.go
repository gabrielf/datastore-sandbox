package neterrors

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func IsTimeoutError(err error) bool {
	// url.Error's implementation of Timeout() doesn't work if the wrapped
	// error is an AppEngine APIError which implements IsTimeout instead
	// of Timeout. Therefore we have to do the unwrapping ourselves.
	if ue, ok := err.(*url.Error); ok {
		return IsTimeoutError(ue.Err)
	}

	// ApiError, CallError
	if t, ok := err.(interface {
		IsTimeout() bool
	}); ok {
		return t.IsTimeout()
	}
	// Context deadline exceeded and URL errors
	if t, ok := err.(interface {
		Timeout() bool
	}); ok {
		return t.Timeout()
	}
	return false
}

func Timeout1(w http.ResponseWriter, r *http.Request) {
	ctx := gaeContext(r)

	for i := 0; ; i++ {
		log.Infof(ctx, "now fetch that thing: %d", i)
		_, err := urlfetch.Client(ctx).Get("https://httpbin.org/delay/10")
		if err != nil {
			log.Infof(ctx, "there is an error: %#v", err)
			log.Infof(ctx, "and it seems to be: %s", err)
			log.Infof(ctx, "is a timeout error: %v", IsTimeoutError(err))
			if ctx.Err() == context.DeadlineExceeded {
				log.Infof(ctx, "the error was a DeadlineExceeded error")
			}
			outputErr(w, err)
			break
		}
	}
}

func Timeout2(w http.ResponseWriter, r *http.Request) {
	ctx := gaeContext(r)

	_, err := urlfetch.Client(ctx).Get("https://httpbin.org/drip?code=200&numbytes=5&duration=10")
	outputErr(w, err)
}

func Timeout3(w http.ResponseWriter, r *http.Request) {
	ctx := gaeContext(r)

	client := urlfetch.Client(ctx)
	client.Timeout = time.Second
	_, err := client.Get("https://httpbin.org/delay/10")
	outputErr(w, err)
}

func Timeout4(w http.ResponseWriter, r *http.Request) {
	ctx := gaeContext(r)

	client := urlfetch.Client(ctx)
	client.Timeout = time.Second
	_, err := client.Get("https://httpbin.org/drip?code=200&numbytes=5&duration=10")
	outputErr(w, err)
}

func Timeout5(w http.ResponseWriter, r *http.Request) {
	ctx := gaeContext(r)

	client := urlfetch.Client(ctx)
	client.Timeout = time.Second
	_, err := client.Get("https://httpbin.org/drip?code=200&numbytes=5&duration=10")
	outputErr(w, err)
}

func Timeout6(w http.ResponseWriter, r *http.Request) {
	ctx := gaeContext(r)

	client := urlfetch.Client(ctx)
	_, err := client.Get("http://localhost:8081/headers/timeout")
	outputErr(w, err)
}

func ConnectionClose1(w http.ResponseWriter, r *http.Request) {
	ctx := gaeContext(r)

	client := urlfetch.Client(ctx)
	_, err := client.Get("http://localhost:8081/close")
	outputErr(w, err)
}

func ConnectionClose2(w http.ResponseWriter, r *http.Request) {
	ctx := gaeContext(r)

	client := urlfetch.Client(ctx)
	_, err := client.Get("http://localhost:8081/headers/close")
	outputErr(w, err)
}

func outputErr(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf("%s\n\n%#v", err.Error(), err), 500)
}

func gaeContext(r *http.Request) context.Context {
	ctx := appengine.NewContext(r)
	ctx, _ = context.WithTimeout(ctx, time.Hour)
	return ctx
}
