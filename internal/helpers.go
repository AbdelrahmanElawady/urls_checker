package internal

import "strings"

// CheckWebsite checks if the website has a prefix https ..?
func checkWebsite(website string) string {
	if !strings.HasPrefix(website, "https://") {
		website = "https://" + website
	}

	return website
}
