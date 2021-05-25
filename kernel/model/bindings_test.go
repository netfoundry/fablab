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

package model

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVariableResolvingModelDefaults(t *testing.T) {
	m := &Model{
		Scope: Scope{
			Defaults: Variables{
				"a": Variables{
					"b": true,
				},
			},
		},
	}

	bindings = Variables{}

	m.init("test")

	val, found := m.GetVariable("a", "b")

	req := require.New(t)
	req.True(found)
	req.Equal(val, true)
}

func TestVariableResolvingModelDefaultsWithDefault(t *testing.T) {
	m := &Model{
		Scope: Scope{
			Defaults: Variables{
				"a": Variables{
					"__default__": "hello",
					"c":           "bye",
				},
			},
		},
	}

	bindings = Variables{}

	m.init("test")

	req := require.New(t)

	val, found := m.GetVariable("a", "b")
	req.True(found)
	req.Equal(val, "hello")

	val, found = m.GetVariable("a", "foo")
	req.True(found)
	req.Equal(val, "hello")

	val, found = m.GetVariable("a", "c")
	req.True(found)
	req.Equal(val, "bye")

	val, found = m.GetVariable("b")
	req.False(found)
}

func TestBindBindingsRequiredToModel(t *testing.T) {
	bValue := "b-value"

	m := &Model{
		Scope: Scope{
			Defaults: Variables{
				"a": Variables{
					"b": true,
				},
			},
		},
	}

	bindings = Variables{
		"a": Variables{
			"b": bValue,
		},
		"c": "c-value",
	}

	m.init("test")

	val, found := m.GetVariable("a", "b")

	req := require.New(t)
	req.True(found)
	req.Equal(val, bValue)
}
