package render

var htmlElm elm = elm{
	tag: "html",
	childs: []elm{
		{
			tag: "head",
			childs: []elm{
				{tag: "meta", attrs: attrsStruct{charset: "UTF-8"}},
			},
		},
	},
}

var styleElm elm = elm{
	tag: "style",
}

