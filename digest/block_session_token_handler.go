package digest

import (
	cdigest "github.com/imfact-labs/currency-model/digest"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/token-model/state"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func PrepareToken(bs *cdigest.BlockSession, st base.State) (string, []mongo.WriteModel, error) {
	switch {
	case state.IsStateDesignKey(st.Key()):
		j, err := handleTokenState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNameToken, j, nil
	case state.IsStateTokenBalanceKey(st.Key()):
		j, err := handleTokenBalanceState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNameTokenBalance, j, nil
	}

	return "", nil, nil
}

func handleTokenState(bs *cdigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if tokenDoc, err := NewTokenDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(tokenDoc),
		}, nil
	}
}

func handleTokenBalanceState(bs *cdigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if tokenBalanceDoc, err := NewTokenBalanceDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(tokenBalanceDoc),
		}, nil
	}
}
