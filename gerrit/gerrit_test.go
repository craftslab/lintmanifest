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

package gerrit

import (
	"testing"
)

const (
	Pass = ""
	Url  = "https://android-review.googlesource.com/"
	User = ""
)

func TestGet(t *testing.T) {
	g := Gerrit{
		Option: "CURRENT_REVISION",
		Pass:   Pass,
		Url:    Url,
		User:   User,
	}

	if _, err := g.Get(-1); err == nil {
		t.Error("FAIL:", err)
	}

	if _, err := g.Get(502075); err != nil {
		t.Error("FAIL:", err)
	}
}

func TestQuery(t *testing.T) {
	g := Gerrit{
		Option: "CURRENT_REVISION",
		Pass:   Pass,
		Url:    Url,
		User:   User,
	}

	if _, err := g.Query("commit:-1", 0); err == nil {
		t.Error("FAIL:", err)
	}

	if _, err := g.Query("commit:b6356a0", 0); err != nil {
		t.Error("FAIL:", err)
	}
}
