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
	"fmt"
	"github.com/openziti/foundation/util/errorz"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"sort"
	"strings"
)

type Scope struct {
	entity           Entity
	Defaults         Variables
	VariableResolver VariableResolver
	Data             Data
	Tags             Tags
	bound            bool
}

func (scope *Scope) initialize(entity Entity, scoped bool) {
	scope.entity = entity
	sort.Strings(scope.Tags)

	if scope.VariableResolver == nil {
		if scoped {
			scope.VariableResolver = entity.GetModel().VarConfig.DefaultScopedVariableResolver
		} else {
			scope.VariableResolver = entity.GetModel().VarConfig.DefaultVariableResolver
		}
	}
}

func (scope *Scope) HasTag(tag string) bool {
	for _, hostTag := range scope.Tags {
		if hostTag == tag {
			return true
		}
	}
	return false
}

func (scope *Scope) WithTags(tags ...string) *Scope {
	scope.Tags = tags
	return scope
}

func (scope *Scope) HasVariable(name string) bool {
	_, found := scope.VariableResolver.Resolve(scope.entity, name)
	return found
}

func (scope *Scope) GetVariable(name string) (interface{}, bool) {
	return scope.VariableResolver.Resolve(scope.entity, name)
}

func (scope *Scope) PutVariable(name string, value interface{}) {
	path := scope.entity.GetModel().VarConfig.VariableNameParser(name)
	scope.Defaults.Put(path, value)
}

func (scope *Scope) GetStringVariable(name string) (string, bool) {
	val, found := scope.VariableResolver.Resolve(scope.entity, name)
	if !found {
		return "", false
	}
	if strVal, ok := val.(string); ok {
		return strVal, true
	}
	return fmt.Sprintf("%v", val), true
}

func (scope *Scope) GetStringVariableOr(name string, defaultValue string) string {
	val, found := scope.GetStringVariable(name)
	if found {
		return val
	}
	return defaultValue
}

func (scope *Scope) GetBoolVariable(name string) (bool, bool) {
	val, found := scope.VariableResolver.Resolve(scope.entity, name)
	if !found {
		return false, false
	}
	if boolVal, ok := val.(bool); ok {
		return boolVal, true
	}
	return strings.EqualFold("true", fmt.Sprintf("%v", val)), true
}

func (scope *Scope) GetVariableOr(name string, defaultValue interface{}) interface{} {
	val, found := scope.GetVariable(name)
	if found {
		return val
	}
	return defaultValue
}

func (scope *Scope) MustVariable(name string) interface{} {
	val, found := scope.GetVariable(name)
	if found {
		return val
	}
	logrus.Fatalf("no value defined for variable %+v", name)
	return nil
}

func (scope *Scope) MustStringVariable(name string) string {
	value := scope.MustVariable(name)
	result, ok := value.(string)
	if !ok {
		logrus.Fatalf("variable [%v] expected to have type string, but was %v", name, reflect.TypeOf(value))
	}
	return result
}

func (scope *Scope) GetRequiredStringVariable(holder errorz.ErrorHolder, name string) string {
	value, found := scope.GetVariable(name)
	if !found {
		holder.SetError(errors.Errorf("missing variable [%s]", name))
		return ""
	}
	result, ok := value.(string)
	if !ok {
		holder.SetError(errors.Errorf("variable [%v] expected to have type string, but was %v", name, reflect.TypeOf(value)))
	}
	return result
}

type Data map[string]interface{}
type Tags []string

type Variables map[string]interface{}

func (v Variables) Put(name []string, newValue interface{}) {
	if len(name) < 1 {
		return
	}

	inputMap := v
	for i := 0; i < (len(name) - 1); i++ {
		key := name[i]
		if value, found := inputMap[key]; found {
			lowerMap, ok := value.(Variables)
			if !ok {
				logrus.Fatalf("path %v overrides a submap", name)
			} else {
				inputMap = lowerMap
			}
		} else {
			lowerMap := Variables{}
			inputMap[key] = lowerMap
			inputMap = lowerMap
		}
	}

	key := name[len(name)-1]
	if val, found := inputMap[key]; found {
		if _, ok := val.(Variables); ok {
			logrus.Fatalf("path %v overrides a submap", name)
		}
	}

	inputMap[key] = newValue
}

func (v Variables) Get(name []string) (interface{}, bool) {
	if len(name) < 1 {
		return nil, false
	}

	inputMap := v
	for i := 0; i < (len(name) - 1); i++ {
		key := name[i]
		if value, found := inputMap[key]; found {
			lowerMap, ok := value.(Variables)
			if !ok {
				return nil, false
			}
			inputMap = lowerMap
		}
	}

	value, found := inputMap[name[len(name)-1]]
	if found {
		if _, ok := value.(Variables); ok {
			return nil, false
		}
	}
	return value, found
}

func (v Variables) getPath(path ...string) []Variables {
	result := []Variables{v}
	current := v
	for _, e := range path {
		next, found := current[e]
		if !found {
			return result
		}
		current, ok := next.(Variables)
		if !ok {
			return result
		}
		result = append(result, current)
	}
	return result
}

func (v Variables) getRelated(name string, path ...string) (interface{}, bool) {
	p := v.getPath(path...)
	for i := len(p) - 1; i >= 0; i-- {
		node := p[i]
		if val, found := node[name]; found {
			return val, true
		}
	}
	return nil, false
}

func (m *Model) IterateScopes(f func(i Entity, path ...string)) {
	f(m, []string{}...)
	for regionId, r := range m.Regions {
		f(r, []string{regionId}...)
		for hostId, h := range r.Hosts {
			f(h, []string{regionId, hostId}...)
			for componentId, c := range h.Components {
				f(c, []string{regionId, hostId, componentId}...)
			}
		}
	}
}

type VariableResolver interface {
	Resolve(entity Entity, name string) (interface{}, bool)
}

func NewScopedVariableResolver(resolver VariableResolver) *ScopedVariableResolver {
	return &ScopedVariableResolver{
		resolver: resolver,
	}
}

type ScopedVariableResolver struct {
	resolver VariableResolver
}

func (self *ScopedVariableResolver) Resolve(entity Entity, name string) (interface{}, bool) {
	entityPath := GetScopedEntityPath(entity)
	prefixedName := entity.GetModel().VarConfig.VariableNamePrefixMapper(entityPath, name)
	return self.resolver.Resolve(entity, prefixedName)
}

func NewCachingVariableResolver(resolver VariableResolver) *CachingVariableResolver {
	return &CachingVariableResolver{
		cache:    map[string]interface{}{},
		resolver: resolver,
	}
}

type CachingVariableResolver struct {
	cache    map[string]interface{}
	resolver VariableResolver
}

func (self *CachingVariableResolver) Resolve(entity Entity, name string) (interface{}, bool) {
	val, found := self.cache[name]
	if found {
		return val, found
	}
	val, found = self.resolver.Resolve(entity, name)
	if found {
		self.cache[name] = val
	}
	return val, found
}

func NewMapVariableResolver(variables Variables) *MapVariableResolver {
	return &MapVariableResolver{
		variables: variables,
	}
}

type MapVariableResolver struct {
	variables Variables
}

func (self *MapVariableResolver) Resolve(entity Entity, name string) (interface{}, bool) {
	path := entity.GetModel().VarConfig.VariableNameParser(name)
	return self.variables.Get(path)
}

type HierarchicalVariableResolver struct{}

func (self HierarchicalVariableResolver) Resolve(entity Entity, name string) (interface{}, bool) {
	path := entity.GetModel().VarConfig.VariableNameParser(name)
	current := entity
	for current != nil {
		if val, found := current.GetScope().Defaults.Get(path); found {
			return val, true
		}
		current = current.GetParentEntity()
	}
	return nil, false
}

type EnvVariableResolver struct{}

func (self EnvVariableResolver) Resolve(entity Entity, name string) (interface{}, bool) {
	key := entity.GetModel().VarConfig.EnvVariableNameMapper(name)
	return os.LookupEnv(key)
}

type ChainedVariableResolver struct {
	resolvers []VariableResolver
}

func (self *ChainedVariableResolver) AppendResolver(resolver VariableResolver) {
	self.resolvers = append(self.resolvers, resolver)
}

func (self *ChainedVariableResolver) Resolve(entity Entity, name string) (interface{}, bool) {
	for _, resolver := range self.resolvers {
		if val, found := resolver.Resolve(entity, name); found {
			return val, true
		}
	}
	return nil, false
}

type CmdLineArgVariableResolver struct{}

func (self CmdLineArgVariableResolver) Resolve(entity Entity, name string) (interface{}, bool) {
	config := entity.GetModel().VarConfig
	key := config.CommandLineVariableNameMapper(name)
	for _, arg := range os.Args {
		for _, prefix := range config.CommandLinePrefixes {
			argPrefix := prefix + key + "="
			if strings.HasPrefix(arg, argPrefix) {
				return strings.TrimPrefix(arg, argPrefix), true
			}
		}
	}
	return nil, false
}
