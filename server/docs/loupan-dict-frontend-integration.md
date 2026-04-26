# 楼盘字典前端对接文档

## 1. 目标

后台楼盘字典页面只维护两类数据：

1. 小区基础信息（`xiao_qu`）
2. 楼盘三层字典（楼栋/单元/房号：`dict_building`、`dict_unit`、`dict_house`）

---

## 2. 鉴权与通用规范

- 接口前缀：`/api`
- 请求头：
  - `x-token: <jwt>`
  - `x-user-id: <userId>`
  - `Content-Type: application/json`
- 统一响应：

```json
{
  "code": 0,
  "data": {},
  "msg": "成功"
}
```

关联约束（重要）：

- `xiao_qu` 与 `dict_building/dict_unit/dict_house` 的关联键是 `community_id`。
- `/api/xiaoqu/edit` 在新增或编辑时，若 `community_id` 为空，后端会自动回填为当前小区ID。
- 发起楼盘字典保存前，仍建议先调用一次 `/api/xiaoqu/edit`，确保 `community_id` 已补齐。

---

## 3. 页面调用顺序（推荐）

## 3.1 列表页

1. `POST /api/xiaoqu/list`：分页查询小区
2. `POST /api/xiaoqu/edit`：新增或更新小区基础信息

## 3.2 编辑楼盘字典弹窗

1. `GET /api/xiaoqu/show?id={xiaoquId}`：拿小区基础信息
2. `GET /api/xiaoqu/dict/tree?id={xiaoquId}`：拿楼栋-单元-房号树
3. 新增/修改走 `POST /api/xiaoqu/dict/upsert`
4. 删除走 `POST /api/xiaoqu/dict/delete`
5. 保存成功后重新调用 `GET /api/xiaoqu/dict/tree` 刷新

---

## 4. 接口明细

### 4.1 小区列表

- 方法：`POST`
- 路径：`/api/xiaoqu/list`
- 请求体：

```json
{
  "page": 1,
  "pageSize": 20,
  "keyword": "留村",
  "cityId": "",
  "areaId": "",
  "districts": []
}
```

---

### 4.2 小区编辑（新增/更新）

- 方法：`POST`
- 路径：`/api/xiaoqu/edit`
- 说明：
  - `ID=0` 新增
  - `ID>0` 更新

```json
{
  "ID": 10,
  "name": "留村东",
  "city": "石家庄",
  "area": "开发区",
  "area_id": 1,
  "position": "天山海世界",
  "districts": "天山海世界",
  "district_ids": "|1|2|",
  "address": "",
  "latitude": "38.0",
  "longitude": "114.0",
  "able": true
}
```

---

### 4.3 楼盘字典树查询

- 方法：`GET`
- 路径：`/api/xiaoqu/dict/tree?id={xiaoquId}`

响应示例：

```json
{
  "xiaoquId": 10,
  "communityId": 12334,
  "buildings": [
    {
      "id": 1,
      "buildingOpenId": "b-uuid-1",
      "name": "1",
      "units": [
        {
          "id": 11,
          "unitOpenId": "u-uuid-1",
          "name": "A",
          "houses": [
            {
              "id": 111,
              "houseOpenId": "h-uuid-1",
              "name": "103"
            }
          ]
        }
      ]
    }
  ]
}
```

---

### 4.4 楼盘字典新增/修改接口

- 方法：`POST`
- 路径：`/api/xiaoqu/dict/upsert`
- 说明：同一请求内支持“楼栋+单元+房号”多条新增/修改混合提交
- 请求体：

```json
{
  "xiaoquId": 10,
  "buildingOps": [
    { "name": "2" },
    { "id": 1, "name": "1号楼" }
  ],
  "unitOps": [
    { "buildingOpenId": "b-uuid-1", "name": "B" },
    { "id": 11, "name": "A单元" }
  ],
  "houseOps": [
    { "unitOpenId": "u-uuid-1", "name": "104" },
    { "id": 111, "name": "103室" }
  ]
}
```

字段说明：

- `buildingOps[]`：楼栋新增/修改；`id` 为空=新增，`id` 有值=修改
- `unitOps[]`：单元新增/修改；`id` 为空=新增（需 `buildingOpenId`），`id` 有值=修改
- `houseOps[]`：房号新增/修改；`id` 为空=新增（需 `unitOpenId`），`id` 有值=修改

必填规则（upsert）：

| 字段 | 场景 | 是否必填 | 说明 |
|---|---|---|---|
| `xiaoquId` | 全部 | 是 | 小区ID |
| `buildingOps[].id` | 楼栋修改 | 是 | 楼栋记录ID |
| `buildingOps[].name` | 楼栋新增/修改 | 是 | 楼栋名称 |
| `buildingOps[].buildingOpenId` | 楼栋新增 | 否 | 不传后端自动生成 |
| `unitOps[].id` | 单元修改 | 是 | 单元记录ID |
| `unitOps[].buildingOpenId` | 单元新增 | 是 | 目标楼栋OpenID |
| `unitOps[].name` | 单元新增/修改 | 是 | 单元名称 |
| `unitOps[].unitOpenId` | 单元新增 | 否 | 不传后端自动生成 |
| `houseOps[].id` | 房号修改 | 是 | 房号记录ID |
| `houseOps[].unitOpenId` | 房号新增 | 是 | 目标单元OpenID |
| `houseOps[].name` | 房号新增/修改 | 是 | 房号名称 |
| `houseOps[].houseOpenId` | 房号新增 | 否 | 不传后端自动生成 |

推荐调用策略（关键）：

1. 如果前端一次新增了“楼栋+单元+房号”，请分 2~3 次调用，不要一次提交。  
2. 第一次：先 upsert 楼栋；成功后用返回树拿到新楼栋 `id`。  
3. 第二次：再 upsert 单元；成功后拿到新单元 `id`。  
4. 第三次：最后 upsert 房号。  

原因：

- 关联是通过 `buildingOpenId` / `unitOpenId` 传递，先保存上层节点能拿到稳定 OpenID，再保存下层更稳妥。

upsert 返回：

- `data` 为最新字典树（同 `dict/tree`），前端可直接替换本地树。

---

### 4.5 楼盘字典删除接口

- 方法：`POST`
- 路径：`/api/xiaoqu/dict/delete`
- 说明：支持同一次请求批量删除多层节点
- 请求体：

```json
{
  "xiaoquId": 10,
  "buildingIds": [8],
  "unitIds": [99],
  "houseIds": [1001]
}
```

必填规则（delete）：

| 字段 | 是否必填 | 说明 |
|---|---|---|
| `xiaoquId` | 是 | 小区ID |
| `buildingIds` | 否 | 楼栋ID数组 |
| `unitIds` | 否 | 单元ID数组 |
| `houseIds` | 否 | 房号ID数组 |

删除顺序（服务端）：

1. 先删 `houseIds`
2. 再删 `unitIds`
3. 最后删 `buildingIds`

delete 返回：

- `data` 为删除后的最新字典树（同 `dict/tree`）。

---

## 5. 删除级联规则

1. 删除楼栋：会级联删除该楼栋下所有单元和房号
2. 删除单元：会级联删除该单元下所有房号
3. 删除房号：仅删除当前房号

---

## 6. 查询级联接口（联动下拉可直接用）

1. 按小区查楼栋  
- `GET /api/base/building?id={xiaoquId}`

2. 按楼栋查单元  
- `GET /api/base/unit?id={buildingOpenId}`

3. 按单元查房号  
- `GET /api/base/house?id={unitOpenId}`

说明：

- 这 3 个是“查询联动”接口，适合表单下拉级联。
- 楼盘字典维护页若使用 `dict/tree` 全量树渲染，也可不调用这 3 个。

联动接口返回结构（统一）：

```json
[
  { "id": "xxx", "name": "xxx" }
]
```

- `/api/base/building` 返回 `buildingOpenId/name`
- `/api/base/unit` 返回 `unitOpenId/name`
- `/api/base/house` 返回 `houseOpenId/name`

---

## 7. 常见错误与提示文案

建议前端将 `msg` 直接展示：

| 场景 | 典型 msg |
|---|---|
| 小区未配置关联键（通常是未先保存小区） | `当前小区未配置community_id，无法维护楼盘字典` |
| 缺少主参数 | `小区ID不能为空` |
| 修改缺少ID | `楼栋ID不存在: xxx` / `单元ID不存在: xxx` / `房号ID不存在: xxx` |
| 新增缺少父标识 | `新增单元必须传buildingOpenId` / `新增房号必须传unitOpenId` |
| 缺少名称 | `楼栋名称不能为空` / `单元名称不能为空` / `房号名称不能为空` |

---

## 8. 前端实现建议

1. 页面本地维护一份 `dict/tree` 数据用于渲染
2. 新增/修改调用 `POST /api/xiaoqu/dict/upsert`，删除调用 `POST /api/xiaoqu/dict/delete`
3. 每次成功后用返回 `data` 直接刷新本地树；必要时再调 `dict/tree` 二次校验
4. 为保存按钮加 `loading`，避免重复提交
5. 后端 `msg` 直接 toast 展示即可
6. 对新增链路采用“楼栋 -> 单元 -> 房号”分步提交，避免父ID缺失

---

## 9. 最小联调清单

1. 新增楼栋后，树里出现新楼栋且有 `id`
2. 新增单元时使用新楼栋 `id`，保存成功
3. 新增房号时使用新单元 `id`，保存成功
4. 删除单元后，其房号全部消失
5. 删除楼栋后，其单元和房号全部消失
6. 新增小区先调用一次 `/api/xiaoqu/edit` 后，再做 upsert/delete
