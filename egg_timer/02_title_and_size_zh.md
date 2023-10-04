---
layout: default
title: 第二章 - 标题与大小
nav_order: 2
parent: 煮蛋计时器
has_children: false
---

# 第二章 - 窗口标题和大小

## 目标

本章节的目标是设置自定义窗口的标题和大小。

![Window with custom title and set size](02_title_and_size.png)

## 大纲

这段代码与[第一章](01_empty_window_zh.md)中的非常相似，我们添加了以下内容：

- 两个额外的包被导入
- 在调用 `app.NewWindow()` 时添加了两个参数

## 代码

```go
package main

import (
  "os"

  "gioui.org/app"
  "gioui.org/unit"
)

func main() {
  go func() {
    // 创建新窗口
    w := app.NewWindow(
      app.Title("Egg timer"),
      app.Size(unit.Dp(400), unit.Dp(600)),
    )

    // 在窗口中监听事件
    for range w.Events() {
    }
    os.Exit(0)
  }()
  app.Main()
}
```

## 注解

第一章中的创建窗口是整个程序基础中的基础，我们希望在这里做一些改进。其中一个改进是确保程序能干净的退出，因此我们导入了 `os` 包并在事件循环后添加了一行 [os.Exit()](https://pkg.go.dev/os?utm_source=gopls#Exit)。一般约定 0 表示正常退出，之后的其他逻辑可以添加其他值。

[gioui.org/unit](https://pkg.go.dev/gioui.org/unit) 实现了独立于设备的单位和值，文档也描述了一些替代选项：

| 类型 | 描述 |
| :--: | :-- |
|  dp  | 设备独立像素 - 独立于底层设备 |
|  sp  | 比例独立像素 - 用于文本大小 |
|  px  | 像素 - 用于实际设备的精度 |

一般来说，`dp` 是最常用的，我们倾向于在尽可能的情况下保证统一（_独立于设备_）的显示大小。因此，在 `app.NewWindow()` 中定义窗口大小时我们使用了 `dp`。

`app.NewWindow(options ...Option)` 的 [options](https://pkg.go.dev/gioui.org/app#Option) 都相容易理解，但请注意以下几点：

- 使用 `app.Size(x, y)` 设置窗口大小。
- 窗口默认可以自由调整大小。如果要限制大小，可以添加：
  - MaxSize
  - MinSize
  - 或同时使用两者以锁定窗口大小
- 需要的话还可以使用全屏
- 如果你正在构建用于安卓的版本，还可以在这里设置 `系统状态栏(Status)` 和 `导航栏(Navigation)` 的颜色。

---

[下一章](03_button_zh.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[在 GitHub 上查看](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }
