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

package gitiles

import (
	"testing"
)

const (
	Pass    = ""
	Project = "platform/build/soong"
	Url     = "https://android.googlesource.com"
	User    = ""
)

func TestHead(t *testing.T) {
	g := Gitiles{
		Pass: Pass,
		Url:  Url,
		User: User,
	}

	if _, err := g.Head(Project, "master"); err != nil {
		t.Error("FAIL:", err)
	}
}

func TestQuery(t *testing.T) {
	g := Gitiles{
		Pass: Pass,
		Url:  Url,
		User: User,
	}

	if _, err := g.Query(Project, "8cf3e5471db04da274965a8e5c0dc3d465f08c5f"); err != nil {
		t.Error("FAIL:", err)
	}

	if _, err := g.Query(Project, "invalid"); err == nil {
		t.Error("FAIL:", err)
	}
}
