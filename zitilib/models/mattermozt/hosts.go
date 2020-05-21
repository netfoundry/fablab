/*
	Copyright 2020 NetFoundry, Inc.

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

package mattermozt

import (
	"fmt"
	"github.com/openziti/fablab/kernel/model"
)

func newHostsFactory() model.Factory {
	return &hostsFactory{}
}

func (self *hostsFactory) Build(m *model.Model) error {
	ctrlType, found := m.GetVariable("mattermozt", "sizing", "ctrl")
	if !found {
		return fmt.Errorf("missing 'mattermozt/sizing/ctrl' variable")
	}
	m.Regions["local"].Hosts["ctrl"].InstanceType = ctrlType.(string)

	terminatorType, found := m.GetVariable("mattermozt", "sizing", "terminator")
	if !found {
		return fmt.Errorf("missing 'mattermozt/sizing/terminator' variable")
	}
	m.Regions["local"].Hosts["terminator"].InstanceType = terminatorType.(string)

	edgeType, found := m.GetVariable("mattermozt", "sizing", "edge")
	if !found {
		return fmt.Errorf("missing 'mattermozt/sizing/edge' variable")
	}
	m.Regions["local"].Hosts["edge"].InstanceType = edgeType.(string)

	serviceType, found := m.GetVariable("mattermozt", "sizing", "service")
	if !found {
		return fmt.Errorf("missing 'mattermozt/sizing/service' variable")
	}
	m.Regions["local"].Hosts["service"].InstanceType = serviceType.(string)

	return nil
}

type hostsFactory struct{}
