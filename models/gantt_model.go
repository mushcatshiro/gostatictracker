package models

import "time"

type GanttEvent struct {
	ID          int64      `json:"id"`
	Start       *time.Time `json:"start"`
	End         *time.Time `json:"end"`
	Group       string     `json:"group"`
	AllDay      bool       `json:"allDay"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	Description string     `json:"description"`
}

func (ge *GanttEvent) ToEvent() Event {
	it := time.Now()
	return Event{
		ID:          ge.ID,
		Start:       ge.Start,
		End:         ge.End,
		InsertTime:  &it,
		Group:       ge.Group,
		AllDay:      ge.AllDay,
		Title:       ge.Title,
		URL:         ge.URL,
		Description: ge.Description,
	}
}
