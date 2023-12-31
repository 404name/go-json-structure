

# [Json解析器(Go语言版)](https://github.com/404name/go-json-structure)
> 使用go语言，完成一个简单的json解析器命令行工具，支持对json的序列化和反序列化，提供简易的交互UI，同时能对外能提供稳定、安全、统一的api接口。

技术栈：`go语言` `cobra命令行框架` `gin Web框架` `git版本管理` 

![](/docs/gson.gif)

1. 提供命令行工具
```
Usage:
  gson [command]

Available Commands:
  server      启动web服务
  init        初始化json数据
  delete      删除内容
  get         获取内容
  insert      添加内容
  update      修改内容
  save        保存内容
  version     显示当前版本
  help        Help about any command

Flags:
      --debug               start with debug mode
  -h, --help                help for gson
      --json string         自定义json数据源(需要满足格式要求)
      --keys string         查询的内容/xxxx/xxx/xxx的形式 (default "/")
      --output string       输出格式[json|yaml|xml] (default "json")
      --outputPath string   输出地址，默认为output.txt (default "./output.txt")
      --value string        增和改设置的json值
      --version string      go-json-structure[BETA] (default "1.0.0")
```

2. 提供WEB接口
```

1. 功能示例:
【初始化】：自定义初始化---默认初始化
【格式化输出】：输出json---输出yaml---输出xml
2. 增删改查: 
接口格式：/{action[get|insert|update|delete]}/{key}/{key}...
【查】/get ---/get/data--- /get/data/details--- /get/data/details/1
【增】/insert/data/test?value=123
【改】/update/data/test?value=[{"data":123},{"test":[1,2]}]
【删】/delete/data/test
```
# 一、需求分析

## 开发进度
0. 完成需求分析 [2023年6月16日]
1. 完成核心逻辑开发[2023年6月17日]
2. 完成单元测试[2023年6月17日]
3. 接入cobra封装成命令行工具[2023年6月18日]
4. 接入gin框架提供对外api接口[2023年6月18日]
5. 接入简易的WEB界面[2023年6月20日]

**1. 基础功能**
- [x] 文件读入json
- [x] 输出json为yaml
- [x] 单元测试覆盖率不低于90%
- [ ] 支持异常处理提示 

**2. api功能**
- [x] 对json的内容增删改查
- [x] json多种读入方式
- [x] json多种输出方式

**3. 附加功能**
- [x] 做成命令行工具(使用cobra实现)
- [ ] 提供UI界面(Go的ui库可以提供)
- [x] 提供api接口(使用gin框架实现) 

# 二、项目介绍

## 0. 项目架构

1. cobra实现命令行框架，允许命令行不同指令进行不同的交互
2. 通过cobra开通一个service指令绑定gin框架启动，就可以完成命令行开启api服务了
3. cobra提供的服务和gin框架服务共用底层的service服务
4. json核心逻辑通过json包实现，采取parser解析器，并且完成100%的单测。


> 此处描述一下技术栈及里面一些点如何实现的
## 1.json对应的数据结构

此处参考java阿里的JSONObject，并且进行了整合，没有区分JSONArray
```
//json树的节点
type JSONObject struct {
	vType int  // null、false、true、数字、字符串、数组、对象 7种类型
	num   float64
	value []byte
	list  []*JSONObject
	obj   map[string]*JSONObject
}
```
## 2.json的解析器parser

解析器核心逻辑就是处理一个string，去除/n空格等内容后递归解析获取得到JSONObject。
- 核心逻辑 string -> JSONObject
- 提供多种序列化、反序列化接口(文件读入、json读入、文件输出、yaml输出) 

## 3.JSONObject对外提供的方法

此处依然参考阿里的JSONObject，假设有一个`var obj JSONObject` 


**基础使用**
- 需要判断值属性使用`obj.Type()`方法 
- 如果是对象需要获取它的子节点使用`obj.Get(key)`方法
- 获取数组的某个节点使用`obj.GetIndex(idx)` 获取值 
- 如果要获取值使用 `obj.Value()` 但这种返回值需要断言
- 获取准确的值使用 `obj.GetString()`(其它类型的例如GetBool、GetNumber也支持) 这种前提是已知属性

**注意**

- 因为json不是强类型的语法，所以后端拿到json其实也不知道具体属性和内容，所以该判断还是得判断
- 部分逻辑情况下后端会知道该读取哪个字段，并且这个字段是那种值是确定的，所以也应该提供直接获取字段的接口

## 4. api及命令行 的接口约定

这里主要是思考如何通过url实现json的增删改查

- 查的话很简单，通过多个key去定位就行
- 增的话只允许【数组、对象】节点去添加元素
- 添加的元素对数组是单个值，对对象是k-v对象，但都可以通过字符串的json写入
- 修改可以允许所有类型的修改，其实也就是通过多个key定位后直接改json

所以可以得出如下的接口定义

【增】/insert/xxx?value=xxx
【删】/delete/xxx/xxx
【改】/update/xxx/xxx?value=xxx
【查】/get/xxx/xxx


# 三、测试内容

1. 支持BOOL值，数值，字符串，数组，对象
2. 构建如下json，并且输出为yaml
```
{
    "basic": {
        "enable": true,
        "ip": "123.123.123.123",
        "port": 389,
        "timeout": 10,
        "basedn": "aaa",
        "fd": -1,
        "maxcnt": 133333333333,
        "dns": ["123.123.123.123", "123.123.123.123"]
    },
    "advance": {
        "dns": [
            {"name":"huanan", "ip": "123.123.123.123"}, 
            {"name":"huabei", "ip": "123.123.123.123"}],
        "portpool": [123,33,4],
        "url": "http://123.123.123.123/main",
        "path": "/etc/sinfors",
        "value": 3.14
    }
}
```
basic.enable
# 四、重构阶段

## 掌握基础设计方法
> 待学习
> 接口和数据结构的设计、掌握一种设计/实现通用数据结构，以及为模块设计API的方法。

## 掌握编码六步法
> 待学习
> 掌握6步法这种基本的编程套路


## 编码风格优化
> 需要考虑API的适用场景是什么，怎么让设计出来的API更好用，更容易调测，更容易扩展，更可靠。

## 考虑安全相关

- [ ]  内存泄露相关
- [ ] api接口安全


 
# 📑 开发日志
> 简易记录下开发日志
<details>
<summary>[v0.1.0] : 完成命令行工具+Web接口 [2022-6-20] </summary>

- 【feat】完成命令行工具，开箱即用
- 【feat】完善WEB端逻辑，提供更多json输出的功能
- 【test】完善单测

</details>
<details>
<summary>[v0.0.2] : cobra+gin框架支持命令行及api操作 [2022-6-18] </summary>

- 【feat】接入cobra框架，支持命令行增删改查json及提供json转换服务
- 【feat】接入gin框架，支持浏览器增删改查json及提供json转换服务
- 【feat】定义增删改查规范
- 【test】完成所有方法的单测

</details>
<details>
<summary>[v0.0.1] : 完成核心逻辑开发&单测覆盖 [2022-6-17] </summary>

- 【feat】完成json核心组件parse解析器功能
- 【feat】完成json的增删改查接口，支持7种类型
- 【feat】支持多种读入方式
- 【test】完成所有方法的单测

</details>
<details>

<summary>[v0.0.0] : 阅读任务内容&需求分析 [2022-6-16] </summary>

- 【需求分析】编写完README文档确认开发任务
- 【需求分析】选择技术栈，确定预期开发功能

</details>

# 📖 学习参考
1. https://zhuanlan.zhihu.com/json-tutorial 从零开始教授如何写一个符合标准的 C 语言 JSON 库
2. https://www.cnblogs.com/tanshaoshenghao/p/14100735.html 手写Json解析器学习心得