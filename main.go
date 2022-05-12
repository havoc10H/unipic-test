package main

import (
	"log"
	"net"

	"github.com/sysdevguru/unipic/config"
	"github.com/sysdevguru/unipic/controller"
	"github.com/sysdevguru/unipic/model"
	"github.com/sysdevguru/unipic/util"

	"github.com/kataras/iris/v12"
)

func init() {
	// generate template html for default page
	templateData := model.Picture{
		Title:       config.Global.Default.Title,
		URL:         config.Global.Default.URL,
		CopyRight:   config.Global.Default.CopyRight,
		Explanation: config.Global.Default.Explanation,
	}
	err := util.ParseTemplate(config.Global.Config.ImageTemPath, templateData)
	if err != nil {
		log.Printf("unipic: generating template html failure: %v\n", err)
	}
}

func main() {
	app := iris.New()
	app.HandleDir("/", "./assets", iris.DirOptions{
		IndexName: "/index.html",
		Gzip:      false,
		ShowList:  false,
	})

	// endpoints
	app.Get("/api/v1/pic", controller.ShowPic)

	// run listener
	listener, err := net.Listen("tcp", config.Global.Config.ServerPort)
	if err != nil {
		log.Printf("unipic: server listening error: %v\n", err)
		return
	}
	app.Run(iris.Listener(listener))
}
