# API Change Log - 2026-04-16

本次文档用于标记这轮后端改动里，Swagger 中哪些接口是“新增”以及哪些接口是“变更”。

规则：
- `[新增]`：本轮新增的接口
- `[变更]`：原有接口本轮修改了入参、返回字段或业务行为

## 新增接口

### 小程序端
- `GET /center/reward/recent`
  - 最近联系过的房源发布人记录
- `POST /center/reward/apply`
  - 发起出房有礼申请
- `POST /center/reward/publisher/list`
  - 发布人查看出房有礼审核列表
- `POST /center/reward/publisher/action`
  - 发布人确认或拒绝出房有礼申请
- `POST /center/house/share`
  - 创建我的房源地图分享链接
- `GET /center/house/share`
  - 通过分享 token 获取地图房源列表
- `POST /center/profile`
  - 小程序设置微信资料

### 后台端
- `POST /api/house/reward/list`
  - 后台成交有礼审核列表
- `POST /api/house/reward/action`
  - 后台操作成交有礼审核状态
- `POST /api/house/contactQuota/grant`
  - 增加经纪人联系方式查看次数
- `POST /api/house/contactQuota/list`
  - 查看联系方式次数流水
- `POST /api/house/batchUpload`
  - Excel 批量上传房源

## 变更接口

### 小程序端
- `GET /center/house/view`
  - 房东房源隐藏门牌号
- `GET /center/house/mobile`
  - 房东房源查看联系方式时扣减查看次数
- `POST /center/house/xiaoquAgg`
  - 新增返佣、房东房源、团队房源筛选逻辑
  - 普通用户后端默认过滤团队房源
- `POST /center/house/xiaoquAggList`
  - 新增返佣、房东房源、团队房源筛选逻辑
  - 普通用户后端默认过滤团队房源
- `POST /center/house/listByXiaoqu`
  - 普通用户后端默认过滤团队房源
- `POST /center/house/my`
  - 返回微信资料、返佣字段、上架额度统计
- `POST /center/house/create`
  - 支持房东房源、返佣金额、团队房源联动字段
- `POST /center/house/edit`
  - 支持房东房源、返佣金额、团队房源联动字段
- `GET /center/options`
  - 筛选项增加房东房源相关配置
- `GET /center/type/options`
  - 筛选项增加房东房源和新增亮点

### 后台端
- `POST /api/house/list`
  - 返回房源发布人微信资料
  - 支持返佣、团队房源相关数据展示
- `POST /api/house/my`
  - 返回微信资料、返佣字段、上架额度统计
- `POST /api/house/create`
  - 支持房东房源、返佣金额、团队房源联动字段
- `POST /api/house/edit`
  - 支持房东房源、返佣金额、团队房源联动字段
- `GET /api/house/options`
  - 筛选项增加房东房源相关配置
- `GET /api/house/type/options`
  - 筛选项增加房东房源和新增亮点
- `GET /api/house/statis/view`
  - 数据中心增加出房有礼申请统计和按手机号汇总能力
- `POST /api/statis/list`
  - 帖子数据返回房源发布人微信资料

### 用户管理
- `POST /user/admin_register`
  - 增加找房超市标识、上架额度、联系方式查看次数
- `POST /user/saler_register`
  - 增加找房超市标识、上架额度、联系方式查看次数
- `PUT /user/setUserInfo`
  - 支持编辑找房超市标识、上架额度、联系方式查看次数、微信资料
- `PUT /user/SetSelfInfo`
  - 支持编辑找房超市标识、上架额度、联系方式查看次数、微信资料

## 使用建议

- 在 Swagger 页面里，通过 `Summary` 前缀直接看：
  - `[新增]`
  - `[变更]`
- 需要核对细节时，再对照本文件查看影响点和字段变化。
