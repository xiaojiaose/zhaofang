# 分享房源地图接口文档（支持筛选）

## 一、概述

分享地图接口支持在分享时传递筛选条件，实现类似 `/center/house/xiaoquAgg` 的筛选能力。筛选参数通过 POST body 传递，与原有 token 参数一起构成完整请求。

---

## 二、接口列表

| # | 接口 | 方法 | 说明 |
|---|------|------|------|
| 1 | /center/house/share | POST | 创建分享链接 |
| 2 | /center/house/share/map | POST | 获取地图点位（支持筛选） |
| 3 | /center/house/share/list | POST | 获取点位下房源列表（支持筛选） |
| 4 | /center/type/options | GET | 获取筛选项数据源 |

---

## 三、接口详情

### 1. 创建分享链接

**请求**
```
POST /center/house/share
Content-Type: application/json
```

**Body 参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| expireDays | int | 否 | 有效期天数，默认 7 天 |

**请求示例**
```json
{
  "expireDays": 7
}
```

**响应**
```json
{
  "success": true,
  "data": {
    "token": "abc123xyz",
    "expireAt": 1704067200
  }
}
```

返回字段说明：
- `token`：分享凭证，用于后续接口调用
- `expireAt`：token 过期时间戳

---

### 2. 获取地图点位（支持筛选）

**请求**
```
POST /center/house/share/map
Content-Type: application/json
```
下面的筛选项 都是从 GET /center/type/options 这个接口来的
**Body 参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| token | string | 是 | 分享 token |
| rentType | string | 否 | 出租类型：整租 / 分整租 / 合租 |
| houseType | string | 否 | 户型：1居 / 2居 / 3居 / 4居+ / 开间 / 主卧 / 次卧 / 暗间 |
| price | int | 否 | 价格区间：1=500以下 / 2=500-1000 / 3=1000-1500 / 4=1500-2000 / 5=2000-2500 / 6=2500-3000 / 7=3000元以上 |
| feature | string | 否 | 特色标签，多个用逗号分隔，如："可短租,有电梯" |
| houseSource | string | 否 | 房源来源：commission=有返佣 / landlord=房东房源 / team=团队房源 |
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页数量，默认 20 |

**请求示例**
```json
{
  "token": "abc123xyz",
  "rentType": "整租",
  "houseType": "2居",
  "price": 3,
  "feature": "有电梯,可短租"
}
```

**响应**
```json
{
  "success": true,
  "data": {
    "list": [
      {
        "xiaoquId": 1,
        "name": "阳光小区",
        "latitude": "31.230416",
        "longitude": "121.473701",
        "num": 5
      }
    ],
    "expireAt": 1704067200,
    "userInfo": {
      "wxNo": "wx123",
      "wxNickName": "张三",
      "headerImg": "https://example.com/avatar.jpg",
      "phone": "13800138000"
    }
  }
}
```

返回字段说明：
- `list`：小区点位数组
  - `xiaoquId`：小区 ID
  - `name`：小区名称
  - `latitude`：纬度
  - `longitude`：经度
  - `num`：应用筛选条件后的房源数量
- `expireAt`：token 过期时间戳
- `userInfo`：分享人信息

> 注意：`num` 是应用筛选条件后的数量，非全部房源数量。筛选条件变化时，点位数量和 `num` 都会动态变化。

---

### 3. 获取点位下房源列表（支持筛选）

**请求**
```
POST /center/house/share/list
Content-Type: application/json
```

**Body 参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| token | string | 是 | 分享 token |
| xiaoquId | int | 是 | 小区 ID |
| rentType | string | 否 | 出租类型 |
| houseType | string | 否 | 户型 |
| price | int | 否 | 价格区间 |
| feature | string | 否 | 特色标签 |
| houseSource | string | 否 | 房源来源 |
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页数量，默认 20 |

**请求示例**
```json
{
  "token": "abc123xyz",
  "xiaoquId": 1,
  "rentType": "整租",
  "page": 1,
  "pageSize": 20
}
```

**响应**
```json
{
  "success": true,
  "data": {
    "list": [
      {
        "id": 1,
        "doorNo": "1号楼 1单元 101室",
        "roomCode": "101",
        "price": 2500,
        "rentType": "整租",
        "houseType": "2居",
        "feature": "有电梯,可短租",
        "attachments": {
          "house": [{"url": "https://example.com/img1.jpg"}]
        },
        "hasPic": true,
        "status": "待出租"
      }
    ],
    "total": 100,
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 4. 获取筛选项（前端下拉菜单数据源）

**请求**
```
GET /center/type/options
```

**响应**
```json
{
  "success": true,
  "data": {
    "rentType": [
      {"name": "整租", "houseType": ["1居", "2居", "3居", "4居+", "开间"], "feature": ["可短租", "包物业", "有电梯", "密码看房", "可办公注册", "协助对接房东", "可带看分佣"]},
      {"name": "分整租", "houseType": ["1居", "2居", "3居", "4居+"], "feature": ["带阳台", "有电梯", "有原卫", "朝南", "有燃气", "可短租", "协助对接房东", "可带看分佣"]},
      {"name": "合租", "houseType": ["2居", "3居", "4居+"], "feature": ["带阳台", "可短租", "有电梯", "有独卫", "可做饭", "纯女生", "协助对接房东", "可带看分佣"]}
    ],
    "price": {
      "1": "500以下", "2": "500-1000元", "3": "1000-1500元",
      "4": "1500-2000元", "5": "2000-2500元", "6": "2500-3000元", "7": "3000元以上"
    },
    "houseSource": {"2": "有返佣", "4": "团队房源"},
    "canViewTeamHouse": true
  }
}
```

返回字段说明：
- `rentType`：出租类型数组，每个类型包含可选的 houseType 和 feature
- `price`：价格区间映射（数字 key 对应 price 参数值）
- `houseSource`：房源来源映射（用于 houseSource 参数）
- `canViewTeamHouse`：当前用户是否有权限查看团队房源

---

## 四、使用流程

1. **创建分享**：`POST /center/house/share`（仅传 expireDays）
2. **打开分享页**：使用 token 调用 `POST /center/house/share/map` 获取地图点位（带上筛选条件）
3. **点击点位**：使用 token + xiaoquId 调用 `POST /center/house/share/list` 获取房源列表（带上筛选条件）

---

## 五、前端实现建议

### 分享 URL 记忆筛选条件

分享时可以将筛选条件编码到分享内容中（如微信分享卡片的 description），打开时解析并填充到筛选项。

```javascript
// 分享时：筛选条件作为上下文传递
const shareContent = {
  token: 'abc123xyz',
  filters: {
    rentType: '整租',
    houseType: '2居',
    price: 3
  }
}

// 打开分享页时：从分享内容恢复筛选项
this.searchForm = { ...shareContent.filters }
```

### 筛选项联动

- `rentType` 变化时，`houseType` 和 `feature` 的可选值联动（参考 `FilterOptions1` 返回的结构）
- 价格区间传的是数字 key（1-7），不是价格本身

---

## 六、注意事项

1. **token 有效期**：默认 7 天，过期后接口返回错误，前端需跳回首页
2. **团队房源**：普通用户分享时，后端自动过滤团队房源；团队权限用户不受影响
3. **筛选条件为空**：表示不过滤，返回该用户所有待出租房源
4. **点位数量变化**：筛选条件变化时，`num` 会动态减少，地图点位数量也可能变化
5. **所有接口均需登录态**：除了 `/center/house/share/map` 和 `/center/house/share/list` 可在未登录状态下访问外，其他接口需要有效的用户会话