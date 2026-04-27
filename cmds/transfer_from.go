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

type TransferFromCommand struct {
	OperationCommand
	Receiver     ccmds.AddressFlag      `arg:"" name:"receiver" help:"token receiver" required:"true"`
	TargetAmount AddressTokenAmountFlag `arg:"" name:"target" help:"target approving" required:"true"`
	receiver     base.Address
}

func (cmd *TransferFromCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *TransferFromCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	receiver, err := cmd.Receiver.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver.String())
	}
	cmd.receiver = receiver

	return nil
}

func (cmd *TransferFromCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("transfer-from operation"))
	var items []token.TransferFromItem
	for i := range cmd.TargetAmount.Address() {
		item := token.NewTransferFromItem(cmd.contract, cmd.receiver, cmd.TargetAmount.Address()[i], cmd.TargetAmount.Amount()[i])
		if err := item.IsValid(nil); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	fact := token.NewTransferFromFact(
		[]byte(cmd.Token), cmd.sender, items, cmd.Currency.CID,
	)

	op := token.NewTransferFrom(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
