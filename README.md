# Revive Action
This Action runs [Revive](https://github.com/mgechev/revive) on your [Go](https://golang.org/) code and adds optional annotations to the check.

## Usage

Checkout
```YAML
- name: Check out code into the Go module directory
  uses: actions/checkout@v1
```
Building from repository
```YAML
- name: Run Revive Action by building from repository
  uses: morphy2k/revive-action@v1
```
Pulling pre-built image
```YAML
- name: Run Revive Action by pulling pre-built image
  uses: docker.pkg.github.com/revive-action/revive-action:v1
```
Configuration
```YAML
  with:
    # Path to your Revive config within the repo (optional)
    config: revive/config.toml
    # Exclude patterns, separated by semicolons (optional)
    exclude: "file.go;foo/bar.go;./foo/bar/..."
  env:
    # GitHub token for annotations (optional)
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
