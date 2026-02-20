package cmds

import (
	"context"

	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/token-model/operation/token"
	"github.com/imfact-labs/token-model/utils"
	"github.com/pkg/errors"
)

type ApproveCommand struct {
	OperationCommand
	Approved1 ccmds.AddressFlag `arg:"" name:"approved" help:"approved account" required:"true"`
	Approved2 ccmds.AddressFlag `arg:"" name:"approved" help:"approved account" required:"true"`
	Amount    ccmds.BigFlag     `arg:"" name:"amount" help:"amount to approve" required:"true"`
	approved1 base.Address
	approved2 base.Address
}

func (cmd *ApproveCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	ccmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *ApproveCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	approved, err := cmd.Approved1.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid approved format, %q", cmd.Approved1.String())
	}
	cmd.approved1 = approved

	approved, err = cmd.Approved2.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid approved format, %q", cmd.Approved2.String())
	}
	cmd.approved2 = approved

	return nil
}

func (cmd *ApproveCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("approve operation"))

	item1 := token.NewApproveItem(cmd.contract,
		cmd.approved1, cmd.Amount.Big, cmd.Currency.CID)

	item2 := token.NewApproveItem(cmd.contract,
		cmd.approved2, cmd.Amount.Big, cmd.Currency.CID)

	fact := token.NewApproveFact(
		[]byte(cmd.Token), cmd.sender, []token.ApproveItem{item1, item2},
	)

	op := token.NewApprove(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
