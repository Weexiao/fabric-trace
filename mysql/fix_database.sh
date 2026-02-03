#!/bin/bash

# 数据库修复脚本 - 添加 dynamic_attributes 列

echo "=========================================="
echo "  Fabric Trace - 数据库修复"
echo "=========================================="
echo ""

# 读取配置
CONFIG_FILE="/home/qikefan/fabric-trace/application/backend/settings/config.yaml"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ 配置文件不存在: $CONFIG_FILE"
    exit 1
fi

# 提取数据库配置
DB_USER=$(grep "user:" "$CONFIG_FILE" | head -1 | awk '{print $2}' | tr -d '"')
DB_PASS=$(grep "password:" "$CONFIG_FILE" | head -1 | awk '{print $2}' | tr -d '"')
DB_HOST=$(grep "host:" "$CONFIG_FILE" | grep -v "#" | head -1 | awk '{print $2}' | tr -d '"')
DB_PORT=$(grep "port:" "$CONFIG_FILE" | head -1 | awk '{print $2}' | tr -d '"')
DB_NAME=$(grep "db:" "$CONFIG_FILE" | head -1 | awk '{print $2}' | tr -d '"')

echo "数据库配置:"
echo "  主机: $DB_HOST:$DB_PORT"
echo "  数据库: $DB_NAME"
echo "  用户: $DB_USER"
echo ""

# 创建临时SQL文件
TEMP_SQL=$(mktemp)
cat > "$TEMP_SQL" << 'SQL_END'
-- 添加 dynamic_attributes 列
ALTER TABLE users ADD COLUMN IF NOT EXISTS dynamic_attributes JSON DEFAULT '{}';

-- 显示表结构
SHOW FULL COLUMNS FROM users;
SQL_END

echo "执行SQL迁移..."
echo ""

# 执行SQL
if [ -z "$DB_PASS" ]; then
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" < "$TEMP_SQL"
else
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" < "$TEMP_SQL"
fi

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 数据库修复成功！"
    echo ""
    echo "现在可以重新尝试注册用户"
else
    echo ""
    echo "❌ 数据库修复失败！"
    echo ""
    echo "请手动执行以下命令:"
    echo "  mysql -u $DB_USER -p $DB_NAME"
    echo "  ALTER TABLE users ADD COLUMN IF NOT EXISTS dynamic_attributes JSON DEFAULT '{}';"
    exit 1
fi

# 清理临时文件
rm -f "$TEMP_SQL"

echo ""
echo "=========================================="

