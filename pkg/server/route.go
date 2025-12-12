package server

import (
	"encoding/json"
	"html/template"

	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {
	// 注册自定义模板函数
	engine.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"seq": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		"js": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

	// 加载 HTML 模板
	engine.LoadHTMLGlob("pkg/tpl/*.tpl")

	// 文章列表
	engine.GET("/articles", GetArticleList)
	// 文章详情
	engine.GET("/articles/:id", GetArticleDetail)
	// 根路径重定向到文章列表
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/articles")
	})
}
