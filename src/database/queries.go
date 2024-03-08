package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/nbd-wtf/go-nostr"
)

func BuildFilterQuery(filter *nostr.Filter, query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if len(filter.Authors) > 0 {
		query = query.Where(squirrel.Eq{"pubkey": filter.Authors})
	}

	if len(filter.IDs) > 0 {
		query = query.Where(squirrel.Eq{"id": filter.IDs})
	}

	if len(filter.Kinds) > 0 {
		query = query.Where(squirrel.Eq{"kind": filter.Kinds})
	}

	if len(filter.Search) > 0 {
		// %% escapes % in printf literals
		query = query.Where("content LIKE ?", fmt.Sprint("%", filter.Search, "%"))
	}

	if !filter.Since.Time().IsZero() {
		query = query.Where(squirrel.GtOrEq{"created_at": filter.Since.Time().Unix()})
	}

	if !filter.Until.Time().IsZero() {
		query = query.Where(squirrel.LtOrEq{"created_at": filter.Until.Time().Unix()})
	}

	query = query.Where(squirrel.Or{squirrel.Eq{"expiration": nil}, squirrel.Lt{"expiration": time.Now().Unix()}})

	if len(filter.Tags) > 0 {
		if etags, ok := filter.Tags["e"]; ok {
			queryEtags := make([]string, 0, len(etags))
			for _, v := range etags {
				queryEtags = append(queryEtags, fmt.Sprint("'", v, ","))
			}

			query = query.Where(fmt.Sprint("etags && ARRAY[", strings.Join(queryEtags, ","), "]"))
		}

		if ptags, ok := filter.Tags["p"]; ok {
			queryPtags := make([]string, 0, len(ptags))
			for _, v := range ptags {
				queryPtags = append(queryPtags, fmt.Sprint("'", v, ","))
			}

			query = query.Where(fmt.Sprint("ptags && ARRAY[", strings.Join(queryPtags, ","), "]"))
		}

		gtags := make(nostr.Tags, 0, len(filter.Tags))
		for k, v := range filter.Tags {
			if k != "e" && k != "p" {
				gtags = append(gtags, v)
			}
		}

		if len(gtags) > 0 {
			queryGtags := make([]string, 0, len(gtags))
			for _, v := range gtags {
				queryGtags = append(queryGtags, fmt.Sprint("'", v, ","))
			}

			query = query.Where(fmt.Sprint("ptags && ARRAY[", strings.Join(queryGtags, ","), "]"))
		}
	}
	return query
}
