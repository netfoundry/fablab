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
	"github.com/pkg/errors"
	"strings"
)

func (m *Model) IsBound() bool {
	return m.bound
}

func (m *Model) GetVariable(name ...string) (interface{}, bool) {
	return m.Variables.Get(name...)
}

func (m *Model) MustVariable(name ...string) interface{} {
	return m.Variables.Must(name...)
}

func (m *Model) GetAction(name string) (Action, bool) {
	action, found := m.actions[name]
	return action, found
}

func (m *Model) SelectRegions(regionSpec string) []*Region {
	var regions []*Region
	for id, region := range m.Regions {
		if regionSpec == "*" || regionSpec == id {
			regions = append(regions, region)
		} else if strings.HasPrefix(regionSpec, "@") {
			for _, tag := range region.Tags {
				if tag == regionSpec[1:] {
					regions = append(regions, region)
				}
			}
		}
	}
	return regions
}

func (m *Model) SelectRegion(regionSpec string) (*Region, error) {
	regions := m.SelectRegions(regionSpec)
	if len(regions) == 1 {
		return regions[0], nil
	} else {
		return nil, errors.Errorf("[%s] matched [%d] regions, expected 1", regionSpec, len(regions))
	}
}

func (m *Model) MustSelectRegion(regionSpec string) *Region {
	region, err := m.SelectRegion(regionSpec)
	if err != nil {
		panic(err)
	}
	return region
}

func (m *Model) SelectHosts(regionSpec, hostSpec string) []*Host {
	var hosts []*Host
	regions := m.SelectRegions(regionSpec)
	for _, region := range regions {
		for id, host := range region.Hosts {
			if hostSpec == "*" || hostSpec == id {
				hosts = append(hosts, host)
			} else if strings.HasPrefix(hostSpec, "@") {
				for _, tag := range host.Tags {
					if tag == hostSpec[1:] {
						hosts = append(hosts, host)
					}
				}
			}
		}
	}
	return hosts
}

func (m *Model) SelectHost(regionSpec, hostSpec string) (*Host, error) {
	hosts := m.SelectHosts(regionSpec, hostSpec)
	if len(hosts) == 1 {
		return hosts[0], nil
	} else {
		return nil, errors.Errorf("[%s, %s] matched [%d] hosts, expected 1", regionSpec, hostSpec, len(hosts))
	}
}

func (m *Model) MustSelectHost(regionSpec, hostSpec string) *Host {
	host, err := m.SelectHost(regionSpec, hostSpec)
	if err != nil {
		panic(err)
	}
	return host
}

func (m *Model) SelectComponents(regionSpec, hostSpec, componentSpec string) []*Component {
	var components []*Component
	hosts := m.SelectHosts(regionSpec, hostSpec)
	for _, host := range hosts {
		for componentId, component := range host.Components {
			if componentSpec == "*" || componentSpec == componentId {
				components = append(components, component)
			} else if strings.HasPrefix(componentSpec, "@") {
				tag := strings.TrimPrefix(componentSpec, "@")
				for _, componentTag := range component.Tags {
					if componentTag == tag {
						components = append(components, component)
					}
				}
			}
		}
	}
	return components
}

func (m *Model) GetAllRegions() []*Region {
	var regions []*Region
	for _, region := range m.Regions {
		regions = append(regions, region)
	}
	return regions
}

func (r *Region) GetAllHosts() []*Host {
	hosts := make([]*Host, 0)
	for _, host := range r.Hosts {
		hosts = append(hosts, host)
	}
	return hosts
}

func (h *Host) HasTag(tag string) bool {
	for _, hostTag := range h.Tags {
		if hostTag == tag {
			return true
		}
	}
	return false
}

func (h *Host) GetComponents(componentSpec string) []*Component {
	var components []*Component
	if strings.HasPrefix(componentSpec, "@") {
		components = h.GetComponentsByTag(strings.TrimPrefix(componentSpec, "@"))
	} else {
		component := h.GetComponent(componentSpec)
		if component != nil {
			components = append(components, component)
		}
	}
	return components
}

func (h *Host) GetComponent(componentId string) *Component {
	component, found := h.Components[componentId]
	if found {
		return component
	}
	return nil
}

func (h *Host) GetComponentsByTag(componentTag string) []*Component {
	var components []*Component
	for _, component := range h.Components {
		for _, tag := range component.Tags {
			if tag == componentTag {
				components = append(components, component)
			}
		}
	}
	return components
}
