package render

import (
	"fmt"
	"strings"
)

type attrsStruct struct {
	class   string
	id      string
	style   string
	href    string
	charset string
	data    string
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
	if a.class != "" {
		attrParts = append(attrParts, fmt.Sprintf(`class="%s"`, a.class))
	}
	if a.id != "" {
		attrParts = append(attrParts, fmt.Sprintf(`id="%s"`, a.id))
	}
	if a.style != "" {
		attrParts = append(attrParts, fmt.Sprintf(`style="%s"`, a.style))
	}
	if a.href != "" {
		attrParts = append(attrParts, fmt.Sprintf(`href="%s"`, a.style))
	}
	if a.charset != "" {
		attrParts = append(attrParts, fmt.Sprintf(`charset="%s"`, a.charset))
	}
	if a.data != "" {
		attrParts = append(attrParts, fmt.Sprintf(`%s`, a.data))
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
