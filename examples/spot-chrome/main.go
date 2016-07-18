package main

import (
	"encoding/base64"
	"github.com/badfortrains/spotcontrol"
	"github.com/gopherjs/gopherjs/js"
)

func setupGlobal() {
	js.Global.Set("spotcontrol", map[string]interface{}{
		"login":      login,
		"loginSaved": loginSaved,
		"convert62":  convert64to62,
	})
}

type controllerWrapper struct {
	controller *spotcontrol.SpircController
}

func (c *controllerWrapper) SendHello(cb *js.Object) {
	go func() {
		err := c.controller.SendHello()
		if err != nil {
			cb.Invoke("Hello failed: " + err.Error())
		} else {
			cb.Invoke(nil)
		}
	}()
}

func (c *controllerWrapper) SendPlay(ident string, cb *js.Object) {
	go func() {
		err := c.controller.SendPlay(ident)
		if err != nil {
			cb.Invoke("Send play failed: " + err.Error())
		} else {
			cb.Invoke(nil)
		}
	}()
}

func (c *controllerWrapper) SendPause(ident string, cb *js.Object) {
	go func() {
		err := c.controller.SendPause(ident)
		if err != nil {
			cb.Invoke("Send pause failed: " + err.Error())
		} else {
			cb.Invoke(nil)
		}
	}()
}

func (c *controllerWrapper) SendVolume(ident string, volume int, cb *js.Object) {
	go func() {
		err := c.controller.SendVolume(ident, volume)
		if err != nil {
			cb.Invoke("Send volume failed: " + err.Error())
		} else {
			cb.Invoke(nil)
		}
	}()
}

func (c *controllerWrapper) LoadTrack(ident string, gids []string, cb *js.Object) {
	go func() {
		err := c.controller.LoadTrack(ident, gids)
		if err != nil {
			cb.Invoke("Load track failed: " + err.Error())
		} else {
			cb.Invoke(nil)
		}
	}()
}

func (c *controllerWrapper) Suggest(term string, cb *js.Object) {
	go func() {
		res, err := c.controller.Suggest(term)
		if err != nil {
			cb.Invoke(nil, "Hello failed: "+err.Error())
		} else {
			cb.Invoke(res)
		}
	}()
}

func (c *controllerWrapper) Search(term string, cb *js.Object) {
	go func() {
		res, err := c.controller.Search(term)
		if err != nil {
			cb.Invoke(nil, "Hello failed: "+err.Error())
		} else {
			cb.Invoke(res)
		}
	}()
}

func (c *controllerWrapper) HandleUpdatesCb(cb func(device string)) {
	c.controller.HandleUpdatesCb(cb)
}

func convert64to62(data64 string) string {
	data, _ := base64.StdEncoding.DecodeString(data64)
	return spotcontrol.ConvertTo62(data)
}

func loginSaved(username, authData string, appkey string, cb *js.Object) {
	go func() {
		key, _ := base64.StdEncoding.DecodeString(appkey)
		data, _ := base64.StdEncoding.DecodeString(authData)
		conn, _ := MakeConn()
		sController, err := spotcontrol.LoginConnectionSaved(username, data, key, "spotcontrol", conn)
		if err != nil {
			cb.Invoke(nil, "", "login failed")
		}
		c := &controllerWrapper{controller: sController}
		cb.Invoke(js.MakeWrapper(c), authData, nil)
	}()
}

func login(username, password, appkey string, cb *js.Object) {
	go func() {
		key, _ := base64.StdEncoding.DecodeString(appkey)
		conn, _ := MakeConn()
		sController, err := spotcontrol.LoginConnection(username, password, key, "spotcontrol", conn)
		if err != nil {
			cb.Invoke(nil, "", "login failed")
		} else {
			authData := sController.SavedCredentials
			c := &controllerWrapper{controller: sController}
			cb.Invoke(js.MakeWrapper(c), base64.StdEncoding.EncodeToString(authData), nil)
		}
	}()
}

func main() {
	setupGlobal()
}
