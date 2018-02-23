package learning

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-errors/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func CreateLogEntry(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	rootKey, _, err := getRootEntity(ctx)
	if err != nil {
		http.Error(w, errors.Wrap(err, 0).ErrorStack(), http.StatusInternalServerError)
		return
	}

	logEntry := RootLog{
		UpdatedAt: time.Now(),
	}
	logEntryKey := datastore.NewIncompleteKey(ctx, "RootLog", rootKey)
	if _, err = datastore.Put(ctx, logEntryKey, &logEntry); err != nil {
		http.Error(w, errors.Wrap(err, 0).ErrorStack(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(logEntry); err != nil {
		http.Error(w, errors.Wrap(err, 0).ErrorStack(), http.StatusInternalServerError)
		return
	}
}

func CreateLogEntryInTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var outerRoot *Root

	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		rootKey, root, err := getRootEntity(ctx)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		logEntry := RootLog{
			UpdatedAt: time.Now(),
		}
		logEntryKey := datastore.NewIncompleteKey(ctx, "RootLog", rootKey)
		root.LogEntries += 1

		if _, err = datastore.Put(ctx, logEntryKey, &logEntry); err != nil {
			return errors.New(err)
		}
		if _, err = datastore.Put(ctx, rootKey, root); err != nil {
			return errors.New(err)
		}

		outerRoot = root
		return nil
	}, &datastore.TransactionOptions{XG: false})

	if err != nil {
		log.Errorf(ctx, errors.Wrap(err, 0).ErrorStack())
		http.Error(w, errors.Wrap(err, 0).ErrorStack(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(outerRoot); err != nil {
		http.Error(w, errors.Wrap(err, 0).ErrorStack(), http.StatusInternalServerError)
		return
	}
}

func getRootEntity(ctx context.Context) (*datastore.Key, *Root, error) {
	root := Root{}
	key := datastore.NewKey(ctx, "Root", "root", 0, nil)
	if err := datastore.Get(ctx, key, &root); err != nil {
		if err == datastore.ErrNoSuchEntity {
			key, err = datastore.Put(ctx, key, &root)
			if err != nil {
				return nil, nil, errors.New(err)
			}
			return key, &root, nil
		} else {
			return nil, nil, errors.New(err)
		}
	}
	return key, &root, nil
}

type Root struct {
	LogEntries int
}

type RootLog struct {
	UpdatedAt time.Time
}
