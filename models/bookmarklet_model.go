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
	it := time.Now()
	return Event{
		InsertTime:  &it,
		Title:       b.Title,
		Description: b.Description,
		Group:       "bookmarklet",
		URL:         b.URL,
	}
}
