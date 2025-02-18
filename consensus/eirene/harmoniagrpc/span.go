package harmoniagrpc

import (
	"context"

	"github.com/zenanet-network/go-zenanet/consensus/eirene/harmonia/span"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/validset"
	"github.com/zenanet-network/go-zenanet/log"

	proto "github.com/maticnetwork/polyproto/heimdall"
	protoutils "github.com/maticnetwork/polyproto/utils"
)

func (h *HarmoniaGRPCClient) Span(ctx context.Context, spanID uint64) (*span.HarmoniaSpan, error) {
	req := &proto.SpanRequest{
		ID: spanID,
	}

	log.Info("Fetching span", "spanID", spanID)

	res, err := h.client.Span(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Info("Fetched span", "spanID", spanID)

	return parseSpan(res.Result), nil
}

func parseSpan(protoSpan *proto.Span) *span.HarmoniaSpan {
	resp := &span.HarmoniaSpan{
		Span: span.Span{
			ID:         protoSpan.ID,
			StartBlock: protoSpan.StartBlock,
			EndBlock:   protoSpan.EndBlock,
		},
		ValidatorSet:      validset.ValidatorSet{},
		SelectedProducers: []validset.Validator{},
		ChainID:           protoSpan.ChainID,
	}

	for _, validator := range protoSpan.ValidatorSet.Validators {
		resp.ValidatorSet.Validators = append(resp.ValidatorSet.Validators, parseValidator(validator))
	}

	resp.ValidatorSet.Proposer = parseValidator(protoSpan.ValidatorSet.Proposer)

	for _, validator := range protoSpan.SelectedProducers {
		resp.SelectedProducers = append(resp.SelectedProducers, *parseValidator(validator))
	}

	return resp
}

func parseValidator(validator *proto.Validator) *validset.Validator {
	return &validset.Validator{
		ID:               validator.ID,
		Address:          protoutils.ConvertH160toAddress(validator.Address),
		VotingPower:      validator.VotingPower,
		ProposerPriority: validator.ProposerPriority,
	}
}
