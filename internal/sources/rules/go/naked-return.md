---
type: rule
description: Do not use naked returns, always specify return values explicitly
globs: ["*.go"]
---

## Good

```go
func ParseConfig(path string) (Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return Config{}, err
    }
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return Config{}, err
    }
    return cfg, nil
}
```

## Bad

```go
func ParseConfig(path string) (cfg Config, err error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return
    }
    if err = json.Unmarshal(data, &cfg); err != nil {
        return
    }
    return
}
```
