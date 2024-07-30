package client

import (
	"testing"

	"github.com/alifakhimi/simple-utils-go/simrest"
	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/test"
)

func TestClient_GetCategoryByAlias(t *testing.T) {
	testReq := test.PreparingTest()
	client, _ := New(testReq.Config)

	type fields struct {
		Client *simrest.Client
	}
	type args struct {
		ctx   *context.Context
		alias string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fetch category by alias",
			fields: fields{
				Client: client.Client,
			},
			args: args{
				ctx:   &context.Context{},
				alias: "اداری",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Client: tt.fields.Client,
			}
			_, err := c.GetCategoryByAlias(tt.args.ctx, tt.args.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetCategoryByAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
