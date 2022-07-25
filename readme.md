# GO BBS 说明文档
<img alt="GO BBS" src="http://go-bbs.com/image/logo.png">

嗯嗯嗯～ 想了很久还是觉得把后端所有的代码都开源出来，希望与大家一起开发一个更加完善的 go bbs 系统。近段时间我也在逐渐完善我们的文档，希望能让所有的人都能轻松上手以及修改自定义属于自己喜欢的系统。
当然，在这里也感谢大家的支持，如果可以的话，希望动动小手点个星星哦 。

现在的 github 地址 https://github.com/gobbscom/go-bbs (原github https://github.com/1920853199/go-blog)

## Go Bbs 官网
http://go-bbs.com/

## Demo
https://nihongdengxia.com

http://clblog.club/

## Go Bbs 是什么?
Go Bbs 是一个基于 Beego 开发的可切换模板的 BBS 社交博客系统。只要把模板文件放在`views/home` 目录下我们就能在后台轻松切换模板。它安装简单便捷，页面简洁优美。前端是HTML+JS+CSS，不需要掌握一些前端技术栈也能轻松自定义页面。

## 安装
1.  把 Go Bbs 项目拉到本地 `git clone https://github.com/gobbscom/go-bbs.git`
2.  进入项目的根目录下执行 `go build -o go-bbs`
3.  **在根目录创建一个配置文件夹conf，将app.conf 内容如 app.conf.example放入conf文件夹中，配置好环境变量 BEEGO_CONFIG_PATH = '配置文件路径'**
4.  执行./go-bbs --install 安装数据库
5.  最后执行 ./go-bbs 访问对应端口即可

> 如果需要启动爬虫任务。在启动的时候 指定 `--t` 为 `true` , 如：`./go-bbs --t ture`
>  管理后台账号密码：admin , 123456

> 大家遇到问题的时候可以翻翻 Issues 列表，说不定有同学会遇到跟你一样的问题并解决了哦


## GO Bbs 模板有哪些?

<table>

<thead>

<tr>

<th> nihongdengxia 模板 </th>

<th> leechan 模板 </th>

<th> goooooooogle 模板 </th>

<th> toutiao 模板 </th>

</tr>

</thead>

<tbody>

<tr>

<td  style="text-align:left;width: 25%;">


![go-bbs 2.0 版本发布，官网以及文档正式上线。后端代码全部开源！！！](https://cdn.learnku.com/uploads/images/202204/21/43046/Rh56egjbGH.png!large)

</td>

<td  style="text-align:left;width: 25%;">


![go-bbs 2.0 版本发布，官网以及文档正式上线。后端代码全部开源！！！](https://cdn.learnku.com/uploads/images/202204/21/43046/4mcA3Cigzs.png!large)

</td>

<td  style="text-align:left;width: 25%;">


![go-bbs 2.0 版本发布，官网以及文档正式上线。后端代码全部开源！！！](https://cdn.learnku.com/uploads/images/202204/21/43046/oYkci4k4c7.png!large)

</td>

<td  style="text-align:left;width: 25%;">


![go-bbs 2.0 版本发布，官网以及文档正式上线。后端代码全部开源！！！](https://cdn.learnku.com/uploads/images/202204/21/43046/bWaWnzLTs6.png!large)

</td>

</tr>

<tr>

<td  style="text-align:left;width: 25%;text-align: center;">

<a  href="https://nihongdengxia.com/"  target="_blank"  title="霓虹灯下社区">Demo</a>

</td>

<td  style="text-align:left;width: 25%;text-align: center;"></td>

<td  style="text-align:left;width: 25%;text-align: center;"></td>

<td  style="text-align:left;width: 25%;text-align: center;">

<a  href="http://clblog.club/"  target="_blank"  title="陈立个人博客">Demo</a>

</td>

</tr>

</tbody>

</table>

## Go Bbs 包含了哪些功能模块?
> 包括但不止下面列出的功能，还有很多这里没有列出来

> 注意：后端代码是全部开源的，所以可以根据自己的需要在对应模板上加上该功能
<table>

<thead>

<tr>

<th  style="width: 40%;"> 功能 </th>

<th> nihongdengxia 模板 </th>

<th> leechan 模板 </th>

<th> goooooooogle 模板 </th>

<th> toutiao 模板 </th>

</tr>

</thead>

<tbody>

<tr>


<td>文章列表/分类</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

</tr>

<tr>

<td>文章评论</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

</tr>

<tr>

<td>文章浏览量记录</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

</tr>

<tr>

<td>文章点赞</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

</tr>

<tr>

<td>热点文章</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

</tr>

<tr>

<td>最新文章</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

</tr>

<tr>

<td>文章推荐</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

</tr>

<tr>

<td>文章Tag</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

<td>✓</td>

</tr>

<tr>

<td>系列文章</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>  <tr>

<td>发布/修改文章(前台)</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>文章同步到社区</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>


<td>图片轮播</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

<td>✓</td>

</tr>

<tr>

<td>广告位</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>


<td>话题分类</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>话题广场</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>创建/参与话题</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>


<td>图片/图文/文字贴子发布</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>评论/点赞/回复贴子</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>全部/推荐列表</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>


<td>书架</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>新增书籍</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>设置书籍状态</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>点赞书架</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>



<td>✗</td>

</tr>

<tr>


<td>登陆注册</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>个人中心</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>修改资料</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>关注/粉丝</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>我的动态/文章</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>个性签名</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>积分累计</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>会员</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>


<td>虎扑热点/赛事爬虫</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

<td>✗</td>

</tr>

<tr>

<td>百度热点爬虫</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

<td>✗</td>

</tr>

<tr>

<td>豆瓣书籍/热点爬虫</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

<td>✗</td>

</tr>

<tr>

<td>知乎热点爬虫</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

<td>✗</td>

</tr>

<tr>

<td>360历史上得今天爬虫</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

<td>✗</td>

</tr>

<tr>

<td>Gocn老版文章爬虫</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

<td>✗</td>

</tr>

<tr>


<td>版权声明</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>关于我们联系我们</td>

<td>✓</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

</tr>

<tr>

<td>友情链接</td>

<td>✓</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

</tr>

<tr>

<td>网站公告</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

<td>✗</td>

</tr>

<tr>

<td>留言</td>

<td>✗</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

</tr>

<tr>


<td>搜索文章/贴子/用户/话题</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>IP限制</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>消息通知</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>精选图片</td>

<td>✓</td>

<td>✗</td>

<td>✗</td>

<td>✗</td>

</tr>

<tr>

<td>热榜大全</td>

<td>✓</td>

<td>✗</td>

<td>✓</td>

<td>✗</td>

</tr>

</tbody>

</table>
