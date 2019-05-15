package ccp
  
import (
	"encoding/json"
        "context"
        "github.com/gorilla/mux"
        "k8s.io/klog"
        "log"
        "net/http"
        //"net/http/httputil"
        "os"
        "os/signal"
        "syscall"
        "time"
	"fmt"
	"io/ioutil"
)
func IngressHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf := make(map[string]string)
	err = json.Unmarshal(b, &buf)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Printf("Map String %v", buf)
}
func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf := make(map[string]string)
	err = json.Unmarshal(b, &buf)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Printf("Map String %v", buf)
}

func PodHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf := make(map[string]string)
	err = json.Unmarshal(b, &buf)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Printf("Map String %v", buf)
}

func EndpointHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf := make(map[string]string)
	err = json.Unmarshal(b, &buf)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Printf("Map String %v", buf)

}


func waitForShutdown(srv *http.Server) {
        interruptChan := make(chan os.Signal, 1)
        signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

        // Block until we receive our signal.
        <-interruptChan

        // Create a deadline to wait for.
        ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
        defer cancel()
        srv.Shutdown(ctx)

        log.Println("Shutting down")
        os.Exit(0)
}

func StartRestServer() {
        // Create Server and Route Handlers
        r := mux.NewRouter()
        //r.HandleFunc("/", handler)
        r.HandleFunc("/api/v1/endpoints", EndpointHandler)
        r.HandleFunc("/api/v1/services", ServiceHandler)
        r.HandleFunc("/api/v1/pods", PodHandler)
        r.HandleFunc("/api/v1/ingresses", IngressHandler)
        srv := &http.Server{
                Handler:      r,
                Addr:         "localhost:8080",
                ReadTimeout:  10 * time.Second,
                WriteTimeout: 10 * time.Second,
        }

        // Start Server
        go func() {
                klog.Info("Starting Server")
                if err := srv.ListenAndServe(); err != nil {
                        klog.Fatal(err)
                }
        }()
        // Graceful Shutdown
        waitForShutdown(srv)
}

func main() {
	StartRestServer()
}
