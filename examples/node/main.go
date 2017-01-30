package main

import (
	"encoding/base64"
	"github.com/badfortrains/spotcontrol"
	Spotify "github.com/badfortrains/spotcontrol/proto"
	"github.com/gopherjs/gopherjs/js"
)

func setupGlobal() {
	js.Module.Get("exports").Set("spotcontrol", map[string]interface{}{
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
		}
	}()
}

func (c *controllerWrapper) SendFrame(frame string, cb *js.Object) {
	go func() {
		err := c.controller.SendJsonFrame(frame)
		if err != nil {
			cb.Invoke("Frame send failed: " + err.Error())
		} else {
			cb.Invoke(nil)
		}
	}()
}

func (c *controllerWrapper) GetTrack(id string, cb *js.Object) {
	go func() {
		track, err := c.controller.GetTrack(id)
		if err != nil {
			cb.Invoke("Frame send failed: " + err.Error())
		} else {
			cb.Invoke(nil, track)
		}
	}()
}

func (c *controllerWrapper) GetAlbum(id string, cb *js.Object) {
	go func() {
		album, err := c.controller.GetAlbum(id)
		if err != nil {
			cb.Invoke("Frame send failed: " + err.Error())
		} else {
			cb.Invoke(nil, album)
		}
	}()
}

func (c *controllerWrapper) GetRootPlaylist(cb *js.Object) {
	go func() {
		result, err := c.controller.GetRootPlaylist()
		if err != nil {
			cb.Invoke("Frame send failed: " + err.Error())
		} else {
			cb.Invoke(nil, result)
		}
	}()
}

func (c *controllerWrapper) GetSuggest(term string, cb *js.Object) {
	go func() {
		result, err := c.controller.Suggest(term)
		if err != nil {
			cb.Invoke("Frame send failed: " + err.Error())
		} else {
			cb.Invoke(nil, result)
		}
	}()
}

func (c *controllerWrapper) SendPlay(ident string, cb *js.Object) {
	go func() {
		err := c.controller.SendPlay(ident)
		if err != nil {
			cb.Invoke("Hello failed: " + err.Error())
		}
	}()
}

func (c *controllerWrapper) SendPause(ident string, cb *js.Object) {
	go func() {
		err := c.controller.SendPause(ident)
		if err != nil {
			cb.Invoke("Hello failed: " + err.Error())
		}
	}()
}

func (c *controllerWrapper) SendVolume(ident string, volume int, cb *js.Object) {
	go func() {
		err := c.controller.SendVolume(ident, volume)
		if err != nil {
			cb.Invoke("Hello failed: " + err.Error())
		}
	}()
}

func (c *controllerWrapper) LoadTrack(ident string, gids []string, cb *js.Object) {
	go func() {
		err := c.controller.LoadTrack(ident, gids)
		if err != nil {
			cb.Invoke("Hello failed: " + err.Error())
		}
	}()
}

func (c *controllerWrapper) HandleUpdatesCbProto(cb *js.Object) {
	c.controller.HandleUpdatesCbProto(func(frame *Spotify.Frame) {
		cb.Invoke(frame)
	})
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
