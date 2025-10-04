package render

import (
	"fmt"
	"strings"
)

type attrsStruct struct {
	class    string
	id       string
	style    string
	href     string
	charset  string
	data     string
	action   string
	method   string
	atype    string
	value    string
	name     string
	required bool
	disabled bool
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
		attrParts = append(attrParts, fmt.Sprintf(`href="%s"`, a.href))
	}
	if a.charset != "" {
		attrParts = append(attrParts, fmt.Sprintf(`charset="%s"`, a.charset))
	}
	if a.data != "" {
		attrParts = append(attrParts, fmt.Sprintf(`%s`, a.data))
	}
	if a.action != "" {
		attrParts = append(attrParts, fmt.Sprintf(`action="%s"`, a.action))
	}
	if a.method != "" {
		attrParts = append(attrParts, fmt.Sprintf(`method="%s"`, a.method))
	}
	if a.atype != "" {
		attrParts = append(attrParts, fmt.Sprintf(`type="%s"`, a.atype))
	}
	if a.value != "" {
		attrParts = append(attrParts, fmt.Sprintf(`value="%s"`, a.value))
	}
	if a.name != "" {
		attrParts = append(attrParts, fmt.Sprintf(`name="%s"`, a.name))
	}
	if a.required {
		attrParts = append(attrParts, fmt.Sprintf(`required`))
	}
	if a.disabled {
		attrParts = append(attrParts, fmt.Sprintf(`disabled`))
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
