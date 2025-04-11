# 二开扩展表迁移说明

这个目录包含了二开扩展的数据库迁移文件，用于管理二开项目中添加的新表。
## 如何使用

### 使用 Flask 命令（推荐）

我们提供了与原始项目类似的 Flask 命令来管理二开扩展表的迁移：

```bash
# 升级数据库到最新版本
flask extend_db upgrade

# 回滚数据库到指定版本
flask extend_db downgrade --revision 版本号

# 查看当前数据库版本
flask extend_db current

# 查看迁移历史
flask extend_db history

# 查看最新的迁移版本
flask extend_db heads
```

### 检查当前数据库版本

要查看当前数据库的版本，请运行：

```bash
flask extend_db current
```

### 降级数据库

如果需要回滚迁移，可以使用：

```bash
flask extend_db downgrade --revision 版本号
```

例如，回滚到 `001_recommended_list_sorted`：

```bash
flask extend_db downgrade --revision 001_recommended_list_sorted
```

## 注意事项

1. 这些迁移文件仅适用于二开扩展的表，不会影响原有系统表。
2. 所有迁移文件都已添加表存在性检查，可以多次运行而不会出错。
3. 如果需要添加新的迁移文件，请确保正确设置前置版本（Revises）值。
4. 使用 Flask 命令时确保在项目 `api` 下运行，而不是根目录下。 