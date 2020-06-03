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
	pass = app.Flag("gitiles-pass", "Gitiles password").String()
	url  = app.Flag("gitiles-url", "Gitiles location").Required().String()
	user = app.Flag("gitiles-user", "Gitiles username").String()
	mode = app.Flag("lint-mode", "Lint mode (async|sync)").Default("sync").String()
	out  = app.Flag("lint-out", "Lint output").Required().String()
	file = app.Flag("manifest-file", "Manifest file").Required().String()
)

func Run() {
	var result []interface{}

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

	if *mode == "async" {
		result, err = LintAsync(projects)
	} else if *mode == "sync" {
		result, err = LintSync(projects)
	} else {
		log.Fatal("mode invalid")
	}

	if err != nil {
		log.Fatal("lint failed: ", err.Error())
	}

	if err := Write(result); err != nil {
		log.Fatal("write failed: ", err.Error())
	}

	log.Println("lint completed.")
}

func LintAsync(projects []interface{}) ([]interface{}, error) {
	result, err := runtime.Run(Routine, projects)
	if err != nil {
		return nil, errors.Wrap(err, "run failed")
	}

	return result, nil
}

func LintSync(projects []interface{}) ([]interface{}, error) {
	var result []interface{}

	for _, val := range projects {
		buf := Routine(val)
		if buf != nil {
			result = append(result, buf)
		}
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
	var buf []string

	if _, err := os.Stat(*out); err == nil {
		return errors.New("file alread exist")
	}

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
