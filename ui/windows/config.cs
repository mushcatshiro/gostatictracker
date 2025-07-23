using System;

namespace QuickBar
{
    public static class Config
    {
        private static string _baseUrl;
        public static void SetBaseUrl(string baseUrl)
        {
            if (string.IsNullOrWhiteSpace(baseUrl))
            {
                throw new ArgumentException("Base URL cannot be null or empty.", nameof(baseUrl));
            }
            // Ensure the base URL ends with a trailing slash for consistent path joining
            _baseUrl = baseUrl.TrimEnd('/') + "/";
        }

        private static string BaseUrlValidated
        {
            get
            {
                if (string.IsNullOrEmpty(_baseUrl))
                {
                    throw new InvalidOperationException("Base URL has not been set in Config. Call SetBaseUrl() first.");
                }
                return _baseUrl;
            }
        }

        public static string SearchEndpoint
        {
            get
            {
                return BaseUrlValidated + "search";
            }
        }

        public static string HealthCheckEndpoint
        {
            get
            {
                return BaseUrlValidated + "healthcheck";
            }
        }

        public static string CreateEventEndpoint
        {
            get
            {
                return BaseUrlValidated + "events";
            }
        }
    }
}