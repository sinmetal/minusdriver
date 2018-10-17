package main

import (
	"context"
	"testing"
)

func TestSpannerStore_ExecuteQuery(t *testing.T) {
	ctx := context.Background()
	s, err := NewSpannerStore(ctx, "gcpug-public-spanner", "merpay-sponsored-instance", "sinmetal")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.ExecuteQuery(ctx, "SELECT 1")
	if err != nil {
		t.Fatal(err)
	}
}
