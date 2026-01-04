# abd开发规范

## 数据类型
### ID
与前端交互时，ID类型加string，如
ToUserID   int64  `json:"to_user_id,string" binding:"required"`
### int的选择
只有极少的数值用int32,如
Status   int32
其他都用int64，如
Seq int64


## 命名规范
### JSON 字段：
全小写，蛇形命名法 (snake_case)。
create_time, user_id, group_info

### Go 结构体/变量：
大驼峰命名法 (PascalCase)。
CreateTime, UserID, GroupInfo
SendMessageReq
SendMessageResp

### Go 专有缩写：
遇到 ID, API, URL 等缩写词，保持全大写。
UserID, APIKey