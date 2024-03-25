package render_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/jimeh/go-render"
)

func ExampleRender_json() {
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

	// Render the object to JSON.
	buf := &bytes.Buffer{}
	err := render.Pretty(buf, "json", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// {
	//   "current": "1.2.2",
	//   "versions": [
	//     {
	//       "version": "1.2.2",
	//       "latest": true,
	//       "stable": true
	//     },
	//     {
	//       "version": "1.2.1",
	//       "latest": false,
	//       "stable": true
	//     },
	//     {
	//       "version": "1.2.0",
	//       "latest": false,
	//       "stable": true
	//     },
	//     {
	//       "version": "1.2.0-rc.0",
	//       "latest": false,
	//       "stable": false
	//     },
	//     {
	//       "version": "1.1.0",
	//       "latest": false,
	//       "stable": true
	//     }
	//   ]
	// }
}

func ExampleRender_yaml() {
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

	// Render the object to YAML.
	buf := &bytes.Buffer{}
	err := render.Pretty(buf, "yaml", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

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

func ExampleRender_xml() {
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

	// Render the object to XML.
	buf := &bytes.Buffer{}
	err := renderer.Pretty(buf, "xml", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// <versions-list>
	//   <current>1.2.2</current>
	//   <version latest="true" stable="true">1.2.2</version>
	//   <version latest="false" stable="true">1.2.1</version>
	//   <version latest="false" stable="true">1.2.0</version>
	//   <version latest="false" stable="false">1.2.0-rc.0</version>
	//   <version latest="false" stable="true">1.1.0</version>
	// </versions-list>
}

func ExampleRender_textFromByteSlice() {
	data := []byte("Hello, World!1")

	// Render the object to XML.
	buf := &bytes.Buffer{}
	err := render.Pretty(buf, "text", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// Hello, World!1
}

func ExampleRender_textFromString() {
	data := "Hello, World!"

	// Render the object to XML.
	buf := &bytes.Buffer{}
	err := render.Pretty(buf, "text", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// Hello, World!
}

func ExampleRender_textFromIOReader() {
	var data io.Reader = strings.NewReader("Hello, World!!!1")

	// Render the object to XML.
	buf := &bytes.Buffer{}
	err := render.Pretty(buf, "text", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// Hello, World!!!1
}

func ExampleRender_textFromWriterTo() {
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

	// Render the object to text.
	buf := &bytes.Buffer{}
	err := render.Pretty(buf, "text", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// 1.2.1 (stable: true, latest: false)
}

type Version struct {
	Version string `json:"version" yaml:"version" xml:",chardata"`
	Latest  bool   `json:"latest" yaml:"latest" xml:"latest,attr"`
	Stable  bool   `json:"stable" yaml:"stable" xml:"stable,attr"`
}

func (v *Version) WriteTo(w io.Writer) (int64, error) {
	s := fmt.Sprintf(
		"%s (stable: %t, latest: %t)", v.Version, v.Stable, v.Latest,
	)
	n, err := w.Write([]byte(s))

	return int64(n), err
}

type OutputList struct {
	Current  string    `json:"current" yaml:"current" xml:"current"`
	Versions []Version `json:"versions" yaml:"versions" xml:"version"`

	XMLName xml.Name `json:"-" yaml:"-" xml:"versions-list"`
}

func (ol *OutputList) String() string {
	buf := &strings.Builder{}

	for _, ver := range ol.Versions {
		if ol.Current == ver.Version {
			buf.WriteString("* ")
		} else {
			buf.WriteString("  ")
		}

		buf.WriteString(ver.Version)
		if !ver.Stable {
			buf.WriteString(" (pre-release)")
		}
		if ver.Latest {
			buf.WriteString(" (latest)")
		}

		buf.WriteByte('\n')
	}

	return buf.String()
}

func ExampleRender_textFromStringer() {
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

	// Render the object to XML.
	buf := &bytes.Buffer{}
	err := render.Pretty(buf, "text", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// * 1.2.2 (latest)
	//   1.2.1
	//   1.2.0
	//   1.2.0-rc.0 (pre-release)
	//   1.1.0
}
