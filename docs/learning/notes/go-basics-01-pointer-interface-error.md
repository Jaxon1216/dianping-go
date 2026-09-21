# Go 基础 01：指针 / 接口 / nil / error / panic

> 来源：第 1 轮前置知识对话（对照 `internal/service/user.go`、`internal/handler/user.go`、`api/v1/`）。
> 读法：每条用 **S 场景 → T 问题 → A 做法 → R 记忆点** 四段，配项目真实代码锚点。
> 面向读者：前端 / C++ 背景，用 TS、Express、C++ 做类比，并标出类比失效的边界。

---

## 核心误区先纠正（贯穿全篇）

> **`nil` = 「空 / 没有」，不是「报错」。**

想通这一句，下面 err、指针、interface 的问题会一起塌方：

```text
nil        = 空 / 什么都没有       ← 它是"没出错"的那一侧
err == nil = 错误盒子是空的 → 正常 ✅
err != nil = 错误盒子里装了错误 → 出错 ❌
```

---

## 1. error 与 `if err != nil`

### S 场景
`internal/service/user.go` 每个方法都返回 `error`，handler 里到处是 `if err != nil`。

### T 问题
`nil` 是空、`err` 是错误，为什么 `err != nil` 才代表「有错」？

### A 做法
`err` 是一个**变量（盒子）**，值有两种：

```text
err = nil            → 盒子空 → 没错误 → 正常
err = &someError{}   → 盒子里装了错误对象 → 出错
```

所以：

```go
// internal/service/user.go  SignCount
} else if err != nil {   // 盒子非空 = 有错
	return 0, err        // 往上抛
}
```

### R 记忆点
- **`if err != nil` ≈ JS 的 `if (err)`**（东西存在 = 有错）。
- 之前拿 `if (!user)` 类比是**错的方向**，别用。
- Go 无 try/catch，「拦截」= 每层手动 `if err != nil`。handler 是翻译官：把 `error` 翻成统一 JSON（见第 4 条）。

---

## 2. 多返回值：按位置匹配，不按名字

### S 场景
`func Me() (*v1.SimpleUser, error)`，`return nil, v1.ErrCanNotGetUser`。

### T 问题
「返回俩？」按顺序还是按名字匹配？

### A 做法
Go 原生支持多返回值，**纯按位置（顺序）匹配，与变量名无关**：

```go
func Me() (*v1.SimpleUser, error) { //  位置1: 指针   位置2: error
	return user, nil                //  user→1       nil→2
}
data, err := userService.Me(...)    //  data 接1     err 接2（叫 foo,bar 也行）
```

### R 记忆点
- 类比 JS 数组解构 `const [a, b] = f()`：位置匹配。
- 铁律约定：**最后一个返回值几乎总是 `error`**，正常填 `nil`，出错填错误。
- 业务错误集中定义在 `api/v1/errors.go`（如 `ErrCanNotGetUser = newError(1103, "获取用户信息失败")`）。

---

## 3. 指针 vs 值，以及 nil 是「空指针」不是「野指针」

### S 场景
`Me` 返回 `*v1.SimpleUser`（指针），并用 `if user == nil` 判空。

### T 问题
为什么用指针不用值？`nil` 是不是空指针 / 野指针？

### A 做法
内存模型：

```text
值    user1 := v1.SimpleUser{...}   变量自己就装整份数据（传参拷一整份）
指针  user2 := &v1.SimpleUser{...}  变量只装地址(箭头)，指向真实数据（传参只拷地址）
```

**关键：只有指针能是 `nil`，值不能。**

```text
值 v1.SimpleUser：永远存在，哪怕零值 {ID:nil,...}，无法表达"没有"
指针 *v1.SimpleUser：要么指向真实数据，要么 = nil（谁都不指）→ 能表达"没有用户"
```

`Me` 的语义是「可能没登录/取不到用户」，这个「没有」只能用指针 `nil` 表达：

```go
user := user_holder.GetUser(ctx) // *v1.SimpleUser，可能 nil
if user == nil {                 // 只有指针能这样判"有没有"
	return nil, v1.ErrCanNotGetUser
}
```

`&`（取地址，值→指针）和 `*`（解引用，指针→值）跟 C++ 一致。见 `createUserWithPhone`：先建值 `user`，再 `return &user`。

### R 记忆点
- **`*T` ≈ TS 的 `T | null`**；看到 `*Something` 先想「可能没有，用前判 `== nil`」。

| | C++ | Go |
|---|---|---|
| 空指针 | `nullptr` | `nil`（合法、可判、安全） |
| 野指针 | `delete` 后仍用 | **Go 没有**（有 GC，不手动释放） |

- `nil` 本身安全；**对 nil 指针取字段才 panic**（= JS `Cannot read property of null`）。判空是保护，漏判就崩。

---

## 4. nil / error 会直接抛给前端吗？—— 不会，handler 拦截并翻译

### S 场景
`service.Me` 返回 `(nil, ErrCanNotGetUser)`。

### T 问题
这个 nil 会直接返回给前端吗？能拦截吗？

### A 做法
不会。service 的返回值只到 handler 为止，handler 拦截并翻译成 JSON：

```go
// internal/handler/user.go  Me
user, err := h.userService.Me(ctx.Request.Context())
if err != nil {                                     // ★ 拦截点
	v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
	return
}
v1.HandleSuccess(ctx, user)
```

```text
前端 GET /user/me
  → handler.Me  调 service.Me() 拿回 (nil, err)
  → if err != nil 拦截 ★
  → v1.HandleError 翻译成统一 JSON（api/v1/v1.go）
  → 前端只收到 { "success": false, "errorMsg": "获取用户信息失败" }
```

### R 记忆点
- service 返回的 `nil`/`error` 是**内部信号**；**前端永远看不到 Go 的 nil / error 对象**，只看到统一 JSON。
- 分层：service 只管返回 `(数据, error)`；handler 是「翻译官」把 error 变 HTTP 码 + JSON。
- Express 类比：service `return nil, err` ≈ `throw`/`return null`；handler `if err != nil` ≈ `if(!user) return res.status(500)...`。差别：Node throw 自动冒泡，Go 手动逐层传。

---

## 5. interface vs struct，以及为什么参数不写 `*`

### S 场景
接口方法 `Login(ctx context.Context, req *v1.LoginReq)`：`ctx` 不加 `*`，`req` 加 `*`。

### T 问题
`ctx` 为什么不加 `*`？是值传递吗？struct 和 interface 到底差在哪？

### A 做法
本质区别：

```text
struct    = 一坨实打实的数据（有字段，占内存）      类比 TS 对象 / C++ struct
interface = 一张"能做什么"的能力清单（只有方法名）  类比 TS interface / C++ 抽象基类
```

要不要 `*`，看类型是哪种：

```text
struct（v1.LoginReq）：一坨数据 → 传值拷一整坨 → 加 * 传指针省拷贝/可改原件
interface（context.Context / error / XxxService）：
   内部天生是"胖指针"（类型信息 + 指向数据的指针），本身已是引用 → 不加 *
基础类型（int/bool/uint64）：太小，拷贝无所谓 → 不加 *
```

interface 变量内部结构：

```text
┌─────────────────────────────┐
│ 类型信息 │ 指向真实数据的指针 ●┼──▶ 真实数据
└─────────────────────────────┘
```

**是值传递吗？** 是。Go 只有值传递（pass-by-value）。但传 interface/指针时，拷贝的「壳」里含指针，仍指向同一份底层数据，效果像共享。

验证（`SignCount(ctx context.Context) (int, error)`）：`context.Context` 是 interface 不加 `*`，`int` 基础类型不加 `*`，`error` 是 interface 不加 `*`。

### R 记忆点
- **一句话规则：struct 用 `*`；interface（`context.Context`/`error`/`XxxService`）和基础类型不用 `*`。**
- 大白话：struct 是「一箱货」给地址(`*`)让人原地取；interface 是「快递单」（单上已写货在哪），直接给单子，不需要「单子的地址」。
- Go 只有值传递；interface/指针拷的是「引用壳」，所以效果像共享。

---

## 6. panic 是什么，会发生什么现象

### S 场景
`service.go` 的 `NewDB` 里连不上库就 `panic(err)`。

### T 问题
panic 是什么？触发后会怎样？和 error 什么关系？

### A 做法
panic = **程序崩溃信号**：当前 goroutine 立刻停止，逐层往上「炸」，直到被 `recover` 接住或炸到顶杀死进程。

```go
// internal/service/service.go  NewDB —— 启动期致命错误，快速失败
if err != nil {
	panic(err)
}
```

两种现象：

```text
A. 没人 recover → 打印 "panic: xxx" + goroutine 栈回溯 → 进程退出(exit status 2)
B. 被 recover 接住 → 不崩。项目里 http.go 用 gin.Default()，自动挂 Recovery 中间件：
   handler 内 panic → Recovery 接住 → 记录栈日志 + 该请求返回 500 → 服务继续跑
```

会**自动**触发 panic 的运行时错误（会踩的坑）：

```go
var p *T = nil; _ = p.Field   // nil 指针解引用
arr[5]                        // 切片越界
m["a"] = 1                    // 往 nil map 写（未 make）
```

### R 记忆点
- **能用 error 就别用 panic**：error = 预期内失败（手动 `return`）；panic = 致命、没法继续（往上炸）。
- 项目里：手动 panic 只用于**启动期致命错误**（`NewDB`）；请求期靠 `gin.Default()` 的 **Recovery** 兜底（单请求 500，服务不崩）。
- 对照：没 recover 的 panic ≈ Node 未捕获异常冒泡到顶后 `process exit`；C++ 空指针/越界是未定义行为，Go 里是明确 panic（更安全但仍崩）。

---

## 一图串联全篇

```text
接口方法签名：ctx(interface,不加*)  req(struct,加*)  →返回 (数据指针 或 nil, error)
        │
        ▼  service 只管返回 (数据, error)
service.Me → (nil, ErrCanNotGetUser)
        │  handler 用 if err != nil 拦截（普通错误走这条）
        ▼
handler 翻译成统一 JSON → 前端只见 { success, errorMsg, data }

（另一条线）致命错误 → panic → 没recover则进程死 / gin Recovery 接住则该请求500
```

共同根：**`nil` = 空 = 没有**（不是报错）。理清它，`err != nil`（有错）、`p == nil`（空指针）、interface 天生带引用（不写 `*`）全部顺了。
