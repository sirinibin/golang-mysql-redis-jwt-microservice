package migrations

import "database/sql"

func Run(db *sql.DB) {
	CreateUserTable(db)
	CreateEmployeeTable(db)
}
