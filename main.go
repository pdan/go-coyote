package main

import (
	"time"

	"coyote/cloudflare"
	"coyote/setting"
)

func main() {
	setting.NewContext()
	c := new(cloudflare.Client)
	c.API = setting.Cfg.ClientAPI
	c.Email = setting.Cfg.ClientEmail
	if c.FetchAll() {

		for {
			c.Run()
			time.Sleep(setting.Cfg.CheckTime * time.Second)

		}
	}

}
