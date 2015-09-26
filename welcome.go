package main

import (
	"github.com/ixchi/foxbot/bot"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var bot *foxbot.Bot
var storage *pluginStorage

func main() {
	foxbot.SetPlugins([]foxbot.Plugin{
		&pluginAdmin{},
		&pluginWatcher{},
		&pluginHelp{},
	})

	storage = &pluginStorage{
		FilePath: "storage.json",
	}

	if storage.Load() != nil {
		if err := storage.Create(); err != nil {
			log.Println(err)
		}
	}

	router := httprouter.New()

	router.GET("/", index)
	router.GET("/rules/:channel", rules)
	router.GET("/settings/:channel", settings)

	go http.ListenAndServe(":8080", router)

	tg := foxbot.LoadBot()
	bot = foxbot.GetBot()
	bot.Start(tg)
}
