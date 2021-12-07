package scanner

import (
	"crypto/md5"
	"fmt"
	"encoding/hex"
	"encoding/json"
	//"time"

	"github.com/weijun-sh/yisheng-mysql/cmd/utils"
	"github.com/urfave/cli/v2"

	"github.com/weijun-sh/yisheng-mysql/params"
	"github.com/weijun-sh/yisheng-mysql/mongodb"
	"github.com/davecgh/go-spew/spew"
)

var (
	// PrepareChainData prepare chain data from mongo
	PrepareChainData = &cli.Command{
		Action:    prepareChainData,
		Name:      "chaindata",
		Usage:     "prepare chain data from mongo",
		ArgsUsage: " ",
		Description: `
prepare chain data from mongo
`,
		Flags: []cli.Flag{
                        utils.ConfigFileFlag,
                },
	}
)

func prepareChainData(ctx *cli.Context) error {
	utils.SetLogger(ctx)
	params.LoadConfig(utils.GetConfigFilePath(ctx))

       //mongo
       params.GetMongodbConfig()
       InitMongodb()

	chainData()
	return nil
}

func chainData() {
	build_hc_card_orders()
	build_reports()
}

func build_reports() {
	// get from reports
	hc, errh := mongodb.Find_reports()
	if errh != nil {
		fmt.Printf("Find_reports err: %v\n", errh)
		return
	}
	var hcCardOrders []interface{} = make([]interface{}, len(hc))
	for _, card := range hc {
		//spew.Printf("card: %#v\n", card)
		var hcCardOrder mongodb.ReportsConfig
		hcCardOrder.Id = card.Id
		hcCardOrder.MemberId = card.MemberId
		//spew.Printf("hcCardOrder: %#v\n", hcCardOrder)
		// get from jj_persons
		jj, errj := mongodb.Find_jj_persons(card.MemberId)
		if errj != nil {
			fmt.Printf("Find_jj_persons err: %v\n", errj)
		} else {
			hcCardOrder.Sex = jj.Sex
			hcCardOrder.Birthday = jj.Birthday
			if jj.UdeskCustomerInfo != "" {
				var phones *mongodb.UdeskCustomerInfoConfig
				var content []*mongodb.CellPhonesConfig
				_ = json.Unmarshal([]byte(jj.UdeskCustomerInfo), &phones)
				inputByte, _ := json.Marshal(phones.CellPhones)
				_ = json.Unmarshal([]byte(inputByte), &content)
				var cellphones []string
				for _, phone := range content {
					cellphones = append(cellphones, md5sum(phone.Content))
				}
				hcCardOrder.Cellphones = cellphones
			}
		}
		// get from jj_persons _id
		jj, errj = mongodb.Find_jj_persons_id(card.JjPersonId)
		if errj != nil {
			fmt.Printf("Find_jj_persons_id err: %v\n", errj)
		} else {
			hcCardOrder.Realname = md5sum(jj.Realname)
		}

		hcCardOrder.YzContents = card.YzContents

		names := getCardType(card.Hxcode)
		hcCardOrder.CardType = names

		hcCardOrder.ReportTime = card.ReportTime
		spew.Printf("\nhcCardOrder: %#v\n", hcCardOrder)
		hcCardOrders = append(hcCardOrders, hcCardOrder)
	}
	mongodb.Insert("chain_reports", hcCardOrders)
}

func md5sum(str string) string {
	m5 := md5.New()
	_,err := m5.Write([]byte(str))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(m5.Sum(nil))
}

func getCardType(hxcode string) []string {
	var names []string
	// get from member_hc_cards
	jj, errj := mongodb.Find_member_hc_cards_hxcode(hxcode)
	if errj != nil {
		fmt.Printf("Find_member_hc_cards err: %v\n", errj)
	} else {
		// get from hc_card_orders
		item, _ := mongodb.Find_hc_card_orders_No(jj.OrderNo)
		var items []*mongodb.ItemsConfig
		_ = json.Unmarshal([]byte(item.Items), &items)
		for _, name := range items {
			names = append(names, name.Name)
		}
	}
	return names
}

func build_hc_card_orders() {
	// get from hc_card_orders
	hc, errh := mongodb.Find_hc_card_orders()
	if errh != nil {
		fmt.Printf("Find_hc_card_orders err: %v\n", errh)
		return
	}
	var hcCardOrders []interface{} = make([]interface{}, len(hc))
	for _, card := range hc {
		var hcCardOrder mongodb.HcCardOrdersConfig
		hcCardOrder.Id = card.Id
		hcCardOrder.MemberId = card.MemberId
		// get from jj_persons
		jj, errj := mongodb.Find_jj_persons(card.MemberId)
		if errj != nil {
			fmt.Printf("Find_jj_persons err: %v\n", errj)
		} else {
			hcCardOrder.Sex = jj.Sex
			hcCardOrder.Realname = md5sum(jj.Realname)
			hcCardOrder.Birthday = jj.Birthday
			if jj.UdeskCustomerInfo != "" {
				var phones *mongodb.UdeskCustomerInfoConfig
				var content []*mongodb.CellPhonesConfig
				_ = json.Unmarshal([]byte(jj.UdeskCustomerInfo), &phones)
				inputByte, _ := json.Marshal(phones.CellPhones)
				_ = json.Unmarshal([]byte(inputByte), &content)
				var cellphones []string
				for _, phone := range content {
					cellphones = append(cellphones, md5sum(phone.Content))
				}
				hcCardOrder.Cellphones = cellphones
			}
		}

		hcCardOrder.OrderNo = card.OrderNo
		var items []*mongodb.ItemsConfig
		_ = json.Unmarshal([]byte(card.Items), &items)
		var names []string
		for _, name := range items {
			names = append(names, name.Name)
		}
		hcCardOrder.CardType = names
		hcCardOrder.Ctime = card.Ctime
		spew.Printf("\n%#v\n", hcCardOrder)
		hcCardOrders = append(hcCardOrders, hcCardOrder)
	}
	mongodb.Insert("chain_hcCardOrders", hcCardOrders)
}

