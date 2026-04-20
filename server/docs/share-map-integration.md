# 个人房源地图对接说明

这个页面的核心是“两步走”：

1. 先拉地图点位
2. 再点位下钻到房源列表

## 接口顺序

### 1. 创建分享链接
- `POST /center/house/share`
- 返回分享 `token`
- `token` 默认 7 天有效

### 2. 获取地图点位
- `GET /center/house/share?token=xxx`
- 返回值里的 `list` 是按 `xiaoqu_id` 聚合后的地图点位
- 每个点位包含：
  - `xiaoquId`
  - `name`
  - `latitude`
  - `longitude`
  - `num`

### 3. 点击地图点位后获取房源列表
- `GET /center/house/share/list?token=xxx&xiaoquId=123`
- 支持分页参数：
  - `page`
  - `pageSize`
- 返回点位下所有房源列表

## 数据规则

- 地图点位坐标来自 `xiao_qu` 表
- 点位下房源列表通过 ZincSearch 查询
- 查询条件固定为：
  - `owner = 分享人`
  - `status = 待出租`
  - 点击点位时再附加 `xiaoqu_id = 当前点位`

## 前端展示建议

- 地图上只画点位，不直接铺房源卡片
- 点击点位时再展示该点位下的房源列表
- 过期 token 返回失败后，前端跳回首页

## 备注

- 如果历史房源索引里缺少 `owner` 字段，需要重新回灌 ZincSearch 后才能完整命中
- 当前后端已把新写入的房源同步带上 `owner` 字段
