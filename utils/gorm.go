package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/sika365/admin-tools/context"
)

var (
	// sortRegex defines the regex pattern for sorting
	sortRegex = regexp.MustCompile(`^(\w+)(?::(asc|desc))?$`)
)

func BuildGormQuery(ctx *context.Context, db *gorm.DB, queryParams url.Values) *gorm.DB {
	// init gorm db
	qb := db

	for field, values := range queryParams {
		// continue if value is empty
		if len(values) == 0 {
			continue
		}

		// extend case if more advanced query is necessary
		// all matches values will reach the default case
		switch field {
		case "search":
			qb = qb.Where("search like ?", values[0])
			for i := 1; i < len(values); i++ {
				qb = qb.Or("search like ?", values[i])
			}
		case "limit":
			qb = qb.Limit(cast.ToInt(values[0]))
		case "offset":
			qb = qb.Offset(cast.ToInt(values[0]))
		case "sort":
			orderByColumns := []clause.OrderByColumn{}
			for _, param := range values {
				param = strings.ToLower(param)
				if sortRegex.MatchString(param) {
					matches := sortRegex.FindStringSubmatch(param)
					orderByColumns = append(orderByColumns,
						clause.OrderByColumn{
							Column: clause.Column{
								Name: matches[1],
							},
							Desc: matches[2] == "desc",
						},
					)
				}
			}
			qb = qb.Order(clause.OrderBy{Columns: orderByColumns})
		case "includes":
		default:
			if len(values) == 1 {
				qb.Where(fmt.Sprintf("%s = ?", field), values)
			} else {
				qb.Where(fmt.Sprintf("%s in ?", field), values)
			}
		}
	}

	return qb
}
