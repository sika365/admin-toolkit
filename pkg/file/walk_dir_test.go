package file

import (
	"net/url"
	"reflect"
	"regexp"
	"testing"

	"github.com/spf13/cast"

	"github.com/sika365/admin-tools/test"
)

func Test_WalkDir(t *testing.T) {
	testReq := test.PreparingTest()
	filters := url.Values(cast.ToStringMapStringSlice(testReq.Meta["filters"]))
	reContentType := regexp.MustCompile(cast.ToString(cast.ToStringMap(testReq.Meta["filters"])["content_types"]))

	type args struct {
		root          string
		maxDepth      int
		reContentType *regexp.Regexp
	}
	tests := []struct {
		name    string
		args    args
		want    Files
		wantErr bool
	}{
		{
			name: "read image (.png|.jpg)",
			args: args{
				root:          cast.ToString(testReq.Meta["root"]),
				maxDepth:      cast.ToInt(filters.Get("max_depth")),
				reContentType: reContentType,
			},
			want: Files{}.AddFiles(
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
			got, err := WalkDir(tt.args.root, tt.args.maxDepth, tt.args.reContentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("walkDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("walkDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
