package render

const editorHtml = `<div class="row-controls">
	<button id="upload-btn" onclick="toggleElement('upload-container')">Upload Content</button>
	<div class="dropdown">
		<button id="assets-btn" onclick="toggleElement('asset-dropdown')">View Assets</button>
		<div id="asset-dropdown" onclick="event.stopPropagation()">
			<ul id="asset-list"></ul>
		</div>
	</div>
	<div id="upload-container" onclick="event.stopPropagation()">
		<form hx-post="/upload" hx-encoding="multipart/form-data" hx-target="#asset-list" hx-swap="beforeend" hx-on::after-request="toggleElement('upload-container')">
			<label style="font-size: 13px; font-weight: 600;">Select Images:</label>
			<input type="file" name="images" multiple accept="image/*" style="margin: 10px 0; display: block;">
			<label style="font-size: 13px; font-weight: 600;">Paste Links (One per line):</label>
			<textarea name="links" rows="5" class="links-input" placeholder="https://..."></textarea>
			<button type="submit" style="width: 100%; background: #007bff; border-color: #007bff;">Upload Assets</button>
		</form>
	</div>
</div>

<div class="row-editor">
<textarea id="vim-textarea" placeholder="Vim Editor Ready... (Esc for Normal Mode, i for Insert)"></textarea>
</div>`

func RenderEditor() string {
	return editorHtml
}
