# Revive Action

This Action runs [Revive](https://github.com/mgechev/revive) on your [Go](https://golang.org/) code and adds optional annotations to the check.

## Usage

Checkout

```YAML
- name: Check out code into the Go module directory
  uses: actions/checkout@v2
```

Use by building from repository

```YAML
- name: Run Revive Action by building from repository
  uses: morphy2k/revive-action@v1
```

Use by pulling pre-built image

```YAML
- name: Run Revive Action by pulling pre-built image
  uses: docker://morphy/revive-action:v1
```

Configuration

```YAML
  with:
    # Path to your Revive config within the repo (optional)
    config: revive/config.toml
    # Exclude patterns, separated by semicolons (optional)
    exclude: "file.go;foo/bar.go;./foo/bar/..."
    # Path pattern (default: ./...)
    path: "./foo/..."
    # Name of the ckeck (default: revive-action)
    name: "Revive"
  env:
    # GitHub token for annotations (optional)
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Screenshots

![Screenshot of annotations](screenshot.png)

_GitHub Annotations_
