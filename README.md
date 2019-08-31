# Revive Action
This Action runs [Revive](https://github.com/mgechev/revive) on your [Go](https://golang.org/) code and adds optional annotations to the check.

## Usage

```YAML
- name: Check out code into the Go module directory
  uses: actions/checkout@v1

- name: Run Revive Action
  use: morphy2k/revive-action@v1
  with:
    # Path to your Revive config within the repo (default: revive.toml)
    config: revive/config.toml
  env:
    # GitHub token needed if you want to use annotations
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
