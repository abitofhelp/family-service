// Copyright (c) 2025 A Bit of Help, Inc.

package resolver

import (
	"context"

	"github.com/abitofhelp/family-service/interface/adapters/graphql/model"
)

// ParentCount is the resolver for the parentCount field.
func (r *familyResolver) ParentCount(ctx context.Context, obj *model.Family) (int, error) {
	return len(obj.Parents), nil
}

// ChildrenCount is the resolver for the childrenCount field.
func (r *familyResolver) ChildrenCount(ctx context.Context, obj *model.Family) (int, error) {
	return len(obj.Children), nil
}