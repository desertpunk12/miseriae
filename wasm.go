//go:build js && wasm

package main

import (
	"cloudflare-worker-boilerplate/pages"
	"cloudflare-worker-boilerplate/utils"
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
	registerTemplRoute("renderCosplays", pages.Cosplays())
	registerTemplRoute("renderResume", pages.Resume())
	registerTemplRoute("renderBase", pages.Base("Base", nil, nil))
	registerTemplRoute("renderBlog", pages.Blog())

	js.Global().Set("renderKV", js.FuncOf(utils.RenderKV))
	js.Global().Set("renderDynamicContent", js.FuncOf(renderDynamicContent))
	fmt.Println("Go: exports set, waiting...")
	<-c
}

func registerTemplRoute(funcName string, page templ.Component) {
	js.Global().Set(funcName, js.FuncOf(func(this js.Value, args []js.Value) any { return utils.RenderToString(page) }))
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
