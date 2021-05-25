/*
	Copyright 2019 NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (m *Model) BindLabel(l *Label) {
	clean := true
	for regionId, region := range m.Regions {
		for hostId, host := range region.Hosts {
			publicIpBinding := fmt.Sprintf("%s_host_%s_public_ip", regionId, hostId)
			if binding, found := l.Bindings[publicIpBinding]; found {
				if publicIp, ok := binding.(string); ok {
					host.PublicIp = publicIp
				}
			} else {
				logrus.Warnf("no binding [%s]", publicIpBinding)
				clean = false
			}

			privateIpBinding := fmt.Sprintf("%s_host_%s_private_ip", regionId, hostId)
			if binding, found := l.Bindings[privateIpBinding]; found {
				if privateIp, ok := binding.(string); ok {
					host.PrivateIp = privateIp
				}
			} else {
				logrus.Warnf("no binding [%s]", privateIpBinding)
				clean = false
			}
		}
	}
	if clean {
		m.bound = true
	}
}

func GetLabel() *Label {
	return label
}

func (label *Label) Has(name ...string) bool {
	_, found := label.Get(name...)
	return found
}

func (label *Label) Must(name ...string) interface{} {
	value, found := label.Get(name...)
	if !found {
		logrus.Fatalf("label binding [%v] not found", name)
	}
	return value
}

func (label *Label) Get(name ...string) (interface{}, bool) {
	if len(name) < 1 {
		return nil, false
	}

	inputMap := label.Bindings
	for i := 0; i < (len(name) - 1); i++ {
		key := name[i]
		if value, found := inputMap[key]; found {
			lowerMap, ok := value.(Variables)
			if !ok {
				return nil, false
			}
			inputMap = lowerMap
		} else {
			return nil, false
		}
	}

	if value, found := inputMap[name[len(name)-1]]; found {
		return value, true
	}
	return nil, false
}

func (label *Label) Save() error {
	return label.SaveAtPath(label.path)
}

func (label *Label) SaveAtPath(path string) error {
	data, err := yaml.Marshal(label)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	labelDir := filepath.Dir(labelPath(path))
	if err := os.MkdirAll(labelDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create label directory [%s] (%s)", labelDir, err)
	}

	if err = ioutil.WriteFile(labelPath(path), data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (label *Label) GetFilePath(fileName string) string {
	return filepath.Join(label.path, fileName)
}

func CreateLabel(instanceId, modelName string) error {
	if err := assertNoLabel(instanceId); err != nil {
		return fmt.Errorf("error with instance path [%s] (%s)", instanceId, err)
	}
	if _, found := GetModel(modelName); !found {
		return fmt.Errorf("no such model [%s]", modelName)
	}
	l := &Label{
		Model: modelName,
		State: Created,
	}
	if err := l.SaveAtPath(instancePath(instanceId)); err != nil {
		return fmt.Errorf("error writing run label [%s] (%s)", instancePath(instanceId), err)
	}
	return nil
}

func ValidateModelName(modelName string) error {
	if _, found := GetModel(modelName); !found {
		return fmt.Errorf("no such model [%s]", modelName)
	}
	return nil
}

func LoadLabelForInstance(instanceId string) (*Label, error) {
	labelPath := instancePath(instanceId)
	return LoadLabel(labelPath)
}

func LoadLabel(path string) (*Label, error) {
	data, err := ioutil.ReadFile(filepath.Join(path, labelFilename))
	if err != nil {
		return nil, err
	}
	l := &Label{Bindings: Variables{}}
	if err = yaml.Unmarshal(data, l); err != nil {
		return nil, err
	}
	l.path = path
	return l, nil
}

func bootstrapLabel() error {
	instancePath := ActiveInstancePath()
	if _, err := os.Stat(labelPath(instancePath)); err != nil {
		if os.IsNotExist(err) {
			logrus.Warnf("no label at instance path [%s]", instancePath)
			return nil
		}
		return fmt.Errorf("unable to stat run label [%s] (%s)", labelPath(instancePath), err)
	}
	if l, err := LoadLabel(instancePath); err == nil {
		label = l
	} else {
		return fmt.Errorf("unable to bootstrap instance label [%s] (%s)", instancePath, err)
	}
	return nil
}

func assertNoLabel(instanceId string) error {
	if _, err := os.Stat(instancePath(instanceId)); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	} else {
		if old, err := LoadLabel(instancePath(instanceId)); err == nil {
			return fmt.Errorf("existing instance [%s] found at [%s]", old.Model, instancePath(instanceId))
		}
		return nil
	}
}

func labelPath(path string) string {
	return filepath.Join(path, labelFilename)
}

type Label struct {
	Model    string        `yaml:"model"`
	State    InstanceState `yaml:"state"`
	Bindings Variables     `yaml:"bindings"`
	path     string
}

type InstanceState int

const (
	Created InstanceState = iota
	Expressed
	Configured
	Kitted
	Distributed
	Activated
	Operating
	Disposed
)

func (instanceState InstanceState) String() string {
	names := [...]string{
		"Created",
		"Expressed",
		"Configured",
		"Kitted",
		"Distributed",
		"Activated",
		"Operating",
		"Disposed",
	}
	if instanceState < Created || instanceState > Disposed {
		return "<<Invalid>>"
	}
	return names[instanceState]
}

const labelFilename = "fablab.yml"
