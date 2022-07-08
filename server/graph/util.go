package graph

import (
	"bytes"
	"encoding/json"
	"io"
)

type ListItem struct {
	Index  int64    `json:"index"`
	Amount string   `json:"amount"`
	Proof  []string `json:"proof"`
}

type List map[string]ListItem

func parseMerkleProofsList(list io.Reader) (List, error) {
	listBz := new(bytes.Buffer)
	listBz.ReadFrom(list)

	var parsedList List

	err := json.Unmarshal(listBz.Bytes(), &parsedList)
	if err != nil {
		return nil, err
	}

	return parsedList, nil
}
