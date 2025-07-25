using System;
using System.Windows.Forms;
using System.Net.Http;

namespace QuickBar
{
    static class Program
    {
        [STAThread]
        static void Main(string[] args)
        {
            // Check for API URL argument
            if (args.Length == 0)
            {
                MessageBox.Show("Please provide the API URL as a command-line argument.", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                return;
            }

            string baseUrl = args[0];

            try
            {
                Config.SetBaseUrl(baseUrl);
            }
            catch (ArgumentException ex)
            {
                MessageBox.Show($"Invalid Base URL provided: {ex.Message}", "Configuration Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                return;
            }

            // Perform health check before starting
            if (!PerformHealthCheck(Config.HealthCheckEndpoint))
            {
                MessageBox.Show("API health check failed. Service will not start.", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                return;
            }

            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new MainForm());
        }

        private static bool PerformHealthCheck(string healthCheckUrl)
        {
            try
            {
                using (HttpClient client = new HttpClient())
                {
                    client.Timeout = TimeSpan.FromSeconds(5);
                    HttpResponseMessage response = client.GetAsync(healthCheckUrl).Result; // .Result to make it synchronous for simple health check
                    return response.IsSuccessStatusCode;
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show($"Health check failed: {ex.Message}", "Health Check Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                return false;
            }
        }
    }
}