package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) getAccountID(r *http.Request) int64 {
	return 5
}

func (s *Server) parseData(dataStr string, obj interface{}) error {

	type Data struct {
		Fields interface{} `json:"fields"`
	}

	data := Data{}

	bytes := []byte(dataStr)

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	return convertObject(data.Fields, &obj)
}

func convertObject(obj1 interface{}, obj2 interface{}) error {
	b, err := json.Marshal(obj1)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, obj2); err != nil {
		return err
	}
	return nil
}
