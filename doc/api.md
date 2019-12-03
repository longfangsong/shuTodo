# API Reference

## 模型

### Todo 项

```json
{
  "id":             id,
  "content":        Todo项内容,
  "due":            到期时间（可空，JS标准格式）,
  "estimate_cost":  预计花费时间（可空，如"2h"，返回时单位为ns）
  "type":           类型（可空，'Homework'/'Coding'/'Report'/'Discussion'）
}
```

## web api

- `GET /ping`

  检查服务是否可用，应该直接返回`pong`。

- `GET /todo`

  获得jwt token指示的学生的所有todo项
  
- `POST /todo`、`PUT /todo`
 
  新增jwt token指示的学生的todo项
  
- `DELETE /todo?id=${id}`
  删除id为id的todo项，这个todo项所有者必须是jwt token指示的学生