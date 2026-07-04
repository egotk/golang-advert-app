package coregrpc

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func GRPCToTimeNullable(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}

	t := ts.AsTime()
	return &t
}
