import os, sys, json, time
from playwright.sync_api import sync_playwright
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
sys.path.append(project_root)
import random
from app import CustomPage

tracking_link = "{tracking_link}"
uuid = "{uuid}"

uuid = sys.argv[1]
tracking_link = sys.argv[2]
ws_url = sys.argv[3]
data = json.loads(sys.argv[4]) if len(sys.argv) > 4 else {}

p = sync_playwright()
u = p.start()
browser = u.chromium.connect_over_cdp(ws_url)
page_org = browser.contexts[0].pages[0] if browser.contexts and browser.contexts[0].pages else browser.new_context().new_page()
page = CustomPage(page_org, data)

print(data)
try:
    page.goto(tracking_link)
except Exception as e:
    print(e)

{scripts_code}
time.sleep(15)
print("200")
u.stop()

