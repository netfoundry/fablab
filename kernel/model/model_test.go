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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetVariable(t *testing.T) {
	m := &Model{
		Scope: Scope{
			Defaults: Variables{
				"a": Variables{
					"b": Variables{
						"c": "oh, wow!",
					},
				},
			},
		},
	}

	value, found := m.Defaults.Get("a", "b", "c")
	assert.True(t, found)
	assert.Equal(t, "oh, wow!", value)

	value, found = m.Defaults.Get("c")
	assert.False(t, found)

	value, found = m.Defaults.Get("d", "e", "f")
	assert.False(t, found)
}
