package render

var listHeader elm = elm{
	tag: "html",
	childs: []elm{
		{
			tag: "head",
			childs: []elm{
				{
					tag: "style",
					innerText: `
					ul li {
						opacity: 0;
						transform: translateY(20px);
						transition: opacity 0.5s ease-out, transform 0.5s ease-out;
					}
					ul li.visible {
						opacity: 1;
						transform: translateY(0);
					}`,
				},
			},
		},
	},
}

var listBody elm = elm{
	tag: "body",
	childs: []elm{
		{tag: "h2"},
	},
}

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
