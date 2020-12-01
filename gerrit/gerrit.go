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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type Gerrit struct {
	Option string
	Pass   string
	Url    string
	User   string
}

func (g Gerrit) Get(id int) (map[string]interface{}, error) {
	_url := g.Url + "/changes/" + strconv.Itoa(id) + "/detail"
	if g.User != "" && g.Pass != "" {
		_url = g.Url + "/a/changes/" + strconv.Itoa(id) + "/detail"
	}

	req, err := http.NewRequest(http.MethodGet, _url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if g.User != "" && g.Pass != "" {
		req.SetBasicAuth(g.User, g.Pass)
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}

	buf := map[string]interface{}{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	return buf, nil
}

func (g Gerrit) Query(search string, start int) (map[string]interface{}, error) {
	_url := g.Url + "/changes/"
	if g.User != "" && g.Pass != "" {
		_url = g.Url + "/a/changes/"
	}

	req, err := http.NewRequest(http.MethodGet, _url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if g.User != "" && g.Pass != "" {
		req.SetBasicAuth(g.User, g.Pass)
	}

	q := req.URL.Query()
	q.Add("o", g.Option)
	q.Add("q", search)
	q.Add("start", strconv.Itoa(start))
	req.URL.RawQuery = q.Encode()

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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}

	buf := []map[string]interface{}{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	if len(buf) == 0 {
		return nil, errors.New("search not matched")
	}

	return buf[0], nil
}
