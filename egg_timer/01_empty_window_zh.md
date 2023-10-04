---
layout: default
title: 第一章 - 空画布
nav_order: 2
parent: 煮蛋计时器
has_children: false
---

# 第一章 - 创建一个空白的窗口

## 目标

本章节的目标是创建一个空白的画布，以便我们之后在其上绘制内容。

![An empty window](01_empty_window.gif)

## 大纲

下面这段代码执行了三个主要操作：

- 导入 Gio
- 创建并调用一个 goroutine，其中包括：
  - 创建一个新窗口，命名为 `w`
  - 启动一个死循环，等待窗口中的事件（在这个示例中永远不会发生任何事件）

就这么多！来看看码：

## 代码

```go
package main

import (
  "gioui.org/app"
)

func main() {
  go func() {
    // 创建新窗口
    w := app.NewWindow()

    // 在窗口中监听事件
    for range w.Events() {
    }
  }()
  app.Main()
}
```

## 注解

这段代码看起来足够简单对吧？尽管如此，还是让我们花些时间来看看其中发生了什么。

1. 我们导入了 `gioui.org/app`（这啥？）

   查看[文档](https://pkg.go.dev/gioui.org/app)后，我们发现：

   > _app 包为带有 GUI 的操作系统的功能提供了一个独立于平台的接口。_

   好消息，Gio 为我们处理了所有与平台相关的事情。我经常在 Windows 和 MacOS 上编码，Gio 可以在这两个平台上正常工作，[GioUI.org](https://gioui.org/#installation) 还列出了更多平台，包括 iOS 和 Android。

   这比你也许能意识到的要深，即使你写的程序今天只支持了单一平台，其实你掌握的技能已经是多平台的了。
   _"我们应该移植到 Mac。"_ &nbsp;搞定了！_"火热的创业公司正在寻找桌面应用开发专家。"_ &nbsp;没问题。_"这里有谁懂 tvOS 吗？"_ &nbsp;你懂！
   _"机长凉了，谁来着陆一下这飞机？"_ &nbsp;好吧，最后一个可能不太现实，但观点仍然成立。Gio 的跨平台简直令人惊叹。

2. 在 goroutine 中的**事件循环**

   - 事件循环指的是 `for range w.Events()` 循环，它监听窗口中的事件，目前我们只是让它监听事件，而不对接收到的事件进行任何操作，之后我们再开始对事件做出响应。

     从 [app.main](https://pkg.go.dev/gioui.org/app#hdr-Main) 中，我们了解到：

     > _因为在某些平台上，Main 也会阻塞，所以窗口的事件循环必须在 goroutine 中运行。_

   - 没有名称的 goroutine（也就是 _匿名函数_）被创建并运行事件循环，由于它在 goroutine 中，它将与程序的其余部分并发运行。

   ```go
   go func {
     // ...
   }()
   ```

   Jeremy Bytes 的[这篇关于匿名函数的文章](https://jeremybytes.blogspot.com/2021/02/go-golang-anonymous-functions-inlining.html)写得很好，其中的内容在许多情况下都很有用，不仅仅在 Gio 中。

3. 通过调用 `app.Main()` 来启动它：
   从 [app.Main 文档](https://pkg.go.dev/gioui.org/app#hdr-Main) 中，我们了解到：
   > _Main 方法必须从程序的主函数中调用，以将主线程的控制权交给需要其的操作系统。_

---

[下一章](02_title_and_size_zh.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }
