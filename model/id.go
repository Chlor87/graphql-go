package model

import (
	"encoding/json"
)

type ID int

func (ID) ImplementsGraphQLType(name string) bool {
	return name == "ID"
}

func (id *ID) UnmarshalGraphQL(input any) error {
	switch t := input.(type) {
	case int32:
		// from query arg
		*id = ID(int(t))
	case float64:
		// from json parsing (golang thingy)
		*id = ID(int(t))
	case int:
		// a normal ID
		*id = ID(t)
	}
	return nil
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(id))
}
