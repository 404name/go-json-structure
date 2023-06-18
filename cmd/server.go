package cmd

import (
	"fmt"
	"strings"

	"github.com/404name/go-json-structure/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// VersionCmd represents the version command
var ServerCmd = &cobra.Command{
	Use:   "service",
	Short: "启动web服务",
	Run: func(cmd *cobra.Command, args []string) {
		// 创建一个默认的 Gin 引擎
		//======如果你打算生成前端静态文件就从嵌入式里面提取，如果不打算，就默认使用嵌入式里面的前端内容，不用提取======
		// checkStaticFiles()
		r := gin.Default()
		r.Use(PathArrayMiddleware())
		// 静态文件服务
		r.GET("/", func(c *gin.Context) {
			html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Welcome to go-json-structure-WEB</title>
	</head>
	<body>
		<h1>欢迎使用 go-json-structure-WEB</h1>
		<p>接口列表如下:</p>
		<hr>
		<h3>1. 增删改查: 可以先<a href="/init">初始化系统demo</a>增删改都是通过多个key去定位节点，参数value填写需要添加或者修改的json</h3>
		<ul>
			<li>【查】<a href="/get">/get</a> 或者 <a href="/get/data">/get/data</a> 或者  <a href="/get/data/details">/get/data/details</a></li> 或者 <a href="/get/data/details/1">/get/data/details/1</a></li>
			<li>【增】<a href="/insert/data/test?value=123">/insert/data/test?value=123</a></li>
			<li>【改】<a href='/update/data/test?value=[{"data":123},{"test":[1,2]}]'>/update/data/test?value=[{"data":123},{"test":[1,2]}]</a></li>
			<li>【删】<a href="/delete/data/test">/delete/data/test</a></li>
			
		</ul>
		<hr>
		<h3>2. 服务列表:</h3>
		<ul>
			<li><a href="/insert">传入string返回json</a></li>
			<li><a href="/insert">传入string输出yaml</a></li>
			<li><a href="/demo">传入文件输出yaml</a></li>
		</ul>
	</body>
	</html>
	`
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, html)
		})
		// 接口服务
		r.GET("/init", demo)
		// 增删改查服务[提供给浏览器直接访问使用，不使用restful(post、put那种)的规范]
		r.GET("/get/*paths", get)
		r.GET("/delete/*paths", delete)
		r.GET("/update/*paths", update)
		r.GET("/insert/*paths", insert)
		// 传入string返回json
		// 传入string输出yaml
		// 传入文件输出yaml
		// 接口服务
		r.GET("/ping", func(c *gin.Context) {
			c.String(200, "Hello, World!"+fmt.Sprintf(`
			Version: %s
			flag: %v
			`,
				version, flag))
		})
		// 启动引擎，监听 8080 端口
		r.Run(":8080")
	},
}

func PathArrayMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取value
		value := c.Query("value")
		c.Set("value", value)

		// 获取路径
		paths := c.Param("paths")
		result := strings.Split(paths, "/")

		var pathArray []string
		for _, part := range result {
			if part != "" {
				pathArray = append(pathArray, part)
			}
		}
		c.Set("keys", pathArray)

		c.Next()
	}
}
func demo(c *gin.Context) {
	c.String(200, service.Demo())
}
func get(c *gin.Context) {
	keys := c.GetStringSlice("keys")
	res, err := service.GetValue(keys)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, res)
}

func update(c *gin.Context) {
	keys := c.GetStringSlice("keys")
	value := c.GetString("value")
	res, err := service.UpdateOrInsertValue(keys, value, true)

	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, res)
}

func insert(c *gin.Context) {
	keys := c.GetStringSlice("keys")
	value := c.GetString("value")
	res, err := service.UpdateOrInsertValue(keys, value, false)

	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, res)
}

func delete(c *gin.Context) {
	keys := c.GetStringSlice("keys")
	res, err := service.DeleteValue(keys)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, res)
}
