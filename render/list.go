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

	htmlBody := bodyElm
	htmlBody.childs = append(htmlBody.childs, elm{tag: "h2", innerText: groupName})
	htmlBody.childs = append(htmlBody.childs, unorderedListElm)

	return buildBaseHtml(listStyleString, htmlBody, listScript)
}

func RenderListHTML(listEntries []models.ListEntry, groupName string) string {
	listHtmlElmTree := buildElmTree(listEntries, groupName)
	htmlString := h(listHtmlElmTree)
	return htmlString
}
