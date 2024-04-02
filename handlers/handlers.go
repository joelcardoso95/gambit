package handlers

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gambit/auth"
	"github.com/gambit/routers"
)

func Handlers(path string, method string, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Processando " + path + " > " + method)

	fmt.Println("Requisição recebida ", request.Headers)

	id := request.PathParameters["id"]
	idn, _ := strconv.Atoi(id)

	isOK, statusCode, user := tokenAuthorization(path, method, headers)
	if !isOK {
		return statusCode, user
	}

	switch path[1:5] {
	case "user":
		return UsersProcess(body, path, method, user, id, request)
	case "prod":
		return ProductProcess(body, path, method, user, idn, request)
	case "stoc":
		return StockProcess(body, path, method, user, idn, request)
	case "addr":
		return AddressProcess(body, path, method, user, idn, request)
	case "cate":
		return CategoryProcess(body, path, method, user, idn, request)
	case "orde":
		return OrdersProcess(body, path, method, user, idn, request)
	}

	return 400, "Method Invalid"
}

func tokenAuthorization(path string, method string, headers map[string]string) (bool, int, string) {
	if (path == "product" && method == "GET") || (path == "category" && method == "GET") {
		return true, 200, ""
	}

	fmt.Println("Validação de Token", headers)

	token := headers["authorization"]
	if len(token) == 0 {
		return false, 401, "Token Requerido"
	}

	success, err, msg := auth.TokenValidate(token)
	if !success {
		if err != nil {
			fmt.Println("Erro ao validar token", err.Error())
			return false, 401, err.Error()
		} else {
			fmt.Println("Erro ao validar token", msg)
			return false, 401, msg
		}
	}

	fmt.Println("Token validado com sucesso")
	return true, 200, msg

}

func UsersProcess(body string, path string, method string, user string, id string, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Method invalid"
}

func ProductProcess(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Method invalid"
}

func CategoryProcess(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	switch method {
	case "POST":
		return routers.InsertCategory(body, user)
	case "PUT":
		return routers.UpdateCategory(body, user, id)
	case "DELETE":
		return routers.DeleteCategory(body, user, id)
	case "GET":
		return routers.SelectCategories(body, request)
	}
	return 400, "Method invalid"
}

func StockProcess(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Method invalid"
}

func AddressProcess(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Method invalid"
}

func OrdersProcess(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Method invalid"
}
