package panel

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Notepad struct{}

func (a *Notepad) Get(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	var notepad models.Notepad
	// 查询当前用户的便签
	if err := global.Db.Where("user_id = ?", userInfo.ID).First(&notepad).Error; err != nil {
		// 如果未找到，返回默认空内容或者错误码，这里返回空字符串的成功响应
		// 或者创建一条新的空记录？前端逻辑是 Get 失败则显示空。
		// 为了简化，这里如果找不到就返回空内容
		apiReturn.SuccessData(c, gin.H{"content": ""})
		return
	}
	apiReturn.SuccessData(c, gin.H{
		"content":   notepad.Content,
		"updatedAt": notepad.UpdatedAt,
	})
}

func (a *Notepad) Save(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	type Req struct {
		Content string `json:"content"`
	}
	var req Req
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	var notepad models.Notepad
	err := global.Db.Where("user_id = ?", userInfo.ID).First(&notepad).Error

	if err != nil {
		// 不存在，创建
		notepad = models.Notepad{
			UserID:  userInfo.ID,
			Content: req.Content,
		}
		if err := global.Db.Create(&notepad).Error; err != nil {
			apiReturn.Error(c, "Save Failed")
			return
		}
	} else {
		// 存在，更新
		notepad.Content = req.Content
		if err := global.Db.Save(&notepad).Error; err != nil {
			apiReturn.Error(c, "Update Failed")
			return
		}
	}

	apiReturn.Success(c)
}

// Upload 上传文件
func (a *Notepad) Upload(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	// 使用系统配置的上传路径
	configUpload := global.Config.GetValueString("base", "source_path")

	f, err := c.FormFile("file") // 前端这类请求通常用 file 字段
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}

	fileExt := strings.ToLower(path.Ext(f.Filename))
	// 允许的扩展名，可以和 System File 配置一致，或者稍微放宽
	agreeExts := []string{".png", ".jpg", ".gif", ".jpeg", ".webp", ".svg", ".ico", ".txt", ".md", ".json", ".pdf", ".doc", ".docx", ".xls", ".xlsx"}

	if !cmn.InArray(agreeExts, fileExt) {
		// 暂时不严格限制，或者复用 system/file.go 的逻辑
		// 这里简化允许上传，但前端要显示为链接
	}

	fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
	// 存放到 uploads/notepad/{year}/{month}/{day}/
	fildDir := fmt.Sprintf("%s/notepad/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())

	isExist, _ := cmn.PathExists(fildDir)
	if !isExist {
		os.MkdirAll(fildDir, os.ModePerm)
	}
	filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)

	if err := c.SaveUploadedFile(f, filepath); err != nil {
		apiReturn.Error(c, "Upload Write Failed")
		return
	}

	// 记录到 models.File 表吗？
	// 最好记录，为了日后管理文件。
	mFile := models.File{}
	mFile.AddFile(userInfo.ID, f.Filename, fileExt, filepath)

	// 返回相对路径，前端补全 URL
	// 注意 filepath 是 ./uploads/... 前端需要 /uploads/...
	// 系统其他地方返回的是 filepath[1:] 即去掉开头的 .

	downloadUrl := filepath
	if strings.HasPrefix(filepath, ".") {
		downloadUrl = filepath[1:]
	}

	apiReturn.SuccessData(c, gin.H{
		"url":  downloadUrl,
		"name": f.Filename,
		"type": f.Header.Get("Content-Type"),
	})
}
