package edge

import (
	"fmt"
	"github.com/netfoundry/fablab/kernel/fablib"
	"github.com/netfoundry/fablab/kernel/fablib/actions/host"
	"github.com/netfoundry/fablab/kernel/model"
	"github.com/sirupsen/logrus"
	"strings"
)

func EdgeRouterEnroll(regionSpec, hostSpec, componentSpec string) model.Action {
	return &edgeRouterEnroll{
		regionSpec:    regionSpec,
		hostSpec:      hostSpec,
		componentSpec: componentSpec,
	}
}

func (enroll *edgeRouterEnroll) Execute(m *model.Model) error {
	hosts := m.GetHosts(enroll.regionSpec, enroll.hostSpec)
	for _, h := range hosts {
		components := h.GetComponents(enroll.componentSpec)
		for _, c := range components {

			_, isEnrolled := c.Data["isEnrolled"]

			if !isEnrolled {
				localJwt := ""
				val, ok := c.Data["localJwt"]

				if !ok || val == nil {
					return fmt.Errorf("component [%s] does not have a local enrollment JWT", c.PublicIdentity)
				}

				localJwt = val.(string)

				if localJwt == "" {
					return fmt.Errorf("could not read local JWT from data component [%s]", c.PublicIdentity)
				}

				remoteJwt := ""
				val, ok = c.Data["remoteJwt"]

				if !ok || val == nil {
					return fmt.Errorf("component [%s] does not have a remote enrollment JWT", c.PublicIdentity)
				}

				remoteJwt = val.(string)

				if remoteJwt == "" {
					return fmt.Errorf("could not read remote JWT from data component [%s]", c.PublicIdentity)
				}

				ssh := fablib.NewSshConfigFactoryImpl(m, h.PublicIp)

				if err := fablib.SendFile(ssh, localJwt, remoteJwt); err != nil {
					return err
				}
				sshConfigFactory := fablib.NewSshConfigFactoryImpl(m, h.PublicIp)
				if output, err := fablib.RemoteExec(sshConfigFactory, "mkdir -p /home/fedora/logs"); err == nil {
					if output != "" {
						logrus.Infof("output [%s]", strings.Trim(output, " \t\r\n"))
					}
				} else {
					return err
				}

				return host.Exec(h, fmt.Sprintf("/home/fedora/fablab/bin/%s enroll /home/fedora/fablab/cfg/%s -j %s > /home/fedora/logs/%s.router.enroll.log 2>&1", c.BinaryName, c.ConfigName, remoteJwt, c.ConfigName)).Execute(m)
			}
		}
	}
	return nil
}

type edgeRouterEnroll struct {
	regionSpec    string
	hostSpec      string
	componentSpec string
}
