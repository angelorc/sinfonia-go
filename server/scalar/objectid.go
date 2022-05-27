package scalar

import (
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"strconv"
)

func MarshalObjectIDScalar(id primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(id.Hex()))
	})
}

func UnmarshalObjectIDScalar(v interface{}) (primitive.ObjectID, error) {
	str, ok := v.(string)
	if !ok {
		return primitive.ObjectID{}, errors.New("ids must be strings")
	}

	oid, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid ObjectID")
	}

	return oid, nil
}
