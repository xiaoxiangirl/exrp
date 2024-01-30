package poa_test

import (
	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"strconv"
)

var (
	bondedStatus        = stakingtypes.Bonded
	unbondedStatus      = stakingtypes.Unbonded
	unbondingStatus     = stakingtypes.Unbonding
	zero                = sdk.ZeroInt()
	DefaultBondedTokens = sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)
)

func (s *IntegrationTestSuite) ConsumeProposalCount() string {
	s.proposalCount = s.proposalCount + 1
	return strconv.Itoa(s.proposalCount)
}

func (s *IntegrationTestSuite) GetCtx() client.Context {
	return s.network.Validators[0].ClientCtx
}

func (s *IntegrationTestSuite) RequireValidator(address string, status *stakingtypes.BondStatus, tokens *sdk.Int) {
	accAddr, _ := sdk.AccAddressFromBech32(address)
	validatorInfo := GetValidator(s.GetCtx(), sdk.ValAddress(accAddr).String())
	if validatorInfo == nil {
		s.Require().True(status == nil)
		s.Require().True(tokens == nil)
	} else {
		s.Require().Equal(*status, validatorInfo.Status)
		s.Require().Equal(*tokens, validatorInfo.Tokens)
	}
}

func (s *IntegrationTestSuite) RequireDelegation(valAddress string, delAddress string, shares sdk.Dec) {
	accAddr, _ := sdk.AccAddressFromBech32(valAddress)
	valAddr := sdk.ValAddress(accAddr).String()
	delegation := GetDelegation(s.GetCtx(), valAddr, delAddress)
	if delegation == nil {
		s.Require().Equal(sdk.ZeroDec(), shares)
	} else {
		s.Require().Equal(delegation.Shares, shares)
	}
}

func (s *IntegrationTestSuite) RequireBondBalance(address string, balance sdk.Int) {
	originalBalance := GetBalance(s.GetCtx(), address, s.cfg.BondDenom)
	expected := sdk.NewCoin(s.cfg.BondDenom, balance)
	s.Require().True(originalBalance.Equal(expected))
}

func (s *IntegrationTestSuite) RequireValidatorSet() struct {
	Contains    func(validator cryptotypes.PubKey)
	NotContains func(validator cryptotypes.PubKey)
} {
	validatorSet := GetValidatorSet(s.GetCtx())
	validatorAddresses := make([]string, 0)
	for _, val := range validatorSet.Validators {
		validatorAddresses = append(validatorAddresses, val.Address)
	}
	return struct {
		Contains    func(pubKey cryptotypes.PubKey)
		NotContains func(pubKey cryptotypes.PubKey)
	}{
		Contains: func(pubKey cryptotypes.PubKey) {
			s.Require().Contains(validatorAddresses, sdk.ConsAddress(pubKey.Address()).String())
		},
		NotContains: func(pubKey cryptotypes.PubKey) {
			s.Require().NotContains(validatorAddresses, sdk.ConsAddress(pubKey.Address()).String())
		},
	}
}
