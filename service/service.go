package service

import (
	"errors"

	"github.com/alifakhimi/simple-service-go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/sika365/admin-tools/config"
	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg"
	"github.com/sika365/admin-tools/pkg/file"
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
		// Migration
	} else if db, err := conf.GetDB("db"); err != nil {
		return err
	} else if h, err := conf.GetHttpServer("main"); err != nil {
		return err
	} else if err := svc.UseContext(h.Echo()); err != nil {
		return err
		// Initializing packages
	} else if err := pkg.
		NewRegistrar().
		Add(file.New(h, db)).
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

func (svc *Service) UseContext(e *echo.Echo) error {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &context.Context{Context: c}
			return next(cc)
		}
	})
	return nil
}
