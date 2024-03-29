package trie

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

var routes = BuildFromMap(map[string]string{
	"":                "root", // as universal prefix
	"/api/user":       "user",
	"/api/user/list":  "usersList",
	"/api/group":      "group",
	"/api/group/list": "groupsList",
	"/api/articles/":  "articles",
})

func Example_routing() {
	var inputs = []string{
		"/api/user/list",
		"/api/user/li",
		"/api/group",
		"/api/unknown",
	}

	for _, inp := range inputs {
		exact, ok := routes.GetByString(inp)
		route := routes.GetAllByString(inp)
		if ok {
			fmt.Printf("%-17s:\thandler %-10s\t(route %v)\n", inp, exact, route)
		} else {
			fmt.Printf("%-17s:\thandler not found\t(route %v)\n", inp, route)
		}
	}

	// Output:
	// /api/user/list   :	handler usersList 	(route [root user usersList])
	// /api/user/li     :	handler not found	(route [root user])
	// /api/group       :	handler group     	(route [root group])
	// /api/unknown     :	handler not found	(route [root])
}

func Example_handlers() {
	tr := &Trie[http.Handler]{}

	tr.PutString("/", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprint(w, "Index page")
	}))
	tr.PutString("/home", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprint(w, "Home page")
	}))
	tr.PutString("/profile", http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprint(w, "Profile page")
	}))

	requests := []*http.Request{
		httptest.NewRequest("GET", "/profile", nil),
		httptest.NewRequest("GET", "/", nil),
		httptest.NewRequest("GET", "/home", nil),
	}

	for _, req := range requests {
		if f, ok := tr.GetByString(req.RequestURI); ok {
			rec := httptest.NewRecorder()
			f.ServeHTTP(rec, req)
			resp, err := io.ReadAll(rec.Result().Body)
			fmt.Println(string(resp), err)
		}
	}

	// Output:
	// Profile page <nil>
	// Index page <nil>
	// Home page <nil>
}
