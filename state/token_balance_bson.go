package state

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/utils/bsonenc"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/token-model/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s TokenBalanceStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  s.Hint().String(),
			"amount": s.Amount,
		},
	)
}

type TokenBalanceStateValueBSONUnmarshaler struct {
	Hint   string `bson:"_hint"`
	Amount string `bson:"amount"`
}

func (s *TokenBalanceStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*s))

	var u TokenBalanceStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(ht)

	big, err := common.NewBigFromString(u.Amount)
	if err != nil {
		return e.Wrap(err)
	}
	s.Amount = big

	return nil
}
