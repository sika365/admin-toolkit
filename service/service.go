package service

import (
	"errors"

	"github.com/alifakhimi/simple-service-go"
	"github.com/sirupsen/logrus"

	"github.com/sika365/admin-tools/config"
	"github.com/sika365/admin-tools/pkg/category"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/pkg/file"
	"github.com/sika365/admin-tools/pkg/image"
	"github.com/sika365/admin-tools/pkg/node"
	"github.com/sika365/admin-tools/pkg/product"
	"github.com/sika365/admin-tools/registrar"
)

type Service struct {
	simple.Simple
}

// New new service
func New(configPath string) simple.Interface {
	if err := simple.ReadConfig(configPath, config.Config()); err != nil {
		logrus.Exit(1)
	}

	s := &Service{
		Simple: simple.NewWithConfig(config.Config().Config),
	}

	if err := s.Simple.Run(s); err != nil {
		logrus.Errorln(err)
	}

	return s
}

// Init call by hook in simple service life-cycle
func (svc *Service) Init() error {
	logrus.Infoln("Initializing service")
	if conf := config.Config(); conf == nil {
		return errors.New("no config instance found")
	} else if client, err := client.New(conf); err != nil {
		return err
	} else if db, err := conf.GetDB("db"); err != nil {
		return err
	} else if h, err := conf.GetHttpServer("main"); err != nil {
		return err
	} else if err := svc.UseContext(h.Echo()); err != nil {
		return err
		// Initializing packages
	} else if err := registrar.
		Add(file.New(h, db, client)).
		Add(image.New(h, db, client)).
		Add(category.New(h, db, client)).
		Add(product.New(h, db, client)).
		Add(node.New(h, db, client)).
		// Add more package
		Init(); err != nil {
		return err
	} else {
		logrus.Infoln("Service initialization done.")
		return nil
	}
	// Viwes
	// views.Init("./views")
}
