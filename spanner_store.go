package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

type SpannerStore struct {
	Client *spanner.Client
}

func NewSpannerStore(ctx context.Context, projectID string, instance string, database string) (*SpannerStore, error) {
	client, err := spanner.NewClient(ctx, fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instance, database))
	if err != nil {
		return nil, err
	}
	return &SpannerStore{client}, nil
}

func (s *SpannerStore) ExecuteQuery(ctx context.Context, sql string) ([]map[string]interface{}, error) {
	iter := s.Client.Single().Query(ctx, spanner.NewStatement(sql))
	defer iter.Stop()

	l := []map[string]interface{}{}
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return l, err
		}
		rm := map[string]interface{}{}
		for _, n := range row.ColumnNames() {
			var v spanner.GenericColumnValue
			if err := row.ColumnByName(n, &v); err != nil {
				return l, errors.Wrapf(err, "failed %s get value", n)
			}
			rm[n] = v.Value
		}
		l = append(l, rm)
	}

	return l, nil
}
