# detour

Windows shortcut replacer tool.

## Usage

### Generating a rule set JSON file

```sh
$ detour gen myrules.json
```

```json
{
  "Rules": [
    {
      "name": "C: -> D:",
      "old": "C:",
      "new": "D:"
    },
    {
      "name": "detour -> shortcut",
      "old": "detour",
      "new": "shortcut"
    },
    {
      "old": "\\bg.",
      "new": "go",
      "regexp": true
    }
  ]
}
```

Note: case insensitive

### execute

```sh
$ detour -v --rule-set myrules.json ./subdir/**/*
```
