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
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/craftslab/lintmanifest/config"
	"github.com/craftslab/lintmanifest/gerrit"
	"github.com/craftslab/lintmanifest/gitiles"
	"github.com/craftslab/lintmanifest/manifest"
	"github.com/craftslab/lintmanifest/runtime"
	"github.com/craftslab/lintmanifest/writer"
)

var (
	app          = kingpin.New("lintmanifest", "Lint Manifest").Version(config.Version + "-build-" + config.Build)
	configFile   = app.Flag("config-file", "Config file, format: .json").Required().String()
	lintMode     = app.Flag("lint-mode", "Lint mode (async|sync)").Default("sync").String()
	lintOut      = app.Flag("lint-out", "Lint output (.json|.txt|.xlsx)").Default("out.json").String()
	manifestFile = app.Flag("manifest-file", "Manifest file").Required().String()
)

var (
	cfg = config.Config{}
)

var (
	localLayout  = "2006-01-02 15:04:05"
	remoteLayout = "Mon Jan _2 15:04:05 2006 -0700"
)

func Run() {
	var err error
	var result []interface{}

	kingpin.MustParse(app.Parse(os.Args[1:]))

	cfg, err = parseConfig(*configFile)
	if err != nil {
		log.Fatal("parse failed: ", err.Error())
	}

	if _, err = os.Stat(*lintOut); err == nil {
		log.Fatal("file exist: ", *lintOut)
	}

	m := manifest.Manifest{}

	if err = m.Load(*manifestFile); err != nil {
		log.Fatal("load failed: ", err.Error())
	}

	projects, err := m.Projects()
	if err != nil {
		log.Fatal("projects failed: ", err.Error())
	}

	log.Println("lint running...")

	if *lintMode == "async" {
		result, err = lintAsync(projects)
		if err != nil {
			log.Fatal("async failed: ", err.Error())
		}
	} else if *lintMode == "sync" {
		result, err = lintSync(projects)
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

	if err = w.Run(result, *lintOut); err != nil {
		log.Fatal("run failed: ", err.Error())
	}

	log.Println("lint completed.")
}

func parseConfig(name string) (config.Config, error) {
	var c config.Config

	fi, err := os.Open(name)
	if err != nil {
		return c, errors.Wrap(err, "open failed")
	}

	defer func() {
		_ = fi.Close()
	}()

	buf, _ := ioutil.ReadAll(fi)
	if err := json.Unmarshal(buf, &c); err != nil {
		return c, errors.Wrap(err, "unmarshal failed")
	}

	return c, nil
}

func lintAsync(projects []interface{}) ([]interface{}, error) {
	result, err := runtime.Run(routine, projects)
	if err != nil {
		return nil, errors.Wrap(err, "run failed")
	}

	return result, nil
}

func lintSync(projects []interface{}) ([]interface{}, error) {
	var result []interface{}

	for _, val := range projects {
		buf := routine(val)
		if buf != nil {
			result = append(result, buf)
		}
	}

	return result, nil
}

func routine(project interface{}) interface{} {
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

	ge := gerrit.Gerrit{
		Option: cfg.GerritConfig.QueryConfig.Option,
		Pass:   cfg.GerritConfig.Pass,
		Url:    cfg.GerritConfig.Url,
		User:   cfg.GerritConfig.User,
	}

	gi := gitiles.Gitiles{
		Pass: cfg.GitilesConfig.Pass,
		Url:  cfg.GitilesConfig.Url,
		User: cfg.GitilesConfig.User,
	}

	commitLocalDate, commitRemoteDate, err := gerritQuery(ge, revision.(string))
	if err != nil {
		return nil
	}

	if commitLocalDate == "" || commitRemoteDate == "" {
		commitLocalDate, commitRemoteDate, err = gitilesQuery(gi, name.(string), revision.(string))
		if err != nil {
			return nil
		}
	}

	headLocalDate, headRemoteDate, headHash, err := gitilesHead(gi, name.(string), upstream.(string))
	if err != nil {
		return nil
	}

	localDate, remoteDate, err := gerritQuery(ge, headHash)
	if err != nil {
		return nil
	}

	if localDate != "" && remoteDate != "" {
		headLocalDate = localDate
		headRemoteDate = remoteDate
	}

	buf := make(map[string]string)
	buf["commitLocalDate"] = commitLocalDate
	buf["commitHash"] = revision.(string)
	buf["commitRemoteDate"] = commitRemoteDate
	buf["commitUrl"] = cfg.GitilesConfig.Url + "/" + name.(string) + "/+/" + revision.(string)
	buf["headLocalDate"] = headLocalDate
	buf["headHash"] = headHash
	buf["headRemoteDate"] = headRemoteDate
	buf["headUrl"] = cfg.GitilesConfig.Url + "/" + name.(string) + "/+/" + headHash
	buf["repoBranch"] = upstream.(string)
	buf["repoName"] = name.(string)

	_, err = gi.Query(name.(string), revision.(string))
	if err != nil {
		buf["reportDetails"] = "Invalid commit in repo branch."
		buf["reportType"] = "ERROR"
		return buf
	}

	if revision.(string) != headHash {
		buf["reportDetails"] = "Commit not found as head of repo branch."
		buf["reportType"] = "WARN"
		return buf
	}

	return nil
}

func gerritQuery(g gerrit.Gerrit, commit string) (localDate, remoteDate string, err error) {
	localDate = ""
	remoteDate = ""

	if buf, err := g.Query("commit:"+commit, 0); err == nil {
		d, _ := time.Parse(localLayout, buf["submitted"].(string))
		localDate = d.Local().Format(localLayout)
		remoteDate = buf["submitted"].(string)
	}

	return localDate, remoteDate, nil
}

func gitilesHead(g gitiles.Gitiles, project, branch string) (localDate, remoteDate, hash string, err error) {
	buf, err := g.Head(project, branch)
	if err != nil {
		return "", "", "", errors.Wrap(err, "head failed")
	}

	if _, ok := buf["log"]; !ok {
		return "", "", "", errors.New("log invalid")
	}

	b := buf["log"].([]interface{})
	if len(b) == 0 {
		return "", "", "", errors.New("list invalid")
	}

	if _, ok := b[0].(map[string]interface{})["commit"]; !ok {
		return "", "", "", errors.New("commit invalid")
	}

	committer := b[0].(map[string]interface{})["committer"]

	d, _ := time.Parse(remoteLayout, committer.(map[string]interface{})["time"].(string))
	localDate = d.Local().Format("2006-01-02 15:04:05")

	remoteDate = committer.(map[string]interface{})["time"].(string)

	hash = b[0].(map[string]interface{})["commit"].(string)

	return localDate, remoteDate, hash, nil
}

func gitilesQuery(g gitiles.Gitiles, project, commit string) (localDate, remoteDate string, err error) {
	buf, err := g.Query(project, commit)
	if err != nil {
		return "", "", errors.Wrap(err, "query failed")
	}

	if _, ok := buf["committer"]; !ok {
		return "", "", errors.New("committer invalid")
	}

	committer := buf["committer"].(map[string]interface{})
	if committer == nil {
		return "", "", errors.New("committer invalid")
	}

	if _, ok := committer["time"]; !ok {
		return "", "", errors.New("time invalid")
	}

	d, _ := time.Parse(remoteLayout, committer["time"].(string))
	localDate = d.Local().Format(localLayout)

	remoteDate = committer["time"].(string)

	return localDate, remoteDate, nil
}
