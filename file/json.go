package file

import (
	"context"
	"encoding/json"
	"io/ioutil"
)

func LoadSql(ctx context.Context, name string) ([]string, error) {
	var sql []string
	if name == "" {
		return sql, nil
	}

	data, err := ioutil.ReadFile(name)
	if err != nil {
		return sql, err
	}
	if err = json.Unmarshal(data, &sql); err != nil {
		return sql, err
	}
	return sql, nil
}
