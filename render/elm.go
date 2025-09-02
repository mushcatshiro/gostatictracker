package render

import (
	"fmt"
	"strings"
)

type attrsStruct struct {
	class string
	id    string
	style string
}

type elm struct {
	tag       string
	attrs     attrsStruct
	innerText string
	childs    []elm
}

func (a *attrsStruct) toString() string {
	var attrParts []string
	if a.class != "" {
		attrParts = append(attrParts, fmt.Sprintf(`class=%s`, a.class))
	}
	if a.id != "" {
		attrParts = append(attrParts, fmt.Sprintf(`id=%s`, a.id))
	}
	if a.style != "" {
		attrParts = append(attrParts, fmt.Sprintf(`style=%s`, a.style))
	}

	return strings.Join(attrParts, " ")
}

func h(e elm) string {
	if len(e.childs) == 0 {
		return fmt.Sprintf(`<%s %s>%s</%s>`, e.tag, e.attrs.toString(), e.innerText, e.tag)
	}
	var childString string
	for _, cElm := range e.childs {
		childString += h(cElm)
	}
	return fmt.Sprintf(
		`<%s %s>%s %s</%s>`, e.tag, e.attrs.toString(), e.innerText, childString, e.tag,
	)
}
