package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CustomerAPI struct {
	base string
}

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const (
	headerContentType = "content-type"
	headerAccept      = "accept"

	contentTypeJSON = "application/json"
)

func (ca *CustomerAPI) GetByID(id string) (*Customer, error) {
	u := fmt.Sprintf("%s/customers/%s", ca.base, id)

	r, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Add(headerAccept, contentTypeJSON)
	r.Header.Add(headerContentType, contentTypeJSON)

	res, err := http.DefaultClient.Do(r)

	if err != nil {
		log.Printf("request customer by id %s failed. reason: %v", id, err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("request customer by id %s returned an status code %d", id, res.StatusCode)
	}

	customer := &Customer{}
	if err = json.NewDecoder(res.Body).Decode(&customer); err != nil {
		return nil, err
	}

	return customer, nil
}
