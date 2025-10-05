package render

import (
	"github.com/mushcatshiro/gostatictracker/models"
)

func buildElmTree(listEntries []models.ListEntry, groupName string) elm {
	unorderedListElm := elm{
		tag:   "ul",
		attrs: attrsStruct{id: "myList"},
	}
	for _, entry := range listEntries {
		entryElm := elm{
			tag:       "li",
			innerText: entry.Title,
		}
		unorderedListElm.childs = append(unorderedListElm.childs, entryElm)
	}

	htmlBody := listBody
	htmlBody.childs = append(htmlBody.childs, elm{tag: "h2", innerText: groupName})
	htmlBody.childs = append(htmlBody.childs, unorderedListElm)
	listHtmlElmTree := listHeader
	listHtmlElmTree.childs = append(listHtmlElmTree.childs, htmlBody)
	listHtmlElmTree.childs = append(listHtmlElmTree.childs, listScript)

	return listHtmlElmTree
}

func RenderListHTML(listEntries []models.ListEntry, groupName string) (string, error) {
	listHtmlElmTree := buildElmTree(listEntries, groupName)
	htmlString := h(listHtmlElmTree)
	return htmlString, nil
}
