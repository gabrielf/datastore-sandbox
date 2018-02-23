package task

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"
)

var AsyncFunc = delay.Func("do-async-stuff", DoAsyncStuff)

func DoAsyncStuff(ctx context.Context, duration time.Duration, times int) error {
	form := url.Values{}
	form.Add("sleep", duration.String())
	form.Add("times", strconv.Itoa(times))

	sleepUrl := "https://datastore-sandbox-1114.appspot.com/sleep"
	if appengine.IsDevAppServer() {
		sleepUrl = "http://localhost:8080/sleep"
	}

	log.Infof(ctx, "Creating HTTP POST to %s", sleepUrl)

	req, err := http.NewRequest("POST", sleepUrl, strings.NewReader(form.Encode()))
	if err != nil {
		log.Errorf(ctx, "Error creating POST request: %v", err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx, _ = context.WithTimeout(ctx, time.Second*600)
	c := urlfetch.Client(ctx)
	res, err := c.Do(req)
	if err != nil {
		log.Errorf(ctx, "Error executing POST request: %v", err.Error())
		return err
	}
	if res.StatusCode != http.StatusOK {
		log.Errorf(ctx, "HTTP error: %v", res.StatusCode)
		return err
	}
	return nil
}

func TriggerSleepTask(w http.ResponseWriter, r *http.Request) {
	sleepMs := r.FormValue("sleep")
	times := r.FormValue("times")

	if sleepMs == "" {
		http.Error(w, "Missing sleep parameter", http.StatusBadRequest)
		return
	}

	ctx := appengine.NewContext(r)
	t := taskqueue.NewPOSTTask("/sleep", map[string][]string{"sleep": {sleepMs}, "times": {times}})
	createdTask, err := taskqueue.Add(ctx, t, "slow-queue")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func TriggerUnstableTask(w http.ResponseWriter, r *http.Request) {
	failCount := r.FormValue("failCount")

	ctx := appengine.NewContext(r)
	t := taskqueue.NewPOSTTask("/unstable", map[string][]string{"failCount": {failCount}})

	createdTask, err := taskqueue.Add(ctx, t, "slow-queue")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func TriggerSleepTaskUsingDelay(w http.ResponseWriter, r *http.Request) {
	duration, err := time.ParseDuration(r.FormValue("sleep"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	times, err := strconv.Atoi(r.FormValue("times"))
	if err != nil {
		times = 1
	}

	t, err := AsyncFunc.Task(duration, times)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := appengine.NewContext(r)
	t, err = taskqueue.Add(ctx, t, "slow-queue")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Sleep(w http.ResponseWriter, r *http.Request) {
	d, err := time.ParseDuration(r.FormValue("sleep"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	times, err := strconv.Atoi(r.FormValue("times"))
	if err != nil {
		times = 1
	}

	ctx := appengine.NewContext(r)
	for i := 1; i <= times; i++ {
		log.Infof(ctx, "Round %d of %d", i, times)
		log.Infof(ctx, "Will sleep for %v", d)
		time.Sleep(d)
		log.Infof(ctx, "Waking up")
	}
	log.Infof(ctx, "Slept for %f seconds", d.Seconds()*float64(times))
}

func Unstable(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	failCount, err := strconv.Atoi(r.FormValue("failCount"))
	if err != nil {
		log.Infof(ctx, "Couldn't parse failCount %s setting to 0", r.FormValue("failCount"))
		failCount = 0
	}
	retryCount, err := strconv.Atoi(r.Header.Get("X-AppEngine-TaskRetryCount"))

	log.Infof(ctx, "failCount: %v", failCount)
	log.Infof(ctx, "retryCount: %v", retryCount)

	log.Infof(ctx, "----------------------------")
	log.Infof(ctx, "Headers")

	for header, values := range r.Header {
		log.Infof(ctx, "%v: %v", header, values)
	}

	if retryCount < failCount {
		log.Infof(ctx, "Failing!")
		http.Error(w, "Failure", http.StatusInternalServerError)
	} else {
		log.Infof(ctx, "Success")
	}
}

func TriggerProtectedTask(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	t := taskqueue.NewPOSTTask("/protected", nil)

	createdTask, err := taskqueue.Add(ctx, t, "slow-queue")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Protected(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Header.Get("X-AppEngine-TaskRetryCount") == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Infof(ctx, "Doing secret stuff")
}

func TriggerParamsTask(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	t := taskqueue.NewPOSTTask("/params", url.Values{"name": []string{"John Doe"}})

	createdTask, err := taskqueue.Add(ctx, t, "slow-queue")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Params(_ http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if err := r.ParseForm(); err != nil {
		log.Errorf(ctx, "Err: %+v", err)
		return
	}

	log.Infof(ctx, "Got parameters: %+v", r.PostForm)
	log.Infof(ctx, "Got headers: %+v", r.Header)
}
