package redis

import (
	"aig-tech-okr/libs/cache"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type User struct{}

func (i *User) ins() string {
	return "user"
}

// 登录状态	------------------------

type LoginUserInfo struct {
	ID     uint   `json:"id"`
	Openid string `json:"openid"`
}

func (i *User) loginToken(token, platform string) (key string) {
	key = fmt.Sprintf("user:login:token:%s:%s", token, platform)
	return
}

func (i *User) SetLoginToken(token, platform string, data LoginUserInfo) error {
	rds := cache.PoolGet(i.ins())

	value, _ := json.Marshal(data)

	err := rds.Set(i.loginToken(token, platform), string(value), 168*time.Hour).Err()
	return err
}

func (i *User) GetLoginToken(token, platform string) (data LoginUserInfo, err error) {
	rds := cache.PoolGet(i.ins())

	value, err := rds.Get(i.loginToken(token, platform)).Result()
	if value == "" {
		err = errors.New("empty cache")
		return
	}

	err = json.Unmarshal([]byte(value), &data)
	return
}

func (i *User) EmptyLoginToken(token, platform string) (err error) {
	rds := cache.PoolGet(i.ins())
	err = rds.Del(i.loginToken(token, platform)).Err()
	return
}

// 登录状态	------------------------
