package controller

import (
	"log"
	"regexp"
	"time"

	"github.com/sysdevguru/unipic/config"
	"github.com/sysdevguru/unipic/model"
	"github.com/sysdevguru/unipic/util"

	"github.com/kataras/iris/v12"
)

// ShowPic returns picture of provided date
func ShowPic(ctx iris.Context) {
	// get and parse date parameter to "YYYY-MM-DD" format
	dateParam := ctx.FormValue("date")

	var dateStr string
	if dateParam == "" {
		date := time.Now()
		dateStr = date.Format("2006-01-02")
	} else {
		date, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			log.Printf("unipic: time coverting error: %v\n", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
		dateStr = date.Format("2006-01-02")
	}

	// picture object for provided date
	pic := model.Picture{
		Date: dateStr,
	}

	// get picture of provided date
	err := pic.GetPic()
	if err != nil {
		log.Printf("unipic: getting picture failure: %v\n", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	// generate html file based on url of video/picture
	rxPat := regexp.MustCompile(`^(http(s?):)([/|.|\w|\s|-])*\.(?:jpg|gif|png)$`)
	if !rxPat.MatchString(pic.URL) {
		// generate html file with returned video information
		err = util.ParseTemplate(config.Global.Config.VideoTemPath, pic)
		if err != nil {
			log.Printf("unipic: generating html failure: %v\n", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
	} else {
		// generate html file with returned picture information
		err = util.ParseTemplate(config.Global.Config.ImageTemPath, pic)
		if err != nil {
			log.Printf("unipic: generating html failure: %v\n", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.ServeFile(config.Global.Config.IndexPath, false)
}
