# Header To Query Plugin

A Traefik plugin that converts HTTP headers to URL query parameters. Supports mapping, renaming, and optionally keeping headers. Handles multiple headers with the same name.

## Installation & Enabling

Add the plugin to your Traefik static configuration:

```yaml
experimental:
  plugins:
    headertoquery:
      moduleName: github.com/zalbiraw/headertoquery
      version: v0.0.3
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
        - headertoquery

  middlewares:
    headertoquery:
      plugin:
        headers:
          - name: SERVICE_TAG
            key: id
          - name: RANK
          - name: GROUP
            keepHeader: true
```

## How It Works

- Each `headers` entry can specify:
  - `name`: The HTTP header to process
  - `key`: (Optional) The query parameter name to use (defaults to the header name)
  - `keepHeader`: (Optional) If `true`, the header is not removed from the request
- If a header appears multiple times, all values are mapped as repeated query parameters (e.g., `?id=1&id=2`).

### Example

Given this configuration:

```yaml
headers:
  - name: SERVICE_TAG
    key: id
  - name: RANK
  - name: GROUP
    keepHeader: true
```

And a request with these headers:

```
SERVICE_TAG: S117
SERVICE_TAG: SPARTAN-117
SERVICE_TAG: 117
RANK: Masterchief
GROUP: UNSC
```

The resulting query string will be:

```
?id=S117&id=SPARTAN-117&id=117&rank=Masterchief&group=UNSC
```

And the resulting headers will be:

```
GROUP: UNSC
```

The `SERVICE_TAG` and `RANK` headers are removed; `GROUP` remains because `keepHeader: true`.

## Development & Testing

Run tests with:

```sh
go test -v
```

---

## Traefik Local Plugin Support

You can use this project as a [Traefik Local Plugin](https://doc.traefik.io/traefik/plugins/local-plugins/). This allows you to develop and test the plugin locally, without needing to publish it to an external registry. Simply reference the plugin's local path in your Traefik configuration for rapid iteration and debugging.

---

## Terraform-Enabled Module

This repository includes a Terraform module (`main.tf`) that provisions all necessary plugin configuration and source code into a Kubernetes `ConfigMap`. This enables you to manage and deploy the plugin as infrastructure-as-code, integrating seamlessly with your Terraform workflows.

### Using Private Repositories with Terraform

If your plugin or configuration files are stored in a private repository, you can securely provide Terraform with access credentials:

- **HTTPS with Personal Access Token:**
  Use a URL like:
  ```hcl
  module "plugin" {
    source = "git::https://<TOKEN>@github.com/username/private-repo.git//module_path"
    # ...
  }
  ```
  Replace `<TOKEN>` with your GitHub/GitLab personal access token. Never commit secrets to version control.

- **SSH Keys:**
  Ensure your Terraform environment has access to the correct SSH key (e.g., via `~/.ssh/id_rsa` or by setting `GIT_SSH_COMMAND`).

- **Environment Variables:**
  Set environment variables (such as `GIT_ASKPASS` or provider-specific variables) before running `terraform init`.

- **Terraform Cloud/Enterprise:**
  Add secrets or environment variables in your workspace settings for secure access.

> **Security Tip:** Always use environment variables or secret managers to handle sensitive information. Avoid hardcoding secrets in `.tf` files.

---

For more details, see the source code and test cases.
