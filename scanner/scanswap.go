package scanner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/anyswap/CrossChain-Bridge/cmd/utils"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/rpc/client"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/urfave/cli/v2"

	ethclient "github.com/jowenshaw/gethclient"
	"github.com/jowenshaw/gethclient/common"
	"github.com/jowenshaw/gethclient/types"

	"github.com/weijun-sh/gethscan/params"
	"github.com/weijun-sh/gethscan/tools"
	"github.com/weijun-sh/gethscan/mongodb"
)

var (
	// ScanSwapCommand scan swaps on eth like blockchain
	ScanSwapCommand = &cli.Command{
		Action:    scanSwap,
		Name:      "scanswap",
		Usage:     "scan cross chain swaps",
		ArgsUsage: " ",
		Description: `
scan cross chain swaps
`,
	}
)

func scanSwap(ctx *cli.Context) error {
	utils.SetLogger(ctx)
	params.LoadConfig(utils.GetConfigFilePath(ctx))

       //mongo
       dbConfig := params.GetMongodbConfig()
       chain = dbConfig.BlockChain
       InitMongodb()

	scanner.initClient()
	scanner.run()
	return nil
}

func addMongodbSwapPost(swap *swapPost) {
       ms := &mongodb.MgoSwap{
               Id:         swap.txid,
               Txid:       swap.txid,
               PairID:     swap.pairID,
               RpcMethod:  swap.rpcMethod,
               SwapServer: swap.swapServer,
               Chain:      chain,
               Timestamp:  uint64(time.Now().Unix()),
       }
       mongodb.AddSwap(ms, false)
}

func addMongodbSwapPendingPost(swap *swapPost) {
       ms := &mongodb.MgoSwap{
               Id:         swap.txid,
               Txid:       swap.txid,
               PairID:     swap.pairID,
               RpcMethod:  swap.rpcMethod,
               SwapServer: swap.swapServer,
               Chain:      chain,
               Timestamp:  uint64(time.Now().Unix()),
       }
       mongodb.AddSwapPending(ms, false)
 }

// InitMongodb init mongodb by config
func InitMongodb() {
       log.Info("InitMongodb")
       dbConfig := params.GetMongodbConfig()
       mongodb.MongoServerInit([]string{dbConfig.DBURL}, dbConfig.DBName, dbConfig.UserName, dbConfig.Password)
}

