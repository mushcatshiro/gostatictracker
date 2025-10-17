package render

var listScript elm = elm{
	tag: "script",
	innerText: `document.addEventListener('DOMContentLoaded', () => {
	const listItems = document.querySelectorAll('#myList li');
	listItems.forEach((item, index) => {
		// Use setTimeout to apply the 'visible' class with a delay
		setTimeout(() => {
			item.classList.add('visible');
		}, index * 100); // 100ms delay between each item
	});
});`,
}

const listStyleString = `ul li {
    opacity: 0;
    transform: translateY(20px);
    transition: opacity 0.5s ease-out, transform 0.5s ease-out;
}
ul li.visible {
    opacity: 1;
    transform: translateY(0);
}`
