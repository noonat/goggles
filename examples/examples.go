package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	port = flag.Int("port", 8080, "Port to listen on")
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	paths := map[string]string{
		"/examples/cube/cube.js":     "bin/cube.js",
		"/examples/cube/cube.js.map": "bin/cube.js.map",
		"/examples/obj/obj.js":       "bin/obj.js",
		"/examples/obj/obj.js.map":   "bin/obj.js.map",
		"/": "src/github.com/noonat/goggles",
	}
	goPath := os.Getenv("GOPATH")
	for route, path := range paths {
		absPath, err := filepath.Abs(filepath.Join(goPath, path))
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		if route == "/" {
			mux.Handle(route, http.FileServer(http.Dir(absPath)))
			continue
		}
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, absPath)
		})
	}

	fmt.Printf("Listening on http://0.0.0.0:%d. Press Ctrl+C to stop.\n", *port)
	if err := http.ListenAndServe(":"+strconv.Itoa(*port), mux); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
