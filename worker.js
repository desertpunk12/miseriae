import "./wasm_exec.js";

const go = new Go();
let wasmInstance;

async function initWasm(env) {
  if (!wasmInstance) {
    // In Cloudflare Workers, WASM modules are often bound as variables or imported
    // Here we assume 'main.wasm' is handled by the bundler/wrangler and available.
    // However, Wrangler 2+ usually imports wasm as a module if configured or we use the import.
    // Let's rely on the module import standard for workers if possible, or fetch it.

    // BUT the standard 'wasm_exec.js' expects to instantiate streaming or from array buffer.
    // We will import the wasm module directly which Wrangler supports.

    // Note: The below import assumes we are using Esbuild/Webpack/Wrangler's module support
    // which effectively makes the .wasm file available.
    // Usually for 'CompiledWasm' rule, we can import it.

    // Let's try to import the wasm file (wrangler handles this)
    const wasm = await import("./main.wasm");

    // 'wasm' export is usually the Module object or we need to instantiate it.
    // Wrangler "rules" -> CompiledWasm means `import mod from './main.wasm'` gives a WebAssembly.Module.

    wasmInstance = await WebAssembly.instantiate(wasm.default, go.importObject);
    // Monitor Go execution
    go.run(wasmInstance);

    // Wait for Go to initialize exports
    let retries = 0;
    while (typeof globalThis.renderIndex !== "function" && retries < 40) {
      await new Promise((r) => setTimeout(r, 50));
      retries++;
    }

    if (typeof globalThis.renderIndex !== "function") {
      console.error("WASM failed to initialize exports in time.");
    }
  }
}

const ROUTES = {
  "/": {
    func: "renderIndex",
    args: ["Wingo Cloduflare Worker App"],
  },
  "/dynamic": {
    func: "renderDynamicContent",
  },
  "/home": {
    func: "renderHome",
  },
  "/cosplays": {
    func: "renderCosplays",
  },
  "/resume": {
    func: "renderResume",
  },
  "/base": {
    func: "renderBase",
  },
  "/blog": {
    func: "renderBlog",
  },
  "/kv": {
    func: "renderKV",
    setup: (env) => {
      // Bind the KV namespace to global scope so Go can find it using js.Global().Get("KV")
      // 'KV_DEMO' must match the binding name in wrangler.toml
      globalThis.KV = env.KV_DEMO;
    },
  },
};

export default {
  async fetch(request, env, ctx) {
    try {
      await initWasm(env);

      // Try to serve static assets first
      if (env.ASSETS) {
        try {
          const assetResponse = await env.ASSETS.fetch(request);
          if (assetResponse.status !== 404) {
            return assetResponse;
          }
        } catch (e) {
          // Fall through to dynamic routes
        }
      }

      const url = new URL(request.url);

      // Handle Google Drive Photo Proxy
      if (url.pathname.startsWith("/gdrivephoto/")) {
        const fileId = url.pathname.replace("/gdrivephoto/", "");
        if (fileId) {
          const driveUrl = `https://drive.google.com/uc?id=${fileId}`;
          const imageResponse = await fetch(driveUrl);
          // Create a new response to allow embedding (CORS/Headers if needed, 
          // though usually direct proxying works fine for basic embedding).
          // We return the response directly.
          return new Response(imageResponse.body, {
            headers: imageResponse.headers,
            status: imageResponse.status,
            statusText: imageResponse.statusText
          });
        }
      }

      // Handle Google Photos Proxy
      if (url.pathname.startsWith("/gphotophoto/")) {
        const shareId = url.pathname.replace("/gphotophoto/", "");
        if (shareId) {
          const photoUrl = `https://photos.app.goo.gl/${shareId}`;
          // Fetch the page with redirect follow to get the final URL
          const pageResponse = await fetch(photoUrl, {
            redirect: "follow",
            headers: {
              "User-Agent":
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
            },
          });

          const html = await pageResponse.text();
          // Look for og:image
          const match = html.match(
            /<meta\s+property="og:image"\s+content="([^"]+)"/,
          );
          if (match && match[1]) {
            let imageUrl = match[1];
            // Remove existing parameters if any (usually after an =)
            // and append =d to get the original/download size
            // or =w10000-h10000 to get a very large version displayed
            // The og:image usually ends with =w...-h...
            // We can replace the suffix or just append if we are careful.
            // Safest: find the last = and replace everything after it with d or w9999

            // Regex to replace the last =... part or append if missing
            // Using =w16383-h16383-no (max size generally) to get full size without forcing download
            const fullSizeParam = "=w16383-h16383-no";

            if (imageUrl.includes("=")) {
              imageUrl = imageUrl.substring(0, imageUrl.lastIndexOf("=")) + fullSizeParam;
            } else {
              imageUrl += fullSizeParam;
            }

            const imageResponse = await fetch(imageUrl);

            // Recreate headers to strip Content-Disposition if present
            const newHeaders = new Headers(imageResponse.headers);
            newHeaders.delete("Content-Disposition");

            return new Response(imageResponse.body, {
              headers: newHeaders,
              status: imageResponse.status,
              statusText: imageResponse.statusText,
            });
          }
        }
      }

      const route = ROUTES[url.pathname];

      if (route) {
        // Run any route-specific setup
        if (route.setup) {
          route.setup(env);
        }

        const funcName = route.func;
        if (typeof globalThis[funcName] !== "function") {
          throw new Error(
            `${funcName} is not defined. WASM may not have initialized correctly.`,
          );
        }

        // Call the function, awaiting it just in case it returns a Promise (like renderKV)
        // If args are provided, pass them; otherwise call without args
        const html = await globalThis[funcName](...(route.args || []));

        return new Response(html, {
          headers: { "Content-Type": "text/html" },
        });
      }

      return new Response("Not Found", { status: 404 });
    } catch (err) {
      console.error(err);
      return new Response(`Error: ${err.message}\nStack: ${err.stack}`, {
        status: 500,
      });
    }
  },
};
