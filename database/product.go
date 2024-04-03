package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

	query := "INSERT INTO products (Prod_Title "

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
	fmt.Println("Insert query produto " + query)

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

func UpdateProduct(product models.Product) error {
	fmt.Println("Realizando atualizacao do produto")

	err := DatabaseConnection()
	if err != nil {
		return err
	}
	defer Database.Close()

	query := "UPDATE products SET "

	query = tools.AdjustQuery(query, "Prod_Title", "S", 0, 0, product.ProdTitle)
	query = tools.AdjustQuery(query, "Prod_Description", "S", 0, 0, product.ProdDescription)
	query = tools.AdjustQuery(query, "Prod_Price", "F", 0, product.ProdPrice, "")
	query = tools.AdjustQuery(query, "Prod_CategoryId", "N", product.ProdCategId, 0, "")
	query = tools.AdjustQuery(query, "Prod_Stock", "N", product.ProdStock, 0, "")
	query = tools.AdjustQuery(query, "Prod_Path", "S", 0, 0, product.ProdPath)

	query += " WHERE Prod_Id = " + strconv.Itoa(product.ProductId)

	_, err = Database.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Produto atualizado com sucesso")
	return nil
}

func DeleteProduct(productId int) error {
	fmt.Println("Realziando exclusao do produto")

	err := DatabaseConnection()
	if err != nil {
		return err
	}
	defer Database.Close()

	query := "DELETE FROM products WHERE Prod_Id = " + strconv.Itoa(productId)

	_, err = Database.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Produto excluido com sucesso")
	return nil
}

func SelectProduct(product models.Product, choice string, page int, pageSize int, orderType string, orderField string) (models.ProductResp, error) {
	fmt.Println("Selecionando produtos", product)
	var response models.ProductResp
	var products []models.Product

	err := DatabaseConnection()
	if err != nil {
		return response, err
	}
	defer Database.Close()

	var query string
	var queryCount string
	var where, limit string

	query = "SELECT Prod_id, Prod_Title, Prod_Description, Prod_CreatedAt, Prod_Updated, Prod_Price, Prod_Path, Prod_CategId, Prod_Stock FROM products"
	query = "SELECT count(*) as registros FROM products"

	switch choice {
	case "P":
		where = " WHERE Prod_Id = " + strconv.Itoa(product.ProductId)
	case "S":
		where = " WHERE UCASE(CONCAT(Prod_Title, Prod_Description)) LIKE '%" + strings.ToUpper(product.ProdSearch) + "%' "
	case "C":
		where = " WHERE Prod_CategoryId = " + strconv.Itoa(product.ProdCategId)
	case "U":
		where = " WHERE UCASE(Prod_Path) LIKE '%" + strings.ToUpper(product.ProdPath) + "%' "
	case "K":
		join := " JOIN category ON Prod_CategoryId = Categ_Id AND Categ_Path LIKE '%" + strings.ToUpper(product.ProdCategPath) + "%' "
		query += join
		queryCount += join
	}

	queryCount += where

	var rows *sql.Rows
	rows, err = Database.Query(queryCount)
	defer rows.Close()

	if err != nil {
		fmt.Println(err.Error())
		return response, err
	}

	rows.Next()
	var register sql.NullInt32
	err = rows.Scan(&register)

	registers := int(register.Int32)

	if page > 0 {
		if registers > pageSize {
			limit = " LIMIT " + strconv.Itoa(pageSize)
			if page > 1 {
				offset := pageSize * (page - 1)
				limit += " OFFSET " + strconv.Itoa(offset)
			}
		} else {
			limit = ""
		}
	}

	// T = Title, 'D' = DATE, P = Price, C = Category, S = Stock
	var orderBy string
	if len(orderField) > 0 {
		switch orderField {
		case "T":
			orderBy = " ORDER BY Prod_Title "
		case "D":
			orderBy = " ORDER BY Prod_CreatedAt "
		case "P":
			orderBy = " ORDER BY Prod_Price "
		case "C":
			orderBy = " ORDER BY Prod_CategoryId "
		case "S":
			orderBy = " ORDER BY Prod_Stock "
		}
		if orderType == "D" {
			orderBy += " DESC"
		}
	}

	query += where + orderBy + limit

	fmt.Println("Query de consulta produtos: ", query)

	rows, err = Database.Query(query)

	for rows.Next() {
		var product models.Product
		var prodId sql.NullInt32
		var prodTitle sql.NullString
		var prodDescription sql.NullString
		var prodCreatedAt sql.NullTime
		var prodUpdated sql.NullTime
		var prodPrice sql.NullFloat64
		var prodPath sql.NullString
		var prodCategoryId sql.NullInt32
		var prodStock sql.NullInt32

		err := rows.Scan(&prodId, &prodTitle, &prodDescription, &prodCreatedAt, &prodUpdated, &prodPrice, &prodPath, &prodCategoryId, &prodStock)
		if err != nil {
			return response, err
		}

		product.ProductId = int(prodId.Int32)
		product.ProdTitle = prodTitle.String
		product.ProdDescription = prodDescription.String
		product.ProdCreatedAt = prodCreatedAt.Time.String()
		product.ProdUptaded = prodUpdated.Time.String()
		product.ProdPrice = prodPrice.Float64
		product.ProdPath = prodPath.String
		product.ProdCategId = int(prodCategoryId.Int32)
		product.ProdStock = int(prodStock.Int32)

		products = append(products, product)
	}

	response.TotalItems = registers
	response.Data = products

	fmt.Println("Produtos consultados com sucesso: ", response.Data)
	return response, nil
}
