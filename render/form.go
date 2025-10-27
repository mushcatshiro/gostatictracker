package render

import (
	"strconv"

	"github.com/mushcatshiro/gostatictracker/common"
	"github.com/mushcatshiro/gostatictracker/models"
)

type formField struct {
	Name     string
	Label    string
	Type     string
	Value    string
	Required bool
	FieldTag string
}

func buildFormSection(m models.Event, viewOnly bool, endpoint string) elm {
	fields := []formField{
		{Name: "group", Label: "Group", Type: "text", Value: m.Group, FieldTag: "input", Required: true},
		{Name: "title", Label: "Title", Type: "text", Value: m.Title, FieldTag: "input", Required: true},
		{Name: "url", Label: "URL", Type: "text", Value: m.URL, FieldTag: "input"},
		{Name: "description", Label: "Description", Type: "", Value: m.Description, FieldTag: "textarea"},
		{Name: "start", Label: "Planned Start Time", Type: "datetime-local", Value: common.ParseTimeToString(m.Start, true), FieldTag: "input", Required: true},
		{Name: "end", Label: "Planned End Time", Type: "datetime-local", Value: common.ParseTimeToString(m.End, true), FieldTag: "input", Required: true},
		{Name: "actualStart", Label: "Actual Start Time", Type: "datetime-local", Value: common.ParseTimeToString(m.ActualStart, true), FieldTag: "input"},
		{Name: "actualEnd", Label: "Actual End Time", Type: "datetime-local", Value: common.ParseTimeToString(m.ActualEnd, true), FieldTag: "input"},
	}

	fieldset := elm{tag: "fieldset"}
	if viewOnly {
		fieldset.attrs = attrsStruct{disabled: true}
	}
	idField := idInputElm
	idField.attrs.value = strconv.Itoa(int(m.ID))
	fieldset.childs = append(fieldset.childs, idField)

	for _, f := range fields {
		formGroup := formGroupElm
		labelE := elm{tag: "label", attrs: attrsStruct{afor: f.Name}, innerText: f.Label}
		fieldE := elm{tag: f.FieldTag, attrs: attrsStruct{
			atype:    f.Type,
			id:       f.Name,
			name:     f.Name,
			required: f.Required,
		}}
		if f.FieldTag == "textarea" {
			fieldE.innerText = f.Value
		} else {
			fieldE.attrs.value = f.Value
		}
		formGroup.childs = []elm{labelE, fieldE}
		fieldset.childs = append(fieldset.childs, formGroup)
	}

	optionFields := []string{"priority", "status"}
	for _, ff := range optionFields {
		f := formGroupElm
		if ff == "priority" {
			pE := priorityElm
			pE[1].childs[int(m.Priority)].attrs.selected = true
			f.childs = pE
		} else {
			sE := statusElm
			sE[1].childs[int(m.Status)].attrs.selected = true
			f.childs = sE
		}
		fieldset.childs = append(fieldset.childs, f)
	}
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

func RenderFormHtml(m models.Event, viewOnly bool, endpoint string) (string, error) {
	formElm := buildFormSection(m, viewOnly, endpoint)
	bd := bodyElm
	bd.childs = append(bd.childs, formElm)

	htmlBase := buildBaseHtml(formStyleString, bd, elm{})
	return h(htmlBase), nil
}
