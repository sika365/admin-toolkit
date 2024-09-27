package woocommerce

import (
	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alifakhimi/simple-utils-go/simscheme"

	"github.com/sika365/admin-tools/pkg/category"
	"github.com/sika365/admin-tools/pkg/node"
	"github.com/sika365/admin-tools/pkg/product"
	"github.com/sika365/admin-tools/registrar"
	"github.com/sika365/admin-tools/service/client"
)

const (
	PackageName = "woocommerce"
)

var (
	defScheme       = simscheme.GetSchema()
	prdRecDoc       = defScheme.AddNewDocumentWithType(&product.ProductRecord{})
	catRecDoc       = defScheme.AddNewDocumentWithType(&category.CategoryRecord{})
	lnodeRecDoc     = defScheme.AddNewDocumentWithType(&node.LocalNode{})
	postDoc         = defScheme.AddNewDocumentWithType(&WpPost{})
	termTaxonomyDoc = defScheme.AddNewDocumentWithType(&WpTermTaxonomy{})
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
	if catPkg, err := registrar.Get(category.PackageName); err != nil {
		return i
	} else if prdPkg, err := registrar.Get(product.PackageName); err != nil {
		return i
	} else if cat, ok := catPkg.(*category.Package); !ok {
		return i
	} else if prd, ok := prdPkg.(*product.Package); !ok {
		return i
	} else if catRepo := cat.GetRepo(); catRepo == nil {
		return i
	} else if catLogic := cat.GetLogic(); catLogic == nil {
		return i
	} else if prdRepo := prd.GetRepo(); prdRepo == nil {
		return i
	} else if prdLogic := prd.GetLogic(); prdLogic == nil {
		return i
	} else if i.repo, i.err = newRepo(client); i.err != nil {
		return i
	} else if i.logic, i.err = newLogic(db, client, i.repo, catLogic, catRepo, prdLogic, prdRepo); i.err != nil {
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
