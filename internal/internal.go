package internal

import (
	"database/sql/driver"
	"fmt"
	"io"
	"os"

	"github.com/istabraq/debit_without_dispense/pkg/config"
	go_ora "github.com/sijms/go-ora"
)

//var query string

func Run(dbParams config.Database, query string) {
	connectionString := "oracle://" + dbParams.UserName + ":" + dbParams.Password + "@" + dbParams.IP + ":" + dbParams.Port + "/" + dbParams.ServiceName
	DB, err := go_ora.NewConnection(connectionString)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	dieOnError("Can't open the driver:", err)
	err = DB.Open()
	dieOnError("Can't open the connection:", err)

	defer DB.Close()
	//query = "SELECT *  FROM v$instance"
	stmt := go_ora.NewStmt(query, DB)
	defer stmt.Close()

	rows, err := stmt.Query(nil)
	dieOnError("Can't query", err)
	defer rows.Close()

	columns := rows.Columns()

	values := make([]driver.Value, len(columns))
	Header(columns)
	for {
		err = rows.Next(values)
		if err != nil {
			break
		}
		Record(values)
	}
	fmt.Println()
	if err != io.EOF {
		dieOnError("Can't Next", err)
	}
}

func Header(columns []string) {
	for _, c := range columns {
		fmt.Printf("%-16s ", c)
	}
	fmt.Println()
}

func Record(values []driver.Value) {
	for _, c := range values {
		fmt.Printf("%-20v", c)
	}

}

func dieOnError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}
