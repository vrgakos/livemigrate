package main

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	lock sync.Mutex
	number int = 0
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	number++
	fmt.Fprintf(w, "Hi, the number is: %d!", number)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}