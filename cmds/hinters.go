package cmds

import (
	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/token-model/operation/token"
	"github.com/imfact-labs/token-model/state"
	"github.com/imfact-labs/token-model/types"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: types.ApproveBoxHint, Instance: types.ApproveBox{}},
	{Hint: types.ApproveInfoHint, Instance: types.ApproveInfo{}},
	{Hint: types.PolicyHint, Instance: types.Policy{}},
	{Hint: types.DesignHint, Instance: types.Design{}},

	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.TokenBalanceStateValueHint, Instance: state.TokenBalanceStateValue{}},

	{Hint: token.RegisterModelHint, Instance: token.RegisterModel{}},
	{Hint: token.MintHint, Instance: token.Mint{}},
	{Hint: token.BurnHint, Instance: token.Burn{}},
	{Hint: token.ApproveHint, Instance: token.Approve{}},
	{Hint: token.ApproveItemHint, Instance: token.ApproveItem{}},
	{Hint: token.TransferHint, Instance: token.Transfer{}},
	{Hint: token.TransferItemHint, Instance: token.TransferItem{}},
	{Hint: token.TransferFromHint, Instance: token.TransferFrom{}},
	{Hint: token.TransferFromItemHint, Instance: token.TransferFromItem{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: token.RegisterModelFactHint, Instance: token.RegisterModelFact{}},
	{Hint: token.MintFactHint, Instance: token.MintFact{}},
	{Hint: token.BurnFactHint, Instance: token.BurnFact{}},
	{Hint: token.ApproveFactHint, Instance: token.ApproveFact{}},
	{Hint: token.TransferFactHint, Instance: token.TransferFact{}},
	{Hint: token.TransferFromFactHint, Instance: token.TransferFromFact{}},
}

func init() {
	Hinters = append(Hinters, ccmds.Hinters...)
	Hinters = append(Hinters, AddedHinters...)

	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, ccmds.SupportedProposalOperationFactHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, AddedSupportedHinters...)
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range Hinters {
		if err := encs.AddDetail(Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}
