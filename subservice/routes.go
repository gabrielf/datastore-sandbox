package app

import (
	"fmt"
	"net/http"

	"github.com/gabrielf/datastore-sandbox/src/learning"
	"github.com/gabrielf/datastore-sandbox/src/task"
)

func init() {
	RegisterRoutes()
}

func RegisterRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from subservice")
	})

	http.HandleFunc("/meta", learning.Meta)
	http.HandleFunc("/echo", learning.Echo)
	http.HandleFunc("/whichServiceForTask", task.WhichServiceDoesATaskRunOn)
}
