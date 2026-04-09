---
title: "Editor IntelliSense"
weight: 7
---

JSON Schemas ship in the nux repository for both config types:

- `schemas/global.schema.json` - global `config.yaml`
- `schemas/project.schema.json` - per-project YAML under `projects/`

Point your YAML language server at these schemas to get completion, validation, and hover docs in your editor.

## 1. Inline modeline

Add a modeline at the top of any config file:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/Drew-Daniels/nux/main/schemas/project.schema.json
```

Use `project.schema.json` for files in `projects/` and `global.schema.json` for `config.yaml`.

## 2. VS Code

In `.vscode/settings.json`, map schema URLs to file patterns:

```json
{
  "yaml.schemas": {
    "https://raw.githubusercontent.com/Drew-Daniels/nux/main/schemas/global.schema.json": [
      ".config/nux/config.yaml",
      "**/nux/config.yaml"
    ],
    "https://raw.githubusercontent.com/Drew-Daniels/nux/main/schemas/project.schema.json": [
      ".config/nux/projects/*.yaml"
    ]
  }
}
```

Adjust paths to match where you keep configs.

## 3. Neovim (yamlls)

If you use `nvim-lspconfig`, add schema mappings to your `yamlls` setup:

```lua
require("lspconfig").yamlls.setup({
  settings = {
    yaml = {
      schemas = {
        ["https://raw.githubusercontent.com/Drew-Daniels/nux/main/schemas/global.schema.json"] = {
          ".config/nux/config.yaml",
          "**/nux/config.yaml",
        },
        ["https://raw.githubusercontent.com/Drew-Daniels/nux/main/schemas/project.schema.json"] = {
          ".config/nux/projects/*.yaml",
        },
      },
    },
  },
})
```
