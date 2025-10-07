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
	htmlBody := bodyElm
	htmlBody.childs = append(htmlBody.childs, elm{tag: "h2", innerText: "Bookmarklet List"})
	htmlBody.childs = append(htmlBody.childs, unorderedListElm)

    return buildBaseHtml(listStyleString, htmlBody, listScript)
}

func renderBookmarkletHtml(bs []models.Bookmarklet) string {
	elm := buildBookmarkletElmTree(bs)
	return h(elm)
}

func RenderBookmarklet(conn *sql.DB) (string, error) {
	var htmlStr string
	bs, err := dbop.GetSpecificGroupEvents(conn, "bookmarklet")
	if err != nil {
		return htmlStr, err
	}
	return renderBookmarkletHtml(bs), nil
}
