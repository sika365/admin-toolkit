package test

import (
	"github.com/alifakhimi/simple-utils-go/simrest"
	"github.com/sirupsen/logrus"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/services/simple"

	"github.com/sika365/admin-tools/config"
	"github.com/sika365/admin-tools/context"
)

type TestRequirements struct {
	Config  *config.ServiceConfig
	Client  *simrest.Client
	Meta    map[string]any
	Context context.Context
}

func PreparingTest() *TestRequirements {
	var (
		testReq = TestRequirements{
			Config: &config.ServiceConfig{},
		}
	)

	if err := simple.ReadConfig("config.test.json", testReq.Config); err != nil {
		logrus.Exit(1)
	}

	testReq.Client, _ = testReq.Config.GetClient("sika365")
	testReq.Meta = testReq.Config.Meta

	return &testReq
}
