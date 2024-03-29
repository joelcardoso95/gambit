package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gambit/models"
	"github.com/gambit/tools"
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

func UpdateCategory(category models.Category) error {
	fmt.Println("Atualizando categoria no banco de dados")

	err := DatabaseConnection()
	if err != nil {
		return err
	}
	defer Database.Close()

	query := "UPDATE category SET "

	if len(category.CategName) > 0 {
		query += " Categ_Name = '" + tools.SkipString(category.CategName) + "'"
	}

	if len(category.CategPath) > 0 {
		if !strings.HasSuffix(query, "SET ") {
			query += ", "
		}
		query += "Categ_Path = '" + tools.SkipString(category.CategPath) + "'"
	}

	query += " WHERE Categ_Id = " + strconv.Itoa(category.CategID)

	_, err = Database.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Categoria Atualizada com sucesso")
	return nil
}
