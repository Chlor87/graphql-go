package model

import (
	"encoding/json"
	"strconv"
)

type ID int

func (ID) ImplementsGraphQLType(name string) bool {
	return name == "ID"
}

func (id *ID) UnmarshalGraphQL(input any) error {
	switch t := input.(type) {
	case int32:
		*id = ID(int(t))
	case string:
		v, err := strconv.Atoi(t)
		if err != nil {
			return err
		}
		*id = ID(v)
	}
	return nil
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(id))
}
