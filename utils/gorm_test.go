package utils

import (
	"net/url"
	"testing"

	"github.com/sika365/admin-tools/context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBuildGormQuery(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		panic(err)
	}

	type args struct {
		ctx         *context.Context
		db          *gorm.DB
		queryParams url.Values
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "with sorts",
			args: args{
				db: db,
				queryParams: url.Values{
					"limit":  []string{"20"},
					"offset": []string{"10"},
					"sort":   []string{"col1", "col2:desc", "col3:asc"},
					"search": []string{"%a%"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildGormQuery(tt.args.ctx, tt.args.db, tt.args.queryParams); (got == nil) != tt.wantErr {
				t.Errorf("BuildGormQuery() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}
