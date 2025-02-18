package harmonia_app

import (
	"context"
	"fmt"
	"math/big"

	"github.com/zenanet-network/go-zenanet/consensus/eirene/harmonia/milestone"

	"github.com/zenanet-network/go-zenanet/log"

	chTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func (h *HarmoniaAppClient) FetchMilestoneCount(_ context.Context) (int64, error) {
	log.Debug("Fetching milestone count")

	res := h.hApp.CheckpointKeeper.GetMilestoneCount(h.NewContext())

	log.Debug("Fetched Milestone Count", "res", int64(res))

	return int64(res), nil
}

func (h *HarmoniaAppClient) FetchMilestone(_ context.Context) (*milestone.Milestone, error) {
	log.Debug("Fetching Latest Milestone")

	res, err := h.hApp.CheckpointKeeper.GetLastMilestone(h.NewContext())
	if err != nil {
		return nil, err
	}

	milestone := toEireneMilestone(res)
	log.Debug("Fetched Latest Milestone", "milestone", milestone)

	return milestone, nil
}

func (h *HarmoniaAppClient) FetchNoAckMilestone(_ context.Context, milestoneID string) error {
	log.Debug("Fetching No Ack Milestone By MilestoneID", "MilestoneID", milestoneID)

	res := h.hApp.CheckpointKeeper.GetNoAckMilestone(h.NewContext(), milestoneID)
	if res {
		log.Info("Fetched No Ack By MilestoneID", "MilestoneID", milestoneID)
		return nil
	}

	return fmt.Errorf("still no-ack milestone exist corresponding to milestoneID: %v", milestoneID)
}

func (h *HarmoniaAppClient) FetchLastNoAckMilestone(_ context.Context) (string, error) {
	log.Debug("Fetching Latest No Ack Milestone ID")

	res := h.hApp.CheckpointKeeper.GetLastNoAckMilestone(h.NewContext())

	log.Debug("Fetched Latest No Ack Milestone ID", "res", res)

	return res, nil
}

func (h *HarmoniaAppClient) FetchMilestoneID(_ context.Context, milestoneID string) error {
	log.Debug("Fetching Milestone ID ", "MilestoneID", milestoneID)

	res := chTypes.GetMilestoneID()

	if res == milestoneID {
		return nil
	}

	return fmt.Errorf("milestone corresponding to milestoneID: %v doesn't exist in harmonia", milestoneID)
}

func toEireneMilestone(hdMilestone *hmTypes.Milestone) *milestone.Milestone {
	return &milestone.Milestone{
		Proposer:   hdMilestone.Proposer.EthAddress(),
		StartBlock: big.NewInt(int64(hdMilestone.StartBlock)),
		EndBlock:   big.NewInt(int64(hdMilestone.EndBlock)),
		Hash:       hdMilestone.Hash.EthHash(),
		EireneChainID: hdMilestone.EireneChainID,
		Timestamp:  hdMilestone.TimeStamp,
	}
}
