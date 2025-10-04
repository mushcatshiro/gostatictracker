package render

func buildBaseHtml(styleString string, bodyElm elm) elm {
	htmlBase := htmlElm

	styleElm = elm{tag: "style", innerText: styleString}
	htmlBase.childs = append(htmlBase.childs, styleElm)

	htmlBase.childs = append(htmlBase.childs, bodyElm)
	return htmlBase
}
