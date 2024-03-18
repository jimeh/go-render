package render_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/jimeh/go-render"
)

type User struct {
	Name  string   `json:"name" yaml:"name" xml:"name"`
	Age   int      `json:"age" yaml:"age" xml:"age"`
	Roles []*Role  `json:"roles" yaml:"roles" xml:"roles"`
	Tags  []string `json:"tags" yaml:"tags" xml:"tags"`

	XMLName xml.Name `json:"-" yaml:"-" xml:"user"`
}

func (u *User) String() string {
	return fmt.Sprintf(
		"%s (%d): %s",
		u.Name, u.Age, strings.Join(u.Tags, ", "),
	)
}

type Role struct {
	Name string `json:"name" yaml:"name" xml:"name"`
	Icon string `json:"icon" yaml:"icon" xml:"icon"`
}

func (r *Role) WriteTo(w io.Writer) (int64, error) {
	s := fmt.Sprintf("%s (%s)", r.Name, r.Icon)
	n, err := w.Write([]byte(s))

	return int64(n), err
}

func ExampleRender_json() {
	data := &User{
		Name: "John Doe",
		Age:  30,
		Roles: []*Role{
			{Name: "admin", Icon: "shield"},
			{Name: "developer", Icon: "keyboard"},
		},
		Tags: []string{"golang", "json", "yaml", "toml"},
	}

	// Render the object to JSON.
	buf := &bytes.Buffer{}
	err := render.Render(buf, "json", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// {
	//   "name": "John Doe",
	//   "age": 30,
	//   "roles": [
	//     {
	//       "name": "admin",
	//       "icon": "shield"
	//     },
	//     {
	//       "name": "developer",
	//       "icon": "keyboard"
	//     }
	//   ],
	//   "tags": [
	//     "golang",
	//     "json",
	//     "yaml",
	//     "toml"
	//   ]
	// }
}

func ExampleRender_yaml() {
	data := &User{
		Name: "John Doe",
		Age:  30,
		Roles: []*Role{
			{Name: "admin", Icon: "shield"},
			{Name: "developer", Icon: "keyboard"},
		},
		Tags: []string{"golang", "json", "yaml", "toml"},
	}

	// Render the object to YAML.
	buf := &bytes.Buffer{}
	err := render.Render(buf, "yaml", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// name: John Doe
	// age: 30
	// roles:
	//   - name: admin
	//     icon: shield
	//   - name: developer
	//     icon: keyboard
	// tags:
	//   - golang
	//   - json
	//   - yaml
	//   - toml
}

func ExampleRender_xml() {
	data := &User{
		Name: "John Doe",
		Age:  30,
		Roles: []*Role{
			{Name: "admin", Icon: "shield"},
			{Name: "developer", Icon: "keyboard"},
		},
		Tags: []string{"golang", "json", "yaml", "toml"},
	}

	// Render the object to XML.
	buf := &bytes.Buffer{}
	err := render.Render(buf, "xml", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// <user>
	//   <name>John Doe</name>
	//   <age>30</age>
	//   <roles>
	//     <name>admin</name>
	//     <icon>shield</icon>
	//   </roles>
	//   <roles>
	//     <name>developer</name>
	//     <icon>keyboard</icon>
	//   </roles>
	//   <tags>golang</tags>
	//   <tags>json</tags>
	//   <tags>yaml</tags>
	//   <tags>toml</tags>
	// </user>
}

func ExampleRender_textViaStringer() {
	// The User struct has a String method which returns a string representation
	// of a user:
	//
	//  func (u *User) String() string {
	//		return fmt.Sprintf(
	//			"%s (%d): %s",
	//			u.Name, u.Age, strings.Join(u.Tags, ", "),
	//		)
	//  }

	data := &User{
		Name: "John Doe",
		Age:  30,
		Roles: []*Role{
			{Name: "admin", Icon: "shield"},
			{Name: "developer", Icon: "keyboard"},
		},
		Tags: []string{"golang", "json", "yaml", "toml"},
	}

	// Render the object to XML.
	buf := &bytes.Buffer{}
	err := render.Render(buf, "text", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// John Doe (30): golang, json, yaml, toml
}

func ExampleRender_textViaWriterTo() {
	// The Role struct has a WriteTo method which writes a string representation
	// of a role to an io.Writer:
	//
	//	func (r *Role) WriteTo(w io.Writer) (int64, error) {
	//		s := fmt.Sprintf("%s (%s)", r.Name, r.Icon)
	//		n, err := w.Write([]byte(s))
	//
	//		return int64(n), err
	//	}

	data := &Role{Name: "admin", Icon: "shield"}

	// Render the object to XML.
	buf := &bytes.Buffer{}
	err := render.Render(buf, "text", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// admin (shield)
}
