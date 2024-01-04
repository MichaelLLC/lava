package ante

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/gogoproto/proto"
	"github.com/lavanet/lava/app"
	testkeeper "github.com/lavanet/lava/testutil/keeper"
	plantypes "github.com/lavanet/lava/x/plans/types"
	spectypes "github.com/lavanet/lava/x/spec/types"
	subsciptiontypes "github.com/lavanet/lava/x/subscription/types"
	"github.com/stretchr/testify/require"
)

func TestNewExpeditedProposalFilterAnteDecorator(t *testing.T) {
	account := authtypes.NewModuleAddressOrBech32Address("cosmos1qypqxpq9qcrsszgjx3ysxf7j8xq9q9qyq9q9q9")
	govAuthority := authtypes.NewModuleAddress(govtypes.ModuleName)
	whitelistedContent := plantypes.PlansDelProposal{}

	tests := []struct {
		name       string
		theMsg     func() sdk.Msg
		shouldFail bool
	}{
		{
			name: "should not fail if the message is in the whitelist, expedited",
			theMsg: func() sdk.Msg {
				proposal, err := v1.NewMsgSubmitProposal(
					[]sdk.Msg{
						&banktypes.MsgSend{},
					},
					sdk.NewCoins(sdk.NewCoin("lava", sdk.NewInt(100))),
					"cosmos1qypqxpq9qcrsszgjx3ysxf7j8xq9q9qyq9q9q9",
					"metadata",
					"title",
					"summary",
					true,
				)
				require.NoError(t, err)

				return proposal
			},
			shouldFail: false,
		},
		{
			name: "should fail if none of the messages are in the whitelist, expedited",
			theMsg: func() sdk.Msg {
				proposal, err := v1.NewMsgSubmitProposal(
					[]sdk.Msg{
						&subsciptiontypes.MsgAutoRenewal{},
					},
					sdk.NewCoins(sdk.NewCoin("lava", sdk.NewInt(100))),
					"cosmos1qypqxpq9qcrsszgjx3ysxf7j8xq9q9qyq9q9q9",
					"metadata",
					"title",
					"summary",
					true,
				)
				require.NoError(t, err)

				return proposal
			},
			shouldFail: true,
		},
		{
			name: "a new msg exec proposal with an expedited message with a whitelisted message should not fail",
			theMsg: func() sdk.Msg {
				proposal, err := v1.NewMsgSubmitProposal(
					[]sdk.Msg{
						&banktypes.MsgSend{},
					},
					sdk.NewCoins(sdk.NewCoin("lava", sdk.NewInt(100))),
					"cosmos1qypqxpq9qcrsszgjx3ysxf7j8xq9q9qyq9q9q9",
					"metadata",
					"title",
					"summary",
					true,
				)
				require.NoError(t, err)

				authMsg := authz.NewMsgExec(
					account,
					[]sdk.Msg{proposal},
				)

				return &authMsg
			},
			shouldFail: false,
		},
		{
			name: "a new msg exec proposal with an expedited message with a non-whitelisted message should fail",
			theMsg: func() sdk.Msg {
				proposal, err := v1.NewMsgSubmitProposal(
					[]sdk.Msg{
						&subsciptiontypes.MsgAutoRenewal{},
					},
					sdk.NewCoins(sdk.NewCoin("lava", sdk.NewInt(100))),
					"cosmos1qypqxpq9qcrsszgjx3ysxf7j8xq9q9qyq9q9q9",
					"metadata",
					"title",
					"summary",
					true,
				)
				require.NoError(t, err)

				authMsg := authz.NewMsgExec(
					account,
					[]sdk.Msg{proposal},
				)

				return &authMsg
			},
			shouldFail: true,
		},
		{
			name: "a v1 proposal that contains a legacy proposal with whitelisted content should not fail, expedited",
			theMsg: func() sdk.Msg {
				anyContent, err := codectypes.NewAnyWithValue(&whitelistedContent)
				require.NoError(t, err)
				submitProposal := v1.NewMsgExecLegacyContent(
					anyContent,
					govAuthority.String(),
				)

				proposal, err := v1.NewMsgSubmitProposal(
					[]sdk.Msg{
						submitProposal,
					},
					sdk.NewCoins(sdk.NewCoin("lava", sdk.NewInt(100))),
					"cosmos1qypqxpq9qcrsszgjx3ysxf7j8xq9q9qyq9q9q9",
					"metadata",
					"title",
					"summary",
					true,
				)
				require.NoError(t, err)

				return proposal
			},
			shouldFail: false,
		},
		{
			name: "a v1 proposal that contains a legacy proposal with not whitelisted content should fail, expedited",
			theMsg: func() sdk.Msg {
				anyContent, err := codectypes.NewAnyWithValue(&plantypes.PlansAddProposal{})
				require.NoError(t, err)
				submitProposal := v1.NewMsgExecLegacyContent(
					anyContent,
					govAuthority.String(),
				)

				proposal, err := v1.NewMsgSubmitProposal(
					[]sdk.Msg{
						submitProposal,
					},
					sdk.NewCoins(sdk.NewCoin("lava", sdk.NewInt(100))),
					"cosmos1qypqxpq9qcrsszgjx3ysxf7j8xq9q9qyq9q9q9",
					"metadata",
					"title",
					"summary",
					true,
				)
				require.NoError(t, err)

				return proposal
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			k, ctx := testkeeper.SpecKeeper(t)
			params := spectypes.DefaultParams()
			params.WhitelistedExpeditedMsgs = []string{
				proto.MessageName(&banktypes.MsgSend{}),
				proto.MessageName(&whitelistedContent),
			} // we whitelist MsgSend proposal

			k.SetParams(ctx, params)

			encodingConfig := app.MakeEncodingConfig()
			clientCtx := client.Context{}.
				WithTxConfig(encodingConfig.TxConfig)

			txBuilder := clientCtx.TxConfig.NewTxBuilder()
			err := txBuilder.SetMsgs(tt.theMsg())
			require.NoError(t, err)

			tx := txBuilder.GetTx()
			anteHandler := NewExpeditedProposalFilterAnteDecorator(k)

			_, err = anteHandler.AnteHandle(ctx, tx, false, func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
				return ctx, nil
			})
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorContains(t, err, "not allowed to be expedited")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
