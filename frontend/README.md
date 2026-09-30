# 黑马点评 - Vue3 重构版

基于 Vue3 + Vite + Element Plus + Pinia 重构的黑马点评前端项目。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Vite** - 下一代前端构建工具
- **Vue Router 4** - 官方路由管理器
- **Pinia** - 状态管理库
- **Element Plus** - UI 组件库
- **Axios** - HTTP 客户端

## 项目结构

```
hmdp/
├── public/                 # 静态资源
│   └── imgs/              # 图片资源
├── src/
│   ├── api/               # API 接口
│   ├── components/        # 公共组件
│   ├── router/            # 路由配置
│   ├── stores/            # Pinia 状态管理
│   ├── utils/             # 工具函数
│   ├── views/             # 页面视图
│   ├── App.vue            # 根组件
│   └── main.js            # 入口文件
├── index.html             # HTML 模板
├── package.json           # 项目依赖
├── vite.config.js         # Vite 配置
└── .eslintrc.cjs          # ESLint 配置
```

## 功能模块

- **首页** - 商户类型展示、热门博客列表
- **商户列表** - 按类型筛选、排序、分页加载
- **商户详情** - 商户信息、优惠券秒杀
- **博客详情** - 图片轮播、点赞、关注
- **博客编辑** - 图片上传、关联商户
- **个人中心** - 我的笔记、关注列表
- **用户主页** - 他人主页、共同关注
- **登录** - 验证码登录、密码登录

## 安装运行

```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build

# 代码检查
npm run lint
```

## 环境要求

- Node.js >= 16
- npm >= 8

## 接口代理

开发环境下，API 请求会代理到 `http://localhost:8081`，可在 `vite.config.js` 中修改。
