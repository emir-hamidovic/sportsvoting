package parser

import (
	"fmt"
	"math/rand"
	"net/http"
)

func setupRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set a User-Agent header to impersonate a browser user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	// Set a random Accept-Language header for each request
	langs := []string{"en-US", "en-GB", "fr-FR", "de-DE", "es-ES", "pt-PT", "it-IT", "ja-JP", "ko-KR", "zh-CN"}
	lang := langs[rand.Intn(len(langs))]
	req.Header.Set("Accept-Language", lang)
	return req, nil
}

func SendRequest(url string) (*http.Response, error) {
	req, err := setupRequest(url)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return res, nil
}
