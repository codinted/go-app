package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tomasen/realip"
)

var mB = make(map[string]MetricsBody)

// MetricsBody is ...
type MetricsBody struct {
	PercentageCPUUsed    int64 `json:"percentage_cpu_used"`
	PercentageMemoryUsed int64 `json:"percentage_memory_used"`
}

//ResultSet is
type ResultSet struct {
	IP        string `json:"ip"`
	MxCPU     int64  `json:"max_cpu"`
	MaxMemory int64  `json:"max_memory"`
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.RemoteAddr)
	fmt.Fprint(w, "{STATUS: UP}")
}

func metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	ip := realip.FromRequest(r)
	var m MetricsBody
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if val, ok := mB[ip]; ok {
		if val.PercentageCPUUsed > m.PercentageCPUUsed {
			m.PercentageCPUUsed = val.PercentageCPUUsed
		}
		if val.PercentageMemoryUsed > m.PercentageMemoryUsed {
			m.PercentageMemoryUsed = val.PercentageMemoryUsed
		}
	}
	mB[ip] = m
	fmt.Fprint(w, 200, "\n")

}

func results(w http.ResponseWriter, r *http.Request) {
	var s []ResultSet
	for ip, key := range mB {
		s = append(s, ResultSet{ip, key.PercentageCPUUsed, key.PercentageMemoryUsed})

	}
	fmt.Println(s)
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func main() {

	http.HandleFunc("/status", status)
	http.HandleFunc("/metrics", metrics)
	http.HandleFunc("/report", results)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
