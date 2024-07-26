package image

import (
	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/registrar"
)

const (
	PackageName           = "image"
	DefaultRoot           = "./images"
	ImageContentTypeRegex = `image/.*`
	ImageBarcodeRegex     = `^(?P<title>\d+)(?:[\.|_].*)?$`
)

type Package struct {
	rest  Rest
	logic Logic
	repo  Repo
	err   error
	//
	h      *simutils.HttpServer
	db     *simutils.DBConnection
	client *client.Client
}

func New(h *simutils.HttpServer, db *simutils.DBConnection, client *client.Client) registrar.Package {
	i := &Package{
		h:      h,
		db:     db,
		client: client,
	}
	if i.repo, i.err = newRepo(); i.err != nil {
		return i
	} else if i.logic, i.err = newLogic(db, client, i.repo); i.err != nil {
		return i
	} else if i.rest, i.err = newRest(h, i.logic); i.err != nil {
		return i
	} else {
		return i
	}
}

func (i *Package) Init() error {
	// update db schema
	if err := i.Migrator(); err != nil {
		return err
	} else {
		return nil
	}
}

func (i *Package) Name() string {
	return PackageName
}

func (i *Package) Error() error {
	return i.err
}

func (i *Package) GetRest() Rest {
	return i.rest
}

func (i *Package) GetLogic() Logic {
	return i.logic
}

func (i *Package) GetRepo() Repo {
	return i.repo
}
