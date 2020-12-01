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

package writer

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	xlsx "github.com/tealeg/xlsx"
)

type Writer struct {
	data []interface{}
	name string
}

func (w Writer) Run(data []interface{}, name string) error {
	var err error

	w.data = data
	w.name = name

	if strings.HasSuffix(w.name, ".json") {
		err = w.writeJson()
	} else if strings.HasSuffix(w.name, ".txt") {
		err = w.writeTxt()
	} else if strings.HasSuffix(w.name, ".xlsx") {
		err = w.writeXlsx()
	} else {
		return errors.New("invalid suffix")
	}

	return err
}

func (w Writer) writeJson() error {
	buf := make(map[string][]interface{})
	buf["lintmanifest"] = w.data

	b, err := json.Marshal(buf)
	if err != nil {
		return errors.Wrap(err, "marshal failed")
	}

	if err := ioutil.WriteFile(w.name, b, 0600); err != nil {
		return errors.Wrap(err, "write failed")
	}

	return nil
}

func (w Writer) writeTxt() error {
	var buf []string

	head := "REPORTTYPE,REPONAME,REPOBRANCH,HEADDATE,HEADHASH,HEADURL,COMMITDATE,COMMITHASH,COMMITURL,REPORTDETAILS"
	buf = append(buf, head)

	for _, val := range w.data {
		if val != nil {
			v := val.(map[string]string)
			buf = append(buf, v["reportType"]+","+v["repoName"]+","+v["repoBranch"]+
				v["headDate"]+","+v["headHash"]+","+v["headUrl"]+","+
				v["commitDate"]+","+v["commitHash"]+","+v["commitUrl"]+",'"+v["reportDetails"]+"'")
		}
	}

	f, err := os.Create(w.name)
	if err != nil {
		return errors.Wrap(err, "create failed")
	}
	defer func() {
		_ = f.Close()
	}()

	b := bufio.NewWriter(f)
	if _, err := b.WriteString(strings.Join(buf, "\n")); err != nil {
		return errors.Wrap(err, "write failed")
	}
	defer func() {
		_ = b.Flush()
	}()

	return nil
}

func (w Writer) writeXlsx() error {
	type R struct {
		ReportType    string
		RepoName      string
		RepoBranch    string
		HeadDate      string
		HeadHash      string
		HeadUrl       string
		CommitDate    string
		CommitHash    string
		CommitUrl     string
		ReportDetails string
	}

	wb := xlsx.NewFile()

	sh, err := wb.AddSheet("lintmanifest")
	if err != nil {
		return errors.Wrap(err, "add failed")
	}

	r := R{
		"REPORTTYPE",
		"REPONAME",
		"REPOBRANCH",
		"HEADDATE",
		"HEADHASH",
		"HEADURL",
		"COMMITDATE",
		"COMMITHASH",
		"COMMITURL",
		"REPORTDETAILS",
	}
	row := sh.AddRow()
	row.WriteStruct(&r, -1)

	for _, val := range w.data {
		if val != nil {
			v := val.(map[string]string)
			r := R{
				v["reportType"],
				v["repoName"],
				v["repoBranch"],
				v["headDate"],
				v["headHash"],
				v["headUrl"],
				v["commitDate"],
				v["commitHash"],
				v["commitUrl"],
				v["reportDetails"],
			}
			row := sh.AddRow()
			row.WriteStruct(&r, -1)
		}
	}

	if err := wb.Save(w.name); err != nil {
		return errors.Wrap(err, "add failed")
	}

	return nil
}
