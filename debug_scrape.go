//go:build ignore

package main

import (
	"fmt"
	scraper "instafix/handlers/scraper"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run debug_scrape.go <postID>")
		fmt.Println("Example: go run debug_scrape.go DDgMSthvmFa")
		os.Exit(1)
	}
	postID := os.Args[1]

	scraper.InitLRU(16)
	scraper.InitDB()
	defer scraper.DB.Close()

	item, err := scraper.GetData(postID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("PostID:       %s\n", item.PostID)
	fmt.Printf("Username:     %s\n", item.Username)
	fmt.Printf("Caption:      %.100s\n", item.Caption)
	fmt.Printf("VideoBlocked: %v\n", item.VideoBlocked)
	fmt.Printf("Medias:       %d\n", len(item.Medias))
	for i, m := range item.Medias {
		fmt.Printf("  [%d] TypeName: %s\n", i, m.TypeName)
		fmt.Printf("      URL:      %s\n", m.URL)
		isVideo := m.TypeName == "GraphVideo" || m.TypeName == "XDTGraphVideo"
		urlLooksLikeVideo := len(m.URL) > 0 && (contains(m.URL, ".mp4") || contains(m.URL, "video"))
		fmt.Printf("      IsVideo:  %v (URL looks like video: %v)\n", isVideo, urlLooksLikeVideo)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
