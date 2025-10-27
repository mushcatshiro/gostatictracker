package render

import (
	"fmt"
	"strings"
)

type attrsStruct struct {
	class       string
	id          string
	style       string
	href        string
	charset     string
	data        string
	action      string
	method      string
	atype       string
	value       string
	name        string
	afor        string
	placeholder string
	required    bool
	disabled    bool
	selected    bool
	x           string
	y           string
	width       string
	height      string
	x1          string
	x2          string
	y1          string
	y2          string
}

type elm struct {
	tag        string
	attrs      attrsStruct
	innerText  string
	childs     []elm
	selfClosed bool
}

func (a *attrsStruct) toString() string {
	var attrParts []string
	attrs := map[string]string{
		"class":   a.class,
		"id":      a.id,
		"style":   a.style,
		"href":    a.href,
		"charset": a.charset,
		"action":  a.action,
		"method":  a.method,
		"type":    a.atype,
		"value":   a.value,
		"name":    a.name,
		"for":     a.afor,
		"x":       a.x,
		"y":       a.y,
		"width":   a.width,
		"x1":      a.x1,
		"x2":      a.x2,
		"y1":      a.y1,
		"y2":      a.y2,
	}
	for key, val := range attrs {
		if val != "" {
			attrParts = append(attrParts, fmt.Sprintf(`%s="%s"`, key, val))
		}
	}
	if a.data != "" {
		attrParts = append(attrParts, a.data)
	}
	if a.required {
		attrParts = append(attrParts, "required")
	}
	if a.disabled {
		attrParts = append(attrParts, "disabled")
	}
	if a.selected {
		attrParts = append(attrParts, "selected")
	}
	return strings.Join(attrParts, " ")
}

func h(e elm) string {
	if len(e.childs) == 0 {
		if e.selfClosed {
			return fmt.Sprintf(`<%s %s/>`, e.tag, e.attrs.toString())
		} else {
			return fmt.Sprintf(`<%s %s>%s</%s>`, e.tag, e.attrs.toString(), e.innerText, e.tag)
		}
	}
	var childString string
	for _, cElm := range e.childs {
		childString += h(cElm)
	}
	if e.selfClosed {
		return fmt.Sprintf(`<%s %s/>`, e.tag, e.attrs.toString())
	} else {
		return fmt.Sprintf(
			`<%s %s>%s %s</%s>`, e.tag, e.attrs.toString(), e.innerText, childString, e.tag,
		)
	}
}
