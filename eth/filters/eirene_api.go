package filters

import (
	"context"
	"errors"

	zenanet "github.com/zenanet-network/go-zenanet"
	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/core/types"
	"github.com/zenanet-network/go-zenanet/params"
	"github.com/zenanet-network/go-zenanet/rpc"
)

// SetChainConfig sets chain config
func (api *FilterAPI) SetChainConfig(chainConfig *params.ChainConfig) {
	api.chainConfig = chainConfig
}

func (api *FilterAPI) GetEireneBlockLogs(ctx context.Context, crit FilterCriteria) ([]*types.Log, error) {
	if api.chainConfig == nil {
		return nil, errors.New("no chain config found. Proper PublicFilterAPI initialization required")
	}

	// get sprint from eirene config
	eireneConfig := api.chainConfig.Eirene

	var filter *EireneBlockLogsFilter
	if crit.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = NewEireneBlockLogsFilter(api.sys.backend, eireneConfig, *crit.BlockHash, crit.Addresses, crit.Topics)
	} else {
		// Convert the RPC block numbers into internal representations
		begin := rpc.LatestBlockNumber.Int64()
		if crit.FromBlock != nil {
			begin = crit.FromBlock.Int64()
		}

		end := rpc.LatestBlockNumber.Int64()
		if crit.ToBlock != nil {
			end = crit.ToBlock.Int64()
		}
		// Construct the range filter
		filter = NewEireneBlockLogsRangeFilter(api.sys.backend, eireneConfig, begin, end, crit.Addresses, crit.Topics)
	}

	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}

	return returnLogs(logs), err
}

// NewDeposits send a notification each time a new deposit received from bridge.
func (api *FilterAPI) NewDeposits(ctx context.Context, crit zenanet.StateSyncFilter) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		stateSyncData := make(chan *types.StateSyncData, 10)
		stateSyncSub := api.events.SubscribeNewDeposits(stateSyncData)

		// nolint: gosimple
		for {
			select {
			case h := <-stateSyncData:
				if h != nil && (crit.ID == h.ID || crit.Contract == h.Contract ||
					(crit.ID == 0 && crit.Contract == common.Address{})) {
					notifier.Notify(rpcSub.ID, h)
				}
			case <-rpcSub.Err():
				stateSyncSub.Unsubscribe()
				return
			case <-notifier.Closed():
				stateSyncSub.Unsubscribe()
				return
			}
		}
	}()

	return rpcSub, nil
}
