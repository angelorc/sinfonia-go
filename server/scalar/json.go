package scalar

import (
	"encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"io"
)

var _ graphql.Marshaler = (*JSON)(nil)
var _ graphql.Unmarshaler = (*JSON)(nil)

type JSON map[string]interface{}

func (j JSON) MarshalGQL(w io.Writer) {
	buf, _ := json.Marshal(j)
	w.Write(buf)
}

func (j JSON) UnmarshalGQL(v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, j)
}

/*func (s JSONScalar) MarshalGQL(w io.Writer) {

}

func (s JSONScalar) UnmarshalGQL(v interface{}) error {
	byteData, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("FAIL WHILE MARSHAL SCHEME")
	}

	if err := json.Unmarshal(byteData, &s); err != nil {
		return fmt.Errorf("FAIL WHILE UNMARSHAL SCHEME")
	}

	return nil
}

func MarshalJSONScalar(j JSONScalar) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		byteData, err := json.Marshal(j)
		if err != nil {
			log.Printf("FAIL WHILE MARSHAL JSON %v\n", string(byteData))
		}
		_, err = w.Write(byteData)
		if err != nil {
			log.Printf("FAIL WHILE WRITE DATA %v\n", string(byteData))
		}
	})
}

func UnmarshalJSONScalar(v interface{}) (JSONScalar, error) {
	byteData, err := json.Marshal(v)
	if err != nil {
		return JSONScalar{}, fmt.Errorf("FAIL WHILE MARSHAL SCHEME")
	}

	tmp := make(JSONScalar)
	if err := json.Unmarshal(byteData, &tmp); err != nil {
		return JSONScalar{}, fmt.Errorf("FAIL WHILE UNMARSHAL SCHEME")
	}

	return tmp, nil
}
*/
