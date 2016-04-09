// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatic(t *testing.T) {
	tests := []struct {
		prefix string
		path   string
		code   int
	}{
		{"/", "/server.go", http.StatusOK},
		{"/", "/nonexisting.go", http.StatusNotFound},

		{"/public", "/public/server.go", http.StatusOK},
		{"/public", "/public/./server.go", http.StatusOK},
		{"/public", "/public/folder/../server.go", http.StatusOK},
		{"/public", "/public/nonexisting.go", http.StatusNotFound},
		{"/public", "/public/../server.go", http.StatusNotFound},
	}

	root := http.Dir(".")

	for idx, test := range tests {
		s := NewServer()
		s.Static(test.prefix, root)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, test.path, nil)

		s.ServeHTTP(w, r)

		if w.Code != test.code {
			t.Errorf("Expected code %d, is %d (test no. %d)", test.code, w.Code, idx)
		}

		if test.code == http.StatusOK && w.Body.Len() == 0 {
			t.Errorf("Expected non-empty body (test no. %d)", idx)
		}
	}
}

func TestRenderer(t *testing.T) {
	server := NewServer()
	locals := &struct{ Title string }{"MyTitle"}

	// Setup renderer with initial template cache
	server.ViewRoot = "/views"
	server.Renderer = NewStdRenderer(".tpl", true)
	server.Renderer.(*stdRenderer).tpl = template.Must(template.New("my.tpl").Parse("{{.Title}}"))

	// Setup route
	server.GetFunc("/myfile", func(res ResponseWriter, req *Request) {
		res.Render("my", locals)
	})

	r, _ := http.NewRequest(http.MethodGet, "/myfile", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK (200), not %s (%d)", http.StatusText(w.Code), w.Code)
	}

	if content := w.Body.String(); content != "MyTitle" {
		t.Errorf("Expected content to be 'MyTitle' not '%s'", content)
	}
}