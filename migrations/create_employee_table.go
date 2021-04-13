package migrations

import (
	"context"
	"database/sql"
	"log"
	"time"
	//_ "github.com/go-sql-driver/mysql"
)

func CreateEmployeeTable(db *sql.DB) (sql.Result, error) {
	query := `CREATE TABLE IF NOT EXISTS employee (
		id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
		created_by int NOT NULL,
		updated_by int NOT NULL,
		name varchar(200) DEFAULT NULL,
		email varchar(200) DEFAULT NULL,
		created_at datetime DEFAULT NULL,
		updated_at datetime DEFAULT NULL
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating employee table", err)

	}
	return res, err
}
