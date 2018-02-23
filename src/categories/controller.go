package categories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type Category struct {
	Key       *datastore.Key `datastore:"-"`
	Ancestors []string
	Name      string
}

func Index(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	if r.Method == "POST" {
		Create(w, r)
	} else {
		Get(w, r)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var result interface{}
	if r.Form.Get("name") != "" {
		result = findByName(ctx, r.Form.Get("name"))
	} else if r.Form.Get("ancestor") != "" {
		result = findByAncestorName(ctx, r.Form.Get("ancestor"))
	} else if r.Form.Get("parent") != "" {
		result = findByParentName(ctx, r.Form.Get("parent"))
	} else {
		result = findAll(ctx)
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	parent := findByName(ctx, r.Form.Get("parent"))

	category := Category{
		Name:      r.Form.Get("name"),
		Ancestors: getAncestorPath(parent),
	}
	categoryKey := datastore.NewIncompleteKey(ctx, "Category", nil)
	var err error
	if category.Key, err = datastore.Put(ctx, categoryKey, &category); err != nil {
		panic(err)
	}

	if err := json.NewEncoder(w).Encode(category); err != nil {
		panic(err)
	}
}

func findAll(ctx context.Context) []Category {
	var categories []Category
	keys, err := datastore.NewQuery("Category").GetAll(ctx, &categories)
	if err != nil {
		panic(err)
	}
	for i, _ := range keys {
		categories[i].Key = keys[i]
	}
	return categories
}

func findByAncestorName(ctx context.Context, ancestorName string) []Category {
	var ancestor = findByName(ctx, ancestorName)

	return findByAncestor(ctx, ancestor)
}

func findByAncestor(ctx context.Context, ancestor *Category) []Category {
	if ancestor == nil {
		return []Category{}
	}

	var categories []Category
	keys, err := datastore.NewQuery("Category").Filter("Ancestors=", ancestor.Key.Encode()).GetAll(ctx, &categories)
	if err != nil {
		panic(err)
	}
	for i, _ := range keys {
		categories[i].Key = keys[i]
	}

	return categories
}

func findByParentName(ctx context.Context, parentName string) []Category {
	var parent = findByName(ctx, parentName)

	ancestors := findByAncestor(ctx, parent)
	var parents []Category
	for i, category := range ancestors {
		log.Infof(ctx, "Checking ancestor: %+v for parentName", category, parentName)

		if category.Ancestors[len(category.Ancestors)-1] == parent.Key.Encode() {
			parents = append(parents, ancestors[i])
		}
	}
	return parents
}

func findByName(ctx context.Context, name string) *Category {
	if name == "" {
		return nil
	}

	var categories []Category
	keys, err := datastore.NewQuery("Category").Filter("Name=", name).GetAll(ctx, &categories)
	if err != nil {
		panic(err)
	}
	for i, _ := range categories {
		categories[i].Key = keys[i]
	}

	return &categories[0]
}

func getAncestorPath(parent *Category) []string {
	if parent == nil {
		return []string{}
	}

	ancestors := []string{}
	ancestors = append(ancestors, parent.Ancestors...)
	ancestors = append(ancestors, parent.Key.Encode())

	return ancestors
}

func TestEventualConsistency(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	now := time.Now()
	name := fmt.Sprintf("Category-%d", now.UnixNano())

	category := Category{
		Name: name,
	}
	categoryKey := datastore.NewIncompleteKey(ctx, "Category", nil)
	var err error
	if category.Key, err = datastore.Put(ctx, categoryKey, &category); err != nil {
		panic(err)
	}

	for {
		var categories []Category
		keys, err := datastore.NewQuery("Category").Filter("Name=", name).GetAll(ctx, &categories)
		if err != nil {
			panic(err)
		}
		if len(keys) == 1 {
			fmt.Fprintf(w, "%dms\n", time.Now().Sub(now).Nanoseconds()/1000/1000)
			return
		}
	}
}
