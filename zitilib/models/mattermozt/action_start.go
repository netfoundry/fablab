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
	"github.com/openziti/fablab/kernel/fablib/actions"
	"github.com/openziti/fablab/kernel/fablib/actions/component"
	"github.com/openziti/fablab/kernel/fablib/actions/semaphore"
	"github.com/openziti/fablab/kernel/model"
	"time"
)

func newStartAction() model.ActionBinder {
	action := &startAction{}
	return action.bind
}

func (a *startAction) bind(m *model.Model) model.Action {
	workflow := actions.Workflow()
	workflow.AddAction(component.Start("@ctrl", "@ctrl", "@ctrl"))
	workflow.AddAction(semaphore.Sleep(2 * time.Second))
	workflow.AddAction(component.Start("@router", "@router", "@router"))
	workflow.AddAction(semaphore.Sleep(2 * time.Second))
	workflow.AddAction(component.Start("@edge-router", "@edge-router", "@edge-router"))
	return workflow
}

type startAction struct{}
