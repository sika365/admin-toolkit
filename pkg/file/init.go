package file

import (
	simutils "github.com/alifakhimi/simple-utils-go"

	"github.com/sika365/admin-tools/registrar"
)

const (
	PackageName = "file"
)

type pkg struct {
	rest  Rest
	logic Logic
	repo  Repo
	err   error
	//
	h  *simutils.HttpServer
	db *simutils.DBConnection
}

func New(h *simutils.HttpServer, db *simutils.DBConnection) registrar.Package {
	i := &pkg{
		h:  h,
		db: db,
	}
	if i.repo, i.err = newRepo(); i.err != nil {
		return i
	} else if i.logic, i.err = newLogic(i.repo, db); i.err != nil {
		return i
	} else if i.rest, i.err = newRest(i.logic, h); i.err != nil {
		return i
	} else {
		return i
	}
}

func (i *pkg) Init() error {
	// update db schema
	if err := i.Migrator(); err != nil {
		return err
	} else {
		return nil
	}
}

func (i *pkg) Name() string {
	return PackageName
}

func (i *pkg) Error() error {
	return i.err
}

func (i *pkg) GetRest() Rest {
	return i.rest
}

func (i *pkg) GetLogic() Logic {
	return i.logic
}

func (i *pkg) GetRepo() Repo {
	return i.repo
}
