package threads

import "context"

type Service interface {
	CreateThread(context.Context, *ReqCreateThreads) error
	SearchByMainHome(context.Context, *ReqSearchByMainHome) (*ThreadsList, error)
}
