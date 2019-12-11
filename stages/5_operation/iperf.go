/*
	Copyright 2019 Netfoundry, Inc.

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

package operation

import (
	"fmt"
	"github.com/netfoundry/fablab/kernel"
	"github.com/netfoundry/fablab/kernel/lib"
	"github.com/sirupsen/logrus"
	"time"
)

func Iperf(seconds int) kernel.OperatingStage {
	return &iperf{seconds: seconds}
}

func (iperf *iperf) Operate(m *kernel.Model) error {
	serverHosts := m.GetHosts("@iperf-server", "@iperf-server")
	clientHosts := m.GetHosts("@iperf-client", "@iperf-client")
	if len(serverHosts) == 1 && len(clientHosts) == 1 {
		serverHost := serverHosts[0]
		clientHost := clientHosts[0]
		sshUser := m.MustVariable("credentials", "ssh", "username").(string)
		go iperf.runServer(serverHost, sshUser)

		time.Sleep(10 * time.Second)

		if err := lib.RemoteKill(sshUser, clientHost.PublicIp, "iperf3"); err != nil {
			return fmt.Errorf("error killing iperf3 clients (%w)", err)
		}

		initiator := m.GetHosts("@initiator", "@initiator")[0]
		iperfCmd := fmt.Sprintf("iperf3 -c %s -p 7002 -t %d --json", initiator.PublicIp, iperf.seconds)
		output, err := lib.RemoteExec(sshUser, clientHost.PublicIp, iperfCmd)
		if err == nil {
			logrus.Infof("iperf3 client completed, output [%s]", output)
		} else {
			logrus.Errorf("iperf3 client failure [%s] (%w)", output, err)
		}

	} else {
		logrus.Warnf("found [%d] server hosts, and [%d] client hosts, skipping", len(serverHosts), len(clientHosts))
	}
	return nil
}

func (iperf *iperf) runServer(h *kernel.Host, sshUser string) {
	if err := lib.RemoteKill(sshUser, h.PublicIp, "iperf3"); err != nil {
		logrus.Errorf("error killing iperf3 clients (%w)", err)
		return
	}

	output, err := lib.RemoteExec(sshUser, h.PublicIp, "iperf3 -s -p 7001 --one-off --json")
	if err == nil {
		logrus.Infof("iperf3 server completed, output [%s]", output)
	} else {
		logrus.Errorf("iperf3 server failure [%s] (%w)", output, err)
	}
}

type iperf struct {
	seconds int
}