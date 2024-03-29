package routers

import (
	"encoding/json"
	"strconv"

	"github.com/gambit/database"
	"github.com/gambit/models"
)

func InsertCategory(body string, User string) (int, string) {
	var category models.Category

	err := json.Unmarshal([]byte(body), &category)
	if err != nil {
		return 400, "Erro ao receber os dados da categoria " + err.Error()
	}

	if len(category.CategName) == 0 {
		return 400, "Erro ao receber o Nome da categoria "
	}

	if len(category.CategPath) == 0 {
		return 400, "Erro ao receber o Path da categoria "
	}

	isAdmin, msg := database.UserIsAdmin(User)
	if !isAdmin {
		return 400, msg
	}

	result, errInsert := database.InsertCategory(category)
	if errInsert != nil {
		return 400, "Erro ao realizar inserção de nova categoria " + category.CategName + " > " + errInsert.Error()
	}

	return 200, "CategID: " + strconv.Itoa(int(result))
}
