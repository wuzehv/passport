# sso实践

## 相关文档
* [gin](https://github.com/gin-gonic/gin)
* [gorm](https://gorm.io/zh_CN/)
* [redisgo](https://github.com/gomodule/redigo)
* [gin-swagger](https://github.com/swaggo/gin-swagger)
* [swag document](https://github.com/swaggo/swag/blob/master/README_zh-CN.md)
* [validator](https://pkg.go.dev/github.com/go-playground/validator/v10)
* [go-ini](https://ini.unknwon.io/docs/intro/getting_started)
* [jwt-go](https://pkg.go.dev/github.com/dgrijalva/jwt-go)

## 开发规范
* 目录采用单数形式
* 表名复数形式

## 使用
1. 配置hosts
```etc
127.0.0.1 client.one.com
127.0.0.1 client.two.com
127.0.0.1 sso.com
```

2. 调整`conf/app.ini`，初始化mysql数据库和表
* 创建数据库

`CREATE DATABASE IF NOT EXISTS passport DEFAULT CHARACTER SET utf8`

* 初始化表结构
```shell
./run.sh init # 或 make init
```

3. 编译主程序，默认监听8099端口

4. 启动两个客户端程序
```shell
# 默认监听8081
./run.sh client # 或 make client1

# 监听8082
./run.sh client -addr=127.0.0.1:8082
```

5. 浏览器访问
```
client.one.com:8081/index
client.two.com:8082/index

# sso中心
sso.com:8099
```

**默认登录信息：**
```
email: admin@gmail.com
password: 123456
```

## TODO
- [ ] SSO
  - [X] 登录逻辑
  - [X] 退出逻辑
  - [X] 签名机制
  - [X] 引入redis做svc接口
  - [X] 登录错误次数限制
  - [X] 考虑jwt使用
  - [X] swagger(http://sso.com:8099/swagger/index.html#/)
  - [ ] IP白名单
  - [ ] sso接口交互不用cookie

- [ ] 后台功能列表
  - [X] 接口管理
  - [X] 用户管理
  - [X] 客户端管理
  - [X] validator详细使用
  - [X] 菜单管理
  - [ ] 权限管理
 
- [ ] 日志
  - [X] 业务日志公共方法
  - [ ] 日志trace id
  - [X] 日志切割
  
- [ ] 部署
  - [X] systemd
  - [ ] nginx upstream(option)
  
- [ ] 对接
  - [ ] php业务对接

## 对接流程
### svc接口
* 请求方式
  * POST
  * `Content-Type: application/x-www-form-urlencoded`


* 请求参数

```json
{
  "domain": "client.one.com",
  "timestamp": "1623680856",
  "token": "6bc4931890225677da85a1cf05ce0fc0",
  "sign": "6BC4931890225677DA85A1CF05CE0FC0"
}
```

* sign签名算法
```php
$s = [
    'domain' => 'client.one.com',
    'token' => '6bc4931890225677da85a1cf05ce0fc0',
    'timestamp' => '1623680856',
    'sign' => '6BC4931890225677DA85A1CF05CE0FC0',
];

// 按key排序
ksort($s);

// 拼接内容
$str = implode("", $s);

// 末尾拼接密钥
$str .= "123456";

echo strtoupper(md5($str)), PHP_EOL;
```


## FAQ
* gin-swagger的结构体名称使用json包修改名称，如果使用gin包的参数绑定，并且是form请求方式，需要与`form`别名一致