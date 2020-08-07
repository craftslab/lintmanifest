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
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Gitiles struct {
}

func (g *Gitiles) Head(url, user, pass, project, branch string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url+"/"+project+"/+log/refs/heads/"+branch+"?format=JSON", nil)
	if err != nil {
		return "", errors.Wrap(err, "request failed")
	}

	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "client failed")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("client failed")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read failed")
	}

	buf := map[string]interface{}{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return "", errors.Wrap(err, "unmarshal failed")
	}

	if _, ok := buf["log"]; !ok {
		return "", errors.New("log invalid")
	}

	b := buf["log"].([]interface{})
	if len(b) == 0 {
		return "", errors.New("list invalid")
	}

	if _, ok := b[0].(map[string]interface{})["commit"]; !ok {
		return "", errors.New("commit invalid")
	}

	return b[0].(map[string]interface{})["commit"].(string), nil
}

func (g *Gitiles) Query(url, user, pass, project, revision string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url+"/"+project+"/+/"+revision+"?format=JSON", nil)
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
		return nil, errors.Wrap(err, "read failed")
	}

	return buf, nil
}
