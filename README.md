# civ - container image verifier

The container image verifier - `civ` checks based on constraints container
images. For this purpose is a config file required which contains the constraint
definitions. The config file must be passed as argument to `civ`. `civ` writes
the results into a separate file.

Currently is `json` and `yaml` supported. As default will be `yaml` used.
Optionally can be specified via the second arg the result file.

`civ config.yaml [ result.yaml ]`

## Constraints

### Labels

#### Exists

Verify if container image `volkerraschek/civ:latest` has label
`org.opencontainers.image.documentation` defined.

```yaml
images:
  volkerraschek/civ:latest:
    labelConstraints:
      org.opencontainers.image.documentation:
        exists: true
```

#### Compare Semantic Versioning

Verify, if the container image `volkerraschek/civ:latest` has label
`org.opencontainers.image.version` defined and has a greater version than
`2.5.7`.

```yaml
images:
  volkerraschek/civ:latest:
    labelConstraints:
      org.opencontainers.image.version:
        compareSemver:
          greaterThan: 2.5.7
```

Alternatively, can `lessThan` and `equal` be used. For example to define a range
of `2.5.7~2.8.4` with `lessThan` and `greaterThan`.

```yaml
images:
  volkerraschek/civ:latest:
    labelConstraints:
      org.opencontainers.image.version:
        compareSemver:
          greaterThan: 2.5.7
          lessThan: 2.8.4
```

#### Compare String

Verify, if the container image `volkerraschek/civ:latest` has label
`org.opencontainers.image.documentation` defined and the value starts with
`https://` and ends with `README.md`.

```yaml
images:
  volkerraschek/civ:latest:
    labelConstraints:
      org.opencontainers.image.documentation:
        compareString:
          hasPrefix: "https://"
          hasSuffix: "README.md"
```

Alternatively, can be `equal` used to compare the value of a label with a
expected value.

#### Count labels with corresponding prefix, suffix or match pattern

No more than 3 labels with the prefix `org.opencontainers` and exactly one
labels with the suffix `version` may be defined for the image
`volkerraschek/civ:latest`.

```yaml
images:
  volkerraschek/civ:latest:
    labelConstraints:
      org.opencontainers%:
        count:
          lowerThan: 4
      %version:
        count:
          equal: 1
```

The functions `lessThan` and `equal` are also available as constraints.
