package routers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gambit/database"
	"github.com/gambit/models"
)

func InsertProduct(body string, user string) (int, string) {
	var product models.Product
	fmt.Println("Dados recebidos: " + body)
	err := json.Unmarshal([]byte(body), &product)
	if err != nil {
		return 400, "Erro ao receber dados: " + err.Error()
	}

	if len(product.ProdTitle) == 0 {
		return 400, "Titulo do produto é obrigatório"
	}

	isAdmin, msg := database.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	result, errInsert := database.InsertProduct(product)
	if errInsert != nil {
		return 400, "Erro ao realizar insert produto: " + errInsert.Error()
	}

	return 200, "{ProductId:  " + strconv.Itoa(int(result)) + "}"
}

func UpdateProduct(body string, user string, productId int) (int, string) {
	var product models.Product

	err := json.Unmarshal([]byte(body), &product)
	if err != nil {
		return 400, "Erro ao receber os dados do produto " + err.Error()
	}

	isAdmin, msg := database.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	product.ProductId = productId
	errUpdate := database.UpdateProduct(product)
	if errUpdate != nil {
		return 400, "Erro ao realizar atualizar produto " + strconv.Itoa(productId) + " > " + errUpdate.Error()
	}

	return 202, "Produto Atualizado"
}

func DeleteProduct(user string, productId int) (int, string) {
	isAdmin, msg := database.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	err := database.DeleteProduct(productId)
	if err != nil {
		return 400, "Erro ao realizar delete d produto " + strconv.Itoa(productId) + " > " + err.Error()
	}

	return 202, "Produto Excluido"

}

func SelectProduct(request events.APIGatewayV2HTTPRequest) (int, string) {
	var product models.Product
	var page, pageSize int
	var orderType, orderField string

	params := request.QueryStringParameters

	page, _ = strconv.Atoi(params["page"])
	pageSize, _ = strconv.Atoi(params["pageSize"])
	orderType = params["orderType"]   // D = DESC . A or Nil = ASC
	orderField = params["orderField"] // T = Title, 'D' = DATE, P = Price, C = Category, S = Stock

	if !strings.Contains("TDPCS", orderField) {
		orderField = ""
	}

	var choice string
	if len(params["prodId"]) > 0 {
		choice = "P"
		product.ProductId, _ = strconv.Atoi(params["prodId"])
	}
	if len(params["search"]) > 0 {
		choice = "S"
		product.ProdSearch, _ = params["search"]
	}
	if len(params["categId"]) > 0 {
		choice = "C"
		product.ProdCategId, _ = strconv.Atoi(params["categId"])
	}
	if len(params["slug"]) > 0 {
		choice = "U"
		product.ProdPath, _ = params["slug"]
	}
	if len(params["slugCateg"]) > 0 {
		choice = "K"
		product.ProdPath, _ = params["slugCateg"]
	}

	fmt.Println("Parametros: ", params)

	result, errSelect := database.SelectProduct(product, choice, page, pageSize, orderType, orderField)
	if errSelect != nil {
		return 400, "Erro ao consultar produtos " + errSelect.Error()
	}

	returnedProduct, errJson := json.Marshal(result)
	if errJson != nil {
		return 400, "Erro ao converter produtos " + errJson.Error()
	}

	return 200, string(returnedProduct)

}
