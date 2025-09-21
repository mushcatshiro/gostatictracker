package server

import (
	"database/sql"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"

	"github.com/mushcatshiro/gostatictracker/dbop"
)

var conn *sql.DB

//go:embed ganttui.html
var ui []byte

//go:embed web/*
var content embed.FS

type groups struct {
	GroupName string `json:"groupName"`
	ActionURL string `json:"actionURL"`
}

func (e *groups) MarshalJSON() ([]byte, error) {
	type Alias groups // Create an alias to avoid recursion
	return json.Marshal((*Alias)(e))
}

func handlerHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleAddEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost || r.Method == http.MethodPut {

		var e dbop.Event
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			fmt.Println(err)
			http.Error(w, "bad request: ", http.StatusBadRequest)
			return
		}
		if e.Start == "" || e.Group == "" {
			http.Error(w, "bad request: missing fields", http.StatusBadRequest)
			return
		}

		var err error
		var ret int
		if r.Method == http.MethodPost {
			_, err = dbop.InsertEvent(conn, e)
			ret = http.StatusCreated
		} else {
			err = dbop.UpdateEvent(conn, e)
			ret = http.StatusOK
		}

		if err != nil {
			fmt.Println(err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(ret)
	} else if r.Method == http.MethodGet {
		group := r.URL.Query().Get("group")
		if group != "" {
			events, err := dbop.GetGanttGroupEvents(conn, group, false)
			if err != nil {
				http.Error(w, "group not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(events); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// TODO create gs obj and encode method
			gs, err := dbop.GetUniqueGroups(conn)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			gstruct := []groups{}
			for _, g := range gs {
				t := groups{GroupName: g, ActionURL: "/events?group=" + url.QueryEscape(g)}
				gstruct = append(gstruct, t)
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(gstruct); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleUngrouppedEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ungrouppedEvents, err := dbop.GetUngrouppedEvents(conn)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		es := []dbop.Event{}
		for _, e := range ungrouppedEvents {
			t := dbop.Event{ID: e.ID, Title: e.Title}
			es = append(es, t)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(es); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handlerUI(w http.ResponseWriter, r *http.Request) {
	w.Write(ui)
}

func Run() {
	var username = flag.String("dbUser", "pgsql", "The postgresql username")
	var password = flag.String("dbPass", "pgsql", "The postgresql password")
	var dbHost = flag.String("dbNost", "localhost", "The postgresql host IP")
	var dbName = flag.String("dbName", "pgsql", "The posgresql database name")
	var serverHost = flag.String("serverHost", "localhost", "The app server host IP")
	var serverPort = flag.String("serverPort", "8080", "The app server host IP")
	flag.Parse()
	var err error

	// if use conn, err := ..., it conn will be a local variable instead of
	// using the global variable hence need to use `=`.
	conn, err = dbop.ConnectDB(*username, *password, *dbHost, *dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()
	dbop.InitDB(conn)

	webFS, err := fs.Sub(content, "web")
	if err != nil {
		log.Fatalf("failed to create sub filesystem: %v", err)
	}
	http.HandleFunc("/events", handleAddEvent)
	http.HandleFunc("/UI", handlerUI)
	http.HandleFunc("/healthcheck", handlerHealthCheck)
	http.HandleFunc("/UIv2", func(w http.ResponseWriter, r *http.Request) {
		// Corrected way to read file from fs.FS
		file, err := webFS.Open("index.html")
		if err != nil {
			http.Error(w, "Could not open index.html", http.StatusInternalServerError)
			return
		}
		defer file.Close() // Ensure the file is closed

		htmlBytes, err := io.ReadAll(file) // Read all content from the file
		if err != nil {
			http.Error(w, "Could not read index.html content", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(htmlBytes)
	})
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Corrected way to read file from fs.FS
		file, err := webFS.Open("test.html")
		if err != nil {
			http.Error(w, "Could not open test.html", http.StatusInternalServerError)
			return
		}
		defer file.Close() // Ensure the file is closed

		htmlBytes, err := io.ReadAll(file) // Read all content from the file
		if err != nil {
			http.Error(w, "Could not read index.html content", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(htmlBytes)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(webFS))))
	fmt.Printf("Server running on http://%s:%s\n", *serverHost, *serverPort)
	addr := *serverHost + ":" + *serverPort
	log.Fatal(http.ListenAndServe(addr, nil))
}
