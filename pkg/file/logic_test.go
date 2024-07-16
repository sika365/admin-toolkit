package file

import (
	"net/url"
	"reflect"
	"regexp"
	"testing"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/test"
	"github.com/sika365/admin-tools/utils"
	"github.com/spf13/cast"
)

func Test_logic_ReadFiles(t *testing.T) {
	testReq := test.PreparingTest()
	repo, _ := newRepo()
	db, _ := testReq.Config.GetDB("db")
	filters := url.Values(cast.ToStringMapStringSlice(testReq.Meta["filters"]))
	maxDepth := cast.ToInt(utils.PopQueryParam[string](filters, "max_depth"))
	reContentType := regexp.MustCompile(utils.PopQueryParam[string](filters, "content_types"))

	type fields struct {
		db   *simutils.DBConnection
		repo Repo
	}
	type args struct {
		ctx           *context.Context
		root          string
		maxDepth      int
		reContentType *regexp.Regexp
		filters       url.Values
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    MapFiles
		wantErr bool
	}{
		{
			name: "read image files",
			fields: fields{
				db:   db,
				repo: repo,
			},
			args: args{
				root:          cast.ToString(testReq.Meta["root"]),
				maxDepth:      maxDepth,
				reContentType: reContentType,
				filters:       filters,
			},
			want: NewMapFiles().FromFiles(
				reContentType,
				"../../samples/images/1-1-1/1234.jpeg",
				"../../samples/images/1-1-2/1235.jpeg",
				"../../samples/images/1236.jpeg",
				"../../samples/images/1237.jpeg",
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logic{
				conn: tt.fields.db,
				repo: tt.fields.repo,
			}
			if got, err := l.ReadFiles(
				tt.args.ctx,
				tt.args.root,
				tt.args.maxDepth,
				tt.args.reContentType,
				tt.args.filters,
			); (err != nil) != tt.wantErr {
				t.Errorf("logic.ReadFiles() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(got.GetKeys(), tt.want.GetKeys()) {
				t.Errorf("walkDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
