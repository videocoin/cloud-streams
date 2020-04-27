package rpc

import "context"

func tracingFilter(ctx context.Context, fullMethodName string) bool {
	methods := []string{"/grpc.health.v1.Health/Check"}

	for _, m := range methods {
		if m == fullMethodName {
			return false
		}
	}

	return true
}

func logrusFilter(fullMethodName string, err error) bool {
	methods := []string{"/grpc.health.v1.Health/Check"}

	for _, m := range methods {
		if m == fullMethodName {
			return false
		}
	}

	return true
}
