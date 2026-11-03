# 发短信执行链：类型、实例、接口值与接收者

## 场景

把 interface、struct、大小写和指针同时讲，容易混淆定义与运行。用户明确需要沿执行顺序观察“每一步传什么、改什么”，而之前没有字段的示例未使用 s，无法展示接收者的作用。

## 已确认的源码与语义

- `type userService struct {...}` 定义类型，`userService{...}` 创建具体值，`&userService{...}` 得到实例指针。
- 普通接口声明列出方法要求；接口变量携带一个具体值。接口实现检查发生在编译期间，不是运行时创建实例的动作。
- `func (s *userService) SendCode(...)` 声明指针类型的方法；调用时 s 携带本次服务实例的指针。Go 按值传参，指针值的复制不会复制整份实例。
- 方法关联类型，接收者确定当前实例。真实 `SendCode` 使用 `s.conf`、`s.rdb`、`s.logger`，请求参数则由 `req.Phone` 提供。
- 首字母大小写控制包外可见性，和 interface / struct 分类、实例创建顺序无关。

## 项目链路

启动：`NewService` → `NewUserService` → `NewUserHandler`，传递公共依赖指针和用户服务接口值。

请求：`POST /user/code` → Handler 填写 req → 通过接口传入 request context 与 `&req` → Service 校验手机号、读取配置、写 Redis、记录日志 → error / nil → HTTP 响应。

当前 Redis key 为 `login:code:手机号`，TTL 为 2 分钟；第三方短信发送尚未接入。

## 教学验证入口

- [逐步执行页](../lessons/0002-sms-execution.html)：显示当前代码、传入/传出值、实例字段变化。
- [独立 Go 示例](../examples/sms-flow/main.go)：模拟短信，打印 p、u 携带的指针与 s 的地址；两次调用展示同一实例计数从 0 到 2。
- [实际方法](../../../internal/service/user.go)、[Handler](../../../internal/handler/user.go)、[Wire 装配](../../../cmd/server/wire/wire_gen.go)。

演示字段 senderName / sentCount 不在真实 userService 中。页面是源码语义演示，不是浏览器执行 Go；学习材料覆盖不等于用户已经独立掌握。

## 骨架语法：返回类型、创建表达式、方法接收者

`func NewUserService(name string) UserService` 按「func → 函数名 → 参数列表 → 返回类型」读。末尾的 UserService 是返回值类型，要求返回的具体值满足该接口；它不表示创建函数自己拥有 SendCode 方法。对应 TypeScript 的签名是 `function newUserService(name: string): UserService`。

`userService{senderName: name}` 按「具体类型名 → 字段名: 字段值」读，创建结构体值；左侧 senderName 是字段，右侧 name 是函数参数。前面加 `&` 得到指向该实例的指针。struct 是定义类型时的关键字，创建表达式用具体类型名。

`func (s *userService) SendCode(phone string) error` 按「func → 接收者 → 方法名 → 参数列表 → 返回类型」读。该方法关联 *userService 类型；调用时 s 携带本次实例指针，phone 携带本次手机号，error 是调用结果的类型。

接口只声明方法要求，具体方法体另行定义。`return p` 能返回为 UserService，是因为 *userService 已满足该接口的方法要求；接口声明本身不会给结构体生成方法。方法实现可以写在创建函数前或后，声明顺序不表示运行顺序，运行从 main 开始。
