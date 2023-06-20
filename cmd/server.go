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
	Use:   "server",
	Short: "启动web服务",
	Run: func(cmd *cobra.Command, args []string) {
		// 创建一个默认的 Gin 引擎
		//======如果你打算生成前端静态文件就从嵌入式里面提取，如果不打算，就默认使用嵌入式里面的前端内容，不用提取======
		// checkStaticFiles()
		r := gin.Default()
		r.Use(PathArrayMiddleware())
		// 静态文件服务
		r.GET("/", Index)
		// 接口服务
		r.GET("/init", Init)
		r.GET("/get/*paths", Get)
		r.GET("/delete/*paths", Delete)
		r.GET("/update/*paths", Update)
		r.GET("/insert/*paths", Insert)
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
		// 获取展示类型
		showType := c.Query("type")
		c.Set("type", showType)
		// 获取json格式
		json := c.Query("json")
		c.Set("json", json)
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
func Init(c *gin.Context) {
	jsonStr := c.GetString("json")
	c.String(200, service.InitJson(jsonStr))
}
func Get(c *gin.Context) {
	keys := c.GetStringSlice("keys")
	outputType := c.GetString("type")
	res, err := service.GetValue(keys, outputType)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, res)
}

func Update(c *gin.Context) {
	keys := c.GetStringSlice("keys")
	value := c.GetString("value")
	res, err := service.UpdateOrInsertValue(keys, value, true)

	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, res)
}

func Insert(c *gin.Context) {
	keys := c.GetStringSlice("keys")
	value := c.GetString("value")
	res, err := service.UpdateOrInsertValue(keys, value, false)

	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, res)
}

func Delete(c *gin.Context) {
	keys := c.GetStringSlice("keys")
	res, err := service.DeleteValue(keys)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.String(200, res)
}

func Index(c *gin.Context) {
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
	
		<h3>1. 功能示例:</h3>
		<ul>
			<li>【初始化】：<a href='/init?json={"data":{"city":"eeeeee","details":[20,20,30],"date":[1,2,3]},"code":"200","msg":"ok"}'>自定义初始化</a>---<a href="/init">默认初始化</a></li>
			<li>【格式化输出】：<a href="/get?type=json">输出json</a>---<a href="/get?type=yaml">输出yaml</a>---<a href="/get?type=xml">输出xml</a></li>
		</ul>
		<hr>
		<h3>2. 增删改查: 增删改都是通过多个key去定位节点，参数value填写需要添加或者修改的json</h3>
		<p>接口格式：/{action[get|insert|update|delete]}/{key}/{key}...</p>
		<ul>
			<li>【查】<a href="/get">/get</a> ---<a href="/get/data">/get/data</a>--- <a href="/get/data/details">/get/data/details</a>--- <a href="/get/data/details/1">/get/data/details/1</a></li>
			<li>【增】<a href="/insert/data/test?value=123">/insert/data/test?value=123</a></li>
			<li>【改】<a href='/update/data/test?value=[{"data":123},{"test":[1,2]}]'>/update/data/test?value=[{"data":123},{"test":[1,2]}]</a></li>
			<li>【删】<a href="/delete/data/test">/delete/data/test</a></li>
			
		</ul>
		<hr>
		<h3>3. 待开发功能</h3>
		<ul>
			<li> 1. 完成简易html页面开发(bootstrap)</li>
			<li> 2. 优化重构代码</li>
			<li> 3. 完善测试 </li>
		</ul>
	</body>
	</html>
	`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}
