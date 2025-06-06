package graphql

import "github.com/slinkeres/ozontask/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostsService      service.Posts
	CommentsService   service.Comments
	CommentsObservers Observers
}