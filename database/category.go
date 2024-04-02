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

func DeleteCategory(id int) error {
	fmt.Println("Excluindo categoria")

	err := DatabaseConnection()
	if err != nil {
		return err
	}
	defer Database.Close()

	query := "DELETE FROM category WHERE Categ_Id = " + strconv.Itoa(id)

	_, err = Database.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Categoria excluida com sucesso")
	return nil
}

func SelectCategories(categId int, slug string) ([]models.Category, error) {
	fmt.Println("Consultando categorias")

	var categories []models.Category

	err := DatabaseConnection()
	if err != nil {
		return categories, err
	}
	defer Database.Close()

	query := "SELECT Categ_Id, Categ_Name, Categ_Path FROM category "

	if categId > 0 {
		query += "WHERE Categ_Id = " + strconv.Itoa(categId)
	} else {
		if len(slug) > 0 {
			query += "WHERE Categ_Path LIKE '%" + slug + "%'"
		}
	}

	fmt.Println("Consultando categorias", query)

	var rows *sql.Rows
	rows, err = Database.Query(query)
	if err != nil {
		return categories, err
	}

	for rows.Next() {
		var category models.Category
		var categId sql.NullInt32
		var categName sql.NullString
		var categPath sql.NullString

		err := rows.Scan(&categId, &categName, &categPath)
		if err != nil {
			return categories, err
		}

		category.CategID = int(categId.Int32)
		category.CategName = categName.String
		category.CategPath = categPath.String

		categories = append(categories, category)
	}

	fmt.Println("Categorias encontradas")
	return categories, nil

}
