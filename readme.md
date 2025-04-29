# Request Template Plugin

A Traefik plugin that transforms HTTP JSON request bodies using [jq](https://stedolan.github.io/jq/) filters. This allows you to flexibly modify, map, or validate incoming requests before they reach your backend service.

## Features
- Apply one or more jq commands to the incoming JSON request body
- Chain transformations: output of one command feeds into the next

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
        commands:
          - .user.message |= "hello, " + .
          - .timestamp = now
```

## How It Works

1. The plugin reads the incoming request body (expects JSON).
2. Each jq command in the `commands` list is applied in order.
   - If any command fails, a 400 error is returned.
   - If all succeed, the final result is marshaled and sent as the response body.
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
commands:
  - .user.message |= "hello, " + .
```

**Output:**
```json
{
  "user": { "message": "hello, world" }
}
```

## Notes
- Only JSON request bodies are processed.
- The plugin uses [gojq](https://github.com/itchyny/gojq) for jq support.
- Errors in JSON or jq filter result in a 400 Bad Request response.

---

For more details, see the source code or open an issue on GitHub.
