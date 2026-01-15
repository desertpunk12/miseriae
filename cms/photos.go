package cms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type MediaItem struct {
	ID          string `json:"id"`
	BaseUrl     string `json:"baseUrl"`
	MimeType    string `json:"mimeType"`
	Description string `json:"description"`
}

type PhotosListResponse struct {
	MediaItems    []MediaItem `json:"mediaItems"`
	NextPageToken string      `json:"nextPageToken"`
}

type Album struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type AlbumsListResponse struct {
	Albums []Album `json:"albums"`
}

// FetchCosplayAlbums fetches all albums and their contents
// Note: This requires the Photos Library API enabled.
// Ideally usage of Google Photos API requires OAuth2 user flow, but for a personal site
// with a long-lived refresh token or API key (public albums), it can work.
// For this implementation, we will assume standard API usage. We might need to upgrade to
// Service Account + Domain Delegation or OAuth Refresh Token flow in `sync.go` later.
func FetchCosplayAlbums(accessToken string) ([]CosplayAlbum, error) {
	// 1. List Albums
	// CAUTION: 'v1/albums' returns only albums created by the app.
	// We might need to use a specific album ID or list all.
	// url := fmt.Sprintf("https://photoslibrary.googleapis.com/v1/albums?key=%s", apiKey)

	// NOTE: Because Google Photos API is strict with Auth, this is a placeholder.
	return nil, fmt.Errorf("FetchCosplayAlbums requires OAuth Access Token, handled in sync.go")
}

// FetchCosplayAlbumDetails fetches details for a specific album using an Access Token
func FetchCosplayAlbumDetails(albumID, accessToken string) (CosplayAlbum, error) {
	// 1. Get Album Metadata
	// Actual implementation with proper request creating
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://photoslibrary.googleapis.com/v1/albums/%s", albumID), nil)

	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CosplayAlbum{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return CosplayAlbum{}, fmt.Errorf("photos api error: %d", resp.StatusCode)
	}

	var googleAlbum Album
	if err := json.NewDecoder(resp.Body).Decode(&googleAlbum); err != nil {
		return CosplayAlbum{}, err
	}

	// 2. List Media Items in Album
	searchUrl := "https://photoslibrary.googleapis.com/v1/mediaItems:search"
	searchBody := fmt.Sprintf(`{"albumId": "%s", "pageSize": "100"}`, albumID)

	req, _ = http.NewRequest("POST", searchUrl, strings.NewReader(searchBody))
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return CosplayAlbum{}, err
	}
	defer resp.Body.Close()

	var list PhotosListResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return CosplayAlbum{}, err
	}

	// 3. Construct CosplayAlbum
	album := CosplayAlbum{
		ID: albumID,
	}

	// Parse Title: "Ahri | League of Legends"
	parts := strings.Split(googleAlbum.Title, "|")
	if len(parts) > 0 {
		album.Title = strings.TrimSpace(parts[0])
	}
	if len(parts) > 1 {
		album.Series = strings.TrimSpace(parts[1])
	}

	// Parse Images
	for i, item := range list.MediaItems {
		// Google Photos Base URLs need parameters to be useful
		// =w1600-h1600 allows high res
		finalUrl := item.BaseUrl + "=w1920-h1080"
		album.Images = append(album.Images, finalUrl)

		// First image is cover
		if i == 0 {
			album.CoverImage = finalUrl
		}
	}

	// We'll leave metadata parsing logic (Description, Photographer)
	// typically this comes from the Album "shareInfo" or we look for a specific
	// "info.txt" or just use the Description of the *First Photo* (Cover)
	// Let's use the Cover Photo Description for metadata
	if len(list.MediaItems) > 0 {
		desc := list.MediaItems[0].Description
		parseMetadataFromDescription(desc, &album)
	}

	return album, nil
}

func parseMetadataFromDescription(desc string, album *CosplayAlbum) {
	lines := strings.Split(desc, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Photographer:") {
			album.Photographer = strings.TrimSpace(strings.TrimPrefix(line, "Photographer:"))
		} else if strings.HasPrefix(line, "Assistant:") {
			album.Assistant = strings.TrimSpace(strings.TrimPrefix(line, "Assistant:"))
		} else if strings.HasPrefix(line, "Location:") {
			album.Location = strings.TrimSpace(strings.TrimPrefix(line, "Location:"))
		} else if strings.HasPrefix(line, "Description:") {
			album.Description = strings.TrimSpace(strings.TrimPrefix(line, "Description:"))
		}
	}
	// Fallback/Cleanup
	if album.Description == "" {
		// If no explicit "Description:" prefix, use remaining lines?
		// For now keep simple.
	}
}
