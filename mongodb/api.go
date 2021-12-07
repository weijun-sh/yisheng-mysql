package mongodb

import (
	"errors"
	"fmt"
	//"strings"
	//"github.com/weijun-sh/yisheng-mysql/log"
	"gopkg.in/mgo.v2/bson"
	//"github.com/davecgh/go-spew/spew"
	//"encoding/json"
)

func Insert(table string, docs []interface{}) (err error) {
	c := deinintCollections(table)
	err = c.Insert(docs...)
	if err != nil {
		fmt.Printf("Insert, docs: %v, err: %v\n", docs, err)
	}
	return nil
}

// card
func Find_hc_card_orders() ([]*hc_card_orders_Config, error) {
	c := deinintCollections("hc_card_orders")
        var res []*hc_card_orders_Config
        err := c.Find(nil).All(&res)
        if err != nil {
                return nil, errors.New("mgo hc_card_orders find failed")
        }
	return res, nil
}

func Find_hc_card_orders_No(orderNo string) (*hc_card_orders_Config, error) {
	c := deinintCollections("hc_card_orders")
        var res hc_card_orders_Config
        err := c.Find(bson.M{"orderNo":orderNo}).One(&res)
        if err != nil {
                return nil, errors.New("mgo hc_card_orders find failed")
        }
	return &res, nil
}

func Find_jj_persons(memberid string) (*jj_persons_Config, error) {
	c := deinintCollections("jj_persons")
        var res *jj_persons_Config
        err := c.Find(bson.M{"memberId":memberid}).One(&res)
        if err != nil {
                return nil, errors.New("mgo jj_persons find failed")
        }
        return res, nil
}

func Find_jj_persons_id(id string) (*jj_persons_Config, error) {
	c := deinintCollections("jj_persons")
        var res *jj_persons_Config
        err := c.Find(bson.M{"_id":id}).One(&res)
        if err != nil {
                return nil, errors.New("mgo jj_persons find failed")
        }
        return res, nil
}

// report
func Find_reports() ([]*reports_Config, error) {
	c := deinintCollections("reports")
        var res []*reports_Config
        err := c.Find(nil).All(&res)
        if err != nil {
                return nil, errors.New("mgo reports find failed")
        }
        return res, nil
}

func Find_member_hc_cards() ([]*member_hc_cards_Config, error) {
	c := deinintCollections("member_hc_cards")
        var res []*member_hc_cards_Config
        err := c.Find(nil).All(&res)
        if err != nil {
                return nil, errors.New("mgo member_hc_cards find failed")
        }
        return res, nil
}

func Find_member_hc_cards_hxcode(hxcode string) (*member_hc_cards_Config, error) {
	c := deinintCollections("member_hc_cards")
        var res member_hc_cards_Config
        err := c.Find(bson.M{"hxcode":hxcode}).One(&res)
        if err != nil {
                return nil, errors.New("mgo member_hc_cards find failed")
        }
        return &res, nil
}
// data
type HcCardOrdersConfig struct {
	Id string `bson:"_id"`
	MemberId string
	Sex string
	Realname string
	Birthday string
	Cellphones []string
	OrderNo string
	CardType []string
	Ctime string
	Hash string
}

type ReportsConfig struct {
	Id string `bson:"_id"`
	MemberId string `bson:"memberId"`
	Sex string `bson:"sex"`
	Realname string `bson:"realname"`
	Birthday string `bson:"birthday"`
	Cellphones []string `bson:"cellphones"`
	YzContents string `bson:"yzContents"`
	CardType []string `bson:"cardType"`
	ReportTime string `bson:"reportTime"`
	Hash string
}

// tables
type hc_card_orders_Config struct {
	Id string `bson:"_id"`
	MemberId string `bson:"memberId"`
	OrderNo string `bson:"orderNo"`
	Items string `bson:"items"`
	Ctime string `bson:"ctime"`
}

type ItemsConfig struct {
	Name string `bson:"name"`// (cardType)
}

type jj_persons_Config struct {
	Id string `bson:"_id"`
	MemberId string `bson:"memberId"`
	Sex string `bson:"sex"`
	Realname string `bson:"realname"`
	Birthday string `bson:"birthday"`
	UdeskCustomerInfo string `bson:"udeskCustomerInfo"`
}

type UdeskCustomerInfoConfig struct {
	CellPhones []CellPhonesConfig `bson:"cellphones"`
}

type CellPhonesConfig struct {
	Content string `bson:"content"`
}

type reports_Config struct {
	Id string `bson:"_id"`
	MemberId string `bson:"memberId"`
	JjPersonId string `bson:"jjPersonId"`
	YzContents string `bson:"yzContents"`
	Hxcode string `bson:"hxcode"`
	ReportTime string `bson:"reportTime"`
}

type member_hc_cards_Config struct {
	Hxcode string `bson:"hxcode"`
	OrderNo string `bson:"orderNo"`
}

