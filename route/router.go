package route

import (
	"github.com/gin-gonic/gin"
	"netspace/handler"
)

func Router() *gin.Engine {
	router := gin.Default()

	// 处理静态资源
	router.Static("/static/", "./static")

	// 不需要验证的借口
	router.GET("/user/signup", handler.SignupHandler)
	router.POST("/user/signup", handler.DoSignupHandler)

	router.GET("/user/sign", handler.SignInHandler)
	router.POST("/user/sign", handler.DoSignInHandler)

	// 加入中间件 用于校验token
	router.Use(handler.HTTPInterceptor())

	// 文件存取接口
	router.GET("/file/upload", handler.UploadHandler)
	router.POST("/file/upload", handler.DoUploadHandler)
	router.GET("/file/upload/suc", handler.UploadSucHandler)
	router.POST("/file/meta", handler.GetFileMetaHandler)

	router.POST("/file/download", handler.FileQueryHandler)
	router.POST("/file/update", handler.Downloadhandler)
	router.POST("/file/delete", handler.FileUpdateMetaHandler)
	router.POST("/file/downloadurl", handler.UploadHandler)

	// 秒传接口
	router.POST("/file/fastupload", handler.TryFastUploadHandler)

	// 分块上传接口
	router.POST("/file/mpupload/init", handler.InitialMultipartUpload)
	router.POST("/file/mpupload/uppart", handler.UploadPartHandler)
	router.POST("/file/mpupload/complete", handler.CompleteUploadHandler)

	// 用户相关接口
	router.POST("/user/info",handler.UserInfoHandler)

	return router
}
