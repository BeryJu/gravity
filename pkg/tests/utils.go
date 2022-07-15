package tests

import (
	"context"
	"encoding/json"
	"time"
)

func MustJSON(in interface{}) string {
	j, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func Context() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	return ctx
}
