package render

import (
	"strconv"

	"github.com/mushcatshiro/gostatictracker/common"
)

var idInputElm elm = elm{
	tag: "input", attrs: attrsStruct{atype: "hidden", id: "id", name: "id"},
}

var formGroupElm elm = elm{
	tag:   "div",
	attrs: attrsStruct{class: "form-group"},
}

var submitButtonElm elm = elm{
	tag: "button",
	attrs: attrsStruct{atype: "submit"},
	innerText: "Submit",
}

var priorityOptions, statusOptions, priorityElm, statusElm []elm

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
	priorityElm = []elm{
		{tag: "label", attrs: attrsStruct{afor: "priority"}, innerText: "Priority"},
		{tag: "select", attrs: attrsStruct{id: "priority", name: "priority"}, childs: priorityOptions},
	}
	statusElm = []elm{
		{tag: "label", attrs: attrsStruct{afor: "status"}, innerText: "Status"},
		{tag: "select", attrs: attrsStruct{id: "status", name: "status"}, childs: priorityOptions},
	}
}

const formStyleString = `body {
	font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	min-height: 100vh;
	margin: 0;
	background-color: #f0f2f5;
	gap: 20px;
	padding: 20px;
	box-sizing: border-box;
}

form {
	background: #ffffff;
	padding: 2rem;
	border-radius: 8px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	width: 100%;
	max-width: 500px; /* Max width for larger screens */
	box-sizing: border-box;
}

form h2 {
	margin-top: 0;
	text-align: center;
	color: #333;
}

.form-group {
	margin-bottom: 1rem;
}

label {
	display: block;
	margin-bottom: 0.5rem;
	font-weight: 600;
	color: #555;
}

input[type="text"],
input[type="datetime-local"],
textarea,
select {
	width: 100%;
	padding: 0.75rem;
	border: 1px solid #ccc;
	border-radius: 4px;
	box-sizing: border-box; /* Ensures padding doesn't affect width */
	font-size: 1rem;
}

textarea {
	height: 120px;
	resize: vertical;
}

button[type="submit"] {
	width: 100%;
	padding: 0.75rem;
	border: none;
	border-radius: 4px;
	background-color: #007bff;
	color: white;
	font-size: 1rem;
	font-weight: bold;
	cursor: pointer;
	transition: background-color 0.2s;
}

button[type="submit"]:hover {
	background-color: #0056b3;
}

fieldset:disabled {
	opacity: 0.6;
	pointer-events: none;
}

fieldset {
	border: none;
	padding: 0;
	margin: 0;
}`
