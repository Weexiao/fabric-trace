#!/usr/bin/env python3

import mysql.connector
import sys
import yaml

def load_config():
    """从YAML配置文件读取数据库配置"""
    try:
        with open('/home/qikefan/fabric-trace/application/backend/settings/config.yaml', 'r') as f:
            config = yaml.safe_load(f)
        return config['mysql']
    except Exception as e:
        print(f"❌ 读取配置文件失败: {e}")
        sys.exit(1)

def fix_database():
    """修复数据库：添加 dynamic_attributes 列"""

    print("=" * 50)
    print("  Fabric Trace - 数据库修复")
    print("=" * 50)
    print()

    # 读取配置
    db_config = load_config()

    print("数据库配置:")
    print(f"  主机: {db_config['host']}:{db_config['port']}")
    print(f"  数据库: {db_config['db']}")
    print(f"  用户: {db_config['user']}")
    print()

    try:
        # 连接数据库
        print("连接数据库...")
        conn = mysql.connector.connect(
            host=db_config['host'],
            port=db_config['port'],
            user=db_config['user'],
            password=db_config['password'],
            database=db_config['db']
        )
        cursor = conn.cursor()
        print("✅ 连接成功")
        print()

        # 检查是否已存在该列
        print("检查表结构...")
        cursor.execute("""
            SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS
            WHERE TABLE_NAME='users' AND COLUMN_NAME='dynamic_attributes'
        """)

        if cursor.fetchone():
            print("ℹ️  列 'dynamic_attributes' 已存在，跳过创建")
        else:
            print("创建列 'dynamic_attributes'...")
            cursor.execute("""
                ALTER TABLE users
                ADD COLUMN dynamic_attributes JSON
            """)
            conn.commit()
            print("✅ 列创建成功")

        print()
        print("验证表结构:")
        print("-" * 50)
        cursor.execute("DESCRIBE users")
        columns = cursor.fetchall()
        for col in columns:
            print(f"  {col[0]:25} | {col[1]:20} | {col[2]}")
        print("-" * 50)
        print()

        # 显示数据
        print("当前用户数据:")
        cursor.execute("SELECT COUNT(*) FROM users")
        count = cursor.fetchone()[0]
        print(f"  总用户数: {count}")
        print()

        cursor.close()
        conn.close()

        print("=" * 50)
        print("✅ 数据库修复成功！")
        print("=" * 50)
        print()
        print("现在可以重新尝试注册用户")

    except mysql.connector.Error as err:
        print(f"❌ 数据库错误: {err}")
        sys.exit(1)
    except Exception as e:
        print(f"❌ 错误: {e}")
        sys.exit(1)

if __name__ == '__main__':
    fix_database()

