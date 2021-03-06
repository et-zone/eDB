package eDB

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var eClient *Client
var table_config string = ""

const (
	TABLE_FIELDS = "TABLEFIELDS"
)

type Client struct {
	tableRows   map[string]string //key=tableName ,val=rowSql
	db          *gorm.DB
	tableFields map[string]string //key=tableName, val=fields
}

type EConfig struct {
	UserName  string `json:"userName"`
	PassWord  string `json:"passWord"`
	Addr      string `json:"addr"`
	Port      int    `json:"port"`
	DB        string `json:"db"`
	TableFile string `json:"tablefile"` // 表的配置文件
}

//初始化gorm
func initOrm(cfg *EConfig) *gorm.DB {

	db, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		cfg.UserName,
		cfg.PassWord,
		cfg.Addr,
		cfg.Port,
		cfg.DB,
	))
	db.SingularTable(true)
	if err != nil {
		panic("orm init error")
	}
	return db

}

func InitClient(config *EConfig) *Client {
	eClient = &Client{
		tableRows:   map[string]string{}, //values
		db:          initOrm(config),
		tableFields: map[string]string{}, //字段名称
	}
	if config.TableFile != "" {
		table_config = config.TableFile
		tblist, err := getTableFieldByJsonFile()
		if err != nil {
			panic(err.Error())
		}
		for key, _ := range tblist {
			eClient.initTableField(key, tblist[key]...)
		}

	}
	return eClient

}

func (cli *Client) Clear(tableName ...string) {
	if len(tableName) == 0 {
		for key, _ := range cli.tableRows {
			cli.tableRows[key] = ""
		}
	} else {
		for _, key := range tableName {
			cli.tableRows[key] = ""
		}
	}

}

func (cli *Client) AddRow(tableName string, row *Row) {
	tmpStr := "("
	for i := 0; i < row.GetSize(); i++ {

		switch reflect.TypeOf(row.GetColumnValues(i)).String() {
		case "string":
			val := fmt.Sprintf("%v", row.GetColumnValues(i))
			if val == "NULL" {
				tmpStr = tmpStr + val
			} else {
				val = strings.ReplaceAll(val, "'", "\\'")
				val = strings.ReplaceAll(val, "\"", "\\\"")
				tmpStr = tmpStr + "'" + val + "'"
			}
		default:
			tmpStr = tmpStr + fmt.Sprintf("%v", row.GetColumnValues(i))

		}
		// fmt.Println(reflect.TypeOf(row.getColumnValues(i)).String())
		if i != row.GetSize()-1 {
			tmpStr += ", "
		}
	}
	tmpStr = tmpStr[:len(tmpStr)] + ")"
	if len(cli.tableRows[tableName]) != 0 {
		tmpStr = "," + tmpStr
	}
	cli.tableRows[tableName] = cli.tableRows[tableName] + tmpStr

}

func (cli *Client) isTableIn(tableName string) bool {
	_, isin := cli.tableFields[tableName]
	return isin

}

func (cli *Client) GetTableNames() []string {
	tbnames := []string{}
	for tbname, _ := range cli.tableFields {
		tbnames = append(tbnames, tbname)
	}
	return tbnames

}

//flush DB All table
func (cli *Client) FlushAll() (err error) {
	sql := ""
	defer func() {
		err := recover()
		if err != nil {
			log.Println("errSql ====> ", sql)

		}
	}()
	for key, _ := range cli.tableRows {

		sql = "INSERT INTO " + key + cli.tableFields[key] + " VALUES " + cli.tableRows[key]
		_, err = cli.db.DB().Exec(sql)

		if err != nil {
			log.Println("FlushAll err  tableName="+key+" sql= ", sql, " err= ", err.Error())
		}
		// d, _ := ret.RowsAffected()
		cli.Clear(key)

	}
	return

}

//事务插入，必须制定哪些数据表,默认不执行插入
func (cli *Client) FlushTx(tableName ...string) (err error) {
	if len(tableName) == 0 {
		return nil
	}

	sql := ""
	defer func() {
		err := recover()
		if err != nil {
			log.Println("errSql ====> ", sql)

		}
	}()

	tx, err := cli.db.DB().Begin()
	if err != nil {
		panic(err.Error())
	}

	for _, key := range tableName {

		sql = "INSERT INTO " + key + cli.tableFields[key] + " VALUES " + cli.tableRows[key]

		_, err = tx.Exec(sql)

		if err != nil {
			log.Println("tableName="+key+" err= ", err.Error())
			tx.Rollback()
			return
		}
		// d, _ := ret.RowsAffected()
	}
	for _, key := range tableName {
		cli.Clear(key)
	}

	err = tx.Commit()
	return

}

//设置字段名称
func (cli *Client) initTableField(tableName string, fields ...string) {
	_, isok := cli.tableFields[tableName]
	if isok {
		return
	}
	if len(fields) == 0 {
		return
	}
	tmpFields := " ("
	for i, field := range fields {
		tmpFields += field
		if i != len(fields)-1 {
			tmpFields += ", "
		}
	}
	tmpFields += ") "
	cli.tableFields[tableName] = tmpFields

}

func getTableFieldByJsonFile() (map[string][]string, error) {
	if table_config == "" {
		return map[string][]string{}, nil
	}
	b, err := ioutil.ReadFile(table_config)
	if err != nil {
		fmt.Println("warning not found file ")
	}

	mafields := map[string][]string{}
	err = json.Unmarshal(b, &mafields)
	if err != nil {
		panic(err.Error())
	}

	return mafields, err
}

func (cli *Client) String() string {
	b, _ := json.Marshal(cli.tableFields)
	return string(b)
}
