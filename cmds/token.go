package cmds

type TokenCommand struct {
	RegisterToken RegisterModelCommand `cmd:"" name:"register-model" help:"register token to contract account"`
	Mint          MintCommand          `cmd:"" name:"mint" help:"mint token to receiver"`
	Burn          BurnCommand          `cmd:"" name:"burn" help:"burn token of target"`
	Approve       ApproveCommand       `cmd:"" name:"approve" help:"approve token to approved account"`
	Transfer      TransferCommand      `cmd:"" name:"transfer" help:"transfer token to receiver"`
	TransferFrom  TransferFromCommand  `cmd:"" name:"transfer-from" help:"transfer token to receiver from target"`
}
