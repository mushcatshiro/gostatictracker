package render

var htmlElm elm = elm{
	tag: "html",
	childs: []elm{
		{
			tag: "head",
			childs: []elm{
				{tag: "meta", attrs: attrsStruct{charset: "UTF-8"}},
			},
		},
	},
}

var bodyElm elm = elm{
	tag: "body",
}

const baseHtmlStyleString = `:root {
	--color-primary: #007bff;
	--color-primary-hover: #0056b3;

	--color-text-primary: #1e293b; /* slate-800 */
	--color-text-secondary: #475569; /* slate-600 */
	--color-text-muted: #64748b;     /* slate-500 */
	--color-text-inverted: #ffffff;

	--color-background-body: #f1f5f9; /* slate-100 */
	--color-background-component: #e2e8f0; /* slate-200 */
	--color-background-card: #ffffff;
	--color-background-overlay: rgba(0, 0, 0, 0.6);

	--color-border: #cbd5e1; /* slate-300 */

	--font-family-sans: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
	--font-family-display: Verdana, Geneva, Tahoma, sans-serif; /* For Gantt Chart */

	--spacing-1: 0.25rem; /* 4px */
	--spacing-2: 0.5rem;
	--spacing-3: 0.75rem;
	--spacing-4: 1rem;
	--spacing-5: 1.25rem;
	--spacing-6: 1.5rem;

	--radius-sm: 0.25rem;
	--radius-md: 0.5rem;
	--border-width: 1px;

	--shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
	--shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
	--shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
	--shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1);
	--shadow-xl: 0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);

	--transition-duration: 0.2s;
	--transition-timing-function: ease-in-out;
	--transition-default: all var(--transition-duration) var(--transition-timing-function);
}

body {
	background-color: var(--color-background-body);
	font-family: var(--font-family-sans);
	color: var(--color-text-primary);
	margin: 0;
	padding: var(--spacing-4);
	box-sizing: border-box;
}

*, *::before, *::after {
	box-sizing: inherit;
}
`

