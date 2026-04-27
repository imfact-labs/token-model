package cmds

import (
	"fmt"
	"strings"

	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/token-model/types"
	"github.com/pkg/errors"
)

type TokenSymbolFlag struct {
	Symbol types.TokenSymbol
}

func (v *TokenSymbolFlag) UnmarshalText(b []byte) error {
	cid := types.TokenSymbol(string(b))
	if err := cid.IsValid(nil); err != nil {
		return fmt.Errorf("invalid token symbol, %q, %w", string(b), err)
	}
	v.Symbol = cid

	return nil
}

func (v *TokenSymbolFlag) String() string {
	return v.Symbol.String()
}

type AddressTokenAmountFlag struct {
	address []base.Address
	amount  []common.Big
}

func (v *AddressTokenAmountFlag) UnmarshalText(b []byte) error {
	arr := strings.SplitN(string(b), "@", -1)
	for i := range arr {
		l := strings.SplitN(arr[i], ",", 2)
		if len(l) != 2 {
			return fmt.Errorf("invalid address-amount, %q", arr[i])
		}

		add, err := base.DecodeAddress(l[0], enc)
		if err != nil {
			return err
		}
		v.address = append(v.address, add)

		b, err := common.NewBigFromString(l[1])
		if err != nil {
			return errors.Wrapf(err, "invalid big string, %q", string(l[1]))
		} else if err := b.IsValid(nil); err != nil {
			return err
		}

		v.amount = append(v.amount, b)
	}

	if len(v.amount) != len(v.address) {
		return errors.Errorf("failed to parse %s", string(b))
	}

	return nil
}

func (v *AddressTokenAmountFlag) Address() []base.Address {
	return v.address
}

func (v *AddressTokenAmountFlag) Amount() []common.Big {
	return v.amount
}
