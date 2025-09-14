package render

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/mushcatshiro/gostatictracker/dbop"
)

func buildElmTree(listEntries []dbop.ListEntry, groupName string) elm {
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

func renderListHTML(listEntries []dbop.ListEntry, groupName string, file *os.File) error {
	listHtmlElmTree := buildElmTree(listEntries, groupName)
	htmlString := h(listHtmlElmTree)
	_, err := file.Write([]byte(htmlString))
	if err != nil {
		return err
	}
	return nil
}

func RenderList(renderTargetPath string, conn *sql.DB) {
	/*
		right now only support flat list, two main option to support nested list
		handle at application level (create a map through for loop perhaps), using
		database CTE. three alternate database design including nested set model
		with left and right pointers to define node position in tree, path
		enumeration by storing full path of ancestors' id or closure table.
	*/

	groups, err := dbop.GetUniqueGroups(conn)
	if err != nil {
		log.Fatalf("Not able to query any list group(s): %v", err)
	}
	for _, group := range groups {
		log.Printf("Processing group: %s", group)
		listEntries, err := dbop.GetListGroupEntries(conn, group)
		if err != nil {
			log.Printf("Failed to get entries for list group %s with %v", group, err)
			continue
		}
		fileName := strings.ReplaceAll(group, " ", "-") + "-list.html"
		file, err := os.Create(renderTargetPath + "/" + fileName) // truncates if exists
		if err != nil {
			log.Printf("Failed to create file %s: %v", fileName, err)
			continue
		}
		defer file.Close()
		err = renderListHTML(listEntries, group, file)
		if err != nil {
			log.Printf("Failed to process %s:\n\t%v", fileName, err)
		}
	}
}
