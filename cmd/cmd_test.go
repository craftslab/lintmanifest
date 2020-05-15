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
)

func TestLint(t *testing.T) {
	projects := make([]interface{}, 1)

	projects[0] = map[string]interface{}{
		"-groups":   "pdk,tradefed",
		"-name":     "platform/build/soong",
		"-path":     "build/soong",
		"-revision": "8cf3e5471db04da274965a8e5c0dc3d465f08c5f",
		"-upstream": "master",
	}

	if _, err := Lint(projects); err != nil {
		t.Error("FAIL:", err)
	}
}

func TestRoutine(t *testing.T) {
	projects := make([]interface{}, 1)

	projects[0] = map[string]interface{}{
		"-groups":   "pdk,tradefed",
		"-name":     "platform/build/soong",
		"-path":     "build/soong",
		"-revision": "8cf3e5471db04da274965a8e5c0dc3d465f08c5f",
		"-upstream": "master",
	}

	if buf := Routine(projects[0]); buf == nil {
		t.Error("FAIL")
	}
}
