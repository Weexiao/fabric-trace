# 各部分功能实现流程泳道图（fabric-trace）

- [x] 覆盖主体：用户、前端（Web/Vue）、后端（Go/Gin）、MySQL、Fabric 链码、文件存储
- [x] 覆盖典型流程：登录鉴权、溯源上链、溯源查询（匿名/鉴权）、附件上传/下载

> 说明：以下泳道图使用 Mermaid `sequenceDiagram`（时序泳道图）表示多主体交互流程，可直接粘贴到支持 Mermaid 的 Markdown 环境中。

---

## 1) 用户登录与 JWT 鉴权流程

```mermaid
sequenceDiagram
    autonumber
    actor U as 用户
    participant FE as 前端（Vue）
    participant BE as 后端（Go/Gin）
    participant DB as MySQL
    participant CC as Fabric链码

    U->>FE: 输入账号/密码并点击登录
    FE->>BE: POST /login (username,password)
    BE->>DB: 查询用户信息/校验密码
    DB-->>BE: 用户记录（或未找到/密码错误）
    alt 校验失败
        BE-->>FE: 登录失败（message）
        FE-->>U: 提示错误
    else 校验成功
        BE->>CC: Query GetUserType(userID)
        CC-->>BE: userType/role
        BE-->>FE: 返回 JWT + 角色信息
        FE->>FE: 保存 JWT（本地存储/状态管理）
        FE-->>U: 跳转到业务页面
    end
```

---

## 2) 溯源信息录入与上链流程（受 JWT 保护）

```mermaid
sequenceDiagram
    autonumber
    actor U as 用户（企业参与方）
    participant FE as 前端（溯源录入页面）
    participant BE as 后端（Go/Gin）
    participant AUTH as JWT中间件
    participant FS as 文件/图片存储
    participant CC as Fabric链码

    U->>FE: 填写溯源码/环节字段，选择图片（可选）
    FE->>BE: POST /uplink (form-data) + Authorization: Bearer JWT
    BE->>AUTH: 校验 JWT、解析 userID/userType
    alt Token 无效/过期
        AUTH-->>BE: 拒绝访问
        BE-->>FE: 401/无权限
        FE-->>U: 提示重新登录
    else Token 有效
        opt 图片上传
            BE->>FS: 保存图片（按 SHA-256 命名/落盘）
            FS-->>BE: 返回图片访问路径/文件名
        end
        BE->>CC: Invoke Uplink/WriteTrace(args...)
        CC-->>BE: 上链交易结果（txId/状态）
        BE-->>FE: 返回上链成功/失败信息
        FE-->>U: 展示提交结果
    end
```

---

## 3) 溯源查询流程（匿名/消费者入口）

```mermaid
sequenceDiagram
    autonumber
    actor U as 用户（消费者/访客）
    participant FE as 前端（溯源查询页）
    participant BE as 后端（Go/Gin）
    participant CC as Fabric链码

    U->>FE: 输入溯源码并点击查询
    FE->>BE: POST /getIndustrialProductInfo (traceCode)
    BE->>CC: Query GetIndustrialProductInfo(traceCode)
    CC-->>BE: 链上溯源信息（当前状态）
    BE-->>FE: 返回查询结果（JSON）
    FE-->>U: 渲染溯源信息（环节/字段/图片链接等）
```

---

## 4) 历史追溯查询流程（受 JWT 保护）

```mermaid
sequenceDiagram
    autonumber
    actor U as 用户（企业参与方）
    participant FE as 前端（历史追溯页面）
    participant BE as 后端（Go/Gin）
    participant AUTH as JWT中间件
    participant CC as Fabric链码

    U->>FE: 输入溯源码并查询历史记录
    FE->>BE: POST /getIndustrialProductHistory + Authorization: Bearer JWT
    BE->>AUTH: 校验 JWT
    alt Token 无效
        BE-->>FE: 401/无权限
        FE-->>U: 提示重新登录
    else Token 有效
        BE->>CC: Query GetIndustrialProductHistory(traceCode)
        CC-->>BE: 该溯源码的链上历史记录集合
        BE-->>FE: 返回历史追溯数据
        FE-->>U: 列表/时间线可视化渲染
    end
```

---

## 5) 附件上传并将 Manifest 上链流程（受 JWT 保护）

```mermaid
sequenceDiagram
    autonumber
    actor U as 用户（企业参与方）
    participant FE as 前端（附件管理）
    participant BE as 后端（Go/Gin）
    participant AUTH as JWT中间件
    participant FS as 链下存储（加密文件）
    participant CC as Fabric链码

    U->>FE: 选择附件并填写溯源码
    FE->>BE: POST /file/upload (file + traceabilityCode) + JWT
    BE->>AUTH: 校验 JWT、确定角色/权限
    alt Token 无效/无权限
        BE-->>FE: 401/拒绝
        FE-->>U: 提示无权限
    else 通过
        BE->>BE: 加密附件（如 AES-GCM）并生成元数据（hash/size/name）
        BE->>FS: 保存加密文件（链下）
        FS-->>BE: 返回 fileID/存储地址
        BE->>CC: Invoke PutFileManifest(traceCode, manifest...)
        CC-->>BE: 上链结果（txId/状态）
        BE-->>FE: 上传成功（fileID + 链上凭证）
        FE-->>U: 刷新文件列表
    end
```

---

## 6) 附件列表与下载流程（受 JWT 保护）

```mermaid
sequenceDiagram
    autonumber
    actor U as 用户（企业参与方）
    participant FE as 前端（附件列表）
    participant BE as 后端（Go/Gin）
    participant AUTH as JWT中间件
    participant CC as Fabric链码
    participant FS as 链下存储

    U->>FE: 查看某溯源码下附件列表
    FE->>BE: POST /file/list (traceabilityCode) + JWT
    BE->>AUTH: 校验 JWT
    alt Token 无效
        BE-->>FE: 401
        FE-->>U: 提示重新登录
    else Token 有效
        BE->>CC: Query GetFileManifestsByTrace(traceCode)
        CC-->>BE: manifests（fileID/名称/hash/时间等）
        BE-->>FE: 返回文件清单
        FE-->>U: 列表展示（下载按钮）
    end

    U->>FE: 点击下载某个附件
    FE->>BE: GET /file/download/:fileID + JWT
    BE->>AUTH: 校验 JWT
    alt Token 无效
        BE-->>FE: 401
        FE-->>U: 提示重新登录
    else Token 有效
        BE->>FS: 读取加密文件（fileID）
        FS-->>BE: 加密文件内容
        BE->>BE: 解密并校验（可选）
        BE-->>FE: 返回文件流（attachment）
        FE-->>U: 浏览器保存/打开文件
    end
```

---

## 在你的文档中引用

如果你的 Markdown 渲染器支持 Mermaid，可以直接把本文件中的某个图复制到目标章节。

如果目标平台不支持 Mermaid：
- 可以用支持 Mermaid 的编辑器导出为 PNG/SVG 再插入；
- 或我可以在仓库里补充一个“离线导出脚本”（不依赖 sudo 的方式）。

