package routers

import (
	"encoding/json"
	"strconv"

	"github.com/gambit/database"
	"github.com/gambit/models"
)

func InsertCategory(body string, user string) (int, string) {
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

	isAdmin, msg := database.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	result, errInsert := database.InsertCategory(category)
	if errInsert != nil {
		return 400, "Erro ao realizar inserção de nova categoria " + category.CategName + " > " + errInsert.Error()
	}

	return 200, "CategID: " + strconv.Itoa(int(result))
}

func UpdateCategory(body string, user string, categoryId int) (int, string) {
	var category models.Category

	err := json.Unmarshal([]byte(body), &category)
	if err != nil {
		return 400, "Erro ao receber os dados da categoria " + err.Error()
	}

	if len(category.CategName) == 0 && len(category.CategPath) == 0 {
		return 400, "Erro ao receber o Path da categoria "
	}

	isAdmin, msg := database.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	category.CategID = categoryId
	errUpdate := database.UpdateCategory(category)
	if errUpdate != nil {
		return 400, "Erro ao realizar atualizar categoria " + strconv.Itoa(categoryId) + " > " + errUpdate.Error()
	}

	return 200, "Categoria Atualizada"
}

func DeleteCategory(body string, user string, id int) (int, string) {
	if id == 0 {
		return 400, "Id Categoria deve ser informado"
	}

	isAdmin, msg := database.UserIsAdmin(user)
	if !isAdmin {
		return 400, msg
	}

	err := database.DeleteCategory(id)
	if err != nil {
		return 400, "Erro ao tentar excluir categoria " + strconv.Itoa(id) + " > " + err.Error()
	}

	return 201, ""

}
