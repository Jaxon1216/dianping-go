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

- 学习总览、用户画像、讲解契约和工具：[docs/learning/README.md](docs/learning/README.md)
- 项目地图、路线和模块状态：[docs/learning/project.md](docs/learning/project.md)
- 主题笔记与索引：[docs/learning/notes.md](docs/learning/notes.md)
- 变更记录：[docs/learning/changelog.md](docs/learning/changelog.md)
