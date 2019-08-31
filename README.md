# Revive Action
This Action runs [Revive](https://github.com/mgechev/revive) on your [Go](https://golang.org/) code and adds optional annotations to the check.

## Usage

```YAML
- name: Check out code into the Go module directory
  uses: actions/checkout@v1

- name: Run Revive Action
  uses: morphy2k/revive-action@v1
  with:
    # Path to your Revive config within the repo (optional)
    config: revive/config.toml
  env:
    # GitHub token for annotations (optional)
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
