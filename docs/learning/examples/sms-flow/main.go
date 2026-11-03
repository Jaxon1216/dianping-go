// 本例只模拟短信发送：打印验证码，不请求短信平台，不访问 Redis。
// 按日志观察 p、u 携带的指针和方法接收者 s 是否指向同一个实例。
// 阅读顺序：接口规定方法 → 定义实例字段 → 创建并返回实例 → 实现方法。
// 实际运行从 main 开始；下面的类型和方法定义不是按书写顺序逐个执行。

// package main：本文件属于 main 包，这是可执行程序的入口包。
package main

// 引入 fmt，下面用它打印文字、类型和地址，方便观察执行过程。
import "fmt"

// ① interface：规定“符合这个接口的东西，要能调用哪些方法”。
// 这里要求有 SendCode，传入 phone 字符串，返回 error。
// 只写方法签名；真正怎么发短信，写在后面的具体方法体里。
//
// 大小写只管其他包能否引用：UserService 大写可导出，userService 小写不导出。
// 是接口还是结构体，看 interface / struct 关键字，不看大小写。
type UserService interface {
	SendCode(phone string) error
}

// ② struct：定义具体实例有哪些字段，可以先按“对象的数据结构”理解。
// 这一段定义的是类型；后面的 userService{...} 才创建具体实例。
type userService struct {
	senderName string // 对象保存的发送方名字，例如“点评短信”。
	sentCount  int    // 对象保存的发送次数；创建时未填写，初始为 0。
}

// ③ 普通创建函数：
// func 函数名(参数列表) 返回类型。
//
// name string：调用者传进来的名字。
// 末尾 UserService：返回值要满足 UserService 接口要求。
// 方法要求约束的是“返回的那个值”，不是创建函数自己。
func NewUserService(name string) UserService {
	fmt.Printf("进入创建函数：name=%q，类型=string\n", name)

	// userService{senderName: name}：创建一个具体实例。
	// senderName 是实例字段，name 是传进来的参数值。
	// 前面的 &：得到这个实例的指针；p 的类型因此是 *userService。
	p := &userService{senderName: name}
	fmt.Printf("创建实例：p=%p，类型=%T，sentCount=%d\n", p, p, p.sentCount)

	// 把实例指针 p 按 UserService 接口类型返回，交给 main 里的 u。
	// 为什么可以？因为 *userService 实现了下面的 SendCode 方法。
	// 编译器会检查方法是否满足要求；接口不会自动给实例生成方法。
	return p
}

// ④ 给 *userService 类型实现接口要求的 SendCode 方法。
// func (接收者名 接收者类型) 方法名(参数列表) 返回类型。
//
// s：本次调用的接收者变量，保存“这次操作哪个实例”的指针，类似 this。
// phone：本次调用传来的手机号，表示“这次发给谁”。
// error：本次方法调用的返回类型。
func (s *userService) SendCode(phone string) error {
	fmt.Printf("进入方法：s=%p，phone=%q，sentCount=%d\n", s, phone, s.sentCount)

	// 通过 s 找到那个实例，把实例里的计数加 1，类似 this.sentCount++。
	s.sentCount++

	// s.senderName、s.sentCount：从实例中读取的数据。
	// phone：这次调用传进来的数据。把它们组合起来，模拟发送短信。
	fmt.Printf("[%s] 模拟向 %s 发验证码 123456；累计 %d 次\n",
		s.senderName, phone, s.sentCount)
	fmt.Println("方法返回 nil：本次没有错误")

	// 返回 nil，表示 error 为空，即本次没有错误。
	return nil
}

// ⑤ 实际运行从这里开始。
func main() {
	fmt.Println("main 调用 NewUserService，传入字符串“点评短信”")

	// 执行创建函数，把“点评短信”传给 name。
	// 创建函数返回 UserService 接口值，所以 u 的类型是 UserService。
	// u 携带的具体值，是刚才创建的那个 *userService 实例指针。
	u := NewUserService("点评短信")
	fmt.Printf("main 拿到 u：静态类型=UserService，动态类型=%T，携带指针=%p\n", u, u)

	fmt.Println("\nmain 第一次调用 u.SendCode，传入第一个手机号")

	// 通过接口调用具体实现：
	// 实例指针交给接收者 s；“13800138000”交给参数 phone。
	// 方法返回的 error 交给 err；非 nil 时进入错误分支。
	// 本例返回 nil，因此不进入错误分支，实例计数从 0 变成 1。
	if err := u.SendCode("13800138000"); err != nil {
		fmt.Println("发送失败：", err)
		return
	}

	fmt.Println("\nmain 第二次调用同一个 u，传入另一个手机号")

	// 这里没有再创建服务，仍然通过同一个 u 调用。
	// s 仍指向同一实例，只换了 phone；所以计数从 1 变成 2。
	if err := u.SendCode("13900139000"); err != nil {
		fmt.Println("发送失败：", err)
		return
	}
}
