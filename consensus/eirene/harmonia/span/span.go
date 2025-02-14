package span

import (
	"github.com/zenanet-network/go-zenanet/consensus/eirene/validset"
)

// Span Bor represents a current bor span
type Span struct {
	ID         uint64 `json:"span_id" yaml:"span_id"`
	StartBlock uint64 `json:"start_block" yaml:"start_block"`
	EndBlock   uint64 `json:"end_block" yaml:"end_block"`
}

// HeimdallSpan represents span from heimdall APIs
type HarmoniaSpan struct {
	Span
	ValidatorSet      validset.ValidatorSet `json:"validator_set" yaml:"validator_set"`
	SelectedProducers []validset.Validator  `json:"selected_producers" yaml:"selected_producers"`
	ChainID           string                `json:"eirene_chain_id" yaml:"eirene_chain_id"`
}
