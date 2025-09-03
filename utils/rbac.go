package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

func CheckRbacAccess(nameRoute string, idUser string) (response bool, err error) {

	apiUrl := fmt.Sprintf("%s/rbac/check-access/%s/%s",
		os.Getenv("SERVICE_RBAC"),
		idUser,
		nameRoute,
	)
	// log.Println("apiUrl: ", apiUrl)

	client := resty.New()
	//client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(apiUrl)
	if err != nil {
		log.Println("check access error: ", err)

		err = RequestError{
			Code:    http.StatusInternalServerError,
			Message: "gagal memeriksa akses ke service RBAC. " + err.Error(),
		}

		return
	}

	if !resp.IsSuccess() {
		err = RequestError{
			Code:    resp.StatusCode(),
			Message: resp.String(),
		}

		return
	}

	//assign response dari rbac true or false ke variable response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return
	}
	// log.Println("response: ", response)

	return
}
