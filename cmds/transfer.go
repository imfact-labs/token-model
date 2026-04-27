package cmds

import (
	"context"

	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/token-model/operation/token"
	"github.com/imfact-labs/token-model/utils"
)

type TransferCommand struct {
	OperationCommand
	ReceiverAmount AddressTokenAmountFlag `arg:"" name:"receiver-amount" help:"receiver token amount (ex: \"<address>,<amount>\") separator @" required:"true"`
}

func (cmd *TransferCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *TransferCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	return nil
}

func (cmd *TransferCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError(utils.ErrStringCreate("transfer operation"))
	var items []token.TransferItem
	for i := range cmd.ReceiverAmount.Address() {
		item := token.NewTransferItem(cmd.contract, cmd.ReceiverAmount.Address()[i], cmd.ReceiverAmount.Amount()[i])
		if err := item.IsValid(nil); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	fact := token.NewTransferFact(
		[]byte(cmd.Token), cmd.sender, items, cmd.Currency.CID,
	)

	op := token.NewTransfer(fact)
	if err := op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
