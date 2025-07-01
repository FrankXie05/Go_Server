import os, requests, sys
from playwright.sync_api import sync_playwright, Page
import random
import time
from faker import Faker

class CustomPage:
    def __init__(self, page: Page, data):
        self.page = page
        self.data = data
        self.faker = Faker('en_US')
    def custom_wait(self):
        self.page.wait_for_timeout(1000)
        self.page.wait_for_load_state('domcontentloaded')
    def goto(self, url):
        self.custom_wait()
        self.page.goto(url, wait_until='domcontentloaded', timeout=60000)
    def fill(self, selector, value, **kvargs):
        self.custom_wait()
        self.page.fill(selector=selector, value=value, **kvargs)
    
    def click(self, selector, **kvargs):
        self.custom_wait()
        self.page.click(selector=selector, **kvargs)

    def select_option(self, selector, index, **kvargs):
        self.custom_wait()
        self.page.select_option(selector=selector, index=index, **kvargs)

    def get_by_text(self, text):
        return self.page.get_by_text(text=text)
    
    def screenshot(self, path):
        return self.page.screenshot(path=path)
    
    def human_scroll(self, target_selector=None):
        total_scroll = random.randint(300, 800)  # 模拟滚动的距离
        step = random.randint(30, 80)            # 每次滚动的步长
        delay = random.randint(50, 120)          # 每步停顿的时间

        for y in range(0, total_scroll, step):
            self.page.evaluate(f"window.scrollBy(0, {step})")
            self.page.wait_for_timeout(delay)

        # 滚动完成后稍等一下
        self.page.wait_for_timeout(random.randint(500, 1000))

        if target_selector:
            self.page.evaluate(f'''
                () => {{
                    const el = document.querySelector("{target_selector}");
                    if (el) el.scrollIntoView({{ behavior: 'auto', block: 'center' }});
                }}
            ''')
            self.page.wait_for_timeout(random.randint(500, 1000))

    def fake_mouse_hover(self, selector):
        self.page.evaluate(f"""
            const el = document.querySelector('{selector}');
            if (el) {{
                el.dispatchEvent(new MouseEvent('mouseover', {{ bubbles: true }}));
                el.dispatchEvent(new MouseEvent('mouseenter', {{ bubbles: true }}));
            }}
        """)

    def weighted_sleep(options):
        """
        参数 options 是一个列表，每个元素是 (等待时间, 权重)
        例如：[(0.5, 5), (1, 3), (3, 1)] 表示更偏向于等待 0.5 秒
        """
        times, weights = zip(*options)
        sleep_time = random.choices(times, weights=weights, k=1)[0]
        time.sleep(sleep_time)

    def clear_input(self,x_path):
        """
        模拟使用回退键清空输入框
        - `x_path`: 输入框的 XPath 选择器
        """
        el = self.page.query_selector(f'xpath={x_path}')
        if not el:
            return
        el.click()
        val = el.input_value()
        for _ in range(len(val)):
            self.page.keyboard.press('Backspace')
            time.sleep(random.uniform(0.05, 0.12))  

    def focus_element(self, x_path):
        """
        模拟用户"点进"输入框，使用JS,让一个DOM元素获取焦点
        """
        self.page.evaluate(f"""
            document.evaluate('{x_path}', docment, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null)
            .singleNodeValue.focus();
        """    )

    def type_with_delay(self,x_path, text, min_delay=120, max_delay=180, error_chance=0.05, pause_chance=0.1, allow_typo=True):
        """
        输入文本，模拟人类输入的延迟和错误
        - `x_path`: 输入框的 XPath 选择器
        - `text`: 要输入的文本
        - `min_delay` 和 `max_delay`: 输入字符之间的延迟时间（毫秒）
        - `error_chance`: 输入错误的概率（0到1之间）
        - `pause_chance`: 输入过程中暂停的概率（0到1之间）
        """

        for char in text:
            delay = random.randint(min_delay, max_delay)
            if allow_typo and random.random() < error_chance:
                wrong_char = random.choice('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789')
                self.page.type(x_path, wrong_char, delay=random.randint(80, 150), timeout=60000)
                self.page.keyboard.press('Backspace')  
                time.sleep(random.uniform(0.1, 0.5)) 

            self.page.type(x_path, char, delay=delay)

            if random.random() < pause_chance:  
                time.sleep(random.uniform(0.5, 2))

    def input_username(self):
        """模拟输入用户名的行为"""
        usename_xpath = '//input[@name="username"]'
        CustomPage.weighted_sleep([(3, 6), (6, 3), (10, 1)])
        self.clear_input(usename_xpath)
        username = self.data['username']
        self.type_with_delay(usename_xpath, username, min_delay=150, max_delay=220)

    def input_password(self):
        """模拟输入密码的行为"""
        password_xpath = '//input[@name="password"]'
        CustomPage.weighted_sleep([(4, 5), (8, 3), (15, 2)])
        self.clear_input(password_xpath)
        password = self.data['password']
        self.type_with_delay(password_xpath, password, min_delay=150, max_delay=220, error_chance=0.02)

    def input_city(self):
        """模拟输入城市的行为"""
        city_xpath = '//input[@name="city"]'
        CustomPage.weighted_sleep([(2, 7), (4, 3), (10, 1)])
        self.clear_input(city_xpath)
        city = self.data['city']
        self.type_with_delay(city_xpath, city, min_delay=120, max_delay=180)

    def input_zipcode(self):
        """模拟输入邮政编码的行为"""
        zipcode_xpath = '//input[@name="postal"]'
        CustomPage.weighted_sleep([(1, 8), (3, 3), (12, 1)])
        self.clear_input(zipcode_xpath)
        zipcode = self.data['zip_code']
        self.type_with_delay(zipcode_xpath, zipcode, min_delay=100, max_delay=150, allow_typo=False)

    def input_email(self):
        """模拟输入邮箱的行为"""
        email_xpath = '//input[@name="email"]'  
        CustomPage.weighted_sleep([(2, 6), (5, 3), (12, 1)])
        self.clear_input(email_xpath)
        email = self.data['email']
        self.type_with_delay(email_xpath, email, min_delay=120, max_delay=180)
    
    def input_address(self):
        """模拟输入地址的行为"""
        address_xpath = '//input[@name="address" or @name="street1" or @name="addressform"]'
        address, street = CustomPage.generate_formatted_address(self=self)
        CustomPage.weighted_sleep([(3, 6), (7, 3), (10, 1)])
        self.clear_input(address_xpath)
        self.type_with_delay(address_xpath, address, min_delay=120, max_delay=180)

    def input_first_name(self):
        """模拟输入名的行为"""
        first_name_xpath = '//input[@name="first_name"]' 
        CustomPage.weighted_sleep([(1, 8), (3, 2), (5, 1)]) 
        self.clear_input(first_name_xpath)
        first_name = self.data['username'].split(' ')[0]
        self.type_with_delay(first_name_xpath, first_name, min_delay=120, max_delay=180)

    def input_last_name(self):
        """模拟输入名的行为"""
        last_name_xpath = '//input[@name="last_name"]' 
        CustomPage.weighted_sleep([(1, 8), (3, 2), (5, 1)]) 
        self.clear_input(last_name_xpath)
        last_name = self.data['username'].split(' ')[-1]
        self.type_with_delay(last_name_xpath, last_name, min_delay=120, max_delay=180)

    def input_phone(self):
        """模拟输入电话号码的行为"""
        phone_xpath = '//input[@name="phone"]' 
        CustomPage.weighted_sleep([(1, 7), (3, 2), (10, 1)]) 
        self.clear_input(phone_xpath)
        phone = self.data['phone_number']
        self.type_with_delay(phone_xpath, phone, min_delay=120, max_delay=180, allow_typo=False)
        

    def generate_formatted_address(self):
        zip_code = self.data['zip_code']
        city = self.data['city']
        province = self.data['province']

        street = Faker('en_US').street_address()
        country_pool = ["US", "USA", "United States", "United States of America"]
        country = random.choice(country_pool)

        # 州缩写替代
        state_map = {
            "Alabama": "al",
            "Alaska": "ak",
            "Arizona": "az",
            "Arkansas": "ar",
            "California": "ca",
            "Colorado": "co",
            "Connecticut": "ct",
            "Delaware": "de",
            "District of Columbia": "dc",
            "Florida": "fl",
            "Georgia": "ga",
            "Hawaii": "hi",
            "Idaho": "id",
            "Illinois": "il",
            "Indiana": "in",
            "Iowa": "ia",
            "Kansas": "ks",
            "Kentucky": "ky",
            "Louisiana": "la",
            "Maine": "me",
            "Maryland": "md",
            "Massachusetts": "ma",
            "Michigan": "mi",
            "Minnesota": "mn",
            "Mississippi": "ms",
            "Missouri": "mo",
            "Montana": "mt",
            "Nebraska": "ne",
            "Nevada": "nv",
            "New Hampshire": "nh",
            "New Jersey": "nj",
            "New Mexico": "nm",
            "New York": "ny",
            "North Carolina": "nc",
            "North Dakota": "nd",
            "Ohio": "oh",
            "Oklahoma": "ok",
            "Oregon": "or",
            "Pennsylvania": "pa",
            "Rhode Island": "ri",
            "South Carolina": "sc",
            "South Dakota": "sd",
            "Tennessee": "tn",
            "Texas": "tx",
            "Utah": "ut",
            "Vermont": "vt",
            "Virginia": "va",
            "Washington": "wa",
            "West Virginia": "wv",
            "Wisconsin": "wi",
            "Wyoming": "wy",
            "Washington, D.C.'":"wa"
        }  
        province_format = random.choice([province, state_map.get(province, province)])

        formats = [
            f"{city}, {province_format}, {zip_code}, {country}",
            f"{city}, {province_format}, {country}",
            f"{zip_code}, {city}, {province_format}, {country}",
            f"{city}, {province_format}, {zip_code}",
            f"{city}, {zip_code}, {province_format}, {country}",
            f"{city},{zip_code}, {country}"
        ]

        return random.choice(formats), street

    def stealth_click(self, visual_selector: str, fallback_selector: str = None, delay_ms_range=(100, 250)):
        """伪装人类点击，优先点击最外层。 防止冒泡点击事件的触发"""
        
        try:
            target = self.page.locator(fallback_selector) if fallback_selector else self.page.locator(visual_selector)

            target.wait_for(state='visible', timeout=5000)
            target.scroll_into_view_if_needed()
            time.sleep(random.uniform(0.2, 0.6))

            target.focus()
            delay = random.randint(*delay_ms_range)
            target.click(delay=delay, force=True)
            time.sleep(random.uniform(0.2, 0.6))
        except Exception as e:
            print(f"Error in stealth_click: {e}")
            if fallback_selector:
                self.page.click(fallback_selector, force=True)
            else:
                self.page.click(visual_selector, force=True)