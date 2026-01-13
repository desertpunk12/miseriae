//go:build js && wasm

package utils

import (
	"bytes"
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/a-h/templ"
)

func RenderToString(c templ.Component) string {
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		return fmt.Sprintf("<div>Error rendering component: %v</div>", err)
	}
	return buf.String()
}

func RenderKV(this js.Value, args []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			defer func() {
				if r := recover(); r != nil {
					reject.Invoke(fmt.Sprintf("Panic in renderKV: %v", r))
				}
			}()

			key := "kv_demo_key"

			// GET
			val, err := KVGet(key)
			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to get KV value: %v", err))
				return
			}

			displayVal := val
			if displayVal == "" {
				displayVal = "(empty - first run?)"
			}

			// SET
			newVal := fmt.Sprintf("Updated at %s from Go WASM", time.Now().Format(time.RFC1123))
			err = KVSet(key, newVal)
			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to set KV value: %v", err))
				return
			}

			// Render simple HTML
			html := fmt.Sprintf(`
				<div style="font-family: sans-serif; padding: 2rem; max-width: 600px; margin: 0 auto;">
					<h1 style="color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 0.5rem;">Cloudflare KV + Go WASM</h1>

					<div style="background: #f8f9fa; border: 1px solid #e9ecef; border-radius: 8px; padding: 1.5rem; margin-top: 1.5rem;">
						<h3 style="margin-top: 0;">Previous Value:</h3>
						<pre style="background: #e9ecef; padding: 0.5rem; border-radius: 4px;">%s</pre>
					</div>

					<div style="background: #d4edda; color: #155724; border: 1px solid #c3e6cb; border-radius: 8px; padding: 1.5rem; margin-top: 1rem;">
						<strong>Success!</strong> Value has been updated.
						<div style="margin-top: 0.5rem;">New Value: %s</div>
					</div>

					<div style="margin-top: 2rem; text-align: center;">
						<p><small>Refresh the page to see the new value cycle through.</small></p>
						<a href="/" style="color: #3498db; text-decoration: none; font-weight: bold;">&larr; Back to Home</a>
					</div>
				</div>
			`, displayVal, newVal)

			resolve.Invoke(html)
		}()

		return nil
	})

	return js.Global().Get("Promise").New(handler)
}
