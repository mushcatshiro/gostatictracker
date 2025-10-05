package render

import (
	"database/sql"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/dbop"
	"github.com/mushcatshiro/gostatictracker/models"
)

func buildBookmarkletElmTree(bs []models.Bookmarklet) elm {
	unorderedListElm := elm{
		tag:   "ul",
		attrs: attrsStruct{id: "myList"},
	}
	for _, b := range bs {
		bElm := elm{
			tag: "li",
			childs: []elm{
				{tag: "a", innerText: b.Title, attrs: attrsStruct{href: b.URL}},
				{tag: "span", innerText: " " + b.Description, attrs: attrsStruct{style: "font-style: italic; color: darkgray;"}},
				{tag: "span", innerText: " " + b.InsertTime.Format(common.TimeLayout)},
			},
		}
		unorderedListElm.childs = append(unorderedListElm.childs, bElm)
	}
	htmlBody := listBody
	htmlBody.childs = append(htmlBody.childs, elm{tag: "h2", innerText: "Bookmarklet List"})
	htmlBody.childs = append(htmlBody.childs, unorderedListElm)
	listHtmlElmTree := listHeader
	listHtmlElmTree.childs = append(listHtmlElmTree.childs, htmlBody)
	listHtmlElmTree.childs = append(listHtmlElmTree.childs, listScript)
	return listHtmlElmTree
}

func RenderBookmarkletHTML(bs []models.Bookmarklet) string {
	elm := buildBookmarkletElmTree(bs)
	return h(elm)
}

func RenderBookmarklet(conn *sql.DB) (string, error) {
	var htmlStr string
	bs, err := dbop.GetSpecificGroupEvents(conn, "bookmarklet")
	if err != nil {
		return htmlStr, err
	}
	return RenderBookmarkletHTML(bs), nil
}
