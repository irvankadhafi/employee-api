package cacher

import (
	"encoding/json"
	redigo "github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

type MultiResponse struct {
	IDs   []uint `json:"ids"`
	Count uint   `json:"count"`
}

// NewMultiResponseFromByte converts interface to multi response model.
func NewMultiResponseFromByte(bt []byte) (mr *MultiResponse, err error) {
	if err := json.Unmarshal(bt, &mr); err != nil {
		log.WithField("bt", string(bt)).Error(err)
		return nil, err
	}

	return
}

func get(client redigo.Conn, key string) (value any, err error) {
	defer WrapCloser(client.Close)

	if err := client.Send("MULTI"); err != nil {
		return nil, err
	}

	if err := client.Send("EXISTS", key); err != nil {
		return nil, err
	}

	if err := client.Send("GET", key); err != nil {
		return nil, err
	}

	res, err := redigo.Values(client.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	val, ok := res[0].(int64)
	if !ok || val <= 0 {
		return nil, ErrKeyNotExist
	}

	return res[1], nil
}

func getHashMember(client redigo.Conn, identifier, key string) (value any, err error) {
	defer func() {
		_ = client.Close()
	}()

	if err := client.Send("MULTI"); err != nil {
		return nil, err
	}

	if err := client.Send("HEXISTS", identifier, key); err != nil {
		return nil, err
	}

	if err := client.Send("HGET", identifier, key); err != nil {
		return nil, err
	}

	res, err := redigo.Values(client.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	val, ok := res[0].(int64)
	if !ok || val <= 0 {
		return nil, ErrKeyNotExist
	}

	return res[1], nil
}
