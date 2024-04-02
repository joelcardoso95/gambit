package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

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
