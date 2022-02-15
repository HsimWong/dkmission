package main

import (
	"dkmission/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	db_instance := utils.NewDatabase()
	//rst, err := db_instance.DbObject.Query("select * from nodes;")
	//utils.Check(err, "query failed")
	//log.Info(rst)
}
