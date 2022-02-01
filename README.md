# WIP: Placekey's Go impl [![ci](https://github.com/ringsaturn/pk/actions/workflows/ci.yml/badge.svg)](https://github.com/ringsaturn/pk/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/ringsaturn/pk.svg)](https://pkg.go.dev/github.com/ringsaturn/pk)

```go
package main

import (
	"fmt"

	"github.com/ringsaturn/pk"
)

func main() {
	fmt.Println(pk.GeoToPlacekey(39.9289, 116.3883))
	// Output: @6qk-v3d-brk
}
```

References:

- <https://github.com/Placekey/placekey-py>
- <https://docs.placekey.io/Placekey_Encoding_Specification_White_Paper.pdf>
