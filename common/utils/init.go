package utils

import (
	"spider/common/env"
)


func InitConfigure(confFilePath string) {
	if err := env.Conf.Load(confFilePath); err != nil {
		panic(err)
	}
}
