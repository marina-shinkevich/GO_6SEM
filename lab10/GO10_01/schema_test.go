package main

import (
	"testing"

	"github.com/graph-gophers/graphql-go"
)

func TestGraphQLSchema(t *testing.T) {
	graphql.MustParseSchema(gqlSchema, &resolver{})
}
