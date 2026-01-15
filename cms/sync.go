//go:build js && wasm

package cms

import (
	"cloudflare-worker-boilerplate/utils"
	"encoding/json"
	"fmt"
)

// SyncContent orchestrates fetching from Drive/Photos and saving to KV
func SyncContent(driveFolderID, driveApiKey, photosApiKey string) (string, error) {
	status := "Starting Sync...\n"

	// 1. Sync Blog Posts
	status += fmt.Sprintf("Fetching Blog Posts from Folder: %s...\n", driveFolderID)
	posts, err := FetchBlogPosts(driveFolderID, driveApiKey)
	if err != nil {
		status += fmt.Sprintf("Error fetching posts: %v\n", err)
	} else {
		status += fmt.Sprintf("Found %d posts.\n", len(posts))

		// Serialize and Store
		postsJSON, _ := json.Marshal(posts)
		if err := utils.KVSet("blog_data", string(postsJSON)); err != nil {
			status += fmt.Sprintf("Error saving blog_data to KV: %v\n", err)
		} else {
			status += "Saved blog_data to KV.\n"
		}
	}

	// 2. Sync Cosplay Albums
	// For now, we assume public albums or API key access (which has limits as discussed)
	// In a real app we'd need an Access Token passed in here.
	// For this demo, assuming 'photosApiKey' is actually an access token or we skip if empty.
	if photosApiKey != "" {
		status += "Fetching Cosplay Albums...\n"
		// TODO: In real usage, photosApiKey here should be an OAuth Access Token
		albums, err := FetchCosplayAlbums(photosApiKey)
		if err != nil {
			status += fmt.Sprintf("Error fetching albums: %v\n", err)
		} else {
			status += fmt.Sprintf("Found %d albums.\n", len(albums))

			albumsJSON, _ := json.Marshal(albums)
			if err := utils.KVSet("cosplay_data", string(albumsJSON)); err != nil {
				status += fmt.Sprintf("Error saving cosplay_data to KV: %v\n", err)
			} else {
				status += "Saved cosplay_data to KV.\n"
			}
		}
	} else {
		status += "Skipping Photos Sync (No API Key/Token provided).\n"
	}

	status += "Sync Complete."
	return status, nil
}
