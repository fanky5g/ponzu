package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var allowedResponseKeys = []string{"data", "error"}

func writeJSON(w io.Writer, resp map[string]interface{}) error {
	enc := json.NewEncoder(w)
	for k := range resp {
		allowed := false
		for _, v := range allowedResponseKeys {
			if k == v {
				allowed = true
				break
			}
		}

		if !allowed {
			return fmt.Errorf(
				"invalid response structure. response must have any of %s keys",
				strings.Join(allowedResponseKeys, ", "),
			)
		}
	}

	err := enc.Encode(resp)
	if err != nil {
		log.Println("Failed to encode data to JSON:", err)
		return err
	}

	return nil
}

// WriteJSONResponse should be used any time you want to communicate
// data back to a foreign client
func WriteJSONResponse(res http.ResponseWriter, statusCode int, response map[string]interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Vary", "Accept-Encoding")

	res.WriteHeader(statusCode)
	if err := writeJSON(res, response); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
