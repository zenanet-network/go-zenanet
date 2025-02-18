package harmonia_app

import (
	"github.com/cosmos/cosmos-sdk/types"

	"github.com/zenanet-network/go-zenanet/log"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/cmd/heimdalld/service"

	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	stateFetchLimit = 50
)

type HarmoniaAppClient struct {
	hApp *app.HarmoniaApp
}

func NewHarmoniaAppClient() *HarmoniaAppClient {
	return &HarmoniaAppClient{
		hApp: service.GetHarmoniaApp(),
	}
}

func (h *HarmoniaAppClient) Close() {
	// Nothing to close as of now
	log.Warn("Shutdown detected, Closing Harmonia App conn")
}

func (h *HarmoniaAppClient) NewContext() types.Context {
	return h.hApp.NewContext(true, abci.Header{Height: h.hApp.LastBlockHeight()})
}
