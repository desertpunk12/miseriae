//go:build js && wasm

package main

import (
	"cloudflare-worker-boilerplate/cms"
	"cloudflare-worker-boilerplate/pages"
	"cloudflare-worker-boilerplate/utils"
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"

	"github.com/a-h/templ"
)

func main() {
	fmt.Println("Go: main started")
	c := make(chan struct{})
	registerTemplRoute("renderIndex", pages.Miseriae())
	registerTemplRoute("renderHome", pages.Miseriae())
	registerTemplRoute("renderResume", pages.Resume())
	registerTemplRoute("renderBase", pages.Base("Base", nil, nil, ""))

	// Dynamic Routes
	js.Global().Set("renderBlog", js.FuncOf(renderBlog))
	js.Global().Set("renderCosplays", js.FuncOf(renderCosplays))

	// CMS Sync
	js.Global().Set("syncContent", js.FuncOf(syncContent))

	js.Global().Set("renderKV", js.FuncOf(utils.RenderKV))
	js.Global().Set("renderDynamicContent", js.FuncOf(renderDynamicContent))
	fmt.Println("Go: exports set, waiting...")
	<-c
}

func registerTemplRoute(funcName string, page templ.Component) {
	js.Global().Set(funcName, js.FuncOf(func(this js.Value, args []js.Value) any { return utils.RenderToString(page) }))
}

func renderBlog(this js.Value, args []js.Value) any {
	// 1. Try to get data from KV
	jsonStr, err := utils.KVGet("blog_data")
	var posts []cms.BlogPost

	// 2. Unmarshal if found
	if err == nil && jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &posts); err != nil {
			fmt.Println("Error unmarshaling blog_data:", err)
		}
	} else {
		fmt.Println("No blog_data found in KV or error:", err)
	}

	// 3. Render
	return utils.RenderToString(pages.Blog(posts))
}

func renderCosplays(this js.Value, args []js.Value) any {
	// 1. Try to get data from KV
	jsonStr, err := utils.KVGet("cosplay_data")
	var albums []cms.CosplayAlbum

	// 2. Unmarshal if found
	if err == nil && jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &albums); err != nil {
			fmt.Println("Error unmarshaling cosplay_data:", err)
		}
	} else {
		fmt.Println("No cosplay_data found in KV or error:", err)
	}

	// 3. Render
	return utils.RenderToString(pages.Cosplays(albums))
}

func syncContent(this js.Value, args []js.Value) any {
	// Args: [driveFolderID, driveApiKey, photosApiKey]
	if len(args) < 2 {
		return "Error: specific driveFolderID and driveApiKey required"
	}
	driveFolderID := args[0].String()
	driveApiKey := args[1].String()
	photosApiKey := ""
	if len(args) > 2 {
		photosApiKey = args[2].String()
	}

	// Run sync in a goroutine? No, we want to return the result content.
	// But sync might take time. We can return a Promise?
	// For simplicity, we'll try to do it synchronously-ish (which blocks JS main thread in Wasm usually)
	// But in standard Cloudflare Workers Wasm, blocking 'main' is okay for a bit.
	// Ideally we return a Promise.

	handler := js.FuncOf(func(this js.Value, pArgs []js.Value) interface{} {
		resolve := pArgs[0]
		reject := pArgs[1]

		go func() {
			status, err := cms.SyncContent(driveFolderID, driveApiKey, photosApiKey)
			if err != nil {
				reject.Invoke(err.Error())
			} else {
				resolve.Invoke(status)
			}
		}()
		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func renderDynamicContent(this js.Value, args []js.Value) any {
	now := time.Now()
	items := []string{
		fmt.Sprintf("Item generated at %s", now.Format(time.TimeOnly)),
		"Another dynamic item",
		"Random Value: " + fmt.Sprint(now.UnixNano()),
	}

	component := pages.DynamicContent(
		"Dynamic Data",
		items,
		now.Format(time.RFC3339),
		now.Format(time.RFC1123),
		now.Format(time.Kitchen),
		"/dynamic",
	)

	return utils.RenderToString(component)
}
