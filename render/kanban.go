package render

import (
	"database/sql"
	"fmt"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/dbop"
)

func RenderKanban(db *sql.DB) (string, error) {
	var htmlString string
	var slElmList []elm

	for idx := common.NOTSTARTED; idx <= common.CANCELLED; idx++ {
		titleName := idx.String()
		listOfEvent, err := dbop.GetKanbanGroup(db, int(idx))
		if err != nil {
			return htmlString, fmt.Errorf("failed to retrieve %s: %v", idx, err)
		}

		slElm := buildOneSwimlane(listOfEvent, titleName)
		slElmList = append(slElmList, slElm)
	}

	bc := boardContainer
	bc.childs = slElmList

	tm := taskModal
	sls := swimlaneScript

	bd := bodyElm
	bd.childs = append(bd.childs, bc)
	bd.childs = append(bd.childs, tm)
	bd.childs = append(bd.childs, sls)

	htmlBase := buildBaseHtml(swimlaneStyleString, bd, elm{})
	htmlString = h(htmlBase)

	return htmlString, nil
}
