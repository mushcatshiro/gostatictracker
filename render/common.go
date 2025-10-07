package render

func buildBaseHtml(styleString string, bodyElm, scriptElm elm) elm {
    // support slice of scriptElm and styleElm instead of current mode
	htmlBase := htmlElm

    styleElm := elm{tag: "style", innerText: styleString}
	htmlBase.childs = append(htmlBase.childs, styleElm)

	htmlBase.childs = append(htmlBase.childs, bodyElm)

    if scriptElm.tag != "" {
        htmlBase.childs = append(htmlBase.childs, scriptElm)
    }
	return htmlBase
}
