package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TokenJSON struct {
	Sub       string
	Event_id  string
	Token_use string
	Scope     string
	Auth_time int
	Iss       string
	Exp       int
	Iat       int
	Client_id string
	Username  string
}

func TokenValidate(token string) (bool, error, string) {
	tokenArray := strings.Split(token, ".")

	if len(tokenArray) != 3 {
		fmt.Println("Token Inválido")
		return false, nil, "Token Inválido"
	}

	userInfo, err := base64.StdEncoding.DecodeString(tokenArray[1])
	if err != nil {
		fmt.Println("Falha ao decodificar token: ", err.Error())
		return false, err, err.Error()
	}

	var jsonToken TokenJSON
	err = json.Unmarshal(userInfo, &jsonToken)
	if err != nil {
		fmt.Println("Falha na transformação para JSON: ", err.Error())
		return false, err, err.Error()
	}

	now := time.Now()
	expirationTime := time.Unix(int64(jsonToken.Exp), 0)

	if expirationTime.Before(now) {
		fmt.Println("Token expirado: ", expirationTime.String())
		return false, err, "Token expirado"
	}

	return true, nil, string(jsonToken.Username)

}
