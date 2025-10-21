package render

import (
	"strconv"

	"github.com/mushcatshiro/gostatictracker/common"
)

var supportedRenderMode = []string{
	"bookmarklet",
	"calendar",
	"form",
	"gantt",
	"kanban",
	"list",
	"searchForm",
}

var idInputElm elm = elm{
	tag: "input", attrs: attrsStruct{atype: "hidden", id: "id", name: "id"},
}

var formGroupElm elm = elm{
	tag:   "div",
	attrs: attrsStruct{class: "form-group"},
}

var submitButtonElm elm = elm{
	tag:       "button",
	attrs:     attrsStruct{atype: "submit"},
	innerText: "Submit",
}

var priorityOptions, statusOptions, priorityElm, statusElm, renderOptions, renderOptionElm []elm

func init() {
	for k, v := range common.PriorityMap {
		priorityOptions = append(
			priorityOptions,
			elm{tag: "option", attrs: attrsStruct{value: strconv.Itoa(int(v))}, innerText: k},
		)
	}
	for k, v := range common.StatusMap {
		statusOptions = append(
			statusOptions,
			elm{tag: "option", attrs: attrsStruct{value: strconv.Itoa(int(v))}, innerText: k},
		)
	}
	for _, rm := range supportedRenderMode {
		renderOptions = append(
			renderOptions,
			elm{tag: "option", attrs: attrsStruct{value: rm}, innerText: rm},
		)
	}
	priorityElm = []elm{
		{tag: "label", attrs: attrsStruct{afor: "priority"}, innerText: "Priority"},
		{tag: "select", attrs: attrsStruct{id: "priority", name: "priority"}, childs: priorityOptions},
	}
	statusElm = []elm{
		{tag: "label", attrs: attrsStruct{afor: "status"}, innerText: "Status"},
		{tag: "select", attrs: attrsStruct{id: "status", name: "status"}, childs: statusOptions},
	}
	renderOptionElm = []elm{
		{tag: "label", attrs: attrsStruct{afor: "render"}, innerText: "Render Mode"},
		{tag: "select", attrs: attrsStruct{id: "render", name: "render"}, childs: renderOptions},
	}
}

const formStyleString = `form {
	background: var(--color-background-card);
	padding: var(--spacing-6);
	border-radius: var(--radius-md);
	box-shadow: var(--shadow-md);
	width: 100%;
	max-width: 500px;
	margin: 2rem auto; /* Center form if it's a standalone page */
}

form h2 {
	margin-top: 0;
	text-align: center;
	color: var(--color-text-primary);
}

.form-group {
	margin-bottom: var(--spacing-4);
}

label {
	display: block;
	margin-bottom: var(--spacing-2);
	font-weight: 600;
	color: var(--color-text-secondary);
}

input[type="text"],
input[type="datetime-local"],
textarea,
select {
	width: 100%;
	padding: var(--spacing-3);
	border: var(--border-width) solid var(--color-border);
	border-radius: var(--radius-sm);
	font-size: 1rem;
	font-family: var(--font-family-sans);
	background-color: var(--color-background-card);
	color: var(--color-text-primary);
	transition: var(--transition-default);
}

input[type="text"]:focus,
input[type="datetime-local"]:focus,
textarea:focus,
select:focus {
	outline: none;
	border-color: var(--color-primary);
	box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.25);
}

textarea {
	min-height: 120px;
	resize: vertical;
}

button[type="submit"] {
	width: 100%;
	padding: var(--spacing-3);
	border: none;
	border-radius: var(--radius-sm);
	background-color: var(--color-primary);
	color: var(--color-text-inverted);
	font-size: 1rem;
	font-weight: bold;
	cursor: pointer;
	transition: var(--transition-default);
}

button[type="submit"]:hover {
	background-color: var(--color-primary-hover);
}

fieldset {
	border: none;
	padding: 0;
	margin: 0;
}

fieldset:disabled {
	opacity: 0.6;
	pointer-events: none;
}`
