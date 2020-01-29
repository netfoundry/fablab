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

package rsync

import (
	"fmt"
	"github.com/netfoundry/fablab/kernel/internal"
	zitilab_bootstrap "github.com/netfoundry/fablab/zitilab/development/bootstrap"
)

func rsync(sourcePath, targetPath string) error {
	rsync := internal.NewProcess(zitilab_bootstrap.RsyncCommand(), "-avz", "-e", zitilab_bootstrap.SshCommand()+" -o StrictHostKeyChecking=no", "--delete", sourcePath, targetPath)
	rsync.WithTail(internal.StdoutTail)
	if err := rsync.Run(); err != nil {
		return fmt.Errorf("rsync failed (%w)", err)
	}
	return nil
}
