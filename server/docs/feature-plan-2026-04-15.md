# Feature Plan 2026-04-15

## Background

This document captures the confirmed implementation plan for the current
house-resource, map-house-finding, reward-audit, and batch-upload changes.
It is based on:

- Current routes in `server/router/center/center.go`
- Current admin routes in `server/router/house_resource/enter.go`
- Current system routes in `server/router/system/enter.go`
- Current schema in `server/docs/zhaofangtest-schema.sql`

## Important Existing Facts

### Route boundaries

- `server/router/center/center.go` is for WeChat Mini Program routes
- `server/router/house_resource/enter.go` is for admin-side house management
- `server/router/system/enter.go` is for admin-side system management

### Existing key interfaces

- `/center/house/view`: house detail
- `/center/house/mobile`: get house contact phone
- `/center/house/xiaoquAgg`: map aggregation
- `/center/house/xiaoquAggList`: map list
- `/center/house/listByXiaoqu`: list by xiaoqu
- `/center/house/my`: my houses
- `/center/favorite/add`, `/center/favorite/del`, `/center/favorite/list`

### Existing data facts

- `house_resources` already contains:
  - `house_type`
  - `rent_type`
  - `price`
  - `feature`
  - `owner`
  - `status`
  - `approval_status`
  - `phone`
  - `follow`
  - `view`
  - `click`
  - `shared`
- `sys_users` already contains:
  - `phone`
  - `header_img`
  - `wx_nick_name`
  - `wx_no`
  - `openid`
- `wx_users` exists in schema, but current business usage should be based on
  `sys_users`

### Existing behavior chain

- `/center/house/mobile` already records click behavior
- `visit_daily` already stores click/view/follow/shared records
- These records are the current base for "recently contacted publishers"

## Front-Prerequisite Checks

These must be clarified before implementing the full feature set.

### 1. WeChat profile persistence

Current business expectation:

- Avatar, nickname, and wechat profile data should come from `sys_users`

Current risk:

- The fields exist, but data is not being successfully written from the WeChat
  side into `sys_users`

Must verify:

- Mini Program login flow
- Authorization flow
- Why profile data is not being persisted into `sys_users`
- How `openid` and login flow currently write user profile fields

Confirmed implementation direction:

- `openid`: keep using current login flow
- `phone`: keep using current WeChat phone authorization flow
- `wx_nick_name` and `header_img`:
  - should be collected on the Mini Program frontend
  - then sent to backend
  - then persisted into `sys_users`
- `wx_no`:
  - should not rely on automatic WeChat acquisition
  - should be filled by user on frontend
  - then persisted into `sys_users.wx_no`

### 2. Contact-record source

Current expectation:

- "Recently contacted publishers" should come from contact/click records
- `/center/house/mobile` is the current contact entry

Must verify:

- Whether current `visit_daily.click` is enough to reconstruct:
  - applicant user
  - house
  - publisher user
- If not enough, add a dedicated contact log table

### 3. Search index sync

Current map and house filtering rely heavily on ZincSearch / ES-style indexing.

Must sync at least:

- `is_team_house`
- `commission_price`
- `house_type`

## Business Rules Confirmed

### Team houses

- Team house = house published by a user with the "find-house-supermarket"
  flag
- A redundant field should be stored on the house record for easier querying
- Suggested field: `is_team_house`

### Landlord house

- `房东房源` is a house type
- It is not a tag

### Reward application source

- "Recently contacted publishers" comes from the contact/click record chain
- `/center/house/mobile` is the current existing entry point

### Contact click quota

- Viewing the contact info for a landlord house costs 1 click quota
- The quota consumed is from the current logged-in user

### Batch upload mode

- Batch upload is sync mode
- For the same user, houses not present in the current Excel upload should be
  automatically taken off shelf

## Data Structure Changes

## 1. `sys_users`

Add fields:

- `is_find_house_supermarket`
- `publish_quota_total`
- `contact_view_quota_total`
- Optional: `can_publish_house`

Used for:

- Whether the user can see team houses
- How many houses the user can keep on shelf
- How many landlord-house contact clicks the user can use

## 2. `house_resources`

Add fields:

- `commission_price`
- `is_team_house`

Extend existing semantics:

- `house_type` adds `房东房源`
- `feature` adds:
  - `协助对接房东`
  - `可带看分佣`

Sync rules:

- On house create/edit, set `is_team_house` from publisher flag
- If user supermarket flag changes, refresh all owned houses'
  `is_team_house`

## 3. House share table

For "my houses map share" page.

Suggested fields:

- `id`
- `token`
- `user_id`
- `expire_at`
- `status`
- `created_at`

Rules:

- Public access by token
- Default valid for 7 days
- Expired token should lead the frontend to home page

## 4. Reward application table

Suggested fields:

- `id`
- `resource_id`
- `apply_user_id`
- `publisher_user_id`
- `apply_user_phone`
- `apply_user_wx_no`
- `publisher_user_phone`
- `publisher_user_wx_no`
- `remark`
- `publisher_confirm_status`
- `audit_status`
- `audit_time`
- `last_operated_at`
- `created_at`
- `updated_at`

Status suggestions:

- `publisher_confirm_status`:
  - `待确认`
  - `已拒绝`
  - `已确认`
- `audit_status`:
  - `未进入审核`
  - `待审核`
  - `审核中`
  - `审核通过待发放`
  - `已发放`
  - `未通过`

## 5. Contact quota records

Recommended minimum design:

- recharge record table
- usage record table

Or simplified:

- one ledger table
- one current balance field on user

Purpose:

- admin can add quota
- landlord-house contact view consumes quota
- admin can query remaining and used quota

## 6. Batch upload record table

Recommended for:

- upload batch traceability
- import result feedback
- failure diagnostics

## Mini Program Changes

Route scope: `server/router/center/center.go`

## 1. Map house finding

Change existing interfaces:

- `/center/options`
- `/center/type/options`
- `/center/house/xiaoquAgg`
- `/center/house/xiaoquAggList`
- `/center/house/listByXiaoqu`

New filters:

- `不限`
- `有返佣`
- `房东房源`
- `团队房源`

Rules:

- Normal users:
  - do not see the team-house filter
  - backend also excludes `is_team_house = true`
- Users with supermarket flag:
  - can see team-house filter
  - can query team houses

## 2. House detail

Change `/center/house/view`:

- landlord house does not show exact address
- return commission amount
- remove dependency on favorite button
- prepare data for complaint entry

## 3. Contact phone

Change `/center/house/mobile`:

- normal house: keep current logic
- landlord house:
  - validate current user's remaining quota first
  - if not enough, return clear error
  - if enough, deduct 1
  - record usage ledger
- continue writing click records for reward-history source

## 4. My house list

Change `/center/house/my`:

- return publisher avatar
- return publisher nickname
- return wechat number
- return commission amount
- return on-shelf count and remaining quota
- remove favorite-related page dependency

Profile fields should ultimately come from `sys_users`, after the persistence
problem is fixed.

## 5. My houses map share

Add interfaces:

- create share token
- query shared map data by token

Rules:

- public access without login
- only current on-shelf houses of the owner are returned
- valid for 7 days
- expired share should be treated as invalid and redirect to home page

## 6. Reward application

Add interfaces:

- recent contacted publishers / houses
- submit application
- my application list
- publisher review list
- publisher confirm / reject

Popup requirements:

- publisher wechat number is displayed from user data
- user only fills remark
- publisher wechat number is not manually entered

## Admin Changes

Route scope:

- `server/router/house_resource/enter.go`
- `server/router/system/enter.go`

## 1. User management

Extend admin user create/edit/list:

- supermarket flag
- on-shelf quota
- contact-view quota

When supermarket flag changes:

- refresh all owned houses' `is_team_house`

## 2. House management

Change:

- `/house/list`
- `/house/my`
- `/house/create`
- `/house/edit`

Add support for:

- landlord house type
- commission amount
- team-house display/filter
- avatar/wechat display on "my houses"
- quota display and shelf control

## 3. Reward review

Add review interfaces:

- publisher reviews applications for own houses
- publisher confirm / reject
- admin views confirmed applications
- admin batch review / payout / reject

## 4. Admin reward list

Only show records with publisher already confirmed.

Fields required:

- applicant phone
- applicant wechat number
- apply time
- publisher phone
- publisher wechat number
- audit status
- last operation time

## 5. Data center

Add:

- reward application count
- search form by phone

Statistics based on phone -> house relation:

- house-post count
- view count
- share count
- contact click count
- reward application count

## 6. Contact quota management

Add admin module:

- add quota
- remark
- ledger list
- search by agent phone
- show remaining and used count

## 7. Batch upload

Add admin-side batch upload interfaces:

- upload Excel
- parse and validate
- batch import
- return result summary

Business rules:

- current login user owns the uploaded houses
- before import, current on-shelf houses of this user are set off shelf
- houses in current Excel are set on shelf
- houses not present in current Excel remain off shelf
- reuse images by same user + same address/building/room identity

## Existing Logic That Must Be Changed

## 1. Favorite flow

- "my favorite" should be removed from Mini Program usage
- `/center/favorite/add`
- `/center/favorite/del`
- `/center/favorite/list`

Need evaluation:

- keep for compatibility only
- or fully deprecate later

## 2. Address visibility

- landlord houses cannot expose exact address in detail response

## 3. Contact phone behavior

- cannot always return phone directly
- landlord houses must go through quota validation first

## 4. Filter options

- `FilterOptions`
- `FilterTypeOptions`

Must be extended on backend, not only frontend static enums.

## Scheduled Job

Add a scheduled task:

- scan reward applications where:
  - publisher has already confirmed
  - backend has not processed
  - more than 48 hours have passed
- auto move them to the target "approved pending payout" style status

This should not be computed only at query time.

## Recommended Delivery Order

1. Fix WeChat profile persistence into `sys_users`
2. Add database fields and models
3. Sync new fields into search index
4. Update Mini Program map/detail/contact behavior
5. Implement reward application and review flow
6. Implement admin user/quota/data-center changes
7. Implement Excel batch upload
