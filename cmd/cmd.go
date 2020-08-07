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
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"

	"lintmanifest/gitiles"
	"lintmanifest/manifest"
	"lintmanifest/runtime"
	"lintmanifest/writer"
)

var (
	app  = kingpin.New("lintmanifest", "Lint Manifest").Author(Author).Version(Version)
	pass = app.Flag("gitiles-pass", "Gitiles password").String()
	url  = app.Flag("gitiles-url", "Gitiles location").Required().String()
	user = app.Flag("gitiles-user", "Gitiles username").String()
	mode = app.Flag("lint-mode", "Lint mode (async|sync)").Default("sync").String()
	out  = app.Flag("lint-out", "Lint output (.json|.txt|.xlsx)").Default("out.txt").String()
	file = app.Flag("manifest-file", "Manifest file").Required().String()
)

func Run() {
	var result []interface{}

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if _, err := os.Stat(*out); err == nil {
		log.Fatal("file exist: ", *out)
	}

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
		if err != nil {
			log.Fatal("async failed: ", err.Error())
		}
	} else if *mode == "sync" {
		result, err = LintSync(projects)
		if err != nil {
			log.Fatal("sync failed: ", err.Error())
		}
	} else {
		log.Fatal("mode invalid")
	}

	if len(result) != 0 {
		log.Println("error/warning found!")
	} else {
		log.Println("no error/warning.")
		log.Println("lint completed.")
		return
	}

	w := writer.Writer{}

	if err := w.Run(result, *out); err != nil {
		log.Fatal("run failed: ", err.Error())
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
	p := project.(map[string]interface{})

	name, ok := p["-name"]
	if !ok {
		return nil
	}

	revision, ok := p["-revision"]
	if !ok {
		return nil
	}

	upstream, ok := p["-upstream"]
	if !ok {
		return nil
	}

	buf := make(map[string]string)

	buf["branch"] = upstream.(string)
	buf["commit"] = revision.(string)
	buf["details"] = ""
	buf["repo"] = name.(string)
	buf["type"] = ""

	g := gitiles.Gitiles{}

	_, err := g.Query(*url, *user, *pass, buf["repo"], buf["commit"])
	if err != nil {
		buf["details"] = "Commit is invalid in branch of repo."
		buf["type"] = "ERROR"
		return buf
	}

	head, err := g.Head(*url, *user, *pass, buf["repo"], buf["branch"])
	if err != nil {
		return nil
	}

	if buf["commit"] != head {
		buf["details"] = "Commit is not the head in branch of repo."
		buf["type"] = "WARN"
		return buf
	}

	return nil
}
