package render

var calendarHeader elm = elm{
	tag: "html",
	childs: []elm{{
			tag: "style",
			innerText: `
			.calendar-container {
				display: grid;
				grid-template-columns: repeat(7, 1fr);
			}
			.day-cell {
				border: 1px solid #ccc,
				padding: 5px;
				min-height: 100px;
			}
			.event-div {
				grid-column-start: var(--start-col);
				grid-row-start: var(--start-row);
				grid-column-end: span var(--span);
				background-color: #4A90E2;
				color: white;
				margin-top: 25px;
				padding: 2px, 5px;
				overflow: hidden;
				z-index: 1;
			}
		`,
	},},
}

var calendarBody elm = elm{
	tag: "body",
}

var calendarContainer elm = elm{
	tag:   "div",
	attrs: attrsStruct{
		class: "calendar-container",
	},
}
