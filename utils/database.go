package utils

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)
//
//sqlite> create table nodes(
//...> node_address text primary key not null,
//...> node_join_time integer not null,
//...> node_current_status text not null);
//sqlite> create table subtasks (
//...> subtask_id text primary key not null,
//...> main_task_id text not null,
//...> main_working_node_address text,
//...> backup_working_node_address text,
//...> main_instance_status text,
//...> backup_instance_status text,
//...> main_instance_result text,
//...> backup_instance_result text);
const databaseSourcePath string = "utils/database.db"

type Database struct {
	DbObject *sql.DB
}

var DB_instance *Database
var once sync.Once

func NewDatabase() *Database {
	once.Do(func() {
		db, err := sql.Open("sqlite3", databaseSourcePath)
		Check(err, "sqlite3 opening failed")
		DB_instance = &Database{DbObject: db}
		//DB_instance.db
	})
	return DB_instance
}

func (receiver Database) name() {

}
//
//func (d *Database) execute(cmd string) error {
//
//	return nil
//}


