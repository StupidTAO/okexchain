package pulsarclient

import (
	"math/rand"
	"os"
	"testing"
	"time"

	appCfg "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/okex/okchain/x/backend"
	"github.com/okex/okchain/x/common"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func TestNewPulsarProducer(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	asyncErrs := make(chan error, 8)
	var err error
	defer func() {
		if len(asyncErrs) != 0 {
			err = <-asyncErrs
		}
		require.Error(t, err)
	}()

	_ = NewPulsarProducer("1.2.3.4:6650", appCfg.DefaultConfig().StreamConfig, logger, &asyncErrs)
	time.Sleep(time.Second * 4)
}

func TestSendMsg(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	asyncErrs := make(chan error, 8)
	go func() {
		for err := range asyncErrs {
			require.NoError(t, err)
		}
	}()
	mp := NewPulsarProducer("localhost:6650", appCfg.DefaultConfig().StreamConfig, logger, &asyncErrs)

	pd := PulsarData{Height: 9, matchResults: nil}
	_, err := mp.SendAllMsg(&pd, logger)
	require.NoError(t, err)
	logger.Info("send zero matchResult")

	marketIdMap["gyl_"+common.NativeToken] = int64(9999)
	results10 := make([]*backend.MatchResult, 0, 10)
	timestamp := time.Now().Unix()
	for i := 0; i < 10; i++ {
		results10 = append(results10, &backend.MatchResult{
			BlockHeight: int64(i),
			Product:     "gyl_" + common.NativeToken,
			Price:       rand.Float64(),
			Quantity:    rand.Float64(),
			Timestamp:   timestamp,
		})
	}

	pds := PulsarData{Height: 11, matchResults: results10}
	_, err = mp.SendAllMsg(&pds, logger)
	if err != nil {
		logger.Info("send 10 matchResult failed")
	}
	require.NoError(t, err)
	logger.Info("send 10 matchResult success")

	results10 = make([]*backend.MatchResult, 0, 10)
	marketIdMap[common.TestToken+common.NativeToken] = int64(10000)
	for i := 0; i < 10; i++ {
		results10 = append(results10, &backend.MatchResult{
			BlockHeight: int64(i),
			Product:     common.TestToken + common.NativeToken,
			Price:       rand.Float64(),
			Quantity:    rand.Float64(),
			Timestamp:   timestamp,
		})
	}

	pds = PulsarData{Height: 11, matchResults: results10}
	_, err = mp.SendAllMsg(&pds, logger)
	if err != nil {
		logger.Info("send 10 matchResult failed")
	}
	require.NoError(t, err)
	logger.Info("send 10 matchResult success")

	/*
	   		results1000 := make([]*backend.MatchResult, 0, 1000)
	   		timestamp = time.Now().Unix()
	   		for i:=0;i < 1000; i++ {
	   			results1000 = append(results1000, &backend.MatchResult{
	   				BlockHeight: int64(i),
	   				Product:     types.TestTokenPair,
	   				Price:       rand.Float64(),
	   				Quantity:    rand.Float64(),
	   				Timestamp:   timestamp,
	   			})
	   		}

	   		pds = PulsarData{Height:12,matchResults:results1000}
	   		_,err = mp.SendAllMsg(&pds,logger)
	   		if err != nil {
	   			logger.Info("send 1000 matchResults failed")
	   			return
	   		}
	   		logger.Info("send 1000 matchResult success")

	   /*
	   		results10000 := make([]*backend.MatchResult, 0, 10000)
	   		timestamp = time.Now().Unix()
	   		for i:=0;i < 10000; i++ {
	   			results10000 = append(results10000, &backend.MatchResult{
	   				BlockHeight: int64(i),
	   				Product:     types.TestTokenPair,
	   				Price:       rand.Float64(),
	   				Quantity:    rand.Float64(),
	   				Timestamp:   timestamp,
	   			})
	   		}

	   		pds = PulsarData{Height:13,matchResults:results10000}
	   		_,err = mp.SendMsg(&pds,logger)
	   		if err != nil {
	   			logger.Info("send 10000 matchResults failed")
	   			return
	   		}
	   		logger.Info("send 10000 matchResult success")*/
}