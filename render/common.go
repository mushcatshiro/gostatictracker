package render

func buildBaseHtml(styleString string, bodyElm, scriptElm elm) elm {
	// support slice of scriptElm and styleElm instead of current mode
	htmlBase := htmlElm

	baseStyleString := baseHtmlStyleString
	styleElm := elm{tag: "style", innerText: baseStyleString + styleString}
	htmlBase.childs = append(htmlBase.childs, styleElm)

	htmlBase.childs = append(htmlBase.childs, bodyElm)

	if scriptElm.tag != "" {
		htmlBase.childs = append(htmlBase.childs, scriptElm)
	}
	return htmlBase
}

func RenderSimpleView(title, textBlock string) string {
	var f, t, tb, e elm
	f.tag = "div"

	if title != "" {
		t.tag = "h2"
		t.innerText = title
		f.childs = append(f.childs, t)
	}
	if textBlock != "" {
		tb.tag = "p"
		tb.innerText = textBlock
		f.childs = append(f.childs, tb)
	}

	simpleElm := buildBaseHtml("", f, e)
	return h(simpleElm)
}
