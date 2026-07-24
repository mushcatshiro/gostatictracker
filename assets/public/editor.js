document.addEventListener('DOMContentLoaded', () => {
  const editor = document.getElementById('vim-textarea');
  const modeStatus = document.getElementById('vim-mode-indicator'); // Optional UI element

  // 1. The State Machine
  let mode = 'insert';
  let keyBuffer = '';
  let bufferTimer = null;

  function setMode(newMode) {
    mode = newMode;
    if (mode === "normal") {
      if (modeStatus) modeStatus.innerText = 'NORMAL';
      editor.classList.add('normal-mode');
    } else {
      if (modeStatus) modeStatus.innerText = 'INSERT';
      editor.classList.remove('normal-mode');
    }
  }

  // Helper: Move the cursor and force the textarea to scroll to it
  function moveCursor(pos) {
    if (mode === 'normal') {
      // Fake a block cursor by highlighting exactly one character
      // We use Math.min to ensure we don't try to select past the end of the text
      const endPos = Math.min(pos + 1, editor.value.length);
      editor.setSelectionRange(pos, endPos);
    } else {
      // Standard collapsed blinking line for Insert mode
      editor.setSelectionRange(pos, pos);
    }
    // editor.setSelectionRange(pos, pos);

    // The Vanilla JS Scroll Hack:
    // Textareas don't have a native 'scrollToCursor' method.
    // Briefly blurring and refocusing the element forces the browser
    // to snap the viewport to the current cursor position.
    editor.blur();
    editor.focus();
  }
  // Helper: Catch multi-key commands like 'gg' or 'gq'
  function pushToBuffer(char) {
    keyBuffer += char;
    clearTimeout(bufferTimer);
    // Clear the buffer if the user types too slowly (1 second window)
    bufferTimer = setTimeout(() => { keyBuffer = ''; }, 1000);
  }
  // Helper: Format text to 80 columns
  function formatTo80Cols(text) {
    // Split by paragraphs to preserve structural spacing
    const paragraphs = text.split('\n\n');
    const formattedParagraphs = paragraphs.map(p => {
      // Flatten the paragraph into a single array of words
      const words = p.replace(/\n/g, ' ').split(/\s+/).filter(w => w.length > 0);
      let currentLine = '';
      const lines = [];
      words.forEach(word => {
        if ((currentLine + word).length > 80) {
          lines.push(currentLine.trim());
          currentLine = word + ' ';
        } else {
          currentLine += word + ' ';
        }
      });
      if (currentLine.trim()) lines.push(currentLine.trim());
        return lines.join('\n');
      });

      return formattedParagraphs.join('\n\n');
  }

  // 2. Event Listener
  editor.addEventListener('keydown', (e) => {

    // Escape always drops back to Normal mode
    if (e.key === 'Escape') {
      setMode('normal');
      keyBuffer = '';
      return;
    }

    // If in Normal Mode, intercept keypresses
    if (mode === 'normal') {
      // Prevent standard typing in normal mode
      e.preventDefault();
      // Switch to Insert Mode
      if (e.key === 'i') {
        setMode('insert');
        keyBuffer = '';
        return;
      }

      const currentPos = editor.selectionStart;
      const text = editor.value;

      // Handle Vim Motions
      switch (e.key) {
        case 'G':
          // Bottom of file
          moveCursor(text.length);
          break;
        case '}': {
          // Next paragraph: find next double-newline
          const nextDoubleNewline = text.indexOf('\n\n', currentPos);
          moveCursor(nextDoubleNewline === -1 ? text.length : nextDoubleNewline + 2);
          break;
        }
        case '{': {
          // Previous paragraph: find previous double-newline
          // We offset the search slightly so we don't just get stuck on the current \n\n
          const searchPos = currentPos > 2 ? currentPos - 2 : 0;
          const prevDoubleNewline = text.lastIndexOf('\n\n', searchPos);
          moveCursor(prevDoubleNewline === -1 ? 0 : prevDoubleNewline);
          break;
        }
        case 'g':
          pushToBuffer('g');
          if (keyBuffer === 'gg') {
            // Top of file
            moveCursor(0);
            keyBuffer = '';
          }
          break;
        case 'q':
          pushToBuffer('q');
          if (keyBuffer === 'gq') {
            // Format the entire text area
            const formatted = formatTo80Cols(text);

            // Using setRangeText preserves undo history better than setting value directly
            editor.setRangeText(formatted, 0, text.length, 'start');

            // Re-orient the cursor
            moveCursor(currentPos);
            keyBuffer = '';
          }
          break;
      }
    }
  });
});
