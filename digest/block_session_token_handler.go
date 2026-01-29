package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-token/state"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func PrepareToken(bs *currencydigest.BlockSession, st base.State) (string, []mongo.WriteModel, error) {
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

func handleTokenState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if tokenDoc, err := NewTokenDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(tokenDoc),
		}, nil
	}
}

func handleTokenBalanceState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if tokenBalanceDoc, err := NewTokenBalanceDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(tokenBalanceDoc),
		}, nil
	}
}
