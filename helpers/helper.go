package helpers

import (
	"log"
	"net/url"
	"os"
	"regexp"
	"time"
)

var Logger *os.File

func LogError(strings string, err error) {
	if err == nil {
		str := time.Now().Local().Format("2006-01-02 15:04:05") + ": " + strings + "\n"
		Logger.Write([]byte(str))
		return
	}
	str := time.Now().String() + ": " + strings + " " + err.Error() + "\n"
	Logger.Write([]byte(str))
}

func InitLogger() *os.File {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	Logger = file
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return Logger
}

func CleanUrl(l string, home_url string) string {
	re := regexp.MustCompile(`^https://www\.naukri\.com/job-listings(?:-[a-zA-Z0-9]+)*(-[0-9]+)?$`)

	str, err := url.Parse(l)
	if err != nil {
		return "" // Invalid URL; return empty string
	}

	// If the URL does not have a hostname, join with the home_url
	if len(str.Hostname()) == 0 {
		ourl, oerr := url.JoinPath(home_url, str.Path)
		if oerr != nil {
			// Log the error and return empty string
			LogError("ourl", oerr)
			return ""
		}
		// Check if the cleaned URL matches the regex
		if re.MatchString(ourl) {
			return ourl
		}
		return ""
	}

	// Ensure the URL uses HTTPS or HTTP
	if str.Scheme != "https" && str.Scheme != "http" {
		str.Scheme = "https"
	}

	// Parse the home URL to extract the hostname
	h, err := url.Parse(home_url)
	if err != nil || str.Hostname() != h.Hostname() {
		return "" // Return empty if hostnames do not match or if home_url is invalid
	}

	// Match the URL against the regex
	if re.MatchString(str.String()) {
		return str.String() // Valid and cleaned URL
	}
	return "" // Return empty if the URL does not match
}
