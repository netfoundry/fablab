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
	initCmd.Flags().StringToStringVarP(&createBindings, "label", "l", nil, "label bindings to include in the model")
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initalizes a fablab project directory",
	Args:  cobra.ExactArgs(0),
	Run:   initialize,
}

var createBindings map[string]string

func initialize(*cobra.Command, []string) {
	var instanceId string

	if model.GetModel() == nil {
		logrus.Fatal("no model configured, exiting")
	}

	if model.GetModel().GetId() == "" {
		logrus.Fatal("no model id provided, exiting")
	}

	if id, err := model.NewInstance(); err == nil {
		instanceId = id
	} else {
		logrus.Fatalf("error creating instance (%v)", err)
	}

	logrus.Infof("allocated new instance [%s]", instanceId)

	if err := model.CreateLabel(instanceId, createBindings); err != nil {
		logrus.Fatalf("unable to create instance label [%s] (%v)", instanceId, err)
	}

	logrus.Infof("initialized fablab project using model [%s]", model.GetModel().GetId())
}
