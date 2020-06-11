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

package subcmd

import (
	"github.com/openziti/fablab/kernel/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "remove instance data from empty or disposed models",
	Args:  cobra.ExactArgs(0),
	Run:   clean,
}

func clean(_ *cobra.Command, _ []string) {
	if err := model.BootstrapInstance(); err != nil {
		logrus.Fatalf("error bootstrapping instance (%v)", err)
	}

	instanceIds, err := model.ListInstances()
	if err != nil {
		logrus.Fatalf("error listing instances (%v)", err)
	}

	activeInstanceId := model.ActiveInstanceId()
	for _, instanceId := range instanceIds {
		if l, err := model.LoadLabelForInstance(instanceId); err == nil {
			if l.State == model.Created || l.State == model.Disposed {
				if err := model.RemoveInstance(instanceId); err != nil {
					logrus.Fatalf("error removing instance [%s] (%v)", instanceId, err)
				}
				if instanceId == activeInstanceId {
					if err := model.ClearActiveInstance(); err != nil {
						logrus.Errorf("error clearing active instance (%v)", err)
					}
				}
				logrus.Infof("removed instance [%s]", instanceId)
			}
		} else {
			logrus.Warnf("error loading label for instance [%s] (%v)", instanceId, err)
			if err := model.RemoveInstance(instanceId); err != nil {
				logrus.Fatalf("error removing instance [%s] (%v)", instanceId, err)
			}
		}
	}
}
