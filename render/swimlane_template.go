package render

var boardContainer elm = elm{
	tag:   "div",
	attrs: attrsStruct{class: "board-container"},
}

var swimlaneCol elm = elm{
	tag:   "div",
	attrs: attrsStruct{class: "swimlane"},
}

var swimlaneTitle elm = elm{
	tag:   "h2",
	attrs: attrsStruct{class: "swimlane-title"},
}

var taskContainer elm = elm{
	tag: "div",
	attrs: attrsStruct{class: "tasks-container"},
}

var taskCard elm = elm{
	tag: "div",
	attrs: attrsStruct{class: "task-card"},
}

var taskModal elm = elm{
	tag:   "div",
	attrs: attrsStruct{class: "modal hidden", id: "taskModal"},
	childs: []elm{
		{
			tag:   "div",
			attrs: attrsStruct{class: "modal-content", id: "modalContent"},
			childs: []elm{
				{
					tag:   "div",
					attrs: attrsStruct{class: "modal-header"},
					childs: []elm{
						{tag: "h3", attrs: attrsStruct{id: "modalTitle"}},
						{tag: "button", innerText: "&times;", attrs: attrsStruct{id: "closeModal"}},
					},
				},
			},
		},
	},
}

var swimlaneScript elm = elm{
	tag: "script",
	innerText: `document.addEventListener('DOMContentLoaded', () => {
	const taskCards = document.querySelectorAll('.task-card');
	const modal = document.getElementById('taskModal');
	const modalContent = document.getElementById('modalContent');
	const modalTitle = document.getElementById('modalTitle');
	const modalDescription = document.getElementById('modalDescription');
	const closeModalButton = document.getElementById('closeModal');

	const showModal = (title, description) => {
		modalTitle.textContent = title;
		modalDescription.textContent = description;
		modal.classList.remove('hidden');
		setTimeout(() => {
			modal.style.opacity = '1';
			modalContent.style.opacity = '1';
			modalContent.style.transform = 'scale(1)';
		}, 10);
	};

	const hideModal = () => {
		modalContent.style.opacity = '0';
		modalContent.style.transform = 'scale(0.95)';
		modal.style.opacity = '0';
		setTimeout(() => {
			modal.classList.add('hidden');
		}, 300);
	};

	taskCards.forEach(card => {
		card.addEventListener('click', () => {
			const title = card.getAttribute('data-title');
			const description = card.getAttribute('data-description');
			showModal(title, description);
		});
	});

	closeModalButton.addEventListener('click', hideModal);

	modal.addEventListener('click', (event) => {
		if (event.target === modal) {
			hideModal();
		}
	});

	document.addEventListener('keydown', (event) => {
		if (event.key === 'Escape' && !modal.classList.contains('hidden')) {
			hideModal();
		}
	});
});`,
}

const swimlaneStyleString = `:root {
	--slate-100: #f1f5f9;
	--slate-200: #e2e8f0;
	--slate-300: #cbd5e1;
	--slate-500: #64748b;
	--slate-600: #475569;
	--slate-700: #334155;
	--slate-800: #1e293b;
	--slate-900: #0f172a;
	--white: #ffffff;
	--black-alpha-60: rgba(0, 0, 0, 0.6);
	--shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
	--shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
	--shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
	--shadow-xl: 0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);
}

body {
	background-color: var(--slate-100);
	font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
	margin: 0;
	padding: 1rem;
}
.board-container {
	display: flex;
	gap: 1rem;
	height: calc(100vh - 2rem);
	overflow-x: auto;
	padding-bottom: 1rem;
}

.swimlane {
	flex-shrink: 0;
	width: 18rem;
	background-color: var(--slate-200);
	border-radius: 0.5rem;
	box-shadow: var(--shadow-md);
	display: flex;
	flex-direction: column;
}

.swimlane-title {
	font-weight: 700;
	font-size: 1.125rem;
	padding: 1rem;
	color: var(--slate-700);
	border-bottom: 1px solid var(--slate-300);
	margin: 0;
}

.tasks-container {
	padding: 1rem;
	display: flex;
	flex-direction: column;
	gap: 1rem;
	overflow-y: auto;
	flex-grow: 1;
}

.tasks-container::-webkit-scrollbar {
	width: 8px;
}
.tasks-container::-webkit-scrollbar-thumb {
	background-color: var(--slate-300);
	border-radius: 4px;
}
.tasks-container::-webkit-scrollbar-track {
	background-color: var(--slate-100);
}

.task-card {
	background-color: var(--white);
	padding: 0.75rem;
	border-radius: 0.375rem;
	box-shadow: var(--shadow);
	cursor: pointer;
	transition: all 0.2s ease-in-out;
}

.task-card:hover {
	background-color: var(--slate-100);
	transform: translateY(-2px);
	box-shadow: var(--shadow-md);
}

.task-card:active {
	transform: scale(0.98);
}

.task-card p {
	margin: 0;
	font-weight: 600;
	color: var(--slate-800);
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.modal {
	position: fixed;
	inset: 0;
	background-color: var(--black-alpha-60);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 50;
	padding: 1rem;
	opacity: 0;
	transition: opacity 0.3s ease;
}

.modal.hidden {
	display: none;
}

.modal-content {
	background-color: var(--white);
	border-radius: 0.5rem;
	box-shadow: var(--shadow-xl);
	width: 100%;
	max-width: 42rem;
	transform: scale(0.95);
	opacity: 0;
	transition: all 0.3s ease;
	padding: 1.5rem;
}

.modal-header {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
}

#modalTitle {
	font-size: 1.5rem;
	font-weight: 700;
	color: var(--slate-800);
	margin: 0;
}

#closeModal {
	font-size: 2rem;
	font-weight: 200;
	line-height: 1;
	color: var(--slate-500);
	background: none;
	border: none;
	cursor: pointer;
	transition: color 0.2s ease;
	padding: 0;
}

#closeModal:hover {
	color: var(--slate-900);
}

#modalDescription {
	margin-top: 1rem;
	color: var(--slate-600);
	white-space: pre-wrap;
	line-height: 1.6;
}

@media (min-width: 640px) {
	body { padding: 1.5rem; }
	.board-container { height: calc(100vh - 3rem); }
	.swimlane { width: 20rem; }
}

@media (min-width: 1024px) {
	body { padding: 2rem; }
	.board-container { height: calc(100vh - 4rem); }
}`
