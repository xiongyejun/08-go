# wechatcmd  [![star this repo](http://github-svg-buttons.herokuapp.com/star.svg?user=liushuchun&repo=wechatcmd&style=flat&background=1081C1)](http://github.com/liushuchun/wechatcmd) [![fork this repo](http://github-svg-buttons.herokuapp.com/fork.svg?user=liushuchun&repo=wechatcmd&style=flat&background=1081C1)](http://github.com/liushuchun/wechatcmd/fork) ![Build](https://camo.githubusercontent.com/46cb8b3469febc6cdb6fbaea2ef1517c396004e7/68747470733a2f2f7472617669732d63692e6f72672f736a77686974776f7274682f676f6c6561726e2e706e673f6272616e63683d6d6173746572)
=================
## 微信命令行版本
开发这个命令行版本，一是为了熟悉微信的接口，二是方便咱们习惯命令行的同学。

现在是中文的支持不是很好，还没有什么特别好的解决方法。

项目还是远未完成，热烈欢迎有兴趣的朋友一起加入开发。

有什么建议可以提issue。谢谢，也欢迎直接提PR。


### 功能特性

1. 用户检索
2. 聊天表情包快捷键
3. 自动聊天
4. Vimer式快捷键让操作丝般顺滑
5. 更加Geek的feel.


### 键盘快捷键


<table>
    <tr><td>Ctrl-n</td><td>下一页</td></tr>
    <tr><td>Ctrl-p</td><td>上一页</td></tr>
    <tr><td>Enter</td><td>输入</td></tr>
    <tr><td>Ctrl-c</td><td>退出</td></tr>
    <tr><td>Ctrl-1</td><td>退出自动回复</td></tr>
    <tr><td>Ctrl-2</td><td>启用自动回复</td></tr>
    <tr><td>Ctrl-3</td><td>机器人自动回复(还没好)</td></tr>
</table>

### 运行bin文件
linux,mac,windows编好的包分别在install 下面的linux/ mac/ win/下。(windows暂时支持的不好，虽然是交叉编译可以运行，但是其UI机制和unix系差的很多，termui支持的并不是很好)

```
git clone git@github.com:liushuchun/wechatcmd.git
cd wechatcmd/install/
进入各自目录
```


### Mac安装

	$ go get -u github.com/liushuchun/wechatcmd


### Linux安装

	$ go get -u github.com/liushuchun/wechatcmd


### 现在实现的界面：

![聊天动态图](https://raw.githubusercontent.com/liushuchun/wechatcmd/master/img/show.gif)

![登陆后图](https://raw.githubusercontent.com/liushuchun/wechatcmd/master/img/wechatcmd-1.png)
![聊天图片](https://raw.githubusercontent.com/liushuchun/wechatcmd/master/img/wechatcmd-2.png)




### 使用

	$ wechatcmd

### 现在完成的功能
- [x] 实现微信登陆
- [x] 实现微信认证
- [x] 实现拉取用户信息
- [x] 同步消息
- [x] 设置自动回复：正在忙，稍后回来，等等。
- [x] 自动更新消息
- [x] 自动回复消息
- [x] 获取其他消息
- [x] 聊天
- [ ] 群聊
- [ ] 读取图片
- [ ] 自动保存消息到本地
- [ ] 表情包的翻译

### roadmap
 -[ ]群内名字解析 @ Haiyang yu 2017-1-17
 -[ ]自动回复消息修复 @liushuchun 2017-1-17
 -[ ]添加表情解析 @liushuchun 2017-1-25
 -[ ]自动添加好友功能 @liushuchun 2017-2-1
 -[ ]好友、公众号分类 @Haiyang yu 2017-2-1
