//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/klauspost/compress/gzhttp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run debug_scrape2.go <postID>")
		os.Exit(1)
	}
	postID := os.Args[1]
	transport := gzhttp.Transport(http.DefaultTransport, gzhttp.TransportAlwaysDecompress(true))
	client := http.Client{Transport: transport}

	// 1. Fetch embed page
	fmt.Println("=== EMBED PAGE ===")
	req, _ := http.NewRequest("GET", "https://www.instagram.com/p/"+postID+"/embed/captioned/", nil)
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Embed fetch error: %v\n", err)
	} else {
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()
		fmt.Printf("Status: %d, Body length: %d\n", res.StatusCode, len(body))
		fmt.Printf("Contains 'shortcode_media': %v\n", bytes.Contains(body, []byte("shortcode_media")))
		fmt.Printf("Contains 'WatchOnInstagram': %v\n", bytes.Contains(body, []byte("WatchOnInstagram")))
		fmt.Printf("Contains 'EmbeddedMediaImage': %v\n", bytes.Contains(body, []byte("EmbeddedMediaImage")))
		fmt.Printf("Contains 'EmbeddedMediaVideo': %v\n", bytes.Contains(body, []byte("EmbeddedMediaVideo")))
	}

	// 2. Fetch GraphQL API
	fmt.Println("\n=== GRAPHQL API ===")
	gqlParams := url.Values{
		"variables": {`{"shortcode":"` + postID + `"}`},
		"doc_id":    {"10015901848480474"},
		"lsd":       {"AVqbxe3J_YA"},
	}
	req2, _ := http.NewRequest("POST", "https://www.instagram.com/api/graphql", strings.NewReader(gqlParams.Encode()))
	req2.Header = http.Header{
		"Content-Type":      {"application/x-www-form-urlencoded"},
		"User-Agent":        {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36"},
		"X-Ig-App-Id":       {"936619743392459"},
		"X-Fb-Lsd":          {"AVqbxe3J_YA"},
		"X-Asbd-Id":         {"129477"},
		"Sec-Fetch-Site":    {"same-origin"},
	}
	res2, err := client.Do(req2)
	if err != nil {
		fmt.Printf("GQL error: %v\n", err)
	} else {
		body2, _ := io.ReadAll(res2.Body)
		res2.Body.Close()
		fmt.Printf("Status: %d\n", res2.StatusCode)
		fmt.Printf("Contains 'require_login': %v\n", bytes.Contains(body2, []byte("require_login")))
		fmt.Printf("Contains 'video_url': %v\n", bytes.Contains(body2, []byte("video_url")))
		fmt.Printf("Contains 'shortcode_media': %v\n", bytes.Contains(body2, []byte("shortcode_media")))
		fmt.Printf("Contains 'xdt_shortcode_media': %v\n", bytes.Contains(body2, []byte("xdt_shortcode_media")))
		// Print first 500 chars
		if len(body2) > 500 {
			fmt.Printf("Body (first 500): %s\n", body2[:500])
		} else {
			fmt.Printf("Body: %s\n", body2)
		}
	}

	// 3. Fetch ?__a=1&__d=dis API
	fmt.Println("\n=== API (?__a=1&__d=dis) ===")
	req3, _ := http.NewRequest("GET", "https://www.instagram.com/p/"+postID+"/?__a=1&__d=dis", nil)
	req3.Header = http.Header{
		"User-Agent":        {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36"},
		"X-Ig-App-Id":       {"936619743392459"},
		"X-Requested-With":  {"XMLHttpRequest"},
		"Sec-Fetch-Site":    {"same-origin"},
	}
	res3, err := client.Do(req3)
	if err != nil {
		fmt.Printf("API error: %v\n", err)
	} else {
		body3, _ := io.ReadAll(res3.Body)
		res3.Body.Close()
		fmt.Printf("Status: %d\n", res3.StatusCode)
		fmt.Printf("Contains 'video_versions': %v\n", bytes.Contains(body3, []byte("video_versions")))
		fmt.Printf("Contains 'items': %v\n", bytes.Contains(body3, []byte("items")))
		fmt.Printf("Contains 'graphql': %v\n", bytes.Contains(body3, []byte("graphql")))
		fmt.Printf("Contains 'require_login': %v\n", bytes.Contains(body3, []byte("require_login")))
		if len(body3) > 500 {
			fmt.Printf("Body (first 500): %s\n", body3[:500])
		} else {
			fmt.Printf("Body: %s\n", body3)
		}
	}
}
