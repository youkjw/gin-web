package v1

import (
	"fmt"
	"gin-web/pkg/app"
	e "gin-web/pkg/error"
	"gin-web/pkg/file"
	"gin-web/pkg/logging"
	"gin-web/pkg/qrcode"
	"gin-web/pkg/setting"
	"gin-web/pkg/upload"
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GenerateQrCode(c *gin.Context) {
	appG := app.Gin{c}
	url := c.PostForm("url")

	valid := validation.Validation{}
	valid.Required(url, "url").Message("url参数必须")

	qc := qrcode.NewQrCode(url, qrcode.HEIGH, qrcode.HEIGH, qr.M, qr.Auto)
	filename, err := qc.Encode(qrcode.GetQrCodeFullPath())
	if err != nil {
		logging.Error(fmt.Sprintf("qrcode encode err:%v", err))
		appG.Response(http.StatusOK, e.ERROR_QRCODE_GENERATE_FAIL, nil)
		return
	}

	imageT := &upload.Image{
		Src: qrcode.GetQrCodeFullPath() + filename,
		Rect:upload.Rect{
			X0:   0,
			Y0:   0,
			X1:   550,
			Y1:   700,
		},
		Pt:upload.Pt{
			X: 125,
			Y: 298,
		},
	}

	bgSrc := qrcode.GetQrCodeFullPath() + "bg.jpg"
	filename, err = imageT.MergeImage("post-" + strconv.FormatInt(time.Now().Unix(), 10), bgSrc)
	if err != nil {
		logging.Error(fmt.Sprintf("image merge err:%v", err))
		appG.Response(http.StatusOK, e.ERROR_QRCODE_GENERATE_FAIL, nil)
		return
	}

	qrFile, _ := file.Open(upload.GetImageFullPath() + filename)
	//qrFile, err := os.Create(upload.GetImageFullPath() + upload.GetImageName("qrcode"))
	defer qrFile.Close()
	drawT := &upload.DrawText{
		FontSrc: setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + "msyhbd.ttc",
		MergeF: qrFile,
		Title:  "Golang Gin 系列文章",
		X:      80,
		Y:      160,
		Size:   20,
	}

	err = drawT.DrawText()
	if err != nil {
		logging.Error(fmt.Sprintf("image draw err:%v", err))
		appG.Response(http.StatusOK, e.ERROR_QRCODE_GENERATE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"qrcode_url" : upload.GetImageFullUrl(filename),
	})
}

func GenerateDraw(c *gin.Context) {
	appG := app.Gin{c}

	qrFile, err := os.Create(upload.GetImageFullPath() + upload.GetImageName("draw"))
	if err != nil {
		logging.Error(fmt.Sprintf("os create err:%v", err))
		appG.Response(http.StatusOK, e.ERROR_QRCODE_GENERATE_FAIL, nil)
		return
	}

	drawT := &upload.DrawText{
		FontSrc: setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + "msyhbd.ttc",
		MergeF: qrFile,
		Title:  "Golang Gin 系列文章",
		X:      80,
		Y:      160,
		Size:   20,
	}

	err = drawT.DrawText()
	if err != nil {
		logging.Error(fmt.Sprintf("image draw err:%v", err))
		appG.Response(http.StatusOK, e.ERROR_QRCODE_GENERATE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}