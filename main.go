package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got request:")
	defer r.Body.Close()
	b, _ := io.ReadAll(r.Body)
	fmt.Println(string(b))
}

type containerFactory struct {
	baseDir string
}

func (c *containerFactory) lxcHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	LxcDir := filepath.Join(c.baseDir, "lxc.monitor."+uuid)
	cgroupEventsFile := "cgroup.events"

	switch r.Method {
	case http.MethodPost:
		// Fake LXC creation
		err := os.MkdirAll(LxcDir, 0744)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(filepath.Join(LxcDir, cgroupEventsFile), []byte("populated 1\nfrozen 0\n"), 0666)
		if err != nil {
			panic(err)
		}
		fmt.Println("Created container ", LxcDir)
	case http.MethodPut:
		// Fake LXC crash
		err := os.WriteFile(filepath.Join(LxcDir, cgroupEventsFile), []byte("populated 0\nfrozen 0\n"), 0666)
		if err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond) // simulate container stopping
		err = os.RemoveAll(filepath.Join(LxcDir))
		if err != nil {
			panic(err)
		}
	case http.MethodDelete:
		// Fake LXC stop
		http.Error(w, "Not implemented yet (Should simulate deactivating restart)", http.StatusBadRequest)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	var port string
	flag.StringVar(&port, "port", "1543", "Port to listen on")
	flag.Parse()

	// Get current directory fr
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executablePath := filepath.Dir(ex)
	fmt.Println(executablePath)

	containers := containerFactory{executablePath}

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/lxc/{uuid}", containers.lxcHandler)

	err = http.ListenAndServe("localhost:"+port, nil)
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
