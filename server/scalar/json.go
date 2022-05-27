package scalar

import (
	"encoding/json"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/angelorc/sinfonia-go/utility"
	"io"
)

func MarshalJSONScalar(s interface{}) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, s.(string))
	})
}

func UnmarshalJSONScalar(str interface{}) (interface{}, error) {
	jsonByte, err := json.Marshal(&str)
	if err != nil {
		return nil, errors.New("field must be valid graphql query")
	}

	jsonString := string(jsonByte)
	if !utility.IsJSON(jsonString) {
		return nil, errors.New("field must be valid graphql query [json]")
	}

	return jsonString, nil
}
