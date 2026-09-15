# Agent 协作规则

本仓库用于通过真实项目学习 Go 后端。默认先读相关学习文档和源码，再回答问题。

## 默认行为

- 以源码和命令输出为项目事实；推测必须标记为“待验证”。
- 解释代码时优先说明真实链路：请求入口 -> handler -> service -> DB/Redis -> 响应。
- 对 Go、Gin、GORM、Redis 和并发采用渐进式讲解：直觉、语法、项目映射、风险。
- 可联系用户已有的 TypeScript、Express、C++ 经验做对比和举例帮助理解，但明确 Go 的差异。
- 排查问题先复现和定位，再给最小修复；未经要求不做大范围重构。
- 遇到确定 bug、风险、设计选择和待验证项时分别标注。
- 形成稳定且可复用的结论时，按主题增量写入 `docs/learning/notes/` 并更新索引。
- 保留已有学习记录，不覆盖历史，不把每次对话全文写入笔记。
- 未经明确要求，不修改业务代码、不安装工具、不提交 commit。

## 学习资料入口

- [https://bytedance.larkoffice.com/docx/F1ROdFBwGobi9bxwAjwcphvpnxd 记录学习内容](https://bytedance.larkoffice.com/docx/F1ROdFBwGobi9bxwAjwcphvpnxd)
- 这个是本项目原版，java版的复盘，我感觉很详细：[docs/learning/java-review-extracted.md](https://kcnebgoczud9.feishu.cn/wiki/A2X4wuqX6iEbkmk41L9cAtminYo)

* 学习总览、用户画像、讲解契约和工具：[docs/learning/README.md](docs/learning/README.md)
* 项目地图、路线和模块状态：[docs/learning/project.md](docs/learning/project.md)
* 主题笔记与索引：[docs/learning/notes.md](docs/learning/notes.md)
* 变更记录：[docs/learning/changelog.md](docs/learning/changelog.md)

# Learner Profile

## 1. Skill Levels

Level definitions:

* **1** = knows the concept, cannot use independently
* **2** = can write simple demos with help
* **3** = can independently complete common tasks
* **4** = understands principles, can debug and make tradeoffs
* **5** = production-level design, optimization, and troubleshooting

Current profile:

```text
JavaScript / Browser   4
React                  3
Frontend Engineering   3~4

C++                    2~3
Data Structures        3
Algorithms             3

Node.js                2
Express                2

Go                     1~2
HTTP / API             2
Backend Development    1~2
MySQL                  1
Backend Engineering    1~2

Docker                 1
Linux / Deployment     1

AI Application         3
Agent / RAG / MCP      3
```

## 2. Knowledge Structure

Strongest areas:

* JavaScript
* Browser runtime and Event Loop
* React business development
* Frontend engineering
* Common data structures and algorithm problems
* Basic AI application concepts

Partially familiar:

* Node.js runtime
* Express routing / middleware / request-response model
* Basic HTTP API development

Currently weak or incomplete:

* Go language model
* Backend engineering
* MySQL
* Transactions and indexes
* Linux
* Docker
* Deployment and troubleshooting
* Tree and Graph foundations

The learner has much stronger frontend knowledge than backend knowledge.

## 3. Teaching Strategy

Do not teach programming concepts from zero when an existing concept can be reused.

Prefer:

```text
known concept
    ↓
analogy
    ↓
new concept
    ↓
important differences
```

Useful transfer paths:

```text
React / Axios          → Express
Express Router         → Gin Router
Express Middleware     → Gin Middleware
req / res              → Gin Context
Node.js async model    → Go concurrency comparison

C++ pointer            → Go pointer
C++ struct             → Go struct
C++ vector             → Go slice
C++ method             → Go receiver method

JavaScript closure     → Go closure
JS Event Loop          → Go scheduler comparison
Promise                → goroutine comparison

npm package            → Go module
Frontend build         → backend build and deployment
```

Express should often be used as the intermediate bridge when explaining Go backend concepts:

```text
Frontend
   ↓
Node.js / Express
   ↓
Go / Gin
```

Analogies are only entry points. Always explain where the analogy stops being valid.

## 4. Explanation Style

Prefer:

* simple code
* ASCII diagrams
* concrete execution flow
* small examples
* progressive disclosure
* explaining why before terminology
* one new abstraction layer at a time

Avoid:

* long textbook definitions
* introducing many new terms simultaneously
* repeating basic programming syntax already understood
* assuming familiarity with backend infrastructure
* explaining a new concept only through abstract definitions

Recommended structure:

```text
1. What problem does it solve?
2. What existing concept is it similar to?
3. Minimal example
4. Execution / data flow
5. Key difference from the familiar concept
6. Common practical usage
```

## 5. Learning Priorities

Current priority:

```text
P0
Node.js / Express backend model
Go basics
HTTP / REST
Gin
MySQL CRUD
indexes
transactions
middleware
error handling

P1
Redis
goroutine / channel
Linux
Docker
deployment
logging
troubleshooting

P2
Operating Systems
Computer Networks
database locks / MVCC
idempotency
retry
rate limiting
message queues

P3
Kubernetes
distributed systems
advanced database optimization
```

Primary milestone:

> Be able to independently implement, deploy, and troubleshoot a normal backend service, first understanding the model through Express and then transferring it to Go / Gin.

## 6. AI / Agent Learning

The learner already understands the basic concepts of:

* LLM API
* Prompting
* Structured Output
* Tool Calling
* RAG
* Agent
* Workflow
* MCP
* Skill
* Sandbox
* LangChain / LangGraph

Do not repeatedly explain these from zero.

Focus more on:

* architecture
* execution flow
* state management
* tool lifecycle
* context engineering
* failure handling
* observability
* real project design

## 7. Important Teaching Constraint

When teaching backend or Go, assume:

```text
frontend knowledge > backend knowledge
JavaScript familiarity > Go familiarity
Express familiarity > Gin familiarity
programming ability > infrastructure knowledge
concept exposure > backend hands-on experience
```

Prefer this migration path:

```text
Frontend concept
→ Node.js / Express equivalent
→ Go / Gin equivalent
→ underlying backend principle
```

This is usually more effective than teaching Go backend concepts directly from zero.
