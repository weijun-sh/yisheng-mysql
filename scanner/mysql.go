package scanner

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"gopkg.in/mgo.v2/bson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/weijun-sh/yisheng-mysql/params"
	"github.com/weijun-sh/yisheng-mysql/mongodb"
)

var (
	tables         = make([]string, 0)
	outputDir      = "./"
)

const (
	driverNameMysql = "mysql"
)

func init() {
	tabs := flag.String("tables", "admin_role_assoc_popedoms,admin_roles,admin_settings,admin_user_assoc_roles,admin_users,after_sale_orders,article_categories,articles,departments,device_categories,devices,doctors,dysms_settings,dysms_templates,ecode_send_logs,email_templates,ezviz_settings,feedbacks,hc_card_orders,hc_cards,hc_services,hospitals,jj_persons,mcode_send_logs,member_groups,member_hc_cards,members,qcloudsms_settings,qcloudsms_templates,qiniu_settings,qqmap_settings,reports,slider_contents,smtp_settings,sp_wenzhen_logs,station_admin,station_jk_devices,station_zn_locks,stations,suppliers,sysconfig,tycard_lq_logs,wxapp_fans,wxapp_subscribe_templates,wxapp_versions,wxapp_visit_logs,wxapps,yizhu_histories,zhuanzhen_logs", "the tables will export data, multi tables separator by comma, default:...")
	//tabs := flag.String("tables", "hc_card_orders,jj_persons", "the tables will export data, multi tables separator by comma, default:...")
	flag.Parse()

	tables = append(tables, strings.Split(*tabs, ",")...)

}

func run() {
	config := params.GetMysqldbConfig()
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", config.UserName, config.Password, config.DBURL, config.DBName)

	count := len(tables)
	ch := make(chan bool, count)

	db, err := sql.Open(driverNameMysql, dataSourceName)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	for _, table := range tables {
		go querySQL(db, table, ch)
		//break
	}

	for i := 0; i < count; i++ {
		<-ch
	}
	fmt.Println("Done!")
}

func querySQL(db *sql.DB, table string, ch chan bool) {
	fmt.Printf("\ntable: %v\n", table)
	rows, err := db.Query(fmt.Sprintf("SELECT * from %s", table))

	if err != nil {
		panic(err)
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("columns: %v %v %v\n", columns[0], columns[1], columns[2])

	//values：一行的所有值,把每一行的各个字段放到values中，values长度==列数
	values := make([]sql.RawBytes, len(columns))
	//print(len(values))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	//fmt.Printf("scanArgs: %v\n", scanArgs)
	//存所有行的内容totalValues
	totalValues := make([]interface{}, 0)
	for rows.Next() {
		ms := bson.M{}

		//把每行的内容添加到scanArgs，也添加到了values
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		for i, v := range values {
			if columns[i] == "id" {
				columns[i] = "_id"
			}
			//vstring := strings.Replace(string(v),"\\","", -1)
			//vstring = strings.Replace(vstring,"\"{","{", -1)
			//vstring = strings.Replace(vstring,"}\n\"","}", -1)
			ms[columns[i]] = string(v)
			fmt.Printf("v: %v\n", string(v))
		}
		//spew.Printf("values: %#v\n", values)
		fmt.Printf("values: %v\n", ms)
		totalValues = append(totalValues, ms)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
	mongodb.Insert(table, totalValues)
	//writeToCSV(outputDir+table+".csv", columns, totalValues)
	ch <- true
}

//writeToCSV
func writeToCSV(file string, columns []string, totalValues [][]string) {
	f, err := os.Create(file)
	// fmt.Println(columns)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	//f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	for i, row := range totalValues {
		//第一次写列名+第一行数据
		if i == 0 {
			w.Write(columns)
			w.Write(row)
		} else {
			w.Write(row)
		}
	}
	w.Flush()
	fmt.Println("处理完毕：", file)
}
