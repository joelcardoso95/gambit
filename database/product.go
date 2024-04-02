package database

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gambit/models"
	"github.com/gambit/tools"
)

func InsertProduct(product models.Product) (int64, error) {
	fmt.Println("Inserindo produto")

	err := DatabaseConnection()
	if err != nil {
		return 0, err
	}
	defer Database.Close()

	query := "INSERT INTO products (Prod_Title, "

	if len(product.ProdDescription) > 0 {
		query += ", Prod_Description"
	}
	if product.ProdPrice > 0 {
		query += ", Prod_Price"
	}
	if product.ProdCategId > 0 {
		query += ", Prod_CategoryId"
	}
	if product.ProdStock > 0 {
		query += ", Prod_Stock"
	}
	if len(product.ProdPath) > 0 {
		query += ", Prod_Path"
	}

	query += ") VALUES ('" + tools.SkipString(product.ProdTitle) + "'"

	if len(product.ProdDescription) > 0 {
		query += ", '" + tools.SkipString(product.ProdDescription) + "'"
	}
	if product.ProdPrice > 0 {
		query += ", " + strconv.FormatFloat(product.ProdPrice, 'e', -1, 64)
	}
	if product.ProdCategId > 0 {
		query += ", " + strconv.Itoa(product.ProdCategId)
	}
	if product.ProdStock > 0 {
		query += ", " + strconv.Itoa(product.ProdStock)
	}
	if len(product.ProdPath) > 0 {
		query += ", '" + tools.SkipString(product.ProdPath) + "'"
	}

	query += ")"

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

	fmt.Println("Produto inserido com sucesso")
	return LastInsertId, nil
}
