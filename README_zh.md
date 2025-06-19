# newc

> **new**-\* **c**onstructor generator

Doc: [English](README.md) | **中文**

----------

一个为 Golang 结构体生成构造器函数(`NewXXX`)代码的自动工具。

## 安装

```bash
go install github.com/GerryLon/newc@latest
```

## 使用方法

在需要生成构造器的结构体上添加一行 `go:generate` 注释。

```go
//go:generate newc
type UserService struct {
	baseService
	userRepository *repositories.UserRepository
	proRepository  *repositories.ProRepository
}
```

然后只需要执行 `go generate ./...` 就能生成下面这样的构造器代码：

```go
// constructor_gen.go

// NewUserService Create a new UserService
func NewUserService(baseService baseService, userRepository *repositories.UserRepository, proRepository *repositories.ProRepository) *UserService {
	return &UserService{
		baseService:    baseService,
		userRepository: userRepository,
		proRepository:  proRepository,
	}
}
```

这里可以[查看更多例子](https://github.com/GerryLon/newc/tree/master/test)

## 使用方式2（无需手动安装）

**推荐在团队合作中采用**

无需手动安装，只需要给结构体添加下面这行注释就行。Go 会在缺失时自动下载这个工具。

```go
//go:generate go run github.com/GerryLon/newc@v0.8.3
```

比如这样：

```go
//go:generate go run github.com/GerryLon/newc@v0.8.3
type UserService struct {
	baseService
	userRepository *repositories.UserRepository
	proRepository  *repositories.ProRepository
}
```

这个方式非常有用，尤其在团队开发中。**就算其他同事没有安装这个工具，这么做也能正常运行，不会影响到其他人的工作**。

## 如何返回结构体的值，而不是引用？

使用 `--value` 参数

```go
//go:generate newc --value
type Config struct {
	debug  bool
}
```

生成代码：

```go
// constructor_gen.go

// NewConfig Create a new Config
func NewConfig(debug bool) Config {
	return Config{
		debug:  debug,
	}
}
```

## 如何在构造器中做初始化操作?

1. 加上 `--init` 参数
2. 为结构体实现一个 `init` 方法

```go
//go:generate newc --init
type Controller struct {
	logger *zap.Logger
	debug  bool
}

func (c *Controller) init() {
	c.logger = c.logger.With(zap.String("tag", "controller-debugger"))
	c.debug = true
}
```

生成代码：

```go
// constructor_gen.go

// NewController Create a new Controller
func NewController(logger *zap.Logger, debug bool) *Controller {
	s := &Controller{
		logger: logger,
		debug:  debug,
	}
	s.init()
	return s
}
```

## 如何在生成构造器时忽略掉一些字段？

给需要忽略的字段添加 `newc:"-"` 标签

```go
type Forbidden struct {
	Msg    string
	Status int    `newc:"-"`
}
```

生成代码：

```go
// NewForbidden Create a new Forbidden
func NewForbidden(msg string) *Forbidden {
	return &Forbidden{
		Msg: msg,
	}
}
```

## 如果你觉得这条注释太长……

一些建议：

1. (推荐）把它加进你的编辑器/IDE的快捷代码片段（code snippest）里
2. ......

## 功能特性与设计理念

**1. 它能让你的代码更容易编写和维护**.

不管是编写还是更新构造器代码，都是一个费力且容易出错的事情，尤其当代码量很大的时候。这些繁琐易错的工作应该交给自动程序来完成，比如这个工具。

同时，这个工具还能完美兼容像[**wire**](https://github.com/google/wire)这种依赖注入工具。如果你的项目中也使用了 **wire**，那你可能非常需要这个工具。**wire** 在 **newc** 的“加持”下会变得更加好用。

**2. 你不需要担心自动生成的代码**.

这个工具在生成代码时会非常小心，会帮你考虑所有代码细节，包括引用依赖、变量命名，甚至还有代码风格。

**3. 它非常适合团队协同工作**.

就算其他同事没有安装这个工具，这么做也不会影响到他们的工作。因为 Go 会在必要时自动安装这个工具。

```go
//go:generate go run github.com/GerryLon/newc@v0.8.3
```

## 赞赏

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://buymeacoffee.com/benn)

![](./doc/donate.png)

## License

MIT
