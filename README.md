# 今日头条自动发布工具

这是一个基于Selenium的自动化工具，用于自动登录今日头条并发布草稿文章。

## 功能特性

- 自动登录今日头条后台
- 支持无头模式运行
- 自动处理草稿发布流程
- 随机间隔发布，模拟人工操作
- 支持代理设置
- 自动处理地理位置权限

## 环境要求

- Python 3.6+
- Chrome浏览器
- ChromeDriver
- Selenium库

## 安装步骤

1. 安装Python 3.6或更高版本
2. 安装依赖库：
   ```bash
   pip install selenium
   ```
3. 下载与Chrome版本匹配的ChromeDriver
4. 将ChromeDriver添加到系统PATH

## 使用方法

1. 修改autoPublish.py中的账号信息：
   ```python
   self.username = "your_username"  # 今日头条账号
   self.password = "your_password"  # 今日头条密码
   ```

2. 运行脚本：
   ```bash
   python autoPublish.py
   ```

3. 程序会自动：
   - 登录今日头条后台
   - 进入草稿箱
   - 逐个发布草稿
   - 随机间隔3-10分钟发布一篇

## 注意事项

- 请勿在代码中明文存储账号信息，建议使用环境变量
- 确保ChromeDriver版本与Chrome浏览器匹配
- 发布间隔时间可自行调整
- 如需使用代理，请取消注释相关代码

## 贡献

欢迎提交Pull Request改进本项目

## 许可证

MIT License
