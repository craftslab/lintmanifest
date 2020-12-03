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

package manifest

import (
	"github.com/clbanning/mxj"
	"github.com/pkg/errors"
)

type Manifest struct {
	manifest mxj.Maps
}

func (m *Manifest) Load(name string) error {
	buf, err := mxj.NewMapsFromXmlFile(name)
	if err != nil {
		return errors.Wrap(err, "load failed")
	}

	if len(buf) != 1 {
		return errors.Wrap(err, "manifest invalid")
	}

	m.manifest = buf

	return nil
}

func (m Manifest) Projects() ([]interface{}, error) {
	if m.manifest == nil || len(m.manifest) != 1 {
		return nil, errors.New("manifest invalid")
	}

	if _, ok := m.manifest[0]["manifest"]; !ok {
		return nil, errors.New("manifest invalid")
	}

	if _, ok := m.manifest[0]["manifest"].(map[string]interface{})["project"]; !ok {
		return nil, errors.New("manifest invalid")
	}

	buf := m.manifest[0]["manifest"].(map[string]interface{})["project"]

	if b, ok := buf.(map[string]interface{}); ok {
		return []interface{}{b}, nil
	}

	return buf.([]interface{}), nil
}
