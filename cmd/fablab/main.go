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

package main

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/fablab/cmd/fablab/subcmd"
	"github.com/openziti/fablab/kernel/model"
	"github.com/openziti/fablab/zitilib"
	_ "github.com/openziti/fablab/zitilib"
	_ "github.com/openziti/fablab/zitilib/models/characterization"
	_ "github.com/openziti/fablab/zitilib/models/development/dilithium"
	_ "github.com/openziti/fablab/zitilib/models/edge"
	_ "github.com/openziti/fablab/zitilib/models/examples"
	_ "github.com/openziti/fablab/zitilib/models/transwarp"
	"github.com/openziti/foundation/transport"
	"github.com/openziti/foundation/transport/quic"
	"github.com/openziti/foundation/transport/tcp"
	"github.com/openziti/foundation/transport/tls"
	"github.com/sirupsen/logrus"
)

func init() {
	pfxlog.Global(logrus.InfoLevel)
	pfxlog.SetPrefix("github.com/openziti/")
	transport.AddAddressParser(quic.AddressParser{})
	transport.AddAddressParser(tls.AddressParser{})
	transport.AddAddressParser(tcp.AddressParser{})
	model.AddBootstrapExtension(&zitilib.Bootstrap{})
}

func main() {
	if err := subcmd.Execute(); err != nil {
		logrus.Fatalf("failure (%v)", err)
	}
}
