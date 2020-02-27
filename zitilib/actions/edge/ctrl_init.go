package edge

import (
	"errors"
	"fmt"
	"github.com/netfoundry/fablab/kernel/fablib"
	"github.com/netfoundry/fablab/kernel/fablib/actions/host"
	"github.com/netfoundry/fablab/kernel/model"
)

func EdgeInit(regionSpec, hostSpec, componentSpec string) model.Action {
	return &edgeInit{
		regionSpec:    regionSpec,
		hostSpec:      hostSpec,
		componentSpec: componentSpec,
	}
}

func (init *edgeInit) Execute(m *model.Model) error {
	hosts := m.GetHosts(init.regionSpec, init.hostSpec)

	username := m.MustVariable("credentials", "edge", "username").(string)
	password := m.MustVariable("credentials", "edge", "password").(string)

	if username == "" {
		return errors.New("variable credentials/edge/username must be a string")
	}

	if password == "" {
		return errors.New("variable credentials/edge/password must be a string")
	}

	for _, h := range hosts {
		components := h.GetComponents(init.componentSpec)
		for _, c := range components {
			sshConfigFactory := fablib.NewSshConfigFactoryImpl(m, h.PublicIp)



			if err := host.Exec(h, fmt.Sprintf("/home/%s/fablab/bin/%s --log-formatter pfxlog edge init /home/%s/fablab/cfg/%s -u %s -p %s > logs/%s.edge.init.log 2>&1", sshConfigFactory.User(), c.BinaryName, sshConfigFactory.User(), c.ConfigName, username, password, c.BinaryName)).Execute(m); err != nil {
				return err
			}

		}
	}

	return nil
}

type edgeInit struct {
	regionSpec    string
	hostSpec      string
	componentSpec string
}
