package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	PORT = "8001"
	OPT  = "dir"
	PATH = "."
)

// liveserve .
// liveserve ~/app/index.html
// liveserve ~/app
func handleArgs() error {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Println("E: Could not parse absolute path", err.Error())
		return err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		log.Println("E: Reading path", err.Error())
		return err
	}

	if info.IsDir() {
		OPT = "dir"
		PATH = absPath
	} else {
		OPT = "file"
		PATH = absPath
	}

	return nil
}

func main() {
	if err := handleArgs(); err != nil {
		return
	}

	// headers to prevent default browser caching
	switch OPT {
	case "file":
		fmt.Println("Serving file", PATH)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			http.ServeFile(w, r, PATH)
		})

	case "dir":
		fmt.Println("Serving dir", PATH)
		fs := http.FileServer(http.Dir(PATH))
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			fs.ServeHTTP(w, r)
		}))
	}

	fmt.Printf("Live at http://localhost:%s/\n", PORT)
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		panic(err)
	}
}
