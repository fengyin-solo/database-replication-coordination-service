# DataSync

数据同步 CDC 管理服务，纯 Go 标准库实现。

服务还包含复制连接器、变更流消费、批次事务、schema 注册、重试调度和游标聚合组件，用于协调跨数据库的增量复制生命周期。

## 实体

- SyncTask（同步任务）
- Source（数据源）
- Target（数据目标）
- ChangeRecord（变更记录）
- Checkpoint（同步断点）

## 运行

```bash
cd origin
go run ./cmd/server
```

## 接口

- POST   /api/tasks          创建任务
- GET    /api/tasks          任务列表
- GET    /api/tasks/{id}     获取任务
- PUT    /api/tasks/{id}     更新任务
- DELETE /api/tasks/{id}     删除任务
- POST   /api/tasks/{id}/start   启动任务
- POST   /api/tasks/{id}/pause   暂停任务
- POST   /api/tasks/{id}/fail    标记失败
- POST   /api/sources        创建源
- GET    /api/sources        源列表
- GET    /api/sources/{id}   获取源
- PUT    /api/sources/{id}   更新源
- DELETE /api/sources/{id}   删除源
- POST   /api/targets        创建目标
- GET    /api/targets        目标列表
- GET    /api/targets/{id}   获取目标
- PUT    /api/targets/{id}   更新目标
- DELETE /api/targets/{id}   删除目标
- POST   /api/records        追加变更记录
- GET    /api/records        变更记录列表
- GET    /api/records/{id}   获取变更记录
- POST   /api/checkpoints    创建断点
- GET    /api/checkpoints    断点列表
- GET    /api/checkpoints/{task_id}  查询断点
- PUT    /api/checkpoints/{task_id}  推进断点
- POST   /api/run            运行同步任务
- POST   /api/record-change  记录变更并推进断点
- GET    /api/stats/overview 统计概览
- GET    /api/stats/tasks-top Top N 任务

## 测试

```bash
go test ./...
```
