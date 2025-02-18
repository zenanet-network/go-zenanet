package eth

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/consensus/eirene"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/clerk"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/harmonia/checkpoint"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/harmonia/milestone"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/harmonia/span"
)

type mockHarmonia struct {
	fetchCheckpoint         func(ctx context.Context, number int64) (*checkpoint.Checkpoint, error)
	fetchCheckpointCount    func(ctx context.Context) (int64, error)
	fetchMilestone          func(ctx context.Context) (*milestone.Milestone, error)
	fetchMilestoneCount     func(ctx context.Context) (int64, error)
	fetchNoAckMilestone     func(ctx context.Context, milestoneID string) error
	fetchLastNoAckMilestone func(ctx context.Context) (string, error)
}

func (m *mockHarmonia) StateSyncEvents(ctx context.Context, fromID uint64, to int64) ([]*clerk.EventRecordWithTime, error) {
	return nil, nil
}
func (m *mockHarmonia) Span(ctx context.Context, spanID uint64) (*span.HarmoniaSpan, error) {
	//nolint:nilnil
	return nil, nil
}
func (m *mockHarmonia) FetchCheckpoint(ctx context.Context, number int64) (*checkpoint.Checkpoint, error) {
	return m.fetchCheckpoint(ctx, number)
}
func (m *mockHarmonia) FetchCheckpointCount(ctx context.Context) (int64, error) {
	return m.fetchCheckpointCount(ctx)
}
func (m *mockHarmonia) FetchMilestone(ctx context.Context) (*milestone.Milestone, error) {
	return m.fetchMilestone(ctx)
}
func (m *mockHarmonia) FetchMilestoneCount(ctx context.Context) (int64, error) {
	return m.fetchMilestoneCount(ctx)
}
func (m *mockHarmonia) FetchNoAckMilestone(ctx context.Context, milestoneID string) error {
	return m.fetchNoAckMilestone(ctx, milestoneID)
}
func (m *mockHarmonia) FetchLastNoAckMilestone(ctx context.Context) (string, error) {
	return m.fetchLastNoAckMilestone(ctx)
}

func (m *mockHarmonia) FetchMilestoneID(ctx context.Context, milestoneID string) error {
	return m.fetchNoAckMilestone(ctx, milestoneID)
}

func (m *mockHarmonia) Close() {}

func TestFetchWhitelistCheckpointAndMilestone(t *testing.T) {
	t.Parallel()

	// create an empty ethHandler
	handler := &ethHandler{}

	// create a mock checkpoint verification function and use it to create a verifier
	verify := func(ctx context.Context, eth *Zenanet, handler *ethHandler, start uint64, end uint64, hash string, isCheckpoint bool) (string, error) {
		return "", nil
	}

	verifier := newEireneVerifier()
	verifier.setVerify(verify)

	// Create a mock harmonia instance and use it for creating a eirene instance
	var harmonia mockHarmonia

	eirene := &eirene.Eirene{HarmoniaClient: &harmonia}

	fetchCheckpointTest(t, &harmonia, eirene, handler, verifier)
	fetchMilestoneTest(t, &harmonia, eirene, handler, verifier)
}

func (b *eireneVerifier) setVerify(verifyFn func(ctx context.Context, eth *Zenanet, handler *ethHandler, start uint64, end uint64, hash string, isCheckpoint bool) (string, error)) {
	b.verify = verifyFn
}

func fetchCheckpointTest(t *testing.T, harmonia *mockHarmonia, eirene *eirene.Eirene, handler *ethHandler, verifier *eireneVerifier) {
	t.Helper()

	var checkpoints []*checkpoint.Checkpoint
	// create a mock fetch checkpoint function
	harmonia.fetchCheckpoint = func(_ context.Context, number int64) (*checkpoint.Checkpoint, error) {
		if len(checkpoints) == 0 {
			return nil, errCheckpoint
		} else if number == -1 {
			return checkpoints[len(checkpoints)-1], nil
		} else {
			return checkpoints[number-1], nil
		}
	}

	// create a background context
	ctx := context.Background()

	_, _, err := handler.fetchWhitelistCheckpoint(ctx, eirene, nil, verifier)
	require.ErrorIs(t, err, errCheckpoint)

	// create 4 mock checkpoints
	checkpoints = createMockCheckpoints(4)

	blockNum, blockHash, err := handler.fetchWhitelistCheckpoint(ctx, eirene, nil, verifier)

	// Check if we have expected result
	require.Equal(t, err, nil)
	require.Equal(t, checkpoints[len(checkpoints)-1].EndBlock.Uint64(), blockNum)
	require.Equal(t, checkpoints[len(checkpoints)-1].RootHash, blockHash)
}

func fetchMilestoneTest(t *testing.T, harmonia *mockHarmonia, eirene *eirene.Eirene, handler *ethHandler, verifier *eireneVerifier) {
	t.Helper()

	var milestones []*milestone.Milestone
	// create a mock fetch checkpoint function
	harmonia.fetchMilestone = func(_ context.Context) (*milestone.Milestone, error) {
		if len(milestones) == 0 {
			return nil, errMilestone
		} else {
			return milestones[len(milestones)-1], nil
		}
	}

	// create a background context
	ctx := context.Background()

	_, _, err := handler.fetchWhitelistMilestone(ctx, eirene, nil, verifier)
	require.ErrorIs(t, err, errMilestone)

	// create 4 mock checkpoints
	milestones = createMockMilestones(4)

	num, hash, err := handler.fetchWhitelistMilestone(ctx, eirene, nil, verifier)

	// Check if we have expected result
	require.Equal(t, err, nil)
	require.Equal(t, milestones[len(milestones)-1].EndBlock.Uint64(), num)
	require.Equal(t, milestones[len(milestones)-1].Hash, hash)
}

func createMockCheckpoints(count int) []*checkpoint.Checkpoint {
	var (
		checkpoints []*checkpoint.Checkpoint = make([]*checkpoint.Checkpoint, count)
		startBlock  int64                    = 257 // any number can be used
	)

	for i := 0; i < count; i++ {
		checkpoints[i] = &checkpoint.Checkpoint{
			Proposer:   common.Address{},
			StartBlock: big.NewInt(startBlock),
			EndBlock:   big.NewInt(startBlock + 255),
			RootHash:   common.Hash{},
			EireneChainID: "2024",
			Timestamp:  uint64(time.Now().Unix()),
		}
		startBlock += 256
	}

	return checkpoints
}

func createMockMilestones(count int) []*milestone.Milestone {
	var (
		milestones []*milestone.Milestone = make([]*milestone.Milestone, count)
		startBlock int64                  = 257 // any number can be used
	)

	for i := 0; i < count; i++ {
		milestones[i] = &milestone.Milestone{
			Proposer:   common.Address{},
			StartBlock: big.NewInt(startBlock),
			EndBlock:   big.NewInt(startBlock + 255),
			Hash:       common.Hash{},
			EireneChainID: "2024",
			Timestamp:  uint64(time.Now().Unix()),
		}
		startBlock += 256
	}

	return milestones
}
