package server

import (
	"github.com/zenanet-network/go-zenanet/common"
	"github.com/zenanet-network/go-zenanet/core/types"
	"github.com/zenanet-network/go-zenanet/internal/cli/server/proto"
	"github.com/zenanet-network/go-zenanet/p2p"

	protobor "github.com/maticnetwork/polyproto/bor"
	protocommon "github.com/maticnetwork/polyproto/common"
	protoutil "github.com/maticnetwork/polyproto/utils"
)

func PeerInfoToPeer(info *p2p.PeerInfo) *proto.Peer {
	return &proto.Peer{
		Id:      info.ID,
		Enode:   info.Enode,
		Enr:     info.ENR,
		Caps:    info.Caps,
		Name:    info.Name,
		Trusted: info.Network.Trusted,
		Static:  info.Network.Static,
	}
}

func ConvertBloomToProtoBloom(bloom types.Bloom) *protobor.Bloom {
	return &protobor.Bloom{
		Bloom: bloom.Bytes(),
	}
}

func ConvertLogsToProtoLogs(logs []*types.Log) []*protobor.Log {
	var protoLogs []*protobor.Log
	for _, log := range logs {
		protoLog := &protobor.Log{
			Address:     protoutil.ConvertAddressToH160(log.Address),
			Topics:      ConvertTopicsToProtoTopics(log.Topics),
			Data:        log.Data,
			BlockNumber: log.BlockNumber,
			TxHash:      protoutil.ConvertHashToH256(log.TxHash),
			TxIndex:     uint64(log.TxIndex),
			BlockHash:   protoutil.ConvertHashToH256(log.BlockHash),
			Index:       uint64(log.Index),
			Removed:     log.Removed,
		}
		protoLogs = append(protoLogs, protoLog)
	}

	return protoLogs
}

func ConvertTopicsToProtoTopics(topics []common.Hash) []*protocommon.H256 {
	var protoTopics []*protocommon.H256
	for _, topic := range topics {
		protoTopics = append(protoTopics, protoutil.ConvertHashToH256(topic))
	}

	return protoTopics
}

func ConvertReceiptToProtoReceipt(receipt *types.Receipt) *protobor.Receipt {
	return &protobor.Receipt{
		Type:              uint64(receipt.Type),
		PostState:         receipt.PostState,
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		Bloom:             ConvertBloomToProtoBloom(receipt.Bloom),
		Logs:              ConvertLogsToProtoLogs(receipt.Logs),
		TxHash:            protoutil.ConvertHashToH256(receipt.TxHash),
		ContractAddress:   protoutil.ConvertAddressToH160(receipt.ContractAddress),
		GasUsed:           receipt.GasUsed,
		EffectiveGasPrice: receipt.EffectiveGasPrice.Int64(),
		BlobGasUsed:       receipt.BlobGasUsed,
		BlockHash:         protoutil.ConvertHashToH256(receipt.BlockHash),
		BlockNumber:       receipt.BlockNumber.Int64(),
		TransactionIndex:  uint64(receipt.TransactionIndex),
	}
}
