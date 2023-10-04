---
layout: default
title: 第五章 - 重构
nav_order: 2
parent: 煮蛋计时器
has_children: false
---

# 第五章 - 重构代码

## 目标

本章节的目标是更好地组织代码。

## 大纲

到目前为止，我们都是通过逐步添加功能来构建程序，这让我们可以从一个空白画布开始，通过改变最少的代码行数来迭代，同时还能取得有意义的进展。

然而如果再继续下去，代码看起来有点难以管理。将所有代码都放在一个大的 `main()` 函数内，可能会让它变得难以理解，也更难以继续构建。因此，我们将对程序进行一些重构，简单地将其分成更小的部分。

> _重构就是以安全而迅速的方式转变代码，对于使其保持成本低廉且易于修改以满足未来需求至关重要。_
> —— [马丁·福勒（Martin Fowler）](https://martinfowler.com/books/refactoring.html)

也就是说，我们现在不会添加新功能，而是为以后更好的功能铺平道路。

## 代码

### 首先是： `main()` 太长了

`main` 函数太长了，功能太多了，最好的做法是 `main()` 只负责启动和控制程序，除此之外委托给其他函数。这是新的 `main` 函数：

```go
func main() {
  go func() {
    // create new window
    w := app.NewWindow(
      app.Title("Egg timer"),
      app.Size(unit.Dp(400), unit.Dp(600)),
    )
    if err := draw(w); err != nil {
      log.Fatal(err)
    }
    os.Exit(0)
  }()
  app.Main()
}
```

现在，在 `main()` 函数内，我们创建了一个窗口 `w`，然后立即将其交给了一个专门处理它的函数 `draw()`。

通过将 `draw()` 的执行结果存储在 `err` 中，我们可以检查执行是否顺利进行，并且可以有序地处理任何错误。

对于错误处理，我们使用了 [os.Exit()](https://pkg.go.dev/os?utm_source=gopls#Exit) 和它的近亲 [log.Fatal(err)](https://pkg.go.dev/log?utm_source=gopls#Fatal)，它们都来自标准库。

如前所述，约定 0 退出码表示正常结束，也就是 `err == nil` 时使用 `os.Exit(0)` 退出，否则我们调用 `log.Fatal(err)`，该函数会打印错误消息并使用 `os.Exit(1)` 退出。

### 之后是： 为 约束(Constraints) 和 尺寸(Dimensions) 创建一个简单的别名

我们之前详细讨论了 **约束** 和 **尺寸**，由于我们经常使用它们，因此定义两个别名 `C` 和 `D` 是很方便的。注意：约束就是上下文的一部分。

```go
type C = layout.Context
type D = layout.Dimensions
```

### 最后是： `draw()` 函数

用一个简化的 `draw()` 展示下它的结构：

```go
func draw(w *app.Window) error {
  // ...

  // listen for events in the window.
  for e := range w.Events() {

    // detect what type of event
    switch e := e.(type) {

    // this is sent when the application should re-render.
    case system.FrameEvent:
        // ...

        // this is sent when the application is closed.
    case system.DestroyEvent:
      return e.Err
    }
  }
  return nil
}
```

与之前一样，我们 `for range w.Events()` 并判断它们的类型。

- 对于 `system.FrameEvent`，我们像以前一样处理
- 我们添加了一个 `system.DestroyEvent` 的新 case，如果窗口正常关闭时返回 _nil_，否则如果窗口因为其他原因关闭时返回 _Err_

## 注解

重构与个人的代码品味有关，这是我的处理方式，如果你有不同的需求，请根据你的程序做正确的事情，主要的目的只是保证代码足够灵活，以支持持续改进和未来需求。祝你好运。

---

[下一章](06_button_low_margin_zh.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }
