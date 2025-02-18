package eirene

import (
	"context"

	"github.com/zenanet-network/go-zenanet/consensus/eirene/clerk"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/harmonia/checkpoint"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/harmonia/milestone"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/harmonia/span"
)

//go:generate mockgen -destination=../../tests/eirene/mocks/IHarmoniaClient.go -package=mocks . IHarmoniaClient
type IHarmoniaClient interface {
	StateSyncEvents(ctx context.Context, fromID uint64, to int64) ([]*clerk.EventRecordWithTime, error)
	Span(ctx context.Context, spanID uint64) (*span.HarmoniaSpan, error)
	FetchCheckpoint(ctx context.Context, number int64) (*checkpoint.Checkpoint, error)
	FetchCheckpointCount(ctx context.Context) (int64, error)
	FetchMilestone(ctx context.Context) (*milestone.Milestone, error)
	FetchMilestoneCount(ctx context.Context) (int64, error)
	FetchNoAckMilestone(ctx context.Context, milestoneID string) error // Fetch the bool value whether milestone corresponding to the given id failed in the Harmonia
	FetchLastNoAckMilestone(ctx context.Context) (string, error)       // Fetch latest failed milestone id
	FetchMilestoneID(ctx context.Context, milestoneID string) error    // Fetch the bool value whether milestone corresponding to the given id is in process in Harmonia
	Close()
}
