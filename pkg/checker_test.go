package pkg

import "testing"

func TestUrlsChecker(t *testing.T) {

	t.Run("test_invalid_url", func(t *testing.T) {
		website := "https://golangz.org"
		err := Check(website)

		if err != nil {
			t.Errorf("Invalid URL")
		}
	})
}
