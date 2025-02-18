package eirene

import (
	"math/big"
	"sort"
	"testing"

	"github.com/maticnetwork/crand"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"

	"github.com/zenanet-network/go-zenanet/common"
	unique "github.com/zenanet-network/go-zenanet/common/set"
	"github.com/zenanet-network/go-zenanet/consensus/eirene/validset"
)

const (
	numVals = 100
)

func TestGetSignerSuccessionNumber_ProposerIsSigner(t *testing.T) {
	t.Parallel()

	validators := buildRandomValidatorSet(numVals)
	validatorSet := validset.NewValidatorSet(validators)
	snap := Snapshot{
		ValidatorSet: validatorSet,
	}

	// proposer is signer
	signerTest := validatorSet.Proposer.Address

	successionNumber, err := snap.GetSignerSuccessionNumber(signerTest)
	if err != nil {
		t.Fatalf("%s", err)
	}

	require.Equal(t, 0, successionNumber)
}

func TestGetSignerSuccessionNumber_SignerIndexIsLarger(t *testing.T) {
	t.Parallel()

	validators := buildRandomValidatorSet(numVals)

	// sort validators by address, which is what NewValidatorSet also does
	sort.Sort(validset.ValidatorsByAddress(validators))

	proposerIndex := 32
	signerIndex := 56
	// give highest ProposerPriority to a particular val, so that they become the proposer
	validators[proposerIndex].VotingPower = 200
	snap := Snapshot{
		ValidatorSet: validset.NewValidatorSet(validators),
	}

	// choose a signer at an index greater than proposer index
	signerTest := snap.ValidatorSet.Validators[signerIndex].Address

	successionNumber, err := snap.GetSignerSuccessionNumber(signerTest)
	if err != nil {
		t.Fatalf("%s", err)
	}

	require.Equal(t, signerIndex-proposerIndex, successionNumber)
}

func TestGetSignerSuccessionNumber_SignerIndexIsSmaller(t *testing.T) {
	t.Parallel()

	validators := buildRandomValidatorSet(numVals)
	proposerIndex := 98
	signerIndex := 11
	// give highest ProposerPriority to a particular val, so that they become the proposer
	validators[proposerIndex].VotingPower = 200
	snap := Snapshot{
		ValidatorSet: validset.NewValidatorSet(validators),
	}

	// choose a signer at an index greater than proposer index
	signerTest := snap.ValidatorSet.Validators[signerIndex].Address

	successionNumber, err := snap.GetSignerSuccessionNumber(signerTest)
	if err != nil {
		t.Fatalf("%s", err)
	}

	require.Equal(t, signerIndex+numVals-proposerIndex, successionNumber)
}

func TestGetSignerSuccessionNumber_ProposerNotFound(t *testing.T) {
	t.Parallel()

	validators := buildRandomValidatorSet(numVals)
	snap := Snapshot{
		ValidatorSet: validset.NewValidatorSet(validators),
	}

	require.Len(t, snap.ValidatorSet.Validators, numVals)

	dummyProposerAddress := randomAddress(toAddresses(validators)...)
	snap.ValidatorSet.Proposer = &validset.Validator{Address: dummyProposerAddress}

	// choose any signer
	signerTest := snap.ValidatorSet.Validators[3].Address

	_, err := snap.GetSignerSuccessionNumber(signerTest)
	require.NotNil(t, err)

	e, ok := err.(*UnauthorizedProposerError)
	require.True(t, ok)
	require.Equal(t, dummyProposerAddress.Bytes(), e.Proposer)
}

func TestGetSignerSuccessionNumber_SignerNotFound(t *testing.T) {
	t.Parallel()

	validators := buildRandomValidatorSet(numVals)
	snap := Snapshot{
		ValidatorSet: validset.NewValidatorSet(validators),
	}

	dummySignerAddress := randomAddress(toAddresses(validators)...)
	_, err := snap.GetSignerSuccessionNumber(dummySignerAddress)
	require.NotNil(t, err)

	e, ok := err.(*UnauthorizedSignerError)
	require.True(t, ok)

	require.Equal(t, dummySignerAddress.Bytes(), e.Signer)
}

// nolint: unparam
func buildRandomValidatorSet(numVals int) []*validset.Validator {
	validators := make([]*validset.Validator, numVals)
	valAddrs := randomAddresses(numVals)

	for i := 0; i < numVals; i++ {
		power := crand.BigInt(big.NewInt(99))
		powerN := power.Int64() + 1

		validators[i] = &validset.Validator{
			Address: valAddrs[i],
			// cannot process validators with voting power 0, hence +1
			VotingPower: powerN,
		}
	}

	// sort validators by address, which is what NewValidatorSet also does
	sort.Sort(validset.ValidatorsByAddress(validators))

	return validators
}

func randomAddress(exclude ...common.Address) common.Address {
	excl := make(map[common.Address]struct{}, len(exclude))

	for _, addr := range exclude {
		excl[addr] = struct{}{}
	}

	r := crand.NewRand()

	for {
		addr := r.Address()
		if _, ok := excl[addr]; ok {
			continue
		}

		return addr
	}
}

func randomAddresses(n int) []common.Address {
	if n <= 0 {
		return []common.Address{}
	}

	addrs := make([]common.Address, 0, n)
	addrsSet := make(map[common.Address]struct{}, n)

	var exist bool

	r := crand.NewRand()

	for {
		addr := r.Address()

		_, exist = addrsSet[addr]
		if !exist {
			addrs = append(addrs, addr)

			addrsSet[addr] = struct{}{}
		}

		if len(addrs) == n {
			return addrs
		}
	}
}

func TestRandomAddresses(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		length := rapid.IntMax(300).AsAny().Draw(t, "length").(int)

		addrs := randomAddresses(length)
		addressSet := unique.New(addrs)

		if len(addrs) != len(addressSet) {
			t.Fatalf("length of unique addresses %d, expected %d", len(addressSet), len(addrs))
		}
	})
}

func toAddresses(vals []*validset.Validator) []common.Address {
	addrs := make([]common.Address, len(vals))

	for i, val := range vals {
		addrs[i] = val.Address
	}

	return addrs
}
