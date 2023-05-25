package trie

import (
	"fmt"
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
