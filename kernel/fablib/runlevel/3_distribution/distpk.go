package distribution

import (
	"fmt"
	"github.com/openziti/fablab/kernel/fablib"
	"github.com/openziti/fablab/kernel/model"
	"github.com/sirupsen/logrus"
)

func DistributeSshKey(hostSpec string) model.DistributionStage {
	return &distSshKey{
		hostSpec: hostSpec,
	}
}

func (self *distSshKey) Distribute(run model.Run) error {
	return run.GetModel().ForEachHost(self.hostSpec, 25, func(host *model.Host) error {
		ssh := fablib.NewSshConfigFactoryImpl(run.GetModel(), host.PublicIp)
		keyPath := fmt.Sprintf("/home/%v/.ssh/id_rsa", ssh.User())

		if _, err := fablib.RemoteExecAll(ssh, fmt.Sprintf("rm -f %v", keyPath)); err == nil {
			logrus.Infof("%s => %s", host.PublicIp, "removing old PK")
		} else {
			return fmt.Errorf("error removing old PK on host [%s] (%w)", host.PublicIp, err)
		}

		if err := fablib.SendFile(ssh, ssh.KeyPath(), keyPath); err != nil {
			logrus.Errorf("[%s] unable to send %s => %s", host.PublicIp, ssh.KeyPath(), keyPath)
			return err
		}
		logrus.Infof("[%s] %s => %s", host.PublicIp, ssh.KeyPath(), keyPath)

		if _, err := fablib.RemoteExecAll(ssh, fmt.Sprintf("chmod 0400 %v", keyPath)); err == nil {
			logrus.Infof("%s => %s", host.PublicIp, "setting pk permissions")
			return nil
		} else {
			return fmt.Errorf("error setting pk permissions on host [%s] (%w)", host.PublicIp, err)
		}
	})
}

type distSshKey struct {
	hostSpec string
}
