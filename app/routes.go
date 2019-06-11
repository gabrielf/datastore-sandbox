package app

import (
	"net/http"

	"github.com/gabrielf/datastore-sandbox/src/categories"
	"github.com/gabrielf/datastore-sandbox/src/learning"
	"github.com/gabrielf/datastore-sandbox/src/neterrors"
	"github.com/gabrielf/datastore-sandbox/src/task"
)

func init() {
	RegisterRoutes()
}

func RegisterRoutes() {
	http.HandleFunc("/", categories.Index)
	http.HandleFunc("/test", categories.TestEventualConsistency)

	http.HandleFunc("/meta", learning.Meta)
	http.HandleFunc("/echo", learning.Echo)
	http.HandleFunc("/log", learning.CreateLogEntry)
	http.HandleFunc("/logtrans", learning.CreateLogEntryInTransaction)

	// Task related routes
	http.HandleFunc("/triggerSleepTask", task.TriggerSleepTask)
	http.HandleFunc("/triggerSleepTaskUsingDelay", task.TriggerSleepTaskUsingDelay)
	http.HandleFunc("/sleep", task.Sleep)
	http.HandleFunc("/triggerUnstableTask", task.TriggerUnstableTask)
	http.HandleFunc("/unstable", task.Unstable)
	http.HandleFunc("/triggerProtectedTask", task.TriggerProtectedTask)
	http.HandleFunc("/protected", task.Protected)
	http.HandleFunc("/triggerParamsTask", task.TriggerParamsTask)
	http.HandleFunc("/params", task.Params)
	http.HandleFunc("/whichServiceForTask", task.WhichServiceDoesATaskRunOn)
	http.HandleFunc("/taskWithETA", task.TaskWithETA)

	// Errors
	http.HandleFunc("/errors/timeout1", neterrors.Timeout1)
	http.HandleFunc("/errors/timeout2", neterrors.Timeout2)
	http.HandleFunc("/errors/timeout3", neterrors.Timeout3)
	http.HandleFunc("/errors/timeout4", neterrors.Timeout4)
	http.HandleFunc("/errors/timeout5", neterrors.Timeout5)
	http.HandleFunc("/errors/timeout6", neterrors.Timeout6)
	http.HandleFunc("/errors/connection_close1", neterrors.ConnectionClose1)
	http.HandleFunc("/errors/connection_close2", neterrors.ConnectionClose2)
}
