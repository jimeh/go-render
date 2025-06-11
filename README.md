<h1 align="center">
  go-render
</h1>

<p align="center">
  <strong>
    A simple and flexible solution to render a value to a <code>io.Writer</code>
    using different formats based on a format string argument.
  </strong>
</p>

<p align="center">
  <a href="https://github.com/jimeh/go-render/releases"><img src="https://img.shields.io/github/v/tag/jimeh/go-render?label=release" alt="GitHub tag (latest SemVer)"></a>
  <a href="https://pkg.go.dev/github.com/jimeh/go-render"><img src="https://img.shields.io/badge/%E2%80%8B-reference-387b97.svg?logo=go&logoColor=white" alt="Go Reference"></a>
  <a href="https://github.com/jimeh/go-render/issues"><img src="https://img.shields.io/github/issues-raw/jimeh/go-render.svg?style=flat&logo=github&logoColor=white" alt="GitHub issues"></a>
  <a href="https://github.com/jimeh/go-render/pulls"><img src="https://img.shields.io/github/issues-pr-raw/jimeh/go-render.svg?style=flat&logo=github&logoColor=white" alt="GitHub pull requests"></a>
  <a href="https://codeclimate.com/github/jimeh/go-render"><img src="https://img.shields.io/codeclimate/coverage/jimeh/go-render.svg?logo=code%20climate" alt="Coverage"></a>
  <a href="https://github.com/jimeh/go-render/blob/main/LICENSE"><img src="https://img.shields.io/github/license/jimeh/go-render.svg?style=flat" alt="License Status"></a>
</p>

Designed around using a custom type/struct to render your output. Thanks to Go's
marshaling interfaces, you get JSON, YAML, and XML support almost for free.
While plain text output is supported by the type implementing `io.Reader`,
`io.WriterTo`, `fmt.Stringer`, and `error` interfaces, or by simply being a type
which can easily be type cast to a byte slice.

Originally intended to easily implement CLI tools which can output their data as
plain text, as well as JSON/YAML with a simple switch of a format string. But it
can just as easily render to any `io.Writer`.

The package is designed to be flexible and extensible with a sensible set of
defaults accessible via package level functions. You can create your own
`Renderer` for custom formats, or create new handlers that support custom
formats.

## Import

```go
import "github.com/jimeh/go-render"
```

## Usage

Basic usage to render a value to various formats into a `io.Writer`:

```go
version := &Version{Version: "1.2.1", Stable: true, Latest: false}

err = render.Pretty(w, "text", version)
// 1.2.1 (stable: true, latest: false)

err = render.Pretty(w, "json", version)
// {
//   "version": "1.2.1",
//   "latest": false,
//   "stable": true
// }

err = render.Compact(w, "json", version)
// {"version":"1.2.1","latest":false,"stable":true}

err = render.Pretty(w, "yaml", version)
// version: 1.2.1
// latest: false
// stable: true

err = render.Pretty(w, "xml", version)
// <version latest="false" stable="true">1.2.1</version>
```

The above assumes the following `Version` struct:

```go
type Version struct {
    Version string `json:"version" yaml:"version" xml:",chardata"`
    Latest  bool   `json:"latest"  yaml:"latest"  xml:"latest,attr"`
    Stable  bool   `json:"stable"  yaml:"stable"  xml:"stable,attr"`
}

func (v *Version) String() string {
    return fmt.Sprintf(
        "%s (stable: %t, latest: %t)", v.Version, v.Stable, v.Latest,
    )
}
```

## Documentation

Please see the
[Go Reference](https://pkg.go.dev/github.com/jimeh/go-render#section-documentation)
for documentation and further examples.

## License

[MIT](https://github.com/jimeh/go-render/blob/main/LICENSE)
