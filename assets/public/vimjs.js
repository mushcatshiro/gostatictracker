<script type="text/javascript">
(function() {
    // 1. Handle Global Variable Redeclaration
    // We attach it to window so it survives the swap or we check for existence.
    window.referenceCount = window.referenceCount || 0;

    // 2. Define functions on window so HTML 'onclick' attributes can find them
    window.toggleElement = window.toggleElement || function(id) {
        const el = document.getElementById(id);
        if (el) el.classList.toggle('visible');
    };

    window.insertMarkdownRef = function(url) {
        const textarea = document.getElementById('vim-textarea');
        if (!textarea) return;

        window.referenceCount++;
        const inlineMarker = `[${window.referenceCount}]`;
        const footerEntry = `\n[${window.referenceCount}]: ${url}\n`;

        textarea.focus();
        const start = textarea.selectionStart;
        const end = textarea.selectionEnd;

        textarea.setRangeText(inlineMarker, start, end, 'end');
        textarea.value = textarea.value.trimEnd() + "\n" + footerEntry;

        const newPos = start + inlineMarker.length;
        textarea.setSelectionRange(newPos, newPos);
        textarea.dispatchEvent(new Event('input', { bubbles: true }));
    };

    // 3. Dependency Loading Logic
    function initVim() {
        if (window.vim) {
            vim.open({
                debug: false,
                element: '#vim-textarea',
                showMsg: function(msg) { console.log('Vim.js:', msg); }
            });
        }
    }

    if (!window.vim) {
        const script = document.createElement('script');
        script.src = "/static/vim.js";
        script.onload = initVim;
        document.head.appendChild(script);
    } else {
        initVim();
    }

    // 4. THE GLOBAL LISTENER (Use addEventListener instead of onclick)
    // We only add this once to the window to avoid the "classList" error.
    if (!window.hasEditorListener) {
        window.addEventListener('click', function(event) {
            const uploadContainer = document.getElementById('upload-container');
            const assetDropdown = document.getElementById('asset-dropdown');

            // Defensive check: if these aren't on the current page, do nothing
            if (!uploadContainer || !assetDropdown) return;

            if (!event.target.matches('button') &&
                !event.target.closest('#upload-container') &&
                !event.target.closest('#asset-dropdown')) {
                uploadContainer.classList.remove('visible');
                assetDropdown.classList.remove('visible');
            }
        });
        window.hasEditorListener = true;
    }

    // 5. text area auto resize
    const autoResize = function(tx) {
      tx.style.height = 'auto';
      tx.style.height = tx.scrollHeight + 'px';
    }
  const txs = document.querySelectorAll('textarea');
  txs.forEach((tx) => {
    tx.addEventListener('input', () => autoResize(tx))
    setTimeout(() => autoResize(tx), 0);
  })
  // 6. focus snap
  const vimTx = document.getElementById('vim-textarea');
  if (vimTx) {
    vimTx.addEventListener('keyup', function(e) {
      // Only trigger on keys that typically cause large jumps
      if (e.key === 'g' || e.key === 'G') {
        // Use a tiny timeout to let Vim.js finish moving the internal cursor first
        setTimeout(() => { vimTx.blur(); vimTx.focus();}, 50);
      }
    });
  }
})();
</script>
