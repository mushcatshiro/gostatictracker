package render

import (
	"fmt"
	"strconv"
)

var gElm elm = elm{
	tag: "g",
}

var rectElm elm = elm{
	tag:        "rect",
	selfClosed: true,
}

var lineElm elm = elm{
	tag:        "line",
	selfClosed: true,
}

var textElm elm = elm{
	tag: "text",
}

var svgElm elm = elm{
	tag:   "svg",
	attrs: attrsStruct{
        class: "gantt",
        id: "gantt-svg",
        style: "padding: 0.5rem;outline-width: 0.5px; outline-style: solid; outline-color: rgb(200, 200, 200); border-radius: 0.5rem;",
    },
}

var todayIndicatorScript elm = elm{
	tag: "script",
	innerText: `document.addEventListener("DOMContentLoaded", function () {
	const today = new Date();
	const todayStr = today.toISOString().split("T")[0];
	const dateGroups = document.getElementById("date-layer").querySelectorAll("g");

	let matched = null;

	dateGroups.forEach(g => {
		const dateAttr = g.getAttribute("data-date");
		if (dateAttr === todayStr) { matched = g;}
	});

	if (matched) {
		const svg = document.getElementById("gantt-svg");
		const highlight = document.createElementNS("http://www.w3.org/2000/svg", "rect");
		const rectStyle = getComputedStyle(matched.querySelector("rect"));
		highlight.setAttribute("x", parseInt(rectStyle.x)+2);
		highlight.setAttribute("y", rectStyle.y);
		highlight.setAttribute("width", parseInt(rectStyle.width)-4);
		highlight.setAttribute("height", 110);  // how to calculate
		highlight.setAttribute("class", "highlight-today");
		highlight.setAttribute("fill", "red")
		highlight.setAttribute("opacity", 0.1)
		highlight.setAttribute("rx", 4);
		svg.appendChild(highlight);
	}});`,
}

const ganttStyleStringP1 = `svg {
	font-family: var(--font-family-display);
}
svg text{
	alignment-baseline: middle;
	dominant-baseline: middle;
	text-anchor: middle;
	font-size: 18px;
	fill: var(--color-text-secondary);
}
svg .header rect{
	width: %dpx;
	height: %dpx;
	fill: var(--color-background-card);
}
svg .year rect{
	y: 0;
}
svg .month rect{
	y: 26;
}
.gantt .rows rect {
	cursor: pointer;
	transition: fill var(--transition-duration) var(--transition-timing-function);
	fill: var(--color-border);
	height: %dpx;
	rx: var(--radius-sm);
	opacity: 0.8;
}`

const ganttStyleStringP2 = `.gantt .rows rect:hover {
	fill: var(--color-text-secondary);
}
.tooltip {
	position: absolute;
	background-color: var(--color-background-overlay);
	color: var(--color-text-inverted);
	padding: var(--spacing-2) var(--spacing-3);
	border-radius: var(--radius-sm);
	font-size: 14px;
	font-family: var(--font-family-sans);
	width: 250px;
	max-height: 100px; /* fixed max height, content will scroll if longer */
	overflow-y: auto; /* Enable vertical scrolling */
	white-space: normal; /* Allow text to wrap within the fixed width */
	pointer-events: none; /* Allows interaction with elements behind it */
	opacity: 0;
	transition: opacity var(--transition-duration) ease, transform var(--transition-duration) ease;
	transform: translate(-50%, -10px); /* Initial offset for smooth appearance */
	z-index: 1000; /* Ensure tooltip is on top */
	box-shadow: var(--shadow-lg);
}
.tooltip.active {
	opacity: 1;
	transform: translate(-50%, 0); /* Move to final position */
}
.date rect {
	y: 52;
}
.rows line {
	stroke: #DDD;
}
.rows text {
	font-size: 12px;
	text-anchor: start;
	fill: var(--color-text-primary);
}`

func buildHeaderGroup(rectX, rectWidth, textX, textY int, textVal, classVal string) elm {
	rectE := rectElm
	rectE.attrs = attrsStruct{x: strconv.Itoa(rectX), width: strconv.Itoa(rectWidth)}
	textE := textElm
	textE.attrs = attrsStruct{x: strconv.Itoa(textX), y: strconv.Itoa(textY)}
	textE.innerText = textVal
	groupElm := gElm
	groupElm.attrs = attrsStruct{class: classVal}
	groupElm.childs = append(groupElm.childs, rectE, textE)
	return groupElm
}

func buildRowGroup(rectX, rectY, rectWidth, lineX1, lineX2, lineY1, lineY2, textX int, textY float32, textVal, classVal string) elm {
	rectE := rectElm
	rectE.attrs = attrsStruct{x: strconv.Itoa(rectX), y: strconv.Itoa(rectY), width: strconv.Itoa(rectWidth)}
	lineE := lineElm
	lineE.attrs = attrsStruct{
		x1: strconv.Itoa(lineX1), x2: strconv.Itoa(lineX2), y1: strconv.Itoa(lineY1), y2: strconv.Itoa(lineY2),
	}
	textE := textElm
	textE.attrs = attrsStruct{x: strconv.Itoa(textX), y: strconv.FormatFloat(float64(textY), 'f', 2, 32)}
	textE.innerText = textVal
	groupElm := gElm
	groupElm.attrs = attrsStruct{class: classVal}
	groupElm.childs = append(groupElm.childs, rectE, lineE, textE)
	return groupElm
}

func buildGanttStyleString(headerRectWidth, headerRectHeight, rowRectHeight int) string {
	ss1t := ganttStyleStringP1
	ss1 := fmt.Sprintf(ss1t, headerRectWidth, headerRectHeight, rowRectHeight)
	ss2 := ganttStyleStringP2
	return ss1 + "\n" + ss2
}
