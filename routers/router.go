package routers

import (
	"gin-web/middleware/jwt"
	"gin-web/pkg/export"
	"gin-web/pkg/qrcode"
	"gin-web/pkg/setting"
	"gin-web/pkg/upload"
	"gin-web/routers/api"
	v1 "gin-web/routers/api/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新增标签
		apiv1.POST("/tags", v1.AddTags)
		//更新标签
		apiv1.PUT("/tags/:id", v1.EditTags)
		//删除标签
		apiv1.DELETE("tags/:id", v1.DeteleTags)
		//导出全部标签
		apiv1.POST("tags/export", v1.ExportTags)
		//导入标签
		apiv1.POST("tags/export_in", v1.ExportInTags)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	{
		r.POST("/api/auth", api.GetAuth)
		r.POST("/api/upload", api.UploadImage)
		//生成二维码
		r.POST("/api/qrcode", v1.GenerateQrCode)
		//生成文字
		r.POST("/api/draw", v1.GenerateDraw)
	}

	//http
	{
		r.GET("/api/get", v1.HttpGet)
	}

	//静态地址
	{
		r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
		r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
		r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))
	}


	return r
}
