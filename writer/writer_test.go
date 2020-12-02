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
	"os"
	"testing"
)

var (
	content = map[string]string{
		"reportType":       "ERROR",
		"repoName":         "project",
		"repoBranch":       "master",
		"headLocalDate":    "2020-01-01 00:00:00",
		"headRemoteDate":   "2020-01-01 00:00:00",
		"headHash":         "8cf3e5471db04da274965a8e5c0dc3d465f08c5f",
		"headUrl":          "http://127.0.0.1",
		"commitLocalDate":  "2020-01-01 00:00:00",
		"commitRemoteDate": "2020-01-01 00:00:00",
		"commitHash":       "8cf3e5471db04da274965a8e5c0dc3d465f08c5f",
		"commitUrl":        "http://127.0.0.1",
		"reportDetails":    "details",
	}
)

func TestWriteJson(t *testing.T) {
	w := Writer{}

	w.data = make([]interface{}, 1)
	w.data[0] = content
	w.name = "out.json"

	err := w.writeJson()
	defer func(name string) { _ = os.Remove(name) }(w.name)
	if err != nil {
		t.Error("FAIL", err)
	}
}

func TestWriteTxt(t *testing.T) {
	w := Writer{}

	w.data = make([]interface{}, 1)
	w.data[0] = content
	w.name = "out.txt"

	err := w.writeTxt()
	defer func(name string) { _ = os.Remove(name) }(w.name)
	if err != nil {
		t.Error("FAIL", err)
	}
}

func TestWriteXlsx(t *testing.T) {
	w := Writer{}

	w.data = make([]interface{}, 1)
	w.data[0] = content
	w.name = "out.xlsx"

	err := w.writeXlsx()
	defer func(name string) { _ = os.Remove(name) }(w.name)
	if err != nil {
		t.Error("FAIL", err)
	}
}
