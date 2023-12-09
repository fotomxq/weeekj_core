package RouterMidWeb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//OutputThemeHtml 输出模版页面
func OutputThemeHtml(c *gin.Context, path string, obj interface{}) {
	c.HTML(http.StatusOK, fmt.Sprint(c.MustGet("themeMark").(string), path), obj)
}