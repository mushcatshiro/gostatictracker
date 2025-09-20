package render

var calendarHeader elm = elm{
	tag: "html",
	childs: []elm{{
		tag: "style",
		innerText: `
			.calendar-container {
				display: grid;
				grid-template-columns: repeat(7, 1fr);
				position: relative;
			}
			.day-cell {
				border: 1px solid #ccc;
				padding: 5px;
				margin: 2px;
				min-height: 100px;
				box-shadow: 0px 0px 3px #CBD5C2;
			}
			.event-div {
				position: absolute;
				grid-column: var(--start-col) / span var(--span);
				grid-row-start: var(--start-row);
				background-color: #4A90E2;
				color: white;
				margin-top: 25px;
				padding: 2px 5px;
				overflow: hidden;
				z-index: 2;
			}
			.weekdays {
				padding: 0.5rem;
			}
		`,
	}},
}

var calendarTitle elm = elm{
	tag: "h2",
	attrs: attrsStruct{
		class: "title",
	},
}

var calendarBody elm = elm{
	tag: "body",
}

var calendarWeekdays []elm = []elm{
	{tag: "div", innerText: "Sunday", attrs: attrsStruct{class: "weekdays"}},
	{tag: "div", innerText: "Monday", attrs: attrsStruct{class: "weekdays"}},
	{tag: "div", innerText: "Tuesday", attrs: attrsStruct{class: "weekdays"}},
	{tag: "div", innerText: "Wednesday", attrs: attrsStruct{class: "weekdays"}},
	{tag: "div", innerText: "Thursday", attrs: attrsStruct{class: "weekdays"}},
	{tag: "div", innerText: "Friday", attrs: attrsStruct{class: "weekdays"}},
	{tag: "div", innerText: "Saturday", attrs: attrsStruct{class: "weekdays"}},
}

var calendarContainer elm = elm{
	tag: "div",
	attrs: attrsStruct{
		class: "calendar-container",
	},
}
