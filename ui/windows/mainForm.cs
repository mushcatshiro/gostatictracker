using System;
using System.Drawing;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Windows.Forms;
using System.Runtime.InteropServices; // Needed for P/Invoke

namespace QuickBar
{
    public class CreateEventPayload
    {
        public string title { get; set; }
        public string group { get; set; }
        public string start { get; set; }
    }

    public partial class MainForm : Form
    {
        private GlobalHotkey _hotkey;
        private readonly TextBox searchTextBox;

        // P/Invoke declarations for moving the form without title bar
        [DllImport("user32.dll")]
        public static extern bool ReleaseCapture();

        [DllImport("user32.dll")]
        public static extern int SendMessage(IntPtr hWnd, int Msg, int wParam, int lParam);

        public const int WM_NCLBUTTONDOWN = 0xA1;
        public const int HT_CAPTION = 0x2;


        public MainForm()
        {
            InitializeComponent();
            this.Hide();
            this.ShowInTaskbar = false;
            searchTextBox = new TextBox();
            _hotkey = new GlobalHotkey(GlobalHotkey.MOD_ALT, Keys.D1, this);
            SetupForm();
            RegisterGlobalHotkey();
        }

        private void SetupForm()
        {
            // Make the form borderless for a "Mac-style" look
            this.FormBorderStyle = FormBorderStyle.None;
            this.ShowInTaskbar = false; // Don't show in taskbar
            this.StartPosition = FormStartPosition.Manual; // We'll position it
            this.TopMost = true; // Keep it on top

            // Set initial size and position (e.g., centered at the top)
            this.Width = 600;
            this.Height = 40;
            CenterFormAtTop();

            // Add the search text box
            searchTextBox.Dock = DockStyle.Fill;
            searchTextBox.Font = new Font("Segoe UI", 16); // A common modern font
            searchTextBox.BorderStyle = BorderStyle.None; // Make it look like part of the form
            searchTextBox.KeyDown += SearchTextBox_KeyDown;
            searchTextBox.LostFocus += SearchTextBox_LostFocus; // Hide when focus is lost
            searchTextBox.CausesValidation = false; // Prevent validation on focus loss
            this.Controls.Add(searchTextBox);

            // Hide the form initially
            // this.Hide();

            // Allow dragging the form by clicking anywhere (since no title bar)
            this.MouseDown += (s, e) =>
            {
                if (e.Button == MouseButtons.Left)
                {
                    ReleaseCapture();
                    SendMessage(this.Handle, WM_NCLBUTTONDOWN, HT_CAPTION, 0);
                }
            };
        }

        private void CenterFormAtTop()
        {
            // Position the form at the top center of the primary screen
            Rectangle screenRect = Screen.PrimaryScreen.WorkingArea;
            this.Location = new Point((screenRect.Width - this.Width) / 2, (screenRect.Height - this.Height) / 2); // 50 pixels from the top
        }

        private void RegisterGlobalHotkey()
        {
            // Register Alt + 1 hotkey
            // Modifiers: MOD_ALT (0x0001)
            // Key: Keys.D1 (D1 represents the '1' key, not Numpad1)
            _hotkey.KeyPressed += new EventHandler(Hotkey_KeyPressed);

            if (!_hotkey.Register())
            {
                MessageBox.Show("Could not register hotkey Alt+1. It might be in use by another application.", "Hotkey Error", MessageBoxButtons.OK, MessageBoxIcon.Warning);
                // Optionally, exit the application if hotkey is crucial
                Application.Exit();
            }
        }

        private void Hotkey_KeyPressed(object sender, EventArgs e)
        {
            // When hotkey is pressed, show the search bar
            this.Show();
            this.Activate(); // Bring to front
            searchTextBox.Focus();
            searchTextBox.Clear();
        }

        private void SearchTextBox_KeyDown(object sender, KeyEventArgs e)
        {
            if (e.KeyCode == Keys.Enter)
            {
                string searchText = searchTextBox.Text.Trim();
                if (!string.IsNullOrEmpty(searchText))
                {
                    PostToApi(searchText);
                }
                this.Hide(); // Hide after pressing Enter
                e.SuppressKeyPress = true; // Prevent the default "ding" sound
            }
            else if (e.KeyCode == Keys.Escape)
            {
                this.Hide(); // Hide on Escape
            }
        }

        private void SearchTextBox_LostFocus(object sender, EventArgs e)
        {
            // Hide the form if the text box loses focus, unless the form is somehow being activated again
            // (e.g., clicking on another part of the form that doesn't steal focus immediately)
            if (!this.ContainsFocus) // Check if any control on this form has focus
            {
                this.Hide();
            }
        }

        private async void PostToApi(string title)
        {
            try
            {
                using (HttpClient client = new HttpClient())
                {
                    DateTime currentTime = DateTime.Now;
                    string formattedTime = currentTime.ToString("MM-dd-yyyy hh:mm");
                    CreateEventPayload payload = new CreateEventPayload {title = title, group = "Unclassified", start = formattedTime};
                    string jsonPayload = JsonSerializer.Serialize(payload);
                    // Create JSON payload
                    StringContent content = new StringContent(jsonPayload, Encoding.UTF8, "application/json");

                    HttpResponseMessage response = await client.PostAsync(Config.CreateEventEndpoint, content);

                    if (response.IsSuccessStatusCode)
                    {
                        // Optionally, read and display API response
                        string responseBody = await response.Content.ReadAsStringAsync();
                        MessageBox.Show("API Call Successful! Response: " + responseBody, "Success", MessageBoxButtons.OK, MessageBoxIcon.Information);
                    }
                    else
                    {
                        string errorResponse = await response.Content.ReadAsStringAsync();
                        MessageBox.Show($"API Call Failed: {response.StatusCode} - {errorResponse}", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show($"An error occurred during API call: {ex.Message}", "API Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }

        protected override void Dispose(bool disposing)
        {
            if (disposing && (_hotkey != null))
            {
                _hotkey.Unregister();
                _hotkey = null;
            }
            base.Dispose(disposing);
        }

        // Dummy InitializeComponent for demonstration. In a real WinForms project,
        // this would be generated by the designer.
        private void InitializeComponent()
        {
            this.SuspendLayout();
            //
            // MainForm
            //
            this.ClientSize = new System.Drawing.Size(284, 261);
            this.Name = "MainForm";
            this.Load += new System.EventHandler(this.MainForm_Load);
            this.ResumeLayout(false);

        }

        private void MainForm_Load(object sender, EventArgs e)
        {
            // Optional: Any initialization on form load
        }
    }
}