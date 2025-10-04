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
		{Name: "insertTime", Label: "Insert Time", Type: "datetime-local", Value: common.ParseTimeToString(m.InsertTime, true), FieldTag: "input"},
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
		formGroup.childs = []elm{
			{tag: "label", attrs: attrsStruct{afor: f.Name}, innerText: f.Label},
			{tag: f.FieldTag, attrs: attrsStruct{
				atype:    f.Type,
				id:       f.Name,
				name:     f.Name,
				value:    f.Value,
				required: f.Required,
			}},
		}
		fieldset.childs = append(fieldset.childs, formGroup)
	}

	optionFields := []string{"priority", "status"}
	for _, ff := range optionFields {
		f := formGroupElm
		if ff == "priority" {
			f.childs = priorityElm
		} else {
			f.childs = statusElm
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

func RenderForm(m models.Event, viewOnly bool, endpoint string) (string, error) {
	formElm := buildFormSection(m, viewOnly, endpoint)
	bd := bodyElm
	bd.childs = append(bd.childs, formElm)

	htmlBase := buildBaseHtml(formStyleString, bd)
	return h(htmlBase), nil
}
