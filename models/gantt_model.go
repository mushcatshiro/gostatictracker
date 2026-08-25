package models

import "time"

func (r *Record) ToGanttRecord() GanttRecord {
	return GanttRecord{
		ID: r.ID,
		Start: r.Start,
		End: r.End,
		Group: r.Group,
		AllDay: r.AllDay,
		Title: r.Title,
		URL: r.URL,
		Description: r.Description,
	}
}

type GanttRecord struct {
	ID          int64      `json:"id"`
	Start       *time.Time `json:"start"`
	End         *time.Time `json:"end"`
	Group       string     `json:"group"`
	AllDay      bool       `json:"allDay"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	Description string     `json:"description"`
}

func (ge *GanttRecord) ToRecord() Record {
	return Record {
		ID:          ge.ID,
		Start:       ge.Start,
		End:         ge.End,
		Group:       ge.Group,
		AllDay:      ge.AllDay,
		Title:       ge.Title,
		URL:         ge.URL,
		Description: ge.Description,
	}
}

type GanttRenderMetadata struct {
	IsDayView            bool
	GroupStartTime       *time.Time
	GroupEndTime         *time.Time
	GroupName            string
	RowTextInRectPadding int     // LR padding for text to be place in rect
	RectToTextMargin     int     // spacing between rect and text
	TextYOffset          float32 // approximate text Y offset from Y (center)
	RowRectMargin        int     // spacing between two rect/row
	RowRectHeight        int
	HeadersOffset        int
	BaseHeaderRectWidth  int // date or week rect width
	HeaderRectHeight     int
	HeaderRectMargin     int
	HeaderTextYOffset    int // somehow it needs +2 after dividing height by 2
	Divisor              int
}
