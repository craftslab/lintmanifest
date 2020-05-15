// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"

	"lintmanifest/gitiles"
	"lintmanifest/manifest"
	"lintmanifest/runtime"
)

var (
	app  = kingpin.New("lintmanifest", "Lint Manifest").Author(Author).Version(Version)
	file = app.Flag("manifest-file", "Manifest file").Required().String()
	pass = app.Flag("gitiles-pass", "Gitiles password").String()
	url  = app.Flag("gitiles-url", "Gitiles location").Required().String()
	user = app.Flag("gitiles-user", "Gitiles username").String()
	out  = app.Flag("lint-out", "Lint output").Required().String()
)

func Run() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	m := manifest.Manifest{}

	if err := m.Load(*file); err != nil {
		log.Fatal("load failed: ", err.Error())
	}

	projects, err := m.Projects()
	if err != nil {
		log.Fatal("projects failed: ", err.Error())
	}

	log.Println("lint running...")

	result, err := Lint(projects)
	if err != nil {
		log.Fatal("lint failed: ", err.Error())
	}

	if err := Write(result); err != nil {
		log.Fatal("write failed: ", err.Error())
	}

	log.Println("lint completed.")
}

func Lint(projects []interface{}) ([]interface{}, error) {
	result, err := runtime.Run(Routine, projects)
	if err != nil {
		return nil, errors.Wrap(err, "run failed")
	}

	return result, nil
}

func Routine(project interface{}) interface{} {
	g := gitiles.Gitiles{}

	buf := project.(map[string]interface{})

	name, ok := buf["-name"]
	if !ok {
		return nil
	}

	path, ok := buf["-path"]
	if !ok {
		path = name
	}

	revision, ok := buf["-revision"]
	if !ok {
		return nil
	}

	_, err := g.Query(*url, *user, *pass, name.(string), revision.(string))
	if err != nil {
		return path
	}

	return nil
}

func Write(data []interface{}) error {
	if _, err := os.Stat(*out); err == nil {
		return errors.New("file alread exist")
	}

	var buf []string

	for _, val := range data {
		if val != nil && val.(string) != "" {
			buf = append(buf, val.(string))
		}
	}

	f, err := os.Create(*out)
	if err != nil {
		return errors.Wrap(err, "create failed")
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err := w.WriteString(strings.Join(buf, "\n")); err != nil {
		return errors.Wrap(err, "write failed")
	}

	w.Flush()

	return nil
}
