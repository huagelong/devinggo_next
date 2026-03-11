---
name: crud-development
description: "使用场景：在 DevingGo-Light 项目中开发新的 CRUD 功能模块。当用户想要为新数据库表创建增删改查接口、或需要手动编写 CRUD 代码（含 API/Controller/Logic/Model 等各层）时使用。包含完整工作流：生成代码 → 完善业务逻辑 → 注册路由 → 运行 make service。"
---

# CRUD 开发工作流

## 完整流程

```bash
# 1. 建表
go run main.go migrate:create -name add_{table}
go run main.go migrate:up
make dao

# 2. 生成代码骨架（推荐）
make gen-crud table={module}_{resource}

# 3. 生成 Service 接口（每次修改 Logic 后必须执行）
make service
```

生成器自动创建：`api/`、`model/req/`、`model/res/`、`controller/`、`logic/` 五层文件。

## 手动编写模板

### API 定义（`api/{module}/{resource}.go`）

```go
type Index{Resource}Req struct {
    g.Meta `path:"/{resource}/index" method:"get" tags:"{分组}" summary:"{中文名}列表" x-permission:"{module}:{resource}:index"`
    model.AuthorHeader
    model.PageListReq
    req.{PascalResource}Search
}
type Index{Resource}Res struct {
    g.Meta `mime:"application/json"`
    page.PageRes
    Items []res.{PascalResource} `json:"items"`
}

type Read{Resource}Req struct {
    g.Meta `path:"/{resource}/read" method:"get" tags:"{分组}" summary:"获取{中文名}详情" x-permission:"{module}:{resource}:read"`
    model.AuthorHeader
    Id int64 `json:"id" v:"required#id不能为空"`
}
type Read{Resource}Res struct {
    g.Meta `mime:"application/json"`
    res.{PascalResource}
}

type Create{Resource}Req struct {
    g.Meta `path:"/{resource}/create" method:"post" tags:"{分组}" summary:"创建{中文名}" x-permission:"{module}:{resource}:create"`
    model.AuthorHeader
    req.{PascalResource}Create
}
type Create{Resource}Res struct{ g.Meta `mime:"application/json"` }

type Update{Resource}Req struct {
    g.Meta `path:"/{resource}/update" method:"put" tags:"{分组}" summary:"更新{中文名}" x-permission:"{module}:{resource}:update"`
    model.AuthorHeader
    Id int64 `json:"id" v:"required#id不能为空"`
    req.{PascalResource}Update
}
type Update{Resource}Res struct{ g.Meta `mime:"application/json"` }

type Delete{Resource}Req struct {
    g.Meta `path:"/{resource}/delete" method:"delete" tags:"{分组}" summary:"删除{中文名}" x-permission:"{module}:{resource}:delete"`
    model.AuthorHeader
    Ids []int64 `json:"ids" v:"required#ids不能为空"`
}
type Delete{Resource}Res struct{ g.Meta `mime:"application/json"` }
```

### Model（`model/req/{table}.go` / `model/res/{table}.go`）

```go
// req
type {PascalResource}Search struct {
    Name   string `json:"name"`
    Status int    `json:"status"`
}
type {PascalResource}Create struct {
    Name   string `json:"name"   v:"required#名称不能为空"`
    Status int    `json:"status"`
    Remark string `json:"remark"`
}
type {PascalResource}Update = {PascalResource}Create

// res
type {PascalResource} struct {
    Id        int64  `json:"id"`
    Name      string `json:"name"`
    Status    int    `json:"status"`
    Remark    string `json:"remark"`
    CreatedAt string `json:"createdAt"`
}
```

### Controller（`controller/{module}/{resource}.go`）

```go
type {resource}Controller struct{ base.BaseController }
var {PascalResource}Controller = {resource}Controller{}

func (c *{resource}Controller) Index(ctx context.Context, in *api.Index{Resource}Req) (out *api.Index{Resource}Res, err error) {
    out = &api.Index{Resource}Res{}
    out.PageRes, out.Items, err = service.{Module}{PascalResource}().GetPageList(ctx, in)
    return
}
// Read / Create / Update / Delete 同理，调对应 service 方法
```

### Logic（`logic/{module}/{table}.go`）

```go
func init() { service.Register{Module}{PascalResource}(New{Module}{PascalResource}()) }

type s{Module}{PascalResource} struct{ base.BaseService }
func New{Module}{PascalResource}() *s{Module}{PascalResource} { return &s{Module}{PascalResource}{} }

func (s *s{Module}{PascalResource}) Model(ctx context.Context) *gdb.Model {
    return dao.{PascalTable}.Ctx(ctx).Hook(hook.Default()).Cache(orm.SetCacheOption(ctx))
}

func (s *s{Module}{PascalResource}) GetPageList(ctx context.Context, req *api.Index{Resource}Req) (pageRes interface{}, items []res.{PascalResource}, err error) {
    var list []*res.{PascalResource}
    m := s.Model(ctx)
    if req.Name != "" { m = m.WhereLike(dao.{PascalTable}.Columns().Name, "%"+req.Name+"%") }
    pageRes, err = orm.NewQuery(m).WithPageListReq(req).ScanAndCount(&list, nil)
    // list → items 转换...
    return
}

func (s *s{Module}{PascalResource}) GetById(ctx context.Context, id int64, out *res.{PascalResource}) error {
    return s.Model(ctx).Where(dao.{PascalTable}.Columns().Id, id).Scan(out)
}

func (s *s{Module}{PascalResource}) Create(ctx context.Context, req *api.Create{Resource}Req) error {
    _, err := s.Model(ctx).Data(do.{PascalTable}{Name: req.Name, Status: req.Status, Remark: req.Remark}).Insert()
    return err
}

func (s *s{Module}{PascalResource}) Update(ctx context.Context, req *api.Update{Resource}Req) error {
    _, err := s.Model(ctx).Data(do.{PascalTable}{Name: req.Name, Status: req.Status, Remark: req.Remark}).
        Where(dao.{PascalTable}.Columns().Id, req.Id).Update()
    return err
}

func (s *s{Module}{PascalResource}) Delete(ctx context.Context, ids []int64) error {
    _, err := s.Model(ctx).WhereIn(dao.{PascalTable}.Columns().Id, ids).Delete()
    return err
}
```

### 注册路由（`router/{module}/router.go`）

```go
group.Group("/{module}", func(group *ghttp.RouterGroup) {
    group.Bind(
        {module}.{PascalResource}Controller,
    ).Middleware(service.Middleware().AdminAuth)
})
```

## 注意事项

- 禁止手动修改 `internal/dao/internal/`、`internal/model/do|entity/`、`modules/*/service/`
- `do.XxxTable{}` 中 `nil` 值字段不写入数据库（避免零值覆盖）
- 软删除由 `hook.Default()` 自动处理

