package scanner

import (
	//"time"

	"github.com/weijun-sh/yisheng-mysql/cmd/utils"
	"github.com/weijun-sh/yisheng-mysql/log"
	"github.com/urfave/cli/v2"

	"github.com/weijun-sh/yisheng-mysql/params"
	"github.com/weijun-sh/yisheng-mysql/mongodb"
)

var (
	// ScanSwapCommand scan swaps on eth like blockchain
	ScanSwapCommand = &cli.Command{
		Action:    scanSwap,
		Name:      "mysql",
		Usage:     "scan cross chain swaps",
		ArgsUsage: " ",
		Description: `
scan cross chain swaps
`,
		Flags: []cli.Flag{
                        utils.ConfigFileFlag,
                },
	}
)

func scanSwap(ctx *cli.Context) error {
	utils.SetLogger(ctx)
	params.LoadConfig(utils.GetConfigFilePath(ctx))

       //mongo
       params.GetMongodbConfig()
       InitMongodb()

	run()
	return nil
}

func add_admin_role_assoc_popedoms(arap *mongodb.Struct_admin_role_assoc_popedoms) {
	//arap.Timestamp = uint64(time.Now().Unix())
       //mongodb.Add_admin_role_assoc_popedoms(arap, false)
}

// InitMongodb init mongodb by config
func InitMongodb() {
       log.Info("InitMongodb")
       dbConfig := params.GetMongodbConfig()
       mongodb.MongoServerInit([]string{dbConfig.DBURL}, dbConfig.DBName, dbConfig.UserName, dbConfig.Password)
}

