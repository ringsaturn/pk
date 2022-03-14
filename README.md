# Unofficial Placekey Go SDK [![ci](https://github.com/ringsaturn/pk/actions/workflows/ci.yml/badge.svg)](https://github.com/ringsaturn/pk/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/ringsaturn/pk.svg)](https://pkg.go.dev/github.com/ringsaturn/pk)

```bash
go install github.com/ringsaturn/pk
```

Cli usage:

```bash
go install github.com/ringsaturn/pk/cmd/placekey@latest

# ToGeo
placekey ToGeo -pk "@627-wbz-tjv"
40.71237820442784 -74.0056425771711

# FromGeo
placekey FromGeo -lat 40.71237820442784 -long -74.0056425771711
@627-wbz-tjv
```

References:

- <https://www.placekey.io>
- <https://docs.placekey.io>
- <https://docs.placekey.io/Placekey_Encoding_Specification_White_Paper.pdf>
- <https://github.com/Placekey/placekey-py>
