package models

import "time"

type Bookmarklet struct {
	InsertTime  *time.Time
	Title       string
	Description string
	Group       string
	URL         string
}

func (b *Bookmarklet) ToEvent() Event {
	return Event{
		Title:       b.Title,
		Description: b.Description,
		Group:       "bookmarklet",
		URL:         b.URL,
	}
}
