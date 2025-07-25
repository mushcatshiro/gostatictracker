package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/mushcatshiro/gostatictracker/dbop"
)

var conn *sql.DB

//go:embed ganttui.html
var ui []byte

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
				http.Error(w, "no groups found", http.StatusNotFound)
			}
			var gstruct []groups
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

func handlerUI(w http.ResponseWriter, r *http.Request) {
	w.Write(ui)
}

func main() {
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

	http.HandleFunc("/events", handleAddEvent)
	http.HandleFunc("/UI", handlerUI)
	http.HandleFunc("/healthcheck", handlerHealthCheck)
	fmt.Printf("Server running on http://%s:%s\n", *serverHost, *serverPort)
	addr := *serverHost + ":" + *serverPort
	log.Fatal(http.ListenAndServe(addr, nil))
}
