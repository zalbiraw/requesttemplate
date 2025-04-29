# Request Template Plugin

A Traefik plugin that transforms HTTP JSON request bodies using a [Go text/template](https://pkg.go.dev/text/template) filter. This allows you to flexibly modify, map, or validate incoming requests before they reach your backend service.

## Features
- Apply a Go text/template transformation to the incoming JSON request body
- Easily map or rewrite fields before forwarding the request

## Installation & Enabling

Add the plugin to your Traefik static configuration:

```yaml
experimental:
  plugins:
    requesttemplate:
      moduleName: github.com/zalbiraw/requesttemplate
      version: v0.0.1
```

## Dynamic Configuration Example

```yaml
http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - requesttemplate

  middlewares:
    requesttemplate:
      plugin:
        template: '{"message": "hello, {{ .user.message }}"}'
```

## How It Works

1. The plugin reads the incoming request body (expects JSON).
2. The Go template in `template` is rendered using the parsed JSON as data.
   - If the template fails, a 400 error is returned.
   - If successful, the result is sent as the new request body.
3. If the request body is empty or not JSON, the request passes through unchanged.

## Example

**Input:**
```json
{
  "user": { "message": "world" }
}
```

**Config:**
```yaml
template: '{"message": "hello, {{ .user.message }}"}'
```

**Output:**
```json
{
  "message": "hello, world"
}
```

## Notes
- Only JSON request bodies are processed.
- The plugin uses Go's [text/template](https://pkg.go.dev/text/template) for templating.
- Errors in JSON or template rendering result in a 400 Bad Request response.

---

For more details, see the source code or open an issue on GitHub.
