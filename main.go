// package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"net"
// 	"net/http"
// 	"net/http/httputil"
// 	"net/url"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// )

// const (
// 	Attempts int = iota
// 	Retry
// )

// // Backend holds the data about a server
// type Backend struct {
// 	URL          *url.URL
// 	Alive        bool
// 	mux          sync.RWMutex
// 	ReverseProxy *httputil.ReverseProxy
// }

// // SetAlive for this backend
// func (b *Backend) SetAlive(alive bool) {
// 	b.mux.Lock()
// 	b.Alive = alive
// 	b.mux.Unlock()
// }

// // IsAlive returns true when backend is alive
// func (b *Backend) IsAlive() (alive bool) {
// 	b.mux.RLock()
// 	alive = b.Alive
// 	b.mux.RUnlock()
// 	return
// }

// // ServerPool holds information about reachable backends
// type ServerPool struct {
// 	backends []*Backend
// 	current  uint64
// }

// // AddBackend to the server pool
// func (s *ServerPool) AddBackend(backend *Backend) {
// 	s.backends = append(s.backends, backend)
// }

// // NextIndex atomically increase the counter and return an index
// func (s *ServerPool) NextIndex() int {
// 	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
// }

// // MarkBackendStatus changes a status of a backend
// func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
// 	for _, b := range s.backends {
// 		if b.URL.String() == backendUrl.String() {
// 			b.SetAlive(alive)
// 			break
// 		}
// 	}
// }

// // GetNextPeer returns next active peer to take a connection
// func (s *ServerPool) GetNextPeer() *Backend {
// 	// loop entire backends to find out an Alive backend
// 	next := s.NextIndex()
// 	l := len(s.backends) + next // start from next and move a full cycle
// 	for i := next; i < l; i++ {
// 		idx := i % len(s.backends)     // take an index by modding
// 		if s.backends[idx].IsAlive() { // if we have an alive backend, use it and store if its not the original one
// 			if i != next {
// 				atomic.StoreUint64(&s.current, uint64(idx))
// 			}
// 			return s.backends[idx]
// 		}
// 	}
// 	return nil
// }

// // HealthCheck pings the backends and update the status
// func (s *ServerPool) HealthCheck() {
// 	for _, b := range s.backends {
// 		status := "up"
// 		alive := isBackendAlive(b.URL)
// 		b.SetAlive(alive)
// 		if !alive {
// 			status = "down"
// 		}
// 		log.Printf("%s [%s]\n", b.URL, status)
// 	}
// }

// // GetAttemptsFromContext returns the attempts for request
// func GetAttemptsFromContext(r *http.Request) int {
// 	if attempts, ok := r.Context().Value(Attempts).(int); ok {
// 		return attempts
// 	}
// 	return 1
// }

// // GetAttemptsFromContext returns the attempts for request
// func GetRetryFromContext(r *http.Request) int {
// 	if retry, ok := r.Context().Value(Retry).(int); ok {
// 		return retry
// 	}
// 	return 0
// }

// // lb load balances the incoming request
// func lb(w http.ResponseWriter, r *http.Request) {
// 	attempts := GetAttemptsFromContext(r)
// 	if attempts > 3 {
// 		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
// 		http.Error(w, "Service not available", http.StatusServiceUnavailable)
// 		return
// 	}

// 	peer := serverPool.GetNextPeer()
// 	if peer != nil {
// 		peer.ReverseProxy.ServeHTTP(w, r)
// 		return
// 	}
// 	http.Error(w, "Service not available", http.StatusServiceUnavailable)
// }

// // isAlive checks whether a backend is Alive by establishing a TCP connection
// func isBackendAlive(u *url.URL) bool {
// 	timeout := 2 * time.Second
// 	conn, err := net.DialTimeout("tcp", u.Host, timeout)
// 	if err != nil {
// 		log.Println("Site unreachable, error: ", err)
// 		return false
// 	}
// 	defer conn.Close()
// 	return true
// }

// // healthCheck runs a routine for check status of the backends every 2 mins
// func healthCheck() {
// 	t := time.NewTicker(time.Minute * 2)
// 	for {
// 		select {
// 		case <-t.C:
// 			log.Println("Starting health check...")
// 			serverPool.HealthCheck()
// 			log.Println("Health check completed")
// 		}
// 	}
// }

// var serverPool ServerPool

// func main() {
// 	var serverList string
// 	var port int
// 	flag.StringVar(&serverList, "backends", "wss://polygon-mainnet.g.alchemy.com/v2/3z5KecNiLnFFeomwszKT2GPiInCZdfiE,wss://eth-mainnet.g.alchemy.com/v2/s8zGtV5Jtr7TPRnOmRW6LaS4I4tMjptX", "Load balanced backends, use commas to separate")
// 	// flag.StringVar(&serverList, "backends", "wss://eth-mainnet.g.alchemy.com/v2/s8zGtV5Jtr7TPRnOmRW6LaS4I4tMjptX", "Load balanced backends, use commas to separate")
// 	flag.IntVar(&port, "port", 3030, "Port to serve")
// 	flag.Parse()

// 	if len(serverList) == 0 {
// 		log.Fatal("Please provide one or more backends to load balance")
// 	}

// 	// parse servers
// 	tokens := strings.Split(serverList, ",")
// 	for _, tok := range tokens {
// 		serverUrl, err := url.Parse(tok)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
// 		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
// 			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
// 			retries := GetRetryFromContext(request)
// 			if retries < 3 {
// 				select {
// 				case <-time.After(10 * time.Millisecond):
// 					ctx := context.WithValue(request.Context(), Retry, retries+1)
// 					proxy.ServeHTTP(writer, request.WithContext(ctx))
// 				}
// 				return
// 			}

// 			// after 3 retries, mark this backend as down
// 			serverPool.MarkBackendStatus(serverUrl, false)

// 			// if the same request routing for few attempts with different backends, increase the count
// 			attempts := GetAttemptsFromContext(request)
// 			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
// 			ctx := context.WithValue(request.Context(), Attempts, attempts+1)
// 			lb(writer, request.WithContext(ctx))
// 		}

// 		serverPool.AddBackend(&Backend{
// 			URL:          serverUrl,
// 			Alive:        true,
// 			ReverseProxy: proxy,
// 		})
// 		log.Printf("Configured server: %s\n", serverUrl)
// 	}

// 	// create http server
// 	server := http.Server{
// 		Addr:    fmt.Sprintf(":%d", port),
// 		Handler: http.HandlerFunc(lb),
// 	}

// 	// start health checking
// 	go healthCheck()

// 	log.Printf("Load Balancer started at :%d\n", port)
// 	if err := server.ListenAndServe(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// // package main

// // import (
// // 	"fmt"
// // 	"io"
// // 	"net/http"
// // 	"sync"
// // 	"time"
// // )

// // // BackendServer represents a backend server.
// // type BackendServer struct {
// // 	URL   string
// // 	Alive bool
// // 	Mutex sync.RWMutex
// // }

// // // BackendServers is a collection of backend servers.
// // type BackendServers []*BackendServer

// // // Next returns the next available backend server.
// // func (servers BackendServers) Next() *BackendServer {
// // 	for {
// // 		for _, server := range servers {
// // 			server.Mutex.RLock()
// // 			alive := server.Alive
// // 			server.Mutex.RUnlock()

// // 			if alive {
// // 				return server
// // 			}
// // 		}
// // 	}
// // }

// // // Update checks the status of each backend server and updates its Alive status.
// // func (servers BackendServers) Update() {
// // 	for _, server := range servers {
// // 		resp, err := http.Head(server.URL)
// // 		if err != nil {
// // 			server.Mutex.Lock()
// // 			server.Alive = false
// // 			server.Mutex.Unlock()
// // 		} else {
// // 			server.Mutex.Lock()
// // 			server.Alive = resp.StatusCode == 200
// // 			server.Mutex.Unlock()
// // 		}
// // 	}
// // }

// // func main() {
// // 	// Create a slice of backend servers.
// // 	servers := BackendServers{
// // 		&BackendServer{URL: "wss://polygon-mainnet.g.alchemy.com/v2/3z5KecNiLnFFeomwszKT2GPiInCZdfiE"},
// // 		&BackendServer{URL: "http://backend2.example.com"},
// // 		&BackendServer{URL: "http://backend3.example.com"},
// // 	}

// // 	// Start a goroutine to periodically check the status of each backend server.
// // 	go func() {
// // 		for {
// // 			servers.Update()
// // 			time.Sleep(time.Second)
// // 		}
// // 	}()

// // 	// Start the HTTP server and balance incoming requests across the backend servers.
// // 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// // 		server := servers.Next()
// // 		if server == nil {
// // 			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
// // 			return
// // 		}

// // 		resp, err := http.Get(server.URL)
// // 		if err != nil {
// // 			http.Error(w, err.Error(), http.StatusInternalServerError)
// // 			return
// // 		}
// // 		defer resp.Body.Close()

// // 		// Forward the response from the backend server to the client.
// // 		if _, err := io.Copy(w, resp.Body); err != nil {
// // 			http.Error(w, err.Error(), http.StatusInternalServerError)
// // 		}
// // 	})

// // 	fmt.Println("Load balancer started on port 8080")
// // 	http.ListenAndServe(":8080", nil)
// // }

package main

import (
	// "fmt"
	"fmt"
	"log"
	"net/http"

	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/websocket"
)

var (
	// WebSocket upgrader
	upgrader = websocket.Upgrader{
		// Allow connections from any origin
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	// Backend servers
	servers = []string{"wss://polygon-mainnet.g.alchemy.com/v2/3z5KecNiLnFFeomwszKT2GPiInCZdfiE", "wss://eth-mainnet.g.alchemy.com/v2/s8zGtV5Jtr7TPRnOmRW6LaS4I4tMjptX"}
	// Current server index
	serverIndex = 0
)

func main() {
	http.HandleFunc("/", handleWebsocket)
	http.ListenAndServe(":8080", nil)
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to a WebSocket
	// conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
	// 	return
	// }
	// defer conn.Close()
	// request := r.FormValue("request")
	// Get the next server to send the request to
	server := servers[serverIndex]
	serverIndex = (serverIndex + 1) % len(servers)

	// Connect to the backend server
	backendConn, err := ethclient.Dial(server)
	if err != nil {
		http.Error(w, "Failed to connect to backend server", http.StatusInternalServerError)
		return
	}
	headers := make(chan *types.Header)
	sub, err := backendConn.SubscribeNewHead(context.Background(), headers)
	fmt.Println(sub)
	if err != nil {
		log.Fatal(err)
	}

	select {
	case err := <-sub.Err():
		log.Fatal(err)
	case header := <-headers:
		block, err := backendConn.BlockByHash(context.Background(), header.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New block:", block.Number().Uint64())
	}

	// defer backendConn.Close()

	// Forward messages between the client and the backend server
	// go func() {
	// 	for {
	// 		_, message, err := conn.ReadMessage()
	// 		if err != nil {
	// 			break
	// 		}
	// 		err = backendConn.WriteMessage(websocket.TextMessage, message)
	// 		if err != nil {
	// 			break
	// 		}
	// 	}
	// }()
	// go func() {
	// 	for {
	// 		_, message, err := backendConn.ReadMessage()
	// 		if err != nil {
	// 			break
	// 		}
	// 		err = conn.WriteMessage(websocket.TextMessage, message)
	// 		if err != nil {
	// 			break
	// 		}
	// 	}
	// }()
}
