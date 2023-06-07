package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"

	"github.com/althea-net/ibc-test-chain/v9/x/icaauth/types"
)

var _ types.QueryServer = Keeper{}

// InterchainAccountFromAddress fetches the interchain account associated with a given connection and owner pair
func (k Keeper) InterchainAccountFromAddress(goCtx context.Context, req *types.QueryInterchainAccountFromAddressRequest) (*types.QueryInterchainAccountFromAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not find account: %s", err)
	}

	addr, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no account found for connectionId %s portID %s",
			req.ConnectionId, portID)
	}

	return types.NewQueryInterchainAccountResponse(addr), nil
}

// InterchainAccountsWithConnection fetches all interchain accounts on a given connection
func (k Keeper) InterchainAccountsWithConnection(goCtx context.Context, req *types.QueryInterchainAccountsWithConnectionRequest) (*types.QueryInterchainAccountsWithConnectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	accounts := k.GetInterchainAccountsWithConnection(ctx, req.ConnectionId)

	return &types.QueryInterchainAccountsWithConnectionResponse{
		InterchainAccounts: accounts,
	}, nil
}

// GetInterchainAccountsWithConnection fetches all the interchain accounts on `connection` as RegisteredInterchainAccounts
func (k Keeper) GetInterchainAccountsWithConnection(ctx sdk.Context, connection string) []*icatypes.RegisteredInterchainAccount {
	var accounts []*icatypes.RegisteredInterchainAccount
	k.IterateInterchainAccounts(ctx, func(key []byte, acc icatypes.RegisteredInterchainAccount) (stop bool) {
		if acc.ConnectionId == connection {
			accounts = append(accounts, &acc)
		}
		return false
	})

	return accounts
}

// IterateInterchainAccounts iterates over all registered interchain account addresses and passes them to the provided callback `cb`
// iteration will end early if `cb` returns true
func (k Keeper) IterateInterchainAccounts(ctx sdk.Context, cb func(key []byte, acc icatypes.RegisteredInterchainAccount) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(icatypes.OwnerKeyPrefix))

	for ; iterator.Valid(); iterator.Next() {
		keySplit := strings.Split(string(iterator.Key()), "/")

		acc := icatypes.RegisteredInterchainAccount{
			ConnectionId:   keySplit[2],
			PortId:         keySplit[1],
			AccountAddress: string(iterator.Value()),
		}
		k := iterator.Key()

		if cb(k, acc) {
			return
		}
	}
}
