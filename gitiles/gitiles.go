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
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Gitiles struct {
}

func (g *Gitiles) Query(url, user, pass, project, revision string) ([]byte, error) {
	if user != "" && pass != "" {
		url = url + "/a/" + project + "/+/" + revision
	} else {
		url = url + "/" + project + "/+/" + revision
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client failed")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("client failed")
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read failed")
	}

	return buf, nil
}
