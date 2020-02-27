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
	"github.com/netfoundry/fablab/kernel/fablib/actions"
	"github.com/netfoundry/fablab/kernel/fablib/actions/component"
	"github.com/netfoundry/fablab/kernel/fablib/actions/host"
	"github.com/netfoundry/fablab/kernel/fablib/actions/semaphore"
	"github.com/netfoundry/fablab/kernel/model"
	"github.com/netfoundry/fablab/zitilib/actions/cli"
	"github.com/netfoundry/fablab/zitilib/actions/edge"
	"path/filepath"
	"time"
)

func newBootstrapAction() model.ActionBinder {
	action := &bootstrapAction{}
	return action.bind
}

func (a *bootstrapAction) bind(m *model.Model) model.Action {
	sshUsername := m.MustVariable("credentials", "ssh", "username").(string)

	workflow := actions.Workflow()

	workflow.AddAction(component.Stop("@ctrl", "@ctrl", "@ctrl"))
	workflow.AddAction(edge.EdgeInit("@ctrl", "@ctrl", "@ctrl"))
	workflow.AddAction(component.Start("@ctrl", "@ctrl", "@ctrl"))
	workflow.AddAction(semaphore.Sleep(2 * time.Second))

	for _, router := range m.GetComponentsByTag("router") {
		cert := fmt.Sprintf("/intermediate/certs/%s-client.cert", router.PublicIdentity)
		workflow.AddAction(cli.Fabric("create", "router", filepath.Join(model.PkiBuild(), cert)))
	}

	for _, h := range m.GetAllHosts() {
		workflow.AddAction(host.Exec(h, fmt.Sprintf("mkdir -p /home/%s/.ziti", sshUsername)))
		workflow.AddAction(host.Exec(h, fmt.Sprintf("rm -f /home/%s/.ziti/identities.yml", sshUsername)))
		workflow.AddAction(host.Exec(h, fmt.Sprintf("ln -s /home/%s/fablab/cfg/remote_identities.yml /home/%s/.ziti/identities.yml", sshUsername, sshUsername)))
	}

	ctrl := m.GetHostByTags("ctrl", "ctrl")
	workflow.AddAction(edge.Login(ctrl))

	workflow.AddAction(edge.CtrlAddRouters("*", "*", "@edge-router"))
	workflow.AddAction(edge.EdgeRouterEnroll("*", "*", "@edge-router"))
	workflow.AddAction(component.Stop("@ctrl", "@ctrl", "@ctrl"))

	return workflow
}

type bootstrapAction struct{}
