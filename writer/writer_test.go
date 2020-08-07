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
		"branch":  "master",
		"commit":  "d297a0797d1c0955ceb923ca38ce0dfb13236337",
		"details": "Details",
		"repo":    "example",
		"type":    "ERROR",
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
