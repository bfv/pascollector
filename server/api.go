package server

import (
	"fmt"
	"net/http"
	"strconv"
)

func startListener() {

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		fmt.Println("stop request received...")
		CloseDB()
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET status")
		fmt.Fprintf(w, "OK")
	})

	http.ListenAndServe(":"+strconv.Itoa(Config.Port), nil)
}
