package main

import (
	"fmt"
	"net/http"
	"sync"
)

// BackendServer represents a backend server.
type BackendServer struct {
	URL   string
	Alive bool
	Mutex sync.RWMutex
}

// BackendServers is a collection of backend servers.
type BackendServers []*BackendServer

// Next returns the next available backend server.
func (servers BackendServers) Next() *BackendServer {
	for {
		for _, server := range servers {
			server.Mutex.RLock()
			alive := server.Alive
			server.Mutex.RUnlock()

			if alive {
				return server
			}
		}
	}
}

// Update checks the status of each backend server and updates its Alive status.
func (servers BackendServers) Update() {
	for _, server := range servers {
		resp, err := http.Head(server.URL)
		if err != nil {
			server.Mutex.Lock()
			server.Alive = false
			server.Mutex.Unlock()
		} else {
			server.Mutex.Lock()
			server.Alive = resp.StatusCode == 200
			server.Mutex.Unlock()
		}
	}
}

func main() {
	// Create a slice of backend servers.
	servers := BackendServers{
		&BackendServer{URL: "http://backend1.example.com"},
		&BackendServer{URL: "http://backend2.example.com"},
		&BackendServer{URL: "http://backend3.example.com"},
	}

	// Start a goroutine to periodically check the status of each backend server.
	go func() {
		for {
			servers.Update()
			time.Sleep(time.Second)
		}
	}()

	// Start the HTTP server and balance incoming requests across the backend servers.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server := servers.Next()
		if server == nil {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		resp, err := http.Get(server.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Forward the response from the backend server to the client.
		if _, err := io.Copy(w, resp.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Load balancer started on port 8080")
	http.ListenAndServe(":8080", nil)
}
