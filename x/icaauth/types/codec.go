package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	// this line is used by starport scaffolding # 1
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	erc20types "github.com/Canto-Network/Canto/v5/x/erc20/types"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	microtxtypes "github.com/althea-net/althea-L1/x/microtx/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	registerLocalAmino(cdc)
	registerForeignAminos(cdc)
}

// Registers the icaauth amino message types
func registerLocalAmino(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterAccount{}, "icaauth/MsgRegisterAccount", nil)
	cdc.RegisterConcrete(&MsgSubmitTx{}, "icaauth/MsgSubmitTx", nil)
}

// Registers any foreign module amino message types used for interchain accounts
// Note: Include any modules here the chain should control which are not running on this chain
func registerForeignAminos(cdc *codec.LegacyAmino) {
	gravitytypes.RegisterCodec(cdc)
	microtxtypes.RegisterCodec(cdc)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registerLocalInterfaces(registry)
	registerForeignInterfaces(registry)
}

// Registers the icaauth protos
func registerLocalInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegisterAccount{},
		&MsgSubmitTx{},
	)

}

// Registers any foreign module protos used for interchain accounts
// Note: Include any modules here the chain should control which are not running on this chain
func registerForeignInterfaces(registry cdctypes.InterfaceRegistry) {
	gravitytypes.RegisterInterfaces(registry)
	microtxtypes.RegisterInterfaces(registry)
	evmtypes.RegisterInterfaces(registry)
	erc20types.RegisterInterfaces(registry)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
