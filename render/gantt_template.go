package render

import "strconv"

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

func buildHeader(rectX, rectWidth, textX, textY int, textVal string) elm {
	rectE := rectElm
	rectE.attrs = attrsStruct{x: strconv.Itoa(rectX), width: strconv.Itoa(rectWidth)}
	textE := textElm
	textE.attrs = attrsStruct{x: strconv.Itoa(textX), y: strconv.Itoa(textY)}
	textE.innerText = textVal
	groupElm := gElm
	groupElm.childs = append(groupElm.childs, rectE, textE)
	return groupElm
}

func buildRow(rectX, rectY, rectWidth, lineX1, lineX2, lineY1, lineY2, textX int, textY float32, textVal string) elm {
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
	groupElm.childs = append(groupElm.childs, rectE, lineE, textE)
	return groupElm
}
