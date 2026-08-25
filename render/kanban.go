package render

import (
	"fmt"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
)

func RenderKanban(molor map[string][]models.Record, groupName string) (string, error) {
	var htmlString string
	var slElmList []elm

	for idx := common.NOTSTARTED; idx <= common.CANCELLED; idx++ {
		titleName := idx.String()
		listOfEvent, ok := molor[titleName]
		if !ok {
			return "", fmt.Errorf("missing title %s", titleName)
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
