package edge

import (
	"errors"
	"github.com/openziti/fablab/kernel/model"
	"github.com/openziti/fablab/zitilib/cli"
	"path/filepath"
)

func Login(hostSelector string) model.Action {
	return &login{
		hostSelector: hostSelector,
	}
}

func (l *login) Execute(m *model.Model) error {
	ctrl, err := m.SelectHost(l.hostSelector)
	if err != nil {
		return err
	}
	username := m.MustVariable("credentials", "edge", "username").(string)
	password := m.MustVariable("credentials", "edge", "password").(string)
	edgeApiBaseUrl := ctrl.PublicIp + ":1280"

	caChain := filepath.Join(model.PkiBuild(), "intermediate", "certs", "intermediate.cert")

	if username == "" {
		return errors.New("variable credentials/edge/username must be a string")
	}

	if password == "" {
		return errors.New("variable credentials/edge/password must be a string")
	}

	_, err = cli.Exec(m, "edge", "login", edgeApiBaseUrl, "-c", caChain, "-u", username, "-p", password)
	return err
}

type login struct {
	hostSelector string
}
