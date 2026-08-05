# mailofly-go

Official **Go** client for the [Mailofly REST API](https://docs.mailofly.com/api).

Requires **Go 1.22+**.

> **Source of truth:** developed in the [mailofly monorepo](https://github.com/godstark82/mailofly) under `packages/go`. This public repo is mirrored automatically on change.

## Install

```bash
go get github.com/teamredevs/mailofly-go@latest
```

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/teamredevs/mailofly-go"
)

func main() {
	client, err := mailofly.New(mailofly.Options{
		APIKey: os.Getenv("MAILOFLY_API_KEY"),
	})
	if err != nil {
		panic(err)
	}

	result, err := client.Compose.Send(map[string]any{
		"account_key": "acc_…",
		"subject":     "Hello",
		"body":        "<p>Hi from Mailofly</p>",
		"recipients":  map[string]any{"emails": []string{"you@example.com"}},
	})
	if err != nil {
		if apiErr, ok := err.(*mailofly.Error); ok {
			fmt.Println(apiErr.Status, apiErr.Err, apiErr.DetailMessage)
		}
		panic(err)
	}
	fmt.Println(result)
}
```

## Docs

- [Go guide](https://docs.mailofly.com/sdks/go)
- [API reference](https://docs.mailofly.com/api)

## Releasing

1. Bump `VERSION` (and this changelog) in the **monorepo** PR.
2. Merge to `main`/`master` → GitHub Action syncs this folder to `teamredevs/mailofly-go`.
3. Publish workflow creates a `v*` git tag for `go get`.

## License

MIT
