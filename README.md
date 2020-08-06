# hugo-search-fast
增加hugo静态博客的搜索功能，可以在廉价的机器上运行，同时支持高性能的搜索。

体验地址：https://www.7benshu.com/search/

体验服务器配置

> 系统： CentOS7.5 64位
>
> 硬件： 2 vCPU 512MB

## 项目目录

两个模块可以单独打开，也可以根目录同时打开，开发工具推荐GoLand

* searchIndex 
  * 主要用于博客内容的索引创建
  * 程序使用Golang语言开发
  * 索引使用的是sonic的搜索引擎
  * 交互使用的是go-sonic组件

* search
  * 主要用于博客内容的搜索
  * 交互使用的是go-sonic组件
  * 基于Go语言开发的Gin Web框架来对外开放搜索API

## 项目架构图

![](https://sevenbooks.oss-cn-hangzhou.aliyuncs.com/postimages/20200805210730.png)

这个图包含了一个使用hugo静态博客的使用全流程

* 写博客-到发布
  * ssh是本地电脑，通过hugo写静态博客
  * 执行发布的脚本，然后页面可以浏览，脚本执行过程如下
    * hugo编译md格式博客编程静态页面
    * 通过git 上传到github上
    * 通过ssh远程发送指令给远程服务器拉取git的资源
    * 通过ssh本地执行searchIndex的脚本创建新增博客的索引内容到服务器上的Sonic搜索引擎上（这个可以在服务器上运行，因为本体验服务器配置太低了，所以放到了本地执行了）
* 搜索博客
  * 进入博客搜索页面
  * 通过页面的搜索功能来访问search程序提供的API,然后返回搜索内容，**此页面需要用户自定义开发**，这块内容会在项目案例中具体描述怎么开发
  * search是通过go-sonic和sonic搜索引擎交互得到匹配内容

## 项目部署



## 项目案例



## 项目依赖

https://github.com/expectedsh/go-sonic

https://github.com/valeriansaliou/sonic

https://github.com/gin-gonic/gin