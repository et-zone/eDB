package eDB

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	tableRows   map[string]string //tableName rowSql
	db          *gorm.DB
	tableFields map[string]string //tableName fields
}

type EConfig struct {
	UserName  string `json:"userName"`
	PassWord  string `json:"passWord"`
	Addr      string `json:"addr"`
	Port      int    `json:"port"`
	DB        string `json:"db"`
	TableFile string `json:"tablefile"`
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
func InitClient(config *EConfig) {
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
			EClient.InitTableField(key, tblist[key]...)
		}

	}

}

func (cli *Client) Clear(tableName string) {
	cli.tableRows[tableName] = ""

}

func (cli *Client) AddRow(tableName string, row *Row) {
	tmpStr := "("
	for i := 0; i < row.getSize(); i++ {

		switch reflect.TypeOf(row.getColumnValues(i)).String() {
		case "string":
			val := fmt.Sprintf("%v", row.getColumnValues(i))
			if val == "NULL" {
				tmpStr = tmpStr + val
			} else {
				val = strings.ReplaceAll(val, "'", "\\'")
				val = strings.ReplaceAll(val, "\"", "\\\"")
				tmpStr = tmpStr + "'" + val + "'"
			}
		default:
			tmpStr = tmpStr + fmt.Sprintf("%v", row.getColumnValues(i))

		}
		// fmt.Println(reflect.TypeOf(row.getColumnValues(i)).String())
		if i != row.getSize()-1 {
			tmpStr += ", "
		}
	}
	tmpStr = tmpStr[:len(tmpStr)] + ")"
	if len(cli.tableRows[tableName]) != 0 {
		tmpStr = "," + tmpStr
	}
	cli.tableRows[tableName] = cli.tableRows[tableName] + tmpStr

}
func (cli *Client) GettableNanme(tableName string) string {

	return cli.tableRows[tableName]

}

func (cli *Client) isTableIn(tableName string) bool {
	_, isin := cli.tableFields[tableName]
	return isin

}

func (cli *Client) Commit() {

	for key, _ := range cli.tableRows {
		// cli.db.Exec(cli.GettableNanme(key)[:len(cli.GettableNanme(key)-1)])
		// fmt.Println(cli.tableRows[key])

		sql := "INSERT INTO " + key + cli.tableFields[key] + " VALUES " + cli.tableRows[key]
		fmt.Println(sql)

		_, err := cli.db.DB().Exec(sql)

		// d, _ := cout.RowsAffected()
		if err != nil {
			fmt.Println(err.Error())
		}

	}

}

//设置字段名称
func (cli *Client) InitTableField(tableName string, fields ...string) {
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
