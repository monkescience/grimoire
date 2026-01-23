---
type: rule
description: Use short (1-2 letter) receiver names based on the type name, never this or self
globs: ["*.go"]
---

## Good

```go
func (s *Server) Start() error
func (s *Server) Stop() error
func (uc *UserController) GetUser(id int64) (*User, error)
```

## Bad

```go
func (this *Server) Start() error
func (self *Server) Stop() error
func (server *Server) Restart() error
```
