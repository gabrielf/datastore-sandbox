package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/gabrielf/datastore-sandbox/app"
	. "github.com/onsi/gomega"
	"google.golang.org/appengine/aetest"
)

func TestWriteLogEntry(t *testing.T) {
	RegisterTestingT(t)

	instance, err := aetest.NewInstance(nil)
	Expect(err).ToNot(HaveOccurred())
	defer instance.Close()

	req, err := instance.NewRequest("GET", "/log", nil)
	Expect(err).ToNot(HaveOccurred())

	res := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(res, req)

	Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
}

func TestWriteLogEntryInTransaction(t *testing.T) {
	RegisterTestingT(t)

	instance, err := aetest.NewInstance(nil)
	Expect(err).ToNot(HaveOccurred())
	defer instance.Close()

	req, err := instance.NewRequest("GET", "/logtrans", nil)
	Expect(err).ToNot(HaveOccurred())

	res := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(res, req)

	Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
	Expect(strings.TrimSpace(res.Body.String())).To(Equal(`{"LogEntries":1}`))
}
