#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
用户登录测试脚本
测试用户名：test
测试密码：123456
"""

import requests
import json

# 服务器配置
BASE_URL = "http://localhost:8000"
LOGIN_URL = f"{BASE_URL}/api/auth/login"
HEALTH_URL = f"{BASE_URL}/api/health"

# 测试数据
test_user = {
    "username": "test",
    "password": "123456"
}

def test_health_check():
    """
    测试服务器健康状态
    """
    print("=== 服务器健康检查 ===")
    try:
        response = requests.get(HEALTH_URL, timeout=5)
        print(f"健康检查状态码: {response.status_code}")
        if response.status_code == 200:
            print("✅ 服务器运行正常")
            return True
        else:
            print("❌ 服务器状态异常")
            return False
    except Exception as e:
        print(f"❌ 无法连接到服务器: {str(e)}")
        return False

def test_login():
    """
    测试用户登录功能
    """
    print("\n=== 用户登录测试 ===")
    print(f"请求URL: {LOGIN_URL}")
    print(f"测试数据: {json.dumps(test_user, ensure_ascii=False, indent=2)}")
    
    try:
        # 发送登录请求
        response = requests.post(
            LOGIN_URL,
            json=test_user,
            headers={"Content-Type": "application/json"},
            timeout=10
        )
        
        print(f"\n响应状态码: {response.status_code}")
        print(f"响应头: {dict(response.headers)}")
        
        # 解析响应
        try:
            response_data = response.json()
            print(f"响应数据: {json.dumps(response_data, ensure_ascii=False, indent=2)}")
        except json.JSONDecodeError:
            print(f"响应内容 (非JSON): {response.text}")
        
        # 判断测试结果
        if response.status_code == 200:
            print("\n✅ 登录测试成功！")
            return True
        elif response.status_code == 401:
            print("\n❌ 登录失败：用户名或密码错误")
            print("💡 提示：请先运行 test_register.py 创建测试用户")
            return False
        elif response.status_code == 403:
            print("\n❌ 登录失败：用户账户未激活")
            return False
        else:
            print(f"\n❌ 登录测试失败，状态码: {response.status_code}")
            return False
            
    except requests.exceptions.ConnectionError:
        print("\n❌ 连接失败！请确保服务器正在运行在 http://localhost:8000")
        return False
    except requests.exceptions.Timeout:
        print("\n❌ 请求超时！")
        return False
    except Exception as e:
        print(f"\n❌ 测试过程中发生错误: {str(e)}")
        return False

def test_invalid_login():
    """
    测试无效登录数据
    """
    print("\n=== 无效登录测试 ===")
    
    # 测试用例
    test_cases = [
        {"name": "错误密码", "data": {"username": "test", "password": "wrongpassword"}},
        {"name": "不存在的用户", "data": {"username": "nonexistentuser", "password": "123456"}},
        {"name": "空用户名", "data": {"username": "", "password": "123456"}},
        {"name": "空密码", "data": {"username": "test", "password": ""}},
        {"name": "缺少用户名", "data": {"password": "123456"}},
        {"name": "缺少密码", "data": {"username": "test"}},
    ]
    
    for test_case in test_cases:
        print(f"\n测试: {test_case['name']}")
        try:
            response = requests.post(
                LOGIN_URL,
                json=test_case['data'],
                headers={"Content-Type": "application/json"},
                timeout=5
            )
            print(f"状态码: {response.status_code}")
            if response.status_code in [400, 401]:
                print("✅ 正确拒绝了无效登录")
            else:
                print("⚠️  意外的响应状态码")
                try:
                    print(f"响应: {response.json()}")
                except:
                    print(f"响应: {response.text}")
        except Exception as e:
            print(f"❌ 测试错误: {str(e)}")

def test_login_flow():
    """
    完整的登录流程测试
    """
    print("\n=== 完整登录流程测试 ===")
    
    # 1. 健康检查
    if not test_health_check():
        return False
    
    # 2. 尝试登录
    login_success = test_login()
    
    # 3. 测试无效登录
    test_invalid_login()
    
    return login_success

if __name__ == "__main__":
    print("开始用户登录测试...")
    print("请确保Go服务器正在运行在端口8000")
    print("请确保已经运行过 test_register.py 创建了测试用户")
    print("-" * 60)
    
    # 执行完整测试流程
    success = test_login_flow()
    
    print("\n" + "="*60)
    if success:
        print("🎉 登录功能测试完成！")
    else:
        print("💥 登录功能测试失败！")
        print("💡 请检查：")
        print("   1. Go服务器是否正在运行")
        print("   2. 数据库是否正确配置")
        print("   3. 是否已运行 test_register.py 创建测试用户")
    print("="*60)