package v1

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/app"
	e "gin-web/pkg/error"
	"gin-web/pkg/export"
	"gin-web/pkg/logging"
	"gin-web/pkg/setting"
	"gin-web/pkg/util"
	"gin-web/service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)


func GetTags (c *gin.Context) {
	appG := app.Gin{c}

	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = 1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS
	data["list"], _ = models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	appG.Response(http.StatusOK, code, data)
}

func AddTags (c *gin.Context) {
	appG := app.Gin{c}

	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	var code int
	if ! valid.HasErrors() {
		if ! models.ExistsTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	} else {
		code = e.INVALID_PARAMS
	}

	appG.Response(http.StatusOK, code, make(map[string]interface{}))
}

func EditTags (c *gin.Context) {
	appG := app.Gin{c}

	id := com.StrTo(c.Param("id")).MustInt()
	name := c.PostForm("name")

	data := make(map[string]interface{})
	data["name"] = name

	valid := validation.Validation{}

	var state int = 1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

		data["state"] = state
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	var code int
	if ! valid.HasErrors() {
		if ! models.ExistsTagByName(name) {
			code = e.SUCCESS
			fmt.Println(id, data)
			models.EditTag(id, data)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	} else {
		code = e.INVALID_PARAMS
	}

	appG.Response(http.StatusOK, code, make(map[string]interface{}))
}

func DeteleTags (c *gin.Context) {
	appG := app.Gin{c}

	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	var code int
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistsTagByID(id) {
			models.DeteleTag(id)
		}
	} else {
		app.MarkErrors(valid.Errors)
	}

	appG.Response(http.StatusOK, code, make(map[string]interface{}))
}

func ExportTags (c *gin.Context) {
	appG := app.Gin{c}

	name := c.PostForm("name")
	state := 1

	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := service.Tag{
		Name: name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url": export.GetExcelFullUrl(filename),
		"export_savel_url": export.GetExcelPath() + filename,
	})
}

func ExportInTags (c *gin.Context) {
	appG := app.Gin{c}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR_IMPORT_FILE_NOT_FOUND, nil)
		return
	}

	tagService := service.Tag{}
	err = tagService.ExportIn(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}