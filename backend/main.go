package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rawdaGastan/urls_checker/internal"
)

type Site struct {
	Url string `json:"url"`
}

var sites = []Site{}
var upgrader = websocket.Upgrader{}

func AddSite(w http.ResponseWriter, r *http.Request) {
	var newSite Site
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the site name")
	}

	err = json.Unmarshal(reqBody, &newSite)
	if err != nil {
		fmt.Fprintf(w, "error in given site")
	}

	sites = append(sites, newSite)
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(newSite)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprint(err))
	}
}

func getAllSites(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(sites)

	if err != nil {
		fmt.Fprintf(w, fmt.Sprint(err))
	}
}

func checkSite(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]
	fmt.Println(url)
	if url == "" {
		fmt.Fprintf(w, "no given site")
	}
	service := internal.NewCheckerService(100)

	noMatchedSites := true
	for _, site := range sites {
		if site.Url == url {
			service.AddSite(url)
			noMatchedSites = false
		}
	}

	if !noMatchedSites {
		service.AddApiOutput()
		service.Start()

		err := json.NewEncoder(w).Encode(service.GetApiOutput())
		if err != nil {
			fmt.Fprintf(w, fmt.Sprint(err))
		}
	} else {
		fmt.Fprintf(w, "no matched site, please add the site first")
	}
}

func checkAllSites(w http.ResponseWriter, r *http.Request) {
	service := internal.NewCheckerService(100)

	for _, site := range sites {
		service.AddSite(site.Url)
	}
	service.AddApiOutput()
	service.Start()

	err := json.NewEncoder(w).Encode(service.GetApiOutput())
	if err != nil {
		fmt.Fprintf(w, fmt.Sprint(err))
	}
}

func checkSiteWithSocket(w http.ResponseWriter, r *http.Request) {

	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer conn.Close()

	// Continuosly read and write message
	for {
		_, website, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed:", err)
			break
		}

		service := internal.NewCheckerService(100)
		service.AddSite(string(website))
		service.AddSocket(conn)
		service.Start()
	}
}

func main() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/site", AddSite).Methods("POST")
	router.HandleFunc("/sites", getAllSites).Methods("GET")
	router.HandleFunc("/report/{url}", checkSite).Methods("GET")
	router.HandleFunc("/reports", checkAllSites).Methods("GET")

	router.HandleFunc("/check", checkSiteWithSocket).Methods("GET")

	fmt.Println("server is running at", 4000)
	err := http.ListenAndServe(":4000", router)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
