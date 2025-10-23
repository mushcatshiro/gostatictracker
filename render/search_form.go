package render

func buildIndexSearchFormSection(groups []string, endpoint string) elm {
	var groupOptions, groupOptionElm []elm
    r := formGroupElm
    r.childs = renderOptionElm

    for _, g := range groups {
        groupOptions = append(
			groupOptions,
			elm{tag: "option", attrs: attrsStruct{value: g}, innerText: g},
		)
    }
    groupOptionElm = []elm{
		{tag: "label", attrs: attrsStruct{afor: "group"}, innerText: "Group"},
		{tag: "select", attrs: attrsStruct{id: "group", name: "group"}, childs: groupOptions},
	}
    g := formGroupElm
    g.childs = groupOptionElm

    fieldset := elm{tag: "fieldset"}
    fieldset.childs = append(fieldset.childs, g, r)
    sb := submitButtonElm
    form := elm{
		tag:   "form",
		attrs: attrsStruct{action: endpoint, method: "post"},
		childs: []elm{
			{tag: "h2", innerText: "Form"},
			fieldset,
			sb,
		},
	}
	return form
}

func RenderIndexSearchFormHtml(groups []string, endpoint string) (string, error) {
    formElm := buildIndexSearchFormSection(groups, endpoint)
	bd := bodyElm
	bd.childs = append(bd.childs, formElm)

	htmlBase := buildBaseHtml(formStyleString, bd, elm{})
	return h(htmlBase), nil
}
