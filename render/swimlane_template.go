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

const swimlaneStyleString = `.board-container {
	display: flex;
	gap: var(--spacing-4);
	height: calc(100vh - (var(--spacing-4) * 2));
	overflow-x: auto;
	padding-bottom: var(--spacing-4);
}

.swimlane {
	flex-shrink: 0;
	width: 18rem;
	background-color: var(--color-background-component);
	border-radius: var(--radius-md);
	box-shadow: var(--shadow-md);
	display: flex;
	flex-direction: column;
}

.swimlane-title {
	font-weight: 700;
	font-size: 1.125rem;
	padding: var(--spacing-4);
	color: var(--color-text-secondary);
	border-bottom: var(--border-width) solid var(--color-border);
	margin: 0;
}

.tasks-container {
	padding: var(--spacing-4);
	display: flex;
	flex-direction: column;
	gap: var(--spacing-4);
	overflow-y: auto;
	flex-grow: 1;
}

/* Custom Scrollbar for Kanban */
.tasks-container::-webkit-scrollbar {
	width: 8px;
}
.tasks-container::-webkit-scrollbar-thumb {
	background-color: var(--color-border);
	border-radius: 4px;
}
.tasks-container::-webkit-scrollbar-track {
	background-color: var(--color-background-body);
}

.task-card {
	background-color: var(--color-background-card);
	padding: var(--spacing-3);
	border-radius: var(--radius-sm);
	box-shadow: var(--shadow);
	cursor: pointer;
	transition: var(--transition-default);
}

.task-card:hover {
	background-color: var(--color-background-body);
	transform: translateY(-2px);
	box-shadow: var(--shadow-md);
}

.task-card:active {
	transform: scale(0.98);
}

.task-card p {
	margin: 0;
	font-weight: 600;
	color: var(--color-text-primary);
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.modal {
	position: fixed;
	inset: 0;
	background-color: var(--color-background-overlay);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 50;
	padding: var(--spacing-4);
	opacity: 0;
	transition: opacity 0.3s ease;
	pointer-events: none;
}

.modal.visible {
	opacity: 1;
	pointer-events: auto;
}

.modal-content {
	background-color: var(--color-background-card);
	border-radius: var(--radius-md);
	box-shadow: var(--shadow-xl);
	width: 100%;
	max-width: 42rem;
	padding: var(--spacing-6);
	transform: scale(0.95);
	opacity: 0;
	transition: all 0.3s ease;
}

.modal.visible .modal-content {
	transform: scale(1);
	opacity: 1;
}

.modal-header {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
}

#modalTitle {
	font-size: 1.5rem;
	font-weight: 700;
	color: var(--color-text-primary);
	margin: 0;
}

#closeModal {
	font-size: 2rem;
	font-weight: 200;
	line-height: 1;
	color: var(--color-text-muted);
	background: none;
	border: none;
	cursor: pointer;
	transition: color var(--transition-duration) ease;
	padding: 0;
}

#closeModal:hover {
	color: var(--color-text-primary);
}

#modalDescription {
	margin-top: var(--spacing-4);
	color: var(--color-text-secondary);
	white-space: pre-wrap;
	line-height: 1.6;
}

@media (min-width: 640px) {
	body { padding: var(--spacing-6); }
	.board-container { height: calc(100vh - (var(--spacing-6) * 2)); }
	.swimlane { width: 20rem; }
}

@media (min-width: 1024px) {
	body { padding: 2rem; }
	.board-container { height: calc(100vh - 4rem); }
}`
