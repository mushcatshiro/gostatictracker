using System;
using System.Windows.Forms;
using System.Runtime.InteropServices;

namespace QuickBar
{
    public class GlobalHotkey : IDisposable
    {
        // P/Invoke declarations
        [DllImport("user32.dll")]
        private static extern bool RegisterHotKey(IntPtr hWnd, int id, int fsModifiers, int vk);

        [DllImport("user32.dll")]
        private static extern bool UnregisterHotKey(IntPtr hWnd, int id);

        // Modifier flags
        public const int MOD_ALT = 0x0001;
        public const int MOD_CONTROL = 0x0002;
        public const int MOD_SHIFT = 0x0004;
        public const int MOD_WIN = 0x0008;

        // Message ID for hotkeys
        public const int WM_HOTKEY = 0x0312;

        private int _modifier;
        private Keys _key;
        private IntPtr _hWnd;
        private int _id;

        // A hidden window to receive hotkey messages
        private MessageWindow _messageWindow;

        public event EventHandler KeyPressed;

        public GlobalHotkey(int modifier, Keys key, Form form)
        {
            _modifier = modifier;
            _key = key;
            // The hWnd passed to RegisterHotKey must be a valid window handle.
            // We'll use a hidden message-only window for this.
            _messageWindow = new MessageWindow(this);
            _hWnd = _messageWindow.Handle; // Get the handle of our hidden message window
            _id = GetHashCode(); // Use a unique ID based on the object's hash code
        }

        // Inner class that creates a hidden window to receive messages
        private class MessageWindow : NativeWindow, IDisposable
        {
            private GlobalHotkey _parent;

            public MessageWindow(GlobalHotkey parent)
            {
                _parent = parent;
                // Create a message-only window. This window is not visible.
                // It's ideal for receiving messages that don't need a UI.
                CreateHandle(new CreateParams());
            }

            protected override void WndProc(ref Message m)
            {
                if (m.Msg == WM_HOTKEY && (int)m.WParam == _parent._id)
                {
                    // Hotkey pressed, raise the event in the parent GlobalHotkey instance
                    _parent.OnKeyPressed();
                }
                base.WndProc(ref m);
            }

            public void Dispose()
            {
                DestroyHandle();
            }
        }

        private void OnKeyPressed()
        {
            if (KeyPressed != null)
            {
                KeyPressed(this, EventArgs.Empty);
            }
        }

        public bool Register()
        {
            // Register the hotkey with the handle of our hidden message window
            return RegisterHotKey(_hWnd, _id, _modifier, (int)_key);
        }

        public bool Unregister()
        {
            // Unregister the hotkey
            return UnregisterHotKey(_hWnd, _id);
        }

        public void Dispose()
        {
            Unregister();
            if (_messageWindow != null)
            {
                _messageWindow.Dispose();
                _messageWindow = null;
            }
            // Suppress finalization, as we've explicitly cleaned up.
            GC.SuppressFinalize(this);
        }

        // Destructor (finalizer) to ensure unregistration if Dispose is not called
        ~GlobalHotkey()
        {
            Dispose();
        }
    }
}