from logging import log
import random
import time
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.chrome.options import Options


def createHeadLessDriver():
    """初始化无头浏览器"""
    chrome_options = Options()
    # chrome_options.add_argument("--headless")
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--disable-dev-shm-usage")
    chrome_options.add_argument("--disable-gpu")
    chrome_options.add_argument("--window-size=1920,1080")
    
    # 如果需要代理
    # chrome_options.add_argument("--proxy-server=http://your-proxy:port")
    chrome_prefs = {
        "profile.default_content_setting_values.geolocation": 2  # 2 表示拒绝， 1 允许
    }
    chrome_options.add_experimental_option("prefs", chrome_prefs)
    

    return webdriver.Chrome(options=chrome_options)

def createChromeDriver():
    chrome_options = Options()
   
    chrome_prefs = {
        "profile.default_content_setting_values.geolocation": 2  # 2 表示拒绝， 1 允许
    }
    chrome_options.add_experimental_option("prefs", chrome_prefs)
    
    return webdriver.Chrome(options=chrome_options)

class ToutiaoAutoPublisher:
    def __init__(self, driver=createChromeDriver()):
        self.driver = driver
        self.wait = WebDriverWait(self.driver, 20)
        
        # 重要提示：请勿在代码中明文存储账号信息
        self.username = "17355605329"  # 建议从环境变量读取
        self.password = "Aa123456"  # 建议从环境变量读取

    def login(self):
        """处理登录流程"""
        self.driver.get("https://mp.toutiao.com/")
        
        # 等待并选择登录方式（示例为账号密码登录）
        self.wait.until(EC.presence_of_element_located(
            (By.XPATH, "//li[@aria-label='账密登录']")
        ))
        self.driver.find_element(By.XPATH, "//li[@aria-label='账密登录']").click()

        # self.wait.until(EC.element_to_be_clickable(
        #     (By.Name, "//div[contains(text(),'密码登录')]"))).click()
        
        # 输入凭据
        self.wait.until(EC.presence_of_element_located(
            (By.NAME, "normal-input"))).send_keys(self.username)
        self.wait.until(EC.presence_of_element_located(
            (By.NAME, "button-input"))).send_keys(self.password)
        
       # 勾选协议框 <span class="web-login-confirm-info__checkbox" role="checkbox" aria-checked="false" tabindex="0" aria-label="协议勾选框"></span>
        self.driver.find_element(By.CLASS_NAME, "web-login-confirm-info__checkbox").click() 
       
        # 提交登录（注意验证码处理）
        self.driver.find_element(By.CLASS_NAME, "web-login-button").click()
        
        # 等待登录成功
        self.wait.until(EC.url_contains("profile_v4/index"))

    # def publish_drafts(self):
    #     """处理草稿发布流程"""
    #     self.driver.get("https://mp.toutiao.com/profile_v4/manage/draft")
    #     self.wait.until(EC.presence_of_element_located(
    #         (By.XPATH, "//div[contains(@class, 'list')]")
    #     ))

  
    #     while True:
    #         # 获取所有草稿项
    #         drafts = self.wait.until(EC.presence_of_all_elements_located(
    #             (By.XPATH, "//div[contains(@class, 'draft-item')]")
    #         ))

    #         if not drafts:
    #             log.info("没有更多草稿了")
    #             # sleep 1hour
    #             time.sleep(3600)
    #             continue
            
    #         for draft in drafts:
    #             try:
    #                 # 点击编辑按钮
    #                 edit_button = draft.find_element(By.XPATH, ".//a[contains(@class, 'op-button') and text()='编辑']")
    #                 edit_button.click()

    #                 # # 等待编辑页面加载
    #                 self.wait.until(EC.presence_of_element_located(
    #                     (By.XPATH, "//div[contains(@class, 'publish-box')]")))

    #                 # 点击发布按钮
    #                 publish_button = self.wait.until(EC.element_to_be_clickable(
    #                     (By.CLASS_NAME, "publish-content")))
    #                 publish_button.click()

    #                 log.info(f"{time.strftime('%Y-%m-%d %H:%M:%S')} 发布成功")
                    
    #                 # 随机间隔等待 5 ~ 10 min
    #                 snooze = 150 + random.random() * 400
    #                 log.info(f"等待 {snooze} 秒")
    #                 time.sleep(snooze)

    #             except Exception as e:
    #                 print(f"处理草稿时出错: {str(e)}")
    #                 continue  # 继续处理下一个草稿
    def publish_drafts(self):
        """处理草稿发布流程"""
        self.driver.get("https://mp.toutiao.com/profile_v4/manage/draft")
        self.wait.until(EC.presence_of_element_located(
            (By.XPATH, "//div[contains(@class, 'list')]")
        ))

        while True:
            # 获取所有草稿项
            drafts = self.wait.until(EC.presence_of_all_elements_located(
                (By.XPATH, "//div[contains(@class, 'draft-item')]")
            ))

            if not drafts:
                log.info("没有更多草稿了")
                # sleep 1hour
                time.sleep(3600)
                continue
            
            for draft in drafts:
                try:
                    # 获取当前窗口句柄
                    current_window = self.driver.current_window_handle
                    editing_window = None

                    # 点击编辑按钮
                    edit_button = draft.find_element(By.XPATH, ".//a[contains(@class, 'op-button') and text()='编辑']")
                    edit_button.click()

                  
                    # 等待新标签页加载完成
                    self.wait.until(EC.number_of_windows_to_be(2))  # 确保有两个窗口
                    
                    # 切换到新标签页
                    for handle in self.driver.window_handles:
                        if handle != current_window:
                            self.driver.switch_to.window(handle)
                            edgeing_window = handle
                            break

                    # 等待编辑页面加载
                    self.wait.until(EC.presence_of_element_located(
                        (By.XPATH, "//div[contains(@class, 'publish-box')]")))

                   
                    # 使用XPath选择按钮
                    # 等待并获取按钮
                    publish_button = self.wait.until(EC.presence_of_element_located(
                        (By.XPATH, "//button[contains(@class, 'publish-content')]")
                    ))

                    # 滚动到按钮
                    self.driver.execute_script("arguments[0].scrollIntoView(true);", publish_button)

                    # 等待遮挡元素消失
                    self.wait.until(EC.invisibility_of_element_located(
                        (By.CLASS_NAME, "byte-drawer-scroll")
                    ))

                    # 等待按钮可点击
                    self.wait.until(EC.element_to_be_clickable(
                        (By.XPATH, "//button[contains(@class, 'publish-content')]")
                    ))

                    # 点击按钮
                    publish_button.click()
                    time.sleep(3)
                    
                    # 切换回原窗口（如果需要继续处理草稿）
                    self.driver.close()
                    self.driver.switch_to.window(current_window)

                    # 随机间隔等待 3 ~ 10 min
                    snooze = 150 + random.random() * 400
                    print(f"等待 {snooze} 秒")
                    time.sleep(snooze)

                   


                except Exception as e:
                    print(f"处理草稿时出错: {str(e)}")
                    continue  # 继续处理下一个草稿

    def run(self):
        try:
            self.login()
            self.publish_drafts()
        except Exception as e:
            print(f"执行出错: {str(e)}")
        finally:
            self.driver.quit()

if __name__ == "__main__":
    publisher = ToutiaoAutoPublisher(createHeadLessDriver())
    publisher.run()
