# go-yaml-trimer

Tiny tool to trim YAML properties.

```yaml
# input.yaml
T:
  ID:
    type: integer
    format: int32
    x-will-be-removed: true
  Tag:
    type: integer
    format: int32
    x-will-be-removed: false
```

After running `yaml-trimer -pattern x-will-* -output output.yaml input.yaml`

```yaml
# output.yaml
T:
  ID:
    type: integer
    format: int32
  Tag:
    type: integer
    format: int32
```

## installation

    go get github.com/suzuken/go-yaml-trimer

## LICENSE

MIT

## Author

suzuken
