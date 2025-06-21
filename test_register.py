#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
用户注册测试脚本
测试用户名：test
测试密码：123456
"""

import requests
import json

# 服务器配置
BASE_URL = "http://localhost:8000"
REGISTER_URL = f"{BASE_URL}/api/auth/register"

# 测试数据
test_user = {
    "username": "test01",
    "password": "123456"
}

def test_register():
    """
    测试用户注册功能
    """
    print("=== 用户注册测试 ===")
    print(f"请求URL: {REGISTER_URL}")
    print(f"测试数据: {json.dumps(test_user, ensure_ascii=False, indent=2)}")
    
    try:
        # 发送注册请求
        response = requests.post(
            REGISTER_URL,
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
        if response.status_code == 201:
            print("\n✅ 注册测试成功！")
            return True
        elif response.status_code == 409:
            print("\n⚠️  用户已存在，这是正常情况")
            return True
        else:
            print(f"\n❌ 注册测试失败，状态码: {response.status_code}")
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

def test_invalid_data():
    """
    测试无效数据的注册请求
    """
    print("\n=== 无效数据测试 ===")
    
    # 测试用例
    test_cases = [
        {"name": "空用户名", "data": {"username": "", "password": "123456"}},
        {"name": "短密码", "data": {"username": "testuser", "password": "123"}},
        {"name": "缺少用户名", "data": {"password": "123456"}},
        {"name": "缺少密码", "data": {"username": "testuser"}},
    ]
    
    for test_case in test_cases:
        print(f"\n测试: {test_case['name']}")
        try:
            response = requests.post(
                REGISTER_URL,
                json=test_case['data'],
                headers={"Content-Type": "application/json"},
                timeout=5
            )
            print(f"状态码: {response.status_code}")
            if response.status_code == 400:
                print("✅ 正确拒绝了无效数据")
            else:
                print("⚠️  意外的响应状态码")
        except Exception as e:
            print(f"❌ 测试错误: {str(e)}")

if __name__ == "__main__":
    print("开始用户注册测试...")
    print("请确保Go服务器正在运行在端口8000")
    print("-" * 50)
    
    # 执行主要测试
    success = test_register()
    
    # 执行无效数据测试
    test_invalid_data()
    
    print("\n" + "="*50)
    if success:
        print("🎉 注册功能测试完成！")
    else:
        print("💥 注册功能测试失败！")
    print("="*50)