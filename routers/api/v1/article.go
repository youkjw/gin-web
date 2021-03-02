package v1

import (
	"gin-web/models"
	"gin-web/pkg/app"
	e "gin-web/pkg/error"
	"gin-web/pkg/setting"
	"gin-web/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"log"
	"net/http"
)

func GetArticle (c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	maps := make(map[string]interface{})

	var data interface{}
	var err error

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if ! models.ExistsArticleByID(id) {
			maps["id"] = id
			data, err = models.GetArticle(maps)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_ARTICLE, nil)
				return
			}

			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		app.MarkErrors(valid.Errors)
	}

	appG.Response(http.StatusOK, code, data)
}

func GetArticles (c *gin.Context) {
	appG := app.Gin{c}

	data := make(map[string]interface{})
	maps := make(map[string]interface{})

	valid := validation.Validation{}

	state := 1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	tagId := 0
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS

		data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
	} else {
		app.MarkErrors(valid.Errors)
	}

	appG.Response(http.StatusOK, code, data)
}

func AddArticle (c *gin.Context) {
	appG := app.Gin{c}

	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	state := com.StrTo(c.DefaultPostForm("state", "1")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if models.ExistsTagByID(tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		app.MarkErrors(valid.Errors)
	}

	appG.Response(http.StatusOK, code, make(map[string]interface{}))
}

func EditArticle (c *gin.Context) {
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")

	var state int = 0
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if models.ExistsArticleByID(id) {
			if models.ExistsTagByID(tagId) {
				data := make(map[string]interface {})
				if tagId > 0 {
					data["tag_id"] = tagId
				}
				if title != "" {
					data["title"] = title
				}
				if desc != "" {
					data["desc"] = desc
				}

				models.EditArticle(id, data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}

func DeleteArticle (c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	var code int
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistsArticleByID(id) {
			models.DeteleArticle(id)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}
