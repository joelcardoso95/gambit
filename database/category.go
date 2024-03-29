package database

import (
	"database/sql"
	"fmt"

	"github.com/gambit/models"
)

func InsertCategory(category models.Category) (int64, error) {
	fmt.Println("Inserindo categoria no banco de dados")

	err := DatabaseConnection()
	if err != nil {
		return 0, err
	}
	defer Database.Close()

	query := "INSERT INTO category (Categ_Name, Categ_Path) VALUES ('" + category.CategName + "','" + category.CategPath + "')"

	var result sql.Result
	result, err = Database.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	LastInsertId, errInsert := result.LastInsertId()
	if errInsert != nil {
		fmt.Println(errInsert.Error())
		return 0, errInsert
	}

	fmt.Println("Categoria inserida com sucesso")
	return LastInsertId, nil
}
