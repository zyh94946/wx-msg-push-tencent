
## 基于腾讯云云函数和API网关实现的企业微信应用消息推送服务

Serverless 云函数目前每月有免费资源使用量40万GBs、免费调用次数100万次

API网关目前开通即送时长12个月100万次免费额度

个人或者低频率使用完全够了，可以通过 GET、POST 方式调用发消息。

对于有服务器、域名资源，通过简单修改也可以直接部署到服务器上。

### 消息效果与限制

目前发送的应用消息类型为图文消息(mpnews)，消息内容支持html标签，不超过666K个字节，效果如下

<img src="https://raw.githubusercontent.com/zyh94946/work-wx-msg-push/main/demo/demo.png" />

不用安装企业微信App，直接通过微信App关注微信插件即可实现在微信App中接收应用消息，还可以选择消息免打扰。

#### 发送应用消息限制

每企业不可超过帐号上限数*30人次/天（注：若调用api一次发给1000人，算1000人次；若企业帐号上限是500人，则每天可发送15000人次的消息）

每应用对同一个成员不可超过30条/分，超过部分会被丢弃不下发

默认已启用重复消息推送检查5分钟内同样内容的消息，不会重复收到，可修改 `EnableDuplicateCheck` `DuplicateCheckInterval` 调整是否开启与时间间隔。

通过Api网关触发，完整请求与响应请[参考](https://cloud.tencent.com/document/product/583/12513)

### 使用方法
#### 注册企业微信
#### 创建应用

登录企业微信web管理 `https://work.weixin.qq.com/`

进入应用管理，创建应用，完成后复制下AgentId，Secret

进入管理工具，素材库，图片，添加图片，上传成功后查看图片地址把 media_id 的值复制下

进入我的企业，把企业ID复制下，进入微信插件，用微信APP扫码关注即可

#### 注册腾讯云账号
#### 开通云函数服务
#### 下载本仓库代码至本地
#### 新建云函数

如果想要绑定域名的话可以选择香港地区，免备案。

选择自定义创建，运行环境Go1，提交方法选择本地上传文件夹选择刚才下载的代码目录

高级配置中增加环境变量

- CORP_ID 企业微信 企业id
- CORP_SECRET 企业微信 应用Secret
- AGENT_ID 企业微信 应用AgentId
- MEDIA_ID 企业微信 图片素材的media_id

触发器配置选择自定义配置，触发方式选择API网关触发

部署云函数

#### 测试云函数

函数服务下点击函数名进入函数管理，函数代码下测试事件选择 Api Gateway 事件模版

path 改成 `/你的函数名称/CORP_SECRET`
queryString 中设置 key 为 title 和 content 的参数
body 中设置 `{"title":"test","content":"test"}`
queryString 和 body 设置一种即可

点击测试调试

#### 设置API网关

进入API网关服务列表，选择配置管理，然后管理API点编辑

增加参数配置，参数名 `SECRET`，参数位置 path，类型 string

路径修改为 `/你的云函数名称/{SECRET}`

然后 立即完成 发布服务

#### 使用

在 基础配置 中复制公网访问地址，想要绑定域名可以在自定义域名中绑定，CNAME指到API网关的二级域名，自定义路径映射 `/` 路径 到 `发布` 环境。

GET方式

`https://service-xx-xx.xx.apigw.tencentcs.com/你的云函数名称/CORP_SECRET?title=消息标题&content=消息内容`

POST方式

```bash
$ curl --location --request POST 'https://service-xx-xx.xx.apigw.tencentcs.com/你的云函数名称/CORP_SECRET' \
--header 'Content-Type: application/json;charset=utf-8' \
--data-raw '{"title":"消息标题","content":"消息内容"}'
```

发送成功状态码返回200，`"Content-Type":"application/json"` body `{"errorCode":0,"errorMessage":""}` 。

### 其它

TODO：

- 使用方法图片补充

