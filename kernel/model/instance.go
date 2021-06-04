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
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func NewInstance() (string, error) {
	workingDir, err := GetWorkingDirectory()
	if err != nil {
		return "", errors.Wrap(err, "unable to get working directory")
	}

	dataDir := filepath.Join(workingDir, ".fablab")

	_, err = os.Stat(dataDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", errors.Wrapf(err, "unable to stat fablab project directory %v", workingDir)
		}

		err = os.MkdirAll(dataDir, 0700)
		if err != nil {
			return "", fmt.Errorf("unable to create fablab project directory [%s] (%w)", workingDir, err)
		}
	}

	return uuid.NewString(), nil
}

func ActiveInstancePath() string {
	return filepath.Join(instancePath(), ".fablab")
}

func instancePath() string {
	wd, err := GetWorkingDirectory()
	if err != nil {
		logrus.Fatalf("unable to get working directory: %v", err)
		return ""
	}
	return wd
}

func GetWorkingDirectory() (string, error) {
	// placeholder so we can add a working directory flag to commands
	return os.Getwd()
}
