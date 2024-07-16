package file

import (
	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sika365/admin-tools/pkg"
)

type filePackage struct {
	rest  Rest
	logic Logic
	repo  Repo
	err   error
	//
	h  *simutils.HttpServer
	db *simutils.DBConnection
}

func New(h *simutils.HttpServer, db *simutils.DBConnection) pkg.Package {
	i := &filePackage{
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

func (i *filePackage) Init() error {
	// update db schema
	if err := MigrateFile(i.db); err != nil {
		return err
	} else {
		return nil
	}
}

func (i *filePackage) Name() string {
	return "image"
}

func (i *filePackage) Error() error {
	return i.err
}

func (i *filePackage) GetRest() Rest {
	return i.rest
}

func (i *filePackage) GetLogic() Logic {
	return i.logic
}

func (i *filePackage) GetRepo() Repo {
	return i.repo
}

// TODO Implement replacer interface
