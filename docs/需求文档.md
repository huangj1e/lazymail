好，这里我直接给你一份**可以丢给 AI 开发执行的完整需求文档（Spec-Driven 风格）**，基于
👉 bubbletea + Go
目标就是你说的：**lazygit 风格 TUI 邮件客户端（支持鼠标）**

---

# 📄 项目名称

**LazyMail（暂定）**

👉 定位：
**终端中的邮件工作台（像 lazygit 一样好用）**

---

# 🎯 一、项目目标

构建一个：

* 可鼠标操作的 TUI 邮件客户端
* 面板式 UI（类似 lazygit）
* 支持 IMAP / SMTP
* 支持多账号
* 可扩展 AI 自动化

---

# 🧠 二、核心设计理念

### 1. 状态驱动 UI（重要）

采用 Bubble Tea 的：

```text
Model -> Update -> View
```

所有 UI 更新必须通过状态变化驱动

---

### 2. 面板化 UI（lazygit 风格）

```text
┌──────────────┬────────────────────────────┬────────────────────────────┐
│ Sidebar      │ Mail List                  │ Mail Viewer                │
│--------------│----------------------------│----------------------------│
│ Inbox        │ > Subject 1                │ From: xxx                  │
│ Sent         │   Subject 2                │ To: xxx                    │
│ Drafts       │   Subject 3                │                            │
│ Archive      │                            │ 邮件正文                    │
└──────────────┴────────────────────────────┴────────────────────────────┘
```

---

### 3. 鼠标 + 键盘双支持

必须支持：

* 鼠标点击选择
* 滚轮滚动
* 键盘导航（vim 风格）

---

# ⚙️ 三、技术栈

### 核心

* Go 1.22+
* Bubble Tea
* Bubbles（组件）
* Lipgloss（样式）

---

### 邮件协议

* IMAP（收邮件）
* SMTP（发邮件）

建议库：

* go-imap
* go-smtp / net/smtp

---

### 本地存储

* SQLite（缓存邮件）

---

# 🧩 四、模块架构

```text
lazymail/
├── cmd/
│   └── lazymail/
│       └── main.go
│
├── internal/
│   ├── tui/            # UI 层
│   │   ├── model.go
│   │   ├── layout.go
│   │   ├── sidebar.go
│   │   ├── maillist.go
│   │   ├── viewer.go
│   │
│   ├── mail/           # 邮件协议层
│   │   ├── imap.go
│   │   ├── smtp.go
│   │
│   ├── store/          # 本地缓存
│   │   └── sqlite.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── domain/         # 数据模型
│   │   └── mail.go
│   │
│   └── app/            # 应用协调层
│       └── service.go
│
└── config.yaml
```

---

# 📦 五、核心数据模型

## 邮件结构

```go
type Mail struct {
    ID        string
    Subject   string
    From      string
    To        []string
    Date      time.Time
    Body      string
    IsRead    bool
    Folder    string
}
```

---

# 🖥 六、UI 组件设计

## 1. Sidebar（左侧）

功能：

* 显示账号
* 显示文件夹（Inbox/Sent等）
* 支持鼠标点击切换

---

## 2. Mail List（中间）

功能：

* 显示邮件列表
* 支持：

  * 上下键
  * 鼠标点击
  * 滚轮
* 当前选中高亮

---

## 3. Mail Viewer（右侧）

功能：

* 显示邮件详情
* 支持滚动
* 支持 HTML 转文本

---

## 4. Status Bar（底部）

```text
[R] Reply  [F] Forward  [D] Delete  [A] Archive  [Q] Quit
```

---

# 🖱 七、交互设计

## 键盘

| 按键    | 行为   |
| ----- | ---- |
| ↑ ↓   | 移动   |
| Enter | 打开   |
| Tab   | 切换面板 |
| q     | 退出   |
| r     | 回复   |
| d     | 删除   |

---

## 鼠标

| 行为    | 功能 |
| ----- | -- |
| 点击    | 选择 |
| 双击    | 打开 |
| 滚轮    | 滚动 |
| hover | 高亮 |

---

# 🌐 八、邮件功能

## 必须实现

* 登录 IMAP
* 拉取邮件列表
* 拉取邮件详情
* 标记已读
* 删除邮件
* 发送邮件（SMTP）

---

## 可选（第二阶段）

* 搜索
* 标签
* 多账号切换
* 本地缓存

---

# 🧠 九、状态管理（关键）

```go
type Model struct {
    ActivePanel string

    Sidebar     SidebarModel
    MailList    MailListModel
    Viewer      ViewerModel

    Mails       []Mail
    SelectedIdx int

    Width  int
    Height int
}
```

---

# 🔄 十、核心事件流

```text
用户操作
   ↓
Bubble Tea Msg
   ↓
Update()
   ↓
修改 Model
   ↓
View()
   ↓
UI 更新
```

---

# 🧪 十一、开发阶段

## Phase 1（必须完成）

* UI 三栏布局
* 模拟数据
* 鼠标 + 键盘
* 切换邮件

---

## Phase 2

* 接入 IMAP
* 显示真实邮件

---

## Phase 3

* SMTP 发信
* 删除 / 标记

---

## Phase 4

* SQLite 缓存
* 搜索

---

## Phase 5（扩展）

* AI 总结邮件
* AI 自动回复
* 与 RPA 系统联动

---

# 🔐 十二、安全要求

* 不允许明文存储密码
* 支持环境变量读取
* 支持 OAuth（未来）

---

# 🚀 十三、扩展方向（非常重要）

未来必须支持：

### 1. AI 插件

```text
- 总结邮件
- 自动生成回复
- 提取任务
```

---

### 2. 自动化（你可以重点做）

```text
邮件 → 触发 RPA
邮件 → 调用 API
邮件 → 自动回复
```

👉 这部分是你最大优势

---

# 🧠 十四、给 AI 的执行指令

你可以直接丢给 AI：

```text
请基于 Go + Bubble Tea 实现一个 TUI 邮件客户端：

要求：
1. 三栏布局（Sidebar / MailList / Viewer）
2. 支持鼠标点击和滚动
3. 支持键盘导航
4. 使用状态驱动架构
5. UI 风格参考 lazygit
6. 第一阶段使用 mock 数据

输出：
- 完整项目结构
- 可运行代码
- 每个模块职责清晰
```

---

# 🧠 最后一句（很关键）

这个项目本质不是“邮件客户端”👇

👉 是一个：

**“终端里的工作流中枢（Mail-driven OS）”**

如果你做对了，可以变成：

* 邮件 + AI
* 邮件 + RPA
* 邮件 + 电商运营

👉 **这是可以做成产品的东西，不只是工具**

---

如果你下一步要，我可以帮你👇

👉 直接生成第一版代码骨架（可运行）
👉 或帮你设计 UI 交互细节（像 lazygit 那种体验级别）
