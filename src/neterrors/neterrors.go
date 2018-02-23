package neterrors

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func Timeout1(w http.ResponseWriter, r *http.Request) {
	ctx := gaeContext(r)

	_, err := urlfetch.Client(ctx).Get("https://httpbin.org/delay/10")
	outputErr(w, err)
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
	ctx, _ = context.WithTimeout(ctx, time.Second)
	return ctx
}
