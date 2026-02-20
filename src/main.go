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
	if len(os.Args) < 2 {
		return nil
	}

	absPath, err := filepath.Abs(os.Args[1])
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

	switch OPT {
	case "file":
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			http.ServeFile(w, r, PATH)
		})

	case "dir":
		fs := http.FileServer(http.Dir(PATH))
		http.Handle("/", fs)

	}

	fmt.Printf("Live at http://localhost:%s/\n", PORT)
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		panic(err)
	}
}
