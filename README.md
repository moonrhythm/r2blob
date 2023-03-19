# r2blob

## Installation

```bash
go get -u github.com/moonrhythm/r2blob
```

## Usage

```go
package main

import (
    "gocloud.dev/blob"

    _ "github.com/moonrhythm/r2blob"
)

func main() {
    bucket, err := blob.OpenBucket(context.Background(), "r2://bucket-name?account=xxx&access_key_id=xxx&access_key_secret=xxx")
    if err != nil {
        panic(err)
    }
    defer bucket.Close()

    // ...
}
