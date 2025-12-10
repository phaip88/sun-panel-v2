package panel

import (
	"sun-panel/api/api_v1/middleware"
	"sun-panel/api/api_v1/panel"

	"github.com/gin-gonic/gin"
)

func InitNotepad(router *gin.RouterGroup) {
	notepad := panel.Notepad{}
	// 使用登录拦截器
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("panel/notepad/get", notepad.Get) // 改用 POST 统一风格？Bookmark 用的 POST. 但 Notepad get 通常 GET. 还是 GET 吧。不过前端之前用 POST 吗？前端没写呢。
		// 统一把 GET 改为 POST？bookmark getList 是 POST。
		// 让 get 也支持 POST 吧，或者保持 GET。
		// 为了省事避免 params 传参问题，很多项目都全是 POST。
		// 但 Get 便签没参数，用 GET 很合适。

		// 修正：在 api/api_v1/panel/notepad.go 里 Get 方法
		// "sun-panel/api/api_v1/common/apiReturn" 通常 SuccessData(c, data)

		// 让我们保持 GET /panel/notepad/get
		// POST /panel/notepad/save
		// POST /panel/notepad/upload

		// Wait, bookmark Use POST for getList.
		// I will use POST for get as well to be safe with this project's style?
		// No, GET is standard. I'll use GET.
	}

	// 重新写：
	r.GET("panel/notepad/get", notepad.Get)
	r.POST("panel/notepad/save", notepad.Save)
	r.POST("panel/notepad/upload", notepad.Upload)
}
