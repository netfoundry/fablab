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
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "name for the new instance")
	createCmd.Flags().StringToStringVarP(&createBindings, "label", "l", nil, "label bindings to include in the model")
	RootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create <model>",
	Short: "create a fablab instance from a model",
	Args:  cobra.ExactArgs(1),
	Run:   create,
}
var createName string
var createBindings map[string]string

func create(_ *cobra.Command, args []string) {
	var instanceId string
	modelName := args[0]

	if err := model.ValidateModelName(modelName); err != nil {
		logrus.Fatalf("unable to create instance (%v)", err)
	}
	if createName != "" {
		if err := model.NewNamedInstance(createName); err == nil {
			instanceId = createName
		} else {
			logrus.Fatalf("error creating named instance [%s] (%v)", createName, err)
		}
	} else {
		if id, err := model.NewInstance(); err == nil {
			instanceId = id
		} else {
			logrus.Fatalf("error creating instance (%v)", err)
		}
	}
	logrus.Infof("allocated new instance [%s]", instanceId)

	if err := model.CreateLabel(instanceId, modelName); err != nil {
		logrus.Fatalf("unable to create instance label [%s] (%v)", instanceId, err)
	}
	if createBindings != nil {
		logrus.Infof("setting label bindings = [%v]", createBindings)
		if l, err := model.LoadLabelForInstance(instanceId); err == nil {
			if l.Bindings == nil {
				l.Bindings = make(model.Bindings)
			}
			for k, v := range createBindings {
				l.Bindings[k] = v
			}
			if err := l.Save(); err != nil {
				logrus.Fatalf("error saving label bindings (%v)", err)
			}
		} else {
			logrus.Fatalf("error loading label (%v)", err)
		}
	}

	_, found := model.GetModel(modelName)
	if !found {
		logrus.Fatalf("no model [%s]", modelName)
	}
	logrus.Infof("using model [%s]", modelName)

	if err := model.SetActiveInstance(instanceId); err != nil {
		logrus.Fatalf("unable to set active instance (%v)", err)
	}
}
