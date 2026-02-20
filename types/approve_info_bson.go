package types

import (
	"github.com/imfact-labs/currency-model/utils/bsonenc"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/token-model/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (a ApproveInfo) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":   a.Hint().String(),
			"account": a.account,
			"amount":  a.amount.String(),
		},
	)
}

type ApproveInfoBSONUnmarshaler struct {
	Hint    string `bson:"_hint"`
	Account string `bson:"account"`
	Amount  string `bson:"amount"`
}

func (a *ApproveInfo) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*a))

	var u ApproveInfoBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return a.unpack(enc, ht, u.Account, u.Amount)
}
