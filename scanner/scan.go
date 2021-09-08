package scanner

// 从Mysql中导出数据到CSV文件。

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	//"github.com/davecgh/go-spew/spew"
)

var (
	tables         = make([]string, 0)
	outputDir      = "./"
	dataSourceName = ""
)

const (
	driverNameMysql = "mysql"
)

func init() {
	port := flag.Int("port", 12345, "the port for mysql,default:12345")
	addr := flag.String("addr", "127.0.0.1", "the address for mysql,default:127.0.0.1")
	user := flag.String("user", "admin", "the username for login mysql,default:admin")
	pwd := flag.String("pwd", "123456", "the password for login mysql by the username,default:123456")
	db := flag.String("db", "testdb", "the port for me to listen on,default:testdb")
	output := flag.String("output", "cvs", "the directory to restored,default:cvs")
	tabs := flag.String("tables", "admin_role_assoc_popedoms,admin_roles,admin_settings,admin_user_assoc_roles,admin_users,after_sale_orders,article_categories,articles,departments,device_categories,devices,doctors,dysms_settings,dysms_templates,ecode_send_logs,email_templates,ezviz_settings,feedbacks,hc_card_orders,hc_cards,hc_services,hospitals,jj_persons,mcode_send_logs,member_groups,member_hc_cards,members,qcloudsms_settings,qcloudsms_templates,qiniu_settings,qqmap_settings,reports,slider_contents,smtp_settings,sp_wenzhen_logs,station_admin,station_jk_devices,station_zn_locks,stations,suppliers,sysconfig,tycard_lq_logs,wxapp_fans,wxapp_subscribe_templates,wxapp_versions,wxapp_visit_logs,wxapps,yizhu_histories,zhuanzhen_logs", "the tables will export data, multi tables separator by comma, default:...")

	flag.Parse()

	tables = append(tables, strings.Split(*tabs, ",")...)

	dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", *user, *pwd, *addr, *port, *db)
	outputDir = fmt.Sprintf("%v/", *output)
}

func test() {

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
		querySQL(db, table, ch)
		break
	}

	//for i := 0; i < count; i++ {
	//	<-ch
	//}
	fmt.Println("Done!")
}

func querySQL(db *sql.DB, table string, ch chan bool) {
	fmt.Printf("table: %v\n", table)
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
	totalValues := make([][]string, 0)
	for rows.Next() {

		//存每一行的内容
		var s []string

		//把每行的内容添加到scanArgs，也添加到了values
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		for _, v := range values {
			s = append(s, string(v))
		}
		//spew.Printf("values: %#v\n", values)
		fmt.Printf("values: %v\n", s)
		totalValues = append(totalValues, s)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
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
