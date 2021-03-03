package trie

import (
	"net/http"
)

/*
   Creation Time: 2021 - Mar - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// ParamStore should be completed by http.ResponseWriter to support dynamic path parameters.
// See the `Writer` type for more.
// This interface can be implemented by custom response writers.
// Example of implementation: to change where and how the parameters are stored and retrieved.
type ParamStore interface {
	Set(key string, value string)
	Get(key string) string
	GetAll() []ParamEntry
}

// GetParam returns the path parameter value based on its key, i.e
// "/hello/:name", the parameter key is the "name".
// For example if a route with pattern of "/hello/:name" is inserted to the `Trie` or handlded by the `Mux`
// and the path "/hello/kataras" is requested through the `Mux#ServeHTTP -> Trie#search`
// then the `GetParam("name")` will return the value of "kataras".
// If not associated value with that key is found then it will return an empty string.
//
// The function will do its job only if the given "w" http.ResponseWriter interface is a `ParamStore`.
func GetParam(w http.ResponseWriter, key string) string {
	if store, ok := w.(ParamStore); ok {
		return store.Get(key)
	}

	return ""
}

// GetParams returns all the available parameters based on the "w" http.ResponseWriter which should be a ParamStore.
//
// The function will do its job only if the given "w" http.ResponseWriter interface is a `ParamStore`.
func GetParams(w http.ResponseWriter) []ParamEntry {
	if store, ok := w.(ParamStore); ok {
		return store.GetAll()
	}

	return nil
}

// SetParam sets manually a parameter to the "w" http.ResponseWriter which should be a ResponseWriter.
// This is not commonly used by the end-developers,
// unless sharing values(string messages only) between handlers is absolutely necessary.
func SetParam(w http.ResponseWriter, key, value string) bool {
	if store, ok := w.(ParamStore); ok {
		store.Set(key, value)
		return true
	}

	return false
}

// ParamEntry holds the Key and the Value of a named path parameter.
type ParamEntry struct {
	Key   string
	Value string
}

// Writer is the muxie's specific ResponseWriter to hold the path parameters.
// Usage: use this to cast a handler's `http.ResponseWriter` and pass it as an embedded parameter to custom response writer
// that will be passed to the next handler in the chain.
type Writer struct {
	http.ResponseWriter
	params []ParamEntry
}

var _ ParamStore = (*Writer)(nil)

// Set implements the `ParamsSetter` which `Trie#search` needs to store the parameters, if any.
// These are decoupled because end-developers may want to use the trie to design a new Mux of their own
// or to store different kind of data inside it.
func (pw *Writer) Set(key, value string) {
	if ln := len(pw.params); cap(pw.params) > ln {
		pw.params = pw.params[:ln+1]
		p := &pw.params[ln]
		p.Key = key
		p.Value = value
		return
	}

	pw.params = append(pw.params, ParamEntry{
		Key:   key,
		Value: value,
	})
}

// Get returns the value of the associated parameter based on its key/name.
func (pw *Writer) Get(key string) string {
	n := len(pw.params)
	for i := 0; i < n; i++ {
		if kv := pw.params[i]; kv.Key == key {
			return kv.Value
		}
	}

	return ""
}

// GetAll returns all the path parameters keys-values.
func (pw *Writer) GetAll() []ParamEntry {
	return pw.params
}

func (pw *Writer) reset(w http.ResponseWriter) {
	pw.ResponseWriter = w
	pw.params = pw.params[0:0]
}
