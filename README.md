# Placekey Go SDK [![Go Reference](https://pkg.go.dev/badge/github.com/ringsaturn/pk.svg)](https://pkg.go.dev/github.com/ringsaturn/pk)

## Install

```bash
go install github.com/ringsaturn/pk
```

## Usage

```go
k, _ := pk.GeoToPlacekey(39.9289, 116.3883)
fmt.Println(k)
// Output: @6qk-v3d-brk
```

```go
lat, long, _ := pk.PlacekeyToGeo("@6qk-v3d-brk")
fmt.Printf("%.3f %.3f \n", lat, long)
// Output: 39.929 116.388
```

More usage examples: <https://pkg.go.dev/github.com/ringsaturn/pk#pkg-examples>

## CLI usage

```bash
go install github.com/ringsaturn/pk/cmd/placekey@latest
```

`ToGeo`:

```bash
placekey ToGeo -pk "@627-wbz-tjv"
```

Output:

```console
40.71237820442784 -74.0056425771711
```

`FromGeo`

```bash
placekey FromGeo -lat 40.71237820442784 -long -74.0056425771711
```

Output:

```console
@627-wbz-tjv
```

## References

- <https://www.placekey.io>
- <https://docs.placekey.io>
- <https://docs.placekey.io/Placekey_Encoding_Specification_White_Paper.pdf>
- <https://github.com/Placekey>
- <https://github.com/Placekey/placekey-py>
