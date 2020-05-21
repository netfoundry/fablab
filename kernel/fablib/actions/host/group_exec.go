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

package host

import (
	"fmt"
	"github.com/openziti/fablab/kernel/fablib"
	"github.com/openziti/fablab/kernel/model"
	"github.com/sirupsen/logrus"
)

func GroupExec(regionSpec, hostSpec, cmd string) model.Action {
	return &groupExec{
		regionSpec: regionSpec,
		hostSpec:   hostSpec,
		cmd:        cmd,
	}
}

func (groupExec *groupExec) Execute(m *model.Model) error {
	hosts := m.GetHosts(groupExec.regionSpec, groupExec.hostSpec)
	for _, h := range hosts {
		sshConfigFactory := fablib.NewSshConfigFactoryImpl(m, h.PublicIp)

		if o, err := fablib.RemoteExec(sshConfigFactory, groupExec.cmd); err != nil {
			logrus.Errorf("output [%s]", o)
			return fmt.Errorf("error executing process [%s] on [%s] (%s)", groupExec.cmd, h.PublicIp, err)
		}
	}
	return nil
}

type groupExec struct {
	regionSpec string
	hostSpec   string
	cmd        string
}
