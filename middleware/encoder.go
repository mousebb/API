package middleware

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

type HtmlReponse struct {
	Glob string
	Body interface{}
}

const (
	textXML  = "text/xml"
	appXML   = "application/xml"
	textYaml = "text/yaml"
	plain    = "text/plain"
	html     = "text/html"
)

// Encode ...
func Encode(r *http.Request, w http.ResponseWriter, obj interface{}) error {
	ct := r.Header.Get("Content-Type")
	var err error

	switch strings.ToLower(ct) {
	case appXML:
		err = toXML(r, w, obj)
	case textXML:
		err = toXML(r, w, obj)
	case textYaml:
		err = toYAML(r, w, obj)
	case plain:
		err = toText(r, w, obj)
	case html:
		err = toHTML(r, w, obj)
	default:
		err = toJSON(r, w, obj)
	}

	return err
}

func toXML(r *http.Request, w http.ResponseWriter, obj interface{}) error {
	w.Header().Set("Content-Type", "text/xml")
	return xml.NewEncoder(w).Encode(obj)
}

func toJSON(r *http.Request, w http.ResponseWriter, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(obj)
}

func toYAML(r *http.Request, w http.ResponseWriter, obj interface{}) error {
	out, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/yaml")

	_, err = w.Write([]byte(out))

	return err
}

func toText(r *http.Request, w http.ResponseWriter, obj interface{}) error {
	var buf bytes.Buffer
	if _, err := fmt.Fprintf(&buf, "%s\n", obj); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write(buf.Bytes())

	return err
}

func toHTML(r *http.Request, w http.ResponseWriter, obj interface{}) error {

	if reflect.TypeOf(obj) != reflect.TypeOf(HtmlReponse{}) {
		return toJSON(r, w, obj)
	}

	hr := obj.(HtmlReponse)

	t, err := template.ParseGlob(hr.Glob)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html")

	return t.Execute(w, hr.Body)
}
