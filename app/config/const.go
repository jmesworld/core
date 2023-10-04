package config

const (
	// AccountAddressPrefix is the prefix of bech32 encoded address
	AccountAddressPrefix = "jmes"

	// AppName is the application name
	AppName = "jmes"

	// CoinType is the LUNA coin type as defined in SLIP44 (https://github.com/satoshilabs/slips/blob/master/slip-0044.md)
	CoinType = 6280

	// BondDenom staking denom
	BondDenom = "ujmes"

	// More denoms
	// Jmes      = "jmes"    // 1 (base denom unit)
	// MilliJmes = "mjmes"   // 10^-3 (milli)
	MicroJmes = BondDenom // 10^-6 (micro)
	// NanoJmes  = "njmes"   // 10^-9 (nano)

	AuthzMsgExec                        = "/cosmos.authz.v1beta1.MsgExec"
	AuthzMsgGrant                       = "/cosmos.authz.v1beta1.MsgGrant"
	AuthzMsgRevoke                      = "/cosmos.authz.v1beta1.MsgRevoke"
	BankMsgSend                         = "/cosmos.bank.v1beta1.MsgSend"
	BankMsgMultiSend                    = "/cosmos.bank.v1beta1.MsgMultiSend"
	DistrMsgSetWithdrawAddr             = "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress"
	DistrMsgWithdrawValidatorCommission = "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"
	DistrMsgFundCommunityPool           = "/cosmos.distribution.v1beta1.MsgFundCommunityPool"
	DistrMsgWithdrawDelegatorReward     = "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
	FeegrantMsgGrantAllowance           = "/cosmos.feegrant.v1beta1.MsgGrantAllowance"
	FeegrantMsgRevokeAllowance          = "/cosmos.feegrant.v1beta1.MsgRevokeAllowance"
	GovMsgVoteWeighted                  = "/cosmos.gov.v1beta1.MsgVoteWeighted"
	GovMsgSubmitProposal                = "/cosmos.gov.v1beta1.MsgSubmitProposal"
	GovMsgDeposit                       = "/cosmos.gov.v1beta1.MsgDeposit"
	GovMsgVote                          = "/cosmos.gov.v1beta1.MsgVote"
	StakingMsgEditValidator             = "/cosmos.staking.v1beta1.MsgEditValidator"
	StakingMsgDelegate                  = "/cosmos.staking.v1beta1.MsgDelegate"
	StakingMsgUndelegate                = "/cosmos.staking.v1beta1.MsgUndelegate"
	StakingMsgBeginRedelegate           = "/cosmos.staking.v1beta1.MsgBeginRedelegate"
	StakingMsgCreateValidator           = "/cosmos.staking.v1beta1.MsgCreateValidator"
	VestingMsgCreateVestingAccount      = "/cosmos.vesting.v1beta1.MsgCreateVestingAccount"
	TransferMsgTransfer                 = "/ibc.applications.transfer.v1.MsgTransfer"
	WasmMsgStoreCode                    = "/cosmwasm.wasm.v1.MsgStoreCode"
	WasmMsgExecuteContract              = "/cosmwasm.wasm.v1.MsgExecuteContract"
	WasmMsgInstantiateContract          = "/cosmwasm.wasm.v1.MsgInstantiateContract"
	WasmMsgMigrateContract              = "/cosmwasm.wasm.v1.MsgMigrateContract"

	// UpgradeName gov proposal name
	Upgrade2_2_0 = "2.2.0"
	Upgrade2_3_0 = "2.3.0"
	Upgrade2_4   = "v2.4"
	Upgrade2_5   = "v2.5"
)
