package harmonia_app

import (
	"context"

	hmTypes "github.com/maticnetwork/heimdall/types"

	"github.com/zenanet-network/go-zenanet/consensus/eirene/harmonia/span"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/validset"
	"github.com/zenanet-network/go-zenanet/log"
)

func (h *HarmoniaAppClient) Span(ctx context.Context, spanID uint64) (*span.HarmoniaSpan, error) {
	log.Info("Fetching span", "spanID", spanID)

	res, err := h.hApp.EireneKeeper.GetSpan(h.NewContext(), spanID)
	if err != nil {
		return nil, err
	}

	log.Info("Fetched span", "spanID", spanID)

	return toSpan(res), nil
}

func toSpan(hdSpan *hmTypes.Span) *span.HarmoniaSpan {
	return &span.HarmoniaSpan{
		Span: span.Span{
			ID:         hdSpan.ID,
			StartBlock: hdSpan.StartBlock,
			EndBlock:   hdSpan.EndBlock,
		},
		ValidatorSet:      toValidatorSet(hdSpan.ValidatorSet),
		SelectedProducers: toValidators(hdSpan.SelectedProducers),
		ChainID:           hdSpan.ChainID,
	}
}

func toValidatorSet(vs hmTypes.ValidatorSet) validset.ValidatorSet {
	return validset.ValidatorSet{
		Validators: toValidatorsRef(vs.Validators),
		Proposer:   toValidatorRef(vs.Proposer),
	}
}

func toValidators(vs []hmTypes.Validator) []validset.Validator {
	newVS := make([]validset.Validator, len(vs))

	for i, v := range vs {
		newVS[i] = toValidator(v)
	}

	return newVS
}

func toValidatorsRef(vs []*hmTypes.Validator) []*validset.Validator {
	newVS := make([]*validset.Validator, len(vs))

	for i, v := range vs {
		if v == nil {
			continue
		}

		newVS[i] = toValidatorRef(v)
	}

	return newVS
}

func toValidatorRef(v *hmTypes.Validator) *validset.Validator {
	return &validset.Validator{
		ID:               v.ID.Uint64(),
		Address:          v.Signer.EthAddress(),
		VotingPower:      v.VotingPower,
		ProposerPriority: v.ProposerPriority,
	}
}

func toValidator(v hmTypes.Validator) validset.Validator {
	return validset.Validator{
		ID:               v.ID.Uint64(),
		Address:          v.Signer.EthAddress(),
		VotingPower:      v.VotingPower,
		ProposerPriority: v.ProposerPriority,
	}
}
