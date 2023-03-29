# Rbac DSL

Common, Concise and Natural DSL for "Subject - Role - Permission - Resource" Management  
基于简洁自然的通用 DSL 提供类 RBAC(Role-based access control) 的 "主体 - 角色 - 权限 - 资源" 的管理功能  
支持配置不同 Model 使用不同 Rbac Engine，默认提供 "DBx / Auth Service" 的引擎实现

Example:
```go
// Definition
For(&Subject{}).Create(&Role{})
For(&Role{}).Create(&Permission{})

// Assignment
Let(&subject).BeA(&role)
Let(&subject).NotBeA(&role)
Let(&subject).NotBeAnyone()

Let(&role).Can(&permission)
Let(&role).Cannot(&permission)
Let(&role).CannotDoAnything()

// Querying
Whether(subject).IsA(role)
Whether(role).Can(permission)
Whether(subject).Can(permission)

// Resource based Permission
Let(&resource).CanBe("get").By(&role)
Whether(&resource).CanBe("get").By(role)
Whether(&resource).CanBe("get").By(subject)
```

## Features

1. 简洁的 DSL
2. 能够灵活替换不同的 Rbac 实现引擎（默认提供 `dbx` 和 auth service 引擎）
3. 支持 temporary role assign & query
4. 支持 Assignment Callbacks

## Concepts and Overview

### In one word:
```
- role has permissions
- subject has the roles
> subject has the permissions through the roles.
```

### Definition and Uniqueness of nouns

1. Subject
    - Someone who can be assigned roles, and who has permissions through the assigned roles.
    - See wiki [RBAC](https://en.wikipedia.org/wiki/Role-based_access_control)
2. Role
    - A job function that groups a series of permissions according to a certain dimension.
    - Also see wiki [RBAC](https://en.wikipedia.org/wiki/Role-based_access_control)
    - Uniquely identified by `name`
3. Permission
    - An action, or an approval of a mode of access to a resource
    - Also see wiki [RBAC](https://en.wikipedia.org/wiki/Role-based_access_control)
    - Uniquely identified by `action( + resource)`
4. Resource
    - Polymorphic association with permissions

### Three steps to use this lib

1. Querying
    - Query if the given role is assigned to the subject
    - Query if the given permission is assigned to the role / subject's roles
2. Assignment
    - assign the given role to the subject
    - assign the given permission to the role
3. Definition
    - the role or permission you want to assign **MUST** be defined before
    - option `AutoDefine` (before assignment) you may need in some cases

**Definition => Assignment => Querying**

## Initialize

1. 初始化 `rbac`
    ```go
    rbac.Init(map[string]rbac.Config{
        "SubjectModel": {
            Engine:               rbac_dbx.Engine{},
            Role:                 "RoleModel",
            Permission:           "PermissionModel",
            AutoDefineRole:       false,
            AutoDefinePermission: false,
        },
       // Auth Engine 的初始化
       "AuthSubject": {
            Engine:               rbac_auth.Engine{},
            Role:                 "AuthRole",
            Permission:           "AuthPermission",
            AutoDefineRole:       false,
            AutoDefinePermission: false,
        },
    })
    ```
2. 模型声明：
    ```go
    type Subject struct {
       rbac.ActAsSubject
    }
    type Role struct {
       rbac.ActAsRole
    }
    type Permission struct {
       rbac.ActAsPermission
    }
    ```

## Definition -- `For().Create()`

如果没有设置 `AutoDefine`，那么如果 role / permission 不存在，
将会在 assign 的时候报错，或者在 query 的时候返回 `false`  
对于 dbx 和 Auth 引擎而言，Definition 就是 create 的动作

使用 `For()` 函数表明这是一个 definition 动作，并发起调用链  
以下演示使用 dbx Engine：
```go
// Role Definition
// 释义：对 User 群体声明一个新的角色 admin
For(&User{}).Create(&Role{Name: "admin"})
// or
For(&User{}).CreateRole(&Role{Name: "admin"})

// Permission Definition
// 释义：对 Role 模型声明一个新的权限 read_data
For(&Role{}).Create(&Permission{Action: "read_data"})
// or
For(&Role{}).CreatePermission(&Permission{Action: "read_data"})
```

## Assignment -- `Let().*`

使用 `Let().*` 祈使句式发起 assign  
其内部流程是：
1. find subject / role / permission
2. 判断是否已经 assign
3. 进行 assign

以下演示使用 dbx Engine：
```go
// Role Assignment
// 释义：赋予 tom 为 admin 角色
Let(&tom).BeA(&admin)
// alias
Let(&tom).BecomeA(&admin)

// Permission Assignment
// 释义：赋予 admin 一个 read_data 的权限
Let(&admin).Can(&read_data)

// Resource based Permission Assignment
// 释义：赋予 admin 能够 read 资源 book 的权限
Let(&book).CanBe("read").By(&admin)

// Cancel Assignment 取消赋予
Let(&tom).NotBeA(&admin)
Let(&admin).Cannot(&read_data)
Let(&book).CannotBe("read").By(&admin)

// Clear Assignment
Let(&tom).NotBeAnyone()
Let(&admin).CannotDoAnything()
Let(&book).CanBe("read").ByNobody()

// Batch Assignment 批量赋予（或者取消赋予）
Let(&tom).BecomeA([]interface{}{&admin, &guest})
Let(&admin).Can([]interface{}{&sing, &jump, &rap, &basketball})
Let(&book).CanBe("read", "write").By(&admin)
```

### Advance

1. Resource based Permission
    ```go
    // 要是用 Resource 赋值语法，Resource Model 需要：
    type Resource struct {
       rbac.ActAsResource
    }
    func (r Resource) NewPermissionForIt(action string, args ...interface{}) rbac.Permission {
       // 返回一个 permission instance，例如：
       return Permission{Action: action, ResourceID: r.ID}
    }
    ```
2. 临时角色赋予（注意：角色必须已 defined）
    ```go
    // 1. 首先，你的 Subject Model 必须实现了这个接口：
    type TemporaryRoleAssignable interface {
    	TemporaryRoleAssign(foundRole interface{}) error
    	TemporaryRoleRemove(foundRole interface{}) error
    	TemporaryRoleClear() error
 	
    	GetTemporaryRoles() []interface{}
    	IsTemporaryRole(foundRole interface{}) bool
    }
    // 2. Usage
    Let(&tom).BeATemp(admin)
    ```
    所谓临时，即不持久化，角色关系只存在于这个 instance 的生命周期内。[实现例子](https://)
3. 替换性 assign：以给予的 roles / permissions 替换当前
    ```go
    Let(&tom).OnlyBeA([]interface{}{&admin})
    Let(&admin).CanOnly([]interface{}{&read_data})
    Let(&book).CanOnlyBe("read").By(&admin)
    ```

## Querying -- `Whether().*`

使用 `Whether` 疑问句式发起查询，结果是 bool 类型，非常简单：
```go
// Role Querying
Whether(tom).IsA(admin)

// Permission Querying
Whether(admin).Can(read_data)
Whether(tom).Can(read_data)

// Resource based Querying
Whether(book).CanBe("read").By(admin)
Whether(book).CanBe("read").By(tom)
```

## Assignment Callbacks

实现了下述 interface 的 Subject / Role 将会自动调用对应回调：
```go
type AssignCallbacks interface {
	AfterAssign(objs []interface{}, args ...interface{})
	AfterCancelAssign(objs []interface{}, args ...interface{})
	AfterClearAssign(args ...interface{})
}
```

## Engines providing default

### `rbac_dbx`

使用注意几点：
1. Definition 和 Assignment 时要用指针，Querying 可不用
2. 支持自动 find，例如：
    ```go
    Let(&User{Name: "Tom"}).BeA(&Role{Name: "admin"})
    ```

### `rbac_auth`

1. 不需要自己写 model
2. 目前只做了 Query
3. 还没有集成 Auth SDK，即你需要自己 request 并拼接好一个 `rbac_auth.Subject{}`
3. Example
    ```go
    Whether(subject).IsA(rbac_auth.Role{Name: "wm240_admin"})
    Whether(subject).IsA("wm240_admin")
    Whether(rbac_auth.Resource{Name: "wm240"}).CanBe("update").By(subject)
    ```
