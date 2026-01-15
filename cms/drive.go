package cms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type DriveFile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
}

type DriveListResponse struct {
	Files []DriveFile `json:"files"`
}

func FetchBlogPosts(folderID, apiKey string) ([]BlogPost, error) {
	// 1. List files in folder
	url := fmt.Sprintf("https://www.googleapis.com/drive/v3/files?q='%s'+in+parents&key=%s", folderID, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google drive api error: %s", string(body))
	}

	var list DriveListResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	var posts []BlogPost

	// 2. Fetch content for each file
	for _, file := range list.Files {
		// Only process text files or google docs
		if strings.Contains(file.MimeType, "folder") {
			continue
		}

		content, err := downloadFileContent(file.ID, file.MimeType, apiKey)
		if err != nil {
			fmt.Printf("Error fetching file %s: %v\n", file.Name, err)
			continue
		}

		post := parseBlogPost(file.ID, content)
		posts = append(posts, post)
	}

	return posts, nil
}

func downloadFileContent(fileID, mimeType, apiKey string) (string, error) {
	var url string
	if strings.Contains(mimeType, "google-apps.document") {
		// Export Google Docs as text/plain for easier parsing, or text/html
		url = fmt.Sprintf("https://www.googleapis.com/drive/v3/files/%s/export?mimeType=text/plain&key=%s", fileID, apiKey)
	} else {
		// Download raw content for other types
		url = fmt.Sprintf("https://www.googleapis.com/drive/v3/files/%s?alt=media&key=%s", fileID, apiKey)
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download error status: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func parseBlogPost(id, content string) BlogPost {
	post := BlogPost{
		ID:          id,
		Title:       "Untitled",
		HTMLContent: "",
	}

	// Simple metadata parser
	// Format:
	// Title: ...
	// Date: ...
	// Tags: ...
	// ---
	// Content...

	lines := strings.Split(content, "\n")
	var bodyLines []string
	inBody := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inBody {
			if trimmed == "---" {
				inBody = true
				continue
			}
			if strings.HasPrefix(trimmed, "Title:") {
				post.Title = strings.TrimSpace(strings.TrimPrefix(trimmed, "Title:"))
			} else if strings.HasPrefix(trimmed, "Date:") {
				post.Date = strings.TrimSpace(strings.TrimPrefix(trimmed, "Date:"))
			} else if strings.HasPrefix(trimmed, "Type:") {
				post.Type = strings.TrimSpace(strings.TrimPrefix(trimmed, "Type:"))
			} else if strings.HasPrefix(trimmed, "Image:") {
				post.ImageURL = strings.TrimSpace(strings.TrimPrefix(trimmed, "Image:"))
			} else if strings.HasPrefix(trimmed, "Summary:") {
				post.Summary = strings.TrimSpace(strings.TrimPrefix(trimmed, "Summary:"))
			} else if strings.HasPrefix(trimmed, "Tags:") {
				tagsRaw := strings.TrimPrefix(trimmed, "Tags:")
				for _, t := range strings.Split(tagsRaw, ",") {
					post.Tags = append(post.Tags, strings.TrimSpace(t))
				}
			}
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	// If no metadata separator found, treat whole file as body
	if !inBody && len(bodyLines) == 0 {
		post.HTMLContent = content
	} else {
		// Convert body lines to HTML paragraphs (very basic markdown-ish)
		// In a real app we'd use a markdown parser library
		var htmlBuilder strings.Builder
		for _, line := range bodyLines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			htmlBuilder.WriteString("<p>" + line + "</p>")
		}
		post.HTMLContent = htmlBuilder.String()
	}

	return post
}
