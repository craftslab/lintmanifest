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
	"testing"

	"github.com/craftslab/lintmanifest/gerrit"
	"github.com/craftslab/lintmanifest/gitiles"
)

func TestParseConfig(t *testing.T) {
	if _, err := parseConfig("foo.json"); err == nil {
		t.Error("FAIL:", err)
	}

	if _, err := parseConfig("../config/config.json"); err != nil {
		t.Error("FAIL:", err)
	}
}

func TestLintAsync(t *testing.T) {
	projects := make([]interface{}, 1)

	projects[0] = map[string]interface{}{
		"-groups":   "pdk,tradefed",
		"-name":     "platform/build/soong",
		"-path":     "build/soong",
		"-revision": "8cf3e5471db04da274965a8e5c0dc3d465f08c5f",
		"-upstream": "master",
	}

	if _, err := lintAsync(projects); err != nil {
		t.Error("FAIL:", err)
	}
}

func TestLintSync(t *testing.T) {
	projects := make([]interface{}, 1)

	projects[0] = map[string]interface{}{
		"-groups":   "pdk,tradefed",
		"-name":     "platform/build/soong",
		"-path":     "build/soong",
		"-revision": "8cf3e5471db04da274965a8e5c0dc3d465f08c5f",
		"-upstream": "master",
	}

	if _, err := lintSync(projects); err != nil {
		t.Error("FAIL:", err)
	}
}

func TestGerritQuery(t *testing.T) {
	g := gerrit.Gerrit{
		Option: "CURRENT_REVISION",
		Pass:   "",
		Url:    "https://android-review.googlesource.com",
		User:   "",
	}

	if _, err := gerritQuery(g, "1514894"); err != nil {
		t.Error("FAIL")
	}

	g = gerrit.Gerrit{
		Option: "CURRENT_REVISION",
		Pass:   "",
		Url:    "https://android-review.googlesource.com",
		User:   "",
	}

	if _, err := gerritQuery(g, "ffbd673b27bd0f60008360e3b8cd34c85886a0a4"); err != nil {
		t.Error("FAIL")
	}
}

func TestGitilesHead(t *testing.T) {
	g := gitiles.Gitiles{
		Pass: "",
		Url:  "https://android.googlesource.com",
		User: "",
	}

	if _, _, err := gitilesHead(g, "platform/build/soong", "foo"); err == nil {
		t.Error("FAIL")
	}

	if _, _, err := gitilesHead(g, "platform/build/soong", "master"); err != nil {
		t.Error("FAIL")
	}
}

func TestGitilesQuery(t *testing.T) {
	g := gitiles.Gitiles{
		Pass: "",
		Url:  "https://android.googlesource.com",
		User: "",
	}

	if _, err := gitilesQuery(g, "platform/build/soong", "foo"); err == nil {
		t.Error("FAIL")
	}

	if _, err := gitilesQuery(g, "platform/build/soong", "8cf3e5471db04da274965a8e5c0dc3d465f08c5f"); err != nil {
		t.Error("FAIL")
	}
}
