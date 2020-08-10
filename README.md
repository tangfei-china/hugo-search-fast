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

> 這里只介绍部署搜索功能的内容
>
> 部署环境介绍：
>
> 本地电脑：Mac 10.13
>
> 服务器：CentOS7 512M

* 服务部署Sonic服务，建议服务器好的可以直接用Docker的方式加载，這里介绍源代码本地交叉编译和部署

  * git下载https://github.com/valeriansaliou/sonic 源代码到本地电脑，因为sonic是rust语言开发，所以需要先装rust环境

  * 本地电脑安装rust环境

    * 如果是中国开发需要配置下数据源，mac 是在 ~/.zshrc， linux 是在 /etc/profile

      ```shell
      export RUSTUP_DIST_SERVER=https://mirrors.sjtug.sjtu.edu.cn/rust-static
      export RUSTUP_UPDATE_ROOT=https://mirrors.sjtug.sjtu.edu.cn/rust-static/rustup
      ```

    * curl https://sh.rustup.rs -sSf | sh

    * 安装完成尝试使用 cargo --version 看看是否有版本出现，验证是否安装完成。

    * 配置crates的国内数据源（限中国）

      ```toml
      [source.crates-io]
      replace-with = 'ustc'
      
      [source.ustc]
      registry = "git://mirrors.ustc.edu.cn/crates.io-index"
      ```

  * 如果服务器比较的好，可以直接在服务器上编译sonic，服务器不好直接跳过此步骤

    ```shell
    #在服务器上下下载源代码，然后安装rust，执行下面的代码来编译安装
    cargo build --release
    ```

  * 使用 cross 来交叉编译项目到指定的目标系统

    * 运行cross的命令需要在sonic的源代码目录下执行

    * 使用cross前提需要本地的电脑有docker的环境，并且是开启的状态

    * ```shell
      #安装cross
      cargo install cross
      ```

    * ```shell
      #通过test 验证是否有编译的环境，是否有报错
      cross test --target x86_64-unknown-linux-gnu
      ```

    * ```shell
      #编译sonic程序 指定目标及其是 x86_64linux
      #自己部署的机器就是编译的目标
      cross build --target x86_64-unknown-linux-gnu
      ```

    * 编译成功后，只需要拷贝两个文件到自己的服务就可以了，一个是配置文件config.cfg在源代码的根目录下的，还有一个就是可执行文件sonic，在target/x86_64-unknown-linux-gnu下

    * 配置config.cfg

      ```toml
      # Sonic
      # Fast, lightweight and schema-less search backend
      # Configuration file
      # Example: https://github.com/valeriansaliou/sonic/blob/master/config.cfg
      
      [server]
      log_level = "info"
      
      # 是本地的一个监听地址
      [channel]
      inet = "0.0.0.0:1491"
      tcp_timeout = 300
      
      #这个是授权的密码 很重要，后面的channel访问需要用到这个密码
      auth_password = "SecretPassword"
      [channel.search]
      
      query_limit_default = 10
      query_limit_maximum = 100
      query_alternates_try = 4
      
      suggest_limit_default = 5
      suggest_limit_maximum = 20
      
      #这个kv的一个存储路径，建议就是在程序目录下，也可以指定你需要的目录
      [store]
      [store.kv]
      path = "data/store/kv/"
      
      retain_word_objects = 1000
      [store.kv.pool]
      inactive_after = 1800
      [store.kv.database]
      flush_after = 900
      
      compress = true
      parallelism = 2
      max_files = 100
      max_compactions = 1
      max_flushes = 1
      write_buffer = 16384
      write_ahead_log = true
      
      #这个fst的一个存储路径，建议就是在程序目录下，也可以指定你需要的目录
      [store.fst]
      path = "data/store/fst/"
      [store.fst.pool]
      inactive_after = 300
      [store.fst.graph]
      consolidate_after = 180
      max_size = 2048
      max_words = 250000
      ```

    * 在执行文件的目录下执行程序，如果运行有报错请看项目提示的内容

      ```shell
      #這里因为配置和执行文件放在一起
      ./sonic -c config.cfg
      ```

  * 程序运行了，可以使用telnet来简单的测试下索引创建和搜索

    详细的访问和介绍请看：https://github.com/valeriansaliou/sonic/blob/master/PROTOCOL.md

    ```shell
    telnet localhost 1491
    START search SecretPassword
    QUERY collection bucket "搜索的关键字" LIMIT(10)
    ```

* 部署search搜索服务

  * 在search的项目文件中找到build.sh脚本，如果是linux的服务器，可以直接执行此脚本，会生成一个main的可执行文件

    ```shell
    #根据服务器的系统不同改变下面的参数
    #GOOS：目标操作系统
    #GOARCH：目标操作系统的架构
    #CGO_ENABLED=0的意思是使用C语言版本的GO编译器，参数配置为0的时候就关闭C语言版本的编译器了
    CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o search
    ```

  * 拷贝可执行文件和配置文件conf.yaml到服务上

  * 配置文件的更新

    ```yaml
    #sonic的服务地址，就是sonic监听的地址，如果和sonic引擎同一台机器上，這里不用改
    sonic_host: localhost
    #sonic的密码，这个是sonic的配置文件中的密码，需要完全一致
    sonic_password: SecretPassword
    #sonic的集合名称，这个是需要和创建索引的searchIndex的collection配置一致
    sonic_collection: demo
    #sonic的桶名称，这个需要和创建索引的searchIndex的bucket配置一致
    sonic_bucket: demo
    sonic_query_limit: 10
    sonic_query_offset: 0
    #服务端口,对外服务API的端口配置
    search_port: 8089
    ```

  * 执行可执行文件

    ```shell
    #执行后会监听服务器的8089的端口服务
    ./search
    
    #如果配置在其他的目录下，需要指定目录地址
    ./search -c="xxxxx"
    ```

  * 这个搜索服务可以用nginx来代理开放

* 部署searchIndex的索引创建程序

  * 此程序可以在服务器上运行，不过会下耗电服务器资源，所以演示博客使用了本地的执行，因为静态博客的文章和服务器上的是一样的。

  * 在searchIndex的项目文件中找到build_mac.sh的脚本文件，创建可执行文件

  * 查看运行脚本内容

    ```shell
    # Mac 下执行
    go build -o search_index
    ```

  * 拷贝可执行文件main和配置文件conf.yaml

  * 配置文件更新conf.yaml

    ```yaml
    #sonic的服务地址,因为本次使用本地运行程序，需要连接远程的sonic的服务器，這里需要配置你的sonic的远程地址
    sonic_host: localhost
    #sonic的密码，这个密码也是sonic的服务器的密码，需要完全一致
    sonic_password: SecretPassword
    #sonic的集合名称，這里需要和search的服务配置的collection名称一致
    sonic_collection: demo
    #sonic的桶名称，這里需要和search的服务配置的bucket的名称一致
    sonic_bucket: demo
    
    #下面都是一些业务需求的配置
    #需要去除的地址和post_path对应的
    #  例如：root_path: /A
    #       post_path: /A/B
    #  结果的文章路径就是   /B
    # root_path是你的huto配置的文章的路径，和post_path的组合使用，是为了后面的搜索内容中的访问地址
    # 例如：post/2020/07/23-2/  这个就是一个文章的地址，目前项目的文章地址配置的是和存储一致的。
    root_path: hugo生成的根路径/content
    #需要检索的静态博客父目录
    post_path: hugo生成的根路径/content/post
    # true 所有博客重新创建索引 false 创建最新博客的索引
    # 如果是false的话会采集到当前的时间前2小时内的博客来创建索引
    index_all_init: false
    ```

  * 在每次写完博客，执行下执行程序，就可以创建索引了

  * 创建索引的规则 - **第一个版本**

    **下面是一个hugo生成的文章模板**

    **关键字的内容是取决于模板中的description内的内容**

    **索引的唯一是采集的title的内容**

    **模板中的author目前是截止位置，后续会更新优化下，可以自定义**

    ```markdown
    ---
    title: "Hugo静态博客搜索功能开发-上集"
    date: 2020-08-05T14:53:08+08:00
    lastmod: 2020-08-05T14:53:08+08:00
    draft: false
    keywords: ["hugo搜索内容","静态博客搜索功能","搜索内容开发"]
    description: "hugo静态博客的搜索功能，自定义搜索功能，使用sonic作为底层搜索引擎，使用gin作为搜索平台，后续准备当开源项目开源出来，让更多的人可以使用，也提出更好的建议和代码。"
    tags: ["hugo"]
    categories: ["网站"]
    author: ""
    
    # You can also close(false) or open(true) something for this content.
    # P.S. comment can only be closed
    
    comment: true
    toc: true
    autoCollapseToc: false
    postMetaInFooter: true
    hiddenFromHomePage: false
    # You can also define another contentCopyright. e.g. contentCopyright: "This is another copyright."
    contentCopyright: false
    reward: true
    mathjax: false
    mathjaxEnableSingleDollar: false
    mathjaxEnableAutoNumber: false
    
    # You unlisted posts you might want not want the header or footer to show
    hideHeaderAndFooter: false
    # You can enable or disable out-of-date content warning for individual post.
    # Comment this out to use the global config.
    #enableOutdatedInfoWarning: false
    
    flowchartDiagrams:
      enable: false
      options: ""
    sequenceDiagrams: 
      enable: false
      options: ""
    ---
    ```
    
  * 运行

    ```shell
    # 如果执行程序和配置在一个目录下，可以不配置-c参数
    ./search_index -c="xxxxxxxx"
    ```

    

* 怎么使用这个搜索服务，请看项目案例的介绍

## 项目案例

> 這里介于演示内容博客的搜索功能添加过程，作为参考
>
> 此博客使用的是hugo静态工具生成的博客，主题是使用的 https://github.com/olOwOlo/hugo-theme-even

* 首先在content/post 文件下创建一个search.html的页面

* 在config.toml配置文件中，配置菜单项，增加一项，配置如下：

  ```toml
  [[menu.main]]
    name = "搜索"
    weight = 8
    identifier = "search"
    url = "/search/"
  ```

* 配置一些静态资源，方便search.html页面的功能开发

  ```toml
    # Link custom CSS and JS assets
  
    #   (relative to /static/css and /static/js respectively)
  
    customCSS = ["index.css"]
    customJS = ["vue.js","index.js","axios.min.js"]
  ```

  这里的customCSS 就是你需要的自定义样式文件

  這里的customJS就是你需要的自定义的脚本文件

  演示的博客使用的是ElementUI和Vue组合来使用。所以引入一些必要的地址，大家也可以直接使用官方开发的地址，那么這里就不用配置了

* 增加静态的样式文件和脚本文件

  <img src="https://sevenbooks.oss-cn-hangzhou.aliyuncs.com/postimages/20200806134624.png" style="zoom:50%;" />

  在themes/even/static文件下增加css和js文件夹，并且把对应的文件放进去，這里提供的是此次演示的一些地址：

  样式地址：https://unpkg.com/element-ui/lib/theme-chalk/index.css

  脚本地址：https://unpkg.com/vue/dist/vue.js

  脚本地址：https://unpkg.com/element-ui/lib/index.js

  如果需要的话可以下载保存起来，然后放到刚才新建的目录下使用。

* 开始写search.html的页面代码，作为参考：

  ```html
  ---
  title: "搜索"
  date: 2019-10-09T16:17:39+08:00
  draft: false
  ---
  
  <link rel="stylesheet" href="/css/index.css">
  <style>
      .el-autocomplete .el-input {
          width: 100%;
      }
      .text {
          font-size: 14px;
      }
      .item {
          margin-bottom: 18px;
      }
      .clearfix:before,
      .clearfix:after {
          display: table;
          content: "";
      }
      .clearfix:after {
          clear: both
      }
      [v-cloak] {
          display: none;
      }
      #app a,
      a:link,
      a:visited,
      a:hover,
      a:active {
         text-decoration: none !important;
         color: black !important;
         border: none;
     }
  </style>
  
  <div id="app" v-cloak>
      <el-card class="box-card">
          <div slot="header" class="clearfix">
              <el-input placeholder="请输入关键字" @keyup.enter.native="search" v-model="searchKey" class="input-with-select">
                  <el-button v-on:click="search" slot="append" icon="el-icon-search"></el-button>
            </el-input>
          </div>
         <div v-if=" status !== '1' " class="text item">
              {{ message }}
          </div>
          <div v-else v-for="item in list" :key="item.Name" class="text item">
             <el-link v-bind:href="item.Path" target="_blank"> {{ item.Name }}</el-link>
          </div>
      </el-card>
  </div>
  
  <!-- import Vue before Element -->
  <script src="/js/vue.js"></script>
  <!-- import JavaScript -->
  <script src="/js/index.js"></script>
  <script src="/js/axios.min.js"></script>
  <script type="module">
      new Vue({
          el: '#app',
          data() {
              return {
                  searchKey: '',
                  list: [],
                  message: '',
                  status: '0'
              };
          },
          methods: {
              search: function () {
                  let that = this
                  // Make a request for a user with a given ID
                  axios.get('https://sh.7benshu.com/search?key=' + this.searchKey)
                      .then(function (response) {
                          // handle success
                          that.status = response.data.status                        
                          if (response.data.status == "1") {
                              that.list = response.data.message
                              if (response.data.message == null) {
                                  that.message = '博主太懒了，没有记录此类文章，请重新输入关键字！'
                                  that.status = '0'
                              }
                          } else {
                              that.message = response.data.message
                          }
                      })
                      .catch(function (error) {
                          that.message = '系统开小差了！'
                          that.status = '0'
                      })
                      .then(function () {
                          // always executed
                      });
              }
          }
      })
  </script>
  ```

  以上代码在Even主题下是测试通过的。如果是其他的主题，大家需要自己测试了

## 项目提示

* 部署后执行出现错误信息：/lib64/libstdc++.so.6: version `CXXABI_1.3.8’ not found

  > ```shell
  > #执行这个命令看看是否不存在,判断是不是gcc版本问题
  > strings /usr/lib64/libstdc++.so.6 | grep CXXABI
  > ```
  >
  > 确定不存在，我们需要手动升级gcc版本，CentOS7 默认是gcc4.8
  >
  > ```sh
  > sudo yum install gmp-devel mpfr-devel libmpc-devel -y 
  > wget ftp://ftp.gnu.org/gnu/gcc/gcc-9.2.0/gcc-9.2.0.tar.xz 
  > xz -d gcc-9.2.0.tar.xz 
  > tar -xf gcc-9.2.0.tar 
  > cd gcc-9.2.0 
  > ./configure --disable-multilib --enable-languages=c,c++ --prefix=$HOME/local 
  > make 
  > make install
  > #这个编辑时间会很长，取决你的机器性能，演示的博客512M的机器2核心的，大概5个小时左右编译好的。
  > ```
  >
  > 配置环境变量 vim /etc/profile
  >
  > ```shell
  > export LD_LIBRARY_PATH=$HOME/local/lib64
  > PATH=$HOME/local/bin:$PATH
  > ```
  >
  > 执行命令：
  >
  > ```shell
  > source /etc/profile
  > ```

* 如果在廉价的机器上编译gcc，那么会出现内存不足

  ```shell
  #如果gcc版本低了，低内存的机器编译，肯定会出现内存不足，所以需要增加交换区，這里根据自己的需要配置交换大小
  dd if=/dev/zero of=/swapfile bs=64M count=60
  mkswap /swapfile
  swapon /swapfile
  
  #这个是编译结束了，可以关闭和删除交换区
  swapoff /swapfile
  rm /swapfile
  
  #這里可以让重启了也可以开启交换区，如果只是编译下，可以不用设置
  /etc/fstab
  /swapfile swap  swap  defaults 0 0
  
  #如果配置了交换区，还是出现内存不足，這里小配置下100，让系统积极的使用交换区
  sysctl vm.swappiness=100
  ```

  

## 项目依赖

https://github.com/expectedsh/go-sonic

https://github.com/valeriansaliou/sonic

https://github.com/gin-gonic/gin

