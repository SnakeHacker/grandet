__author__ = 'mickey'

from selenium import webdriver
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By

import tushare as ts
import psycopg2
import time
import yaml
import argparse
import socket
import os
socket.setdefaulttimeout(120)


def get_concept(stocks, conn, debug):
    desired_capabilities = DesiredCapabilities.CHROME
    desired_capabilities["pageLoadStrategy"] = "none"
    chrome_options = Options()
    if not debug:
        chrome_options.add_argument('--headless')
        chrome_options.add_argument('--no-sandbox')
        chrome_options.add_argument('--disbale-dev-shm-usage')
    browser = webdriver.Chrome(options=chrome_options)
    browser.implicitly_wait(30)

    url = "http://www.iwencai.com/stockpick/search?\
        tid=stockpick&qs=sl_box_main_ths&w="

    cur = conn.cursor()

    for stock in stocks:
        print(stock)
        browser.get('{}{}'.format(url, stock))
        time.sleep(5)
        locator = (By.XPATH, '//*[@class="em alignCenter split"]')
        WebDriverWait(browser, 120, 1).until(
            EC.presence_of_element_located(locator))

        try:
            moreBtn = browser.find_element_by_xpath('//*[@class="em alignCenter split"]\
                /a[contains(@class, "ml5 moreSplit fr")]')
            moreBtn.click()
            time.sleep(2)
        except Exception:
            # no moreBtn
            pass

        concepts = browser.find_elements_by_xpath('//*[@class="em alignCenter split"]/\
            span/a')
        for concept in concepts:
            print(concept.text)
            cur.execute("INSERT INTO concept_details \
                (ts_code, concept_name) VALUES(%s, %s)", (stock, concept.text))
            conn.commit()
        time.sleep(1)
    cur.close()
    return


def connect_db(conf):
    conn = psycopg2.connect(database=conf['db']['database'],
                            user=conf['db']['username'],
                            password=conf['db']['password'],
                            host=conf['db']['host'],
                            port=conf['db']['port'])

    return conn


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--debug', type=bool, default=False)
    parser.add_argument('-conf', type=str)
    args = parser.parse_args()

    yaml_path = os.path.join(os.getcwd(), args.conf)
    f = open(yaml_path)
    conf = yaml.load(f, Loader=yaml.FullLoader)

    pro = ts.pro_api(os.getenv("TUSHARE_TOKEN"))
    stocks = pro.stock_basic()

    tsCodes = stocks['ts_code'].to_list()
    print(tsCodes)

    conn = connect_db(conf)
    get_concept(tsCodes, conn, args.debug)
    conn.close()
