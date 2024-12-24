package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

type PingResp struct {
	Status string `json:"status_code"`
	Ret    struct {
		ResponseCode string      `json:"resp_code,omitempty"`
		ResponseMsg  string      `json:"resp_msg,omitempty"`
		BuildDate    string      `json:"build_date"`
		Version      string      `json:"version"`
		RunMode      string      `json:"run_mode"`
		ServiceName  string      `json:"service_name"`
		StartTime    time.Time   `json:"start_time,omitempty"`
		Info         interface{} `json:"info,omitempty"`
	} `json:"ret,omitempty"`
}

var glSysStatus PingResp

func pingHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(glSysStatus)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HelloWorld)
	r.HandleFunc("/ping", pingHandler)
	http.Handle("/", r)

	fmt.Println("Server is running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
