package render_test

import (
	"encoding/xml"
	"io"
	"os"
	"strings"

	"github.com/jimeh/go-render"
)

//nolint:lll
func ExampleCompact_json() {
	type Version struct {
		Version string `json:"version" yaml:"version" xml:",chardata"`
		Latest  bool   `json:"latest" yaml:"latest" xml:"latest,attr"`
		Stable  bool   `json:"stable" yaml:"stable" xml:"stable,attr"`
	}

	type OutputList struct {
		Current  string    `json:"current" yaml:"current" xml:"current"`
		Versions []Version `json:"versions" yaml:"versions" xml:"version"`

		XMLName xml.Name `json:"-" yaml:"-" xml:"versions-list"`
	}

	data := &OutputList{
		Current: "1.2.2",
		Versions: []Version{
			{Version: "1.2.2", Stable: true, Latest: true},
			{Version: "1.2.1", Stable: true},
			{Version: "1.2.0", Stable: true},
			{Version: "1.2.0-rc.0", Stable: false},
			{Version: "1.1.0", Stable: true},
		},
	}

	err := render.Compact(os.Stdout, "json", data)
	if err != nil {
		panic(err)
	}

	// Output:
	// {"current":"1.2.2","versions":[{"version":"1.2.2","latest":true,"stable":true},{"version":"1.2.1","latest":false,"stable":true},{"version":"1.2.0","latest":false,"stable":true},{"version":"1.2.0-rc.0","latest":false,"stable":false},{"version":"1.1.0","latest":false,"stable":true}]}
}

func ExampleCompact_yaml() {
	type Version struct {
		Version string `json:"version" yaml:"version" xml:",chardata"`
		Latest  bool   `json:"latest" yaml:"latest" xml:"latest,attr"`
		Stable  bool   `json:"stable" yaml:"stable" xml:"stable,attr"`
	}

	type OutputList struct {
		Current  string    `json:"current" yaml:"current" xml:"current"`
		Versions []Version `json:"versions" yaml:"versions" xml:"version"`

		XMLName xml.Name `json:"-" yaml:"-" xml:"versions-list"`
	}

	data := &OutputList{
		Current: "1.2.2",
		Versions: []Version{
			{Version: "1.2.2", Stable: true, Latest: true},
			{Version: "1.2.1", Stable: true},
			{Version: "1.2.0", Stable: true},
			{Version: "1.2.0-rc.0", Stable: false},
			{Version: "1.1.0", Stable: true},
		},
	}

	err := render.Compact(os.Stdout, "yaml", data)
	if err != nil {
		panic(err)
	}

	// Output:
	// current: 1.2.2
	// versions:
	//   - version: 1.2.2
	//     latest: true
	//     stable: true
	//   - version: 1.2.1
	//     latest: false
	//     stable: true
	//   - version: 1.2.0
	//     latest: false
	//     stable: true
	//   - version: 1.2.0-rc.0
	//     latest: false
	//     stable: false
	//   - version: 1.1.0
	//     latest: false
	//     stable: true
}

//nolint:lll
func ExampleCompact_xml() {
	type Version struct {
		Version string `json:"version" yaml:"version" xml:",chardata"`
		Latest  bool   `json:"latest" yaml:"latest" xml:"latest,attr"`
		Stable  bool   `json:"stable" yaml:"stable" xml:"stable,attr"`
	}

	type OutputList struct {
		Current  string    `json:"current" yaml:"current" xml:"current"`
		Versions []Version `json:"versions" yaml:"versions" xml:"version"`

		XMLName xml.Name `json:"-" yaml:"-" xml:"versions-list"`
	}

	data := &OutputList{
		Current: "1.2.2",
		Versions: []Version{
			{Version: "1.2.2", Stable: true, Latest: true},
			{Version: "1.2.1", Stable: true},
			{Version: "1.2.0", Stable: true},
			{Version: "1.2.0-rc.0", Stable: false},
			{Version: "1.1.0", Stable: true},
		},
	}

	// Create a new renderer that supports XML in addition to the default JSON,
	// Text, and YAML formats.
	renderer := render.NewWith("json", "text", "xml", "yaml")

	err := renderer.Compact(os.Stdout, "xml", data)
	if err != nil {
		panic(err)
	}

	// Output:
	// <versions-list><current>1.2.2</current><version latest="true" stable="true">1.2.2</version><version latest="false" stable="true">1.2.1</version><version latest="false" stable="true">1.2.0</version><version latest="false" stable="false">1.2.0-rc.0</version><version latest="false" stable="true">1.1.0</version></versions-list>
}

func ExampleCompact_textFromByteSlice() {
	data := []byte("Hello, World!1")

	err := render.Compact(os.Stdout, "text", data)
	if err != nil {
		panic(err)
	}

	// Output:
	// Hello, World!1
}

func ExampleCompact_textFromString() {
	data := "Hello, World!"

	err := render.Compact(os.Stdout, "text", data)
	if err != nil {
		panic(err)
	}

	// Output:
	// Hello, World!
}

func ExampleCompact_textFromIOReader() {
	var data io.Reader = strings.NewReader("Hello, World!!!1")

	err := render.Compact(os.Stdout, "text", data)
	if err != nil {
		panic(err)
	}

	// Output:
	// Hello, World!!!1
}

func ExampleCompact_textFromWriterTo() {
	// The Version struct has a WriteTo method which writes a string
	// representation of a version to an io.Writer:
	//
	//	func (v *Version) WriteTo(w io.Writer) (int64, error) {
	//		s := fmt.Sprintf(
	//			"%s (stable: %t, latest: %t)", v.Version, v.Stable, v.Latest,
	//		)
	//		n, err := w.Write([]byte(s))
	//
	//		return int64(n), err
	//	}

	data := &Version{Version: "1.2.1", Stable: true, Latest: false}

	err := render.Compact(os.Stdout, "text", data)
	if err != nil {
		panic(err)
	}

	// Output:
	// 1.2.1 (stable: true, latest: false)
}

func ExampleCompact_textFromStringer() {
	// The User struct has a String method which returns a string representation
	// of a user:
	//
	//	func (ol *OutputList) String() string {
	//		buf := &strings.Builder{}
	//
	//		for _, ver := range ol.Versions {
	//			if ol.Current == ver.Version {
	//				buf.WriteString("* ")
	//			} else {
	//				buf.WriteString("  ")
	//			}
	//
	//			buf.WriteString(ver.Version)
	//			if !ver.Stable {
	//				buf.WriteString(" (pre-release)")
	//			}
	//			if ver.Latest {
	//				buf.WriteString(" (latest)")
	//			}
	//
	//			buf.WriteByte('\n')
	//		}
	//
	//		return buf.String()
	//	}

	data := &OutputList{
		Current: "1.2.2",
		Versions: []Version{
			{Version: "1.2.2", Stable: true, Latest: true},
			{Version: "1.2.1", Stable: true},
			{Version: "1.2.0", Stable: true},
			{Version: "1.2.0-rc.0", Stable: false},
			{Version: "1.1.0", Stable: true},
		},
	}

	err := render.Compact(os.Stdout, "text", data)
	if err != nil {
		panic(err)
	}

	// Output:
	// * 1.2.2 (latest)
	//   1.2.1
	//   1.2.0
	//   1.2.0-rc.0 (pre-release)
	//   1.1.0
}
