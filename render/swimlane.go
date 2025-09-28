package render

import (
	"fmt"
	"strings"

	"github.com/mushcatshiro/gostatictracker/models"
)

func buildOneSwimlane(listOfEvents []models.Event, title string) elm {
	var listOfTask []elm
	for _, e := range listOfEvents {
		attr := e.ToDataMap()
		var d []string
		for k, v := range attr {
			d = append(d, fmt.Sprintf(`data-%s="%s"`, k, v))
		}
		t := taskCard
		t.attrs.data = strings.Join(d, "\n")
		t.childs = append(t.childs, elm{tag: "p", innerText: attr["title"]})
		listOfTask = append(listOfTask, t)
	}
	tc := taskContainer
	tc.childs = listOfTask

	st := swimlaneTitle
	st.innerText = title

	sl := swimlaneCol
	sl.childs = append(sl.childs, st)
	sl.childs = append(sl.childs, tc)
	return sl
}
