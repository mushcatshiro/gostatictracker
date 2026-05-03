package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
	"github.com/mushcatshiro/gostatictracker/render"
)

func handlePostForm(r *http.Request) (models.Event, error) {
	var e models.Event
	priority, err := common.ParsePriority(r.PostFormValue("priority"))
	if err != nil {
		return e, errors.New("Bad Request")
	}
	status, err := common.ParseStatus(r.PostFormValue("status"))
	if err != nil {
		return e, errors.New("Bad Request")
	}
	e = models.Event{
		Start:       common.ParseStringDate(r.PostFormValue("start"), true, false),
		End:         common.ParseStringDate(r.PostFormValue("end"), true, false),
		ActualStart: common.ParseStringDate(r.PostFormValue("actualStart"), true, false),
		ActualEnd:   common.ParseStringDate(r.PostFormValue("actualEnd"), true, false),
		Group:       r.PostFormValue("group"),
		Title:       r.PostFormValue("title"),
		URL:         r.PostFormValue("url"),
		Description: r.PostFormValue("description"),
		Priority:    priority,
		Status:      status,
	}
	return e, nil
}

func (s *Server) handleEventFormView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := r.Method
		nextURL := r.URL.Query().Get("next")
		if nextURL == "" {
			nextURL = "/"
		}
		switch m {
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Bad Request: Failed to parse form", http.StatusBadRequest)
				return
			}
			e, err := handlePostForm(r)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			}
			_, err = dbop.InsertEvent(s.db, e)
			if err != nil {
				fmt.Printf("%v\n", err)
				http.Error(w, "Bad Request: Failed to insert to database", http.StatusBadRequest)
				return
			}
			http.Redirect(w, r, nextURL, http.StatusSeeOther)
		case http.MethodGet:
			id := r.URL.Query().Get("id")

			var e models.Event
			if id != "" {
				iid, err := strconv.Atoi(id)
				if err != nil {
					http.Error(w, "Bad Request: Unexpected id", http.StatusBadRequest)
					return
				}
				e, err = dbop.ReadEventById(s.db, int64(iid))
				if err != nil {
					http.Error(w, fmt.Sprintf("Bad Request: event %s not found", id), http.StatusBadRequest)
					return
				}
			}
			htmlString, err := render.RenderFormHtml(e, false, "/eventForm")
			if err != nil {
				http.Error(w, "Bad Request: failed to render view", http.StatusBadRequest)
			}
			io.WriteString(w, htmlString)
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
