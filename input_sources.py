import sys
import argparse
import re
import json
from datetime import datetime
import requests


def parsePurchase_smbc(log_file):
    # MEMO: 元号が変わるとバグる
    purchase_list = []
    for line in log_file:
        if not line.startswith('H'):
            continue  # drop the line if it's beginning is not year

        entry = line[:-1].split(',')
        year = int(entry[0][1:3]) + 1988
        date = re.sub(r'^H\d\d', str(year), entry[0])  # 和暦を西暦に
        date = datetime.strptime(date, '%Y.%m.%d')
        desc = entry[3][1:-1]  # eliminate double-quotes from head / tail
        # check income or outgo
        if entry[1] != '':
            price = entry[1]
        else:
            price = '-' + entry[2]
        purchase_list.append((date, desc, int(price), line.replace('\n', '')))
    return purchase_list


def parsePurchase_pitapa(log_file):
    # TODO: プリペイド利用をうまく処理する
    purchase_list = []
    for line in log_file:
        if not line.startswith('20'):
            continue  # drop the line if it's beginning is not year

        entry = line[:-1].split(',')
        date = datetime.strptime(entry[0], '%Y/%m/%d')
        desc = entry[3]
        price = entry[4]
        purchase_list.append((date, desc, int(price), line.replace('\n', '')))
    return purchase_list


def parsePurchase_jpbank(log_file):
    purchase_list = []
    for line in log_file:
        if not line.startswith('20'):
            continue  # drop the line if it's beginning is not year

        entry = line[:-1].split(',')
        date = datetime.strptime(entry[0], '%Y/%m/%d')
        desc = entry[1]
        if entry[6] != '':
            desc += ' ' + entry[6]
        price = entry[2]
        purchase_list.append((date, desc, int(price), line.replace('\n', '')))
    return purchase_list


def parsePurchase_visa(log_file):
    purchase_list = []
    for line in log_file:
        if not line.startswith('20'):
            continue  # drop the line if it's beginning is not year

        entry = line[:-1].split(',')
        date = datetime.strptime(entry[0], '%Y/%m/%d')
        desc = entry[1]
        if entry[6] != '':
            desc += ' ' + entry[6]
        price = entry[2]
        if price == '':
            # TODO: たまにあるけどどうしよう？
            print('empty price is found: ' + line.replace('\n', ''))
            continue
        purchase_list.append((date, desc, int(price), line.replace('\n', '')))
    return purchase_list


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Parse a purchase log file.')
    parser.add_argument('filename', metavar='file', type=str,
                        help='file name that is parsed')
    args = parser.parse_args()

    # check file type (jpbank/pitapa/visa/...)
    input_file = open(args.filename, 'r', encoding='cp932')
    headline = input_file.readline()
    purchase_items = None
    if 'ご利用日,入場時刻,出場時刻,ご利用内容,ご利用額（円）,備考' in headline:
        purchase_items = parsePurchase_pitapa(input_file)
        source = 'pitapa'
    elif 'ＪＰＢＡＮＫＶＩＳＡ' in headline:
        purchase_items = parsePurchase_jpbank(input_file)
        source = 'jpbank'
    elif 'ＳＭＢＣＣＡＲＤ' in headline:
        purchase_items = parsePurchase_visa(input_file)
        source = 'visa'
    elif '"年月日（和暦）","お引出し","お預入れ","お取り扱い内容","残高"' \
            in headline:
        purchase_items = parsePurchase_smbc(input_file)
        source = 'smbc'
    else:
        sys.exit('unknown file type. abort.')

    for item in purchase_items:
        trans = {
            "type": source,
            "time": item[0].isoformat() + "+09:00",  # timezone ?
            "price": item[2],
            "content": item[1],
            "raw": item[3],
            "temporary": True
        }
        data = json.dumps(trans, ensure_ascii=False)
        print(data)
        response = requests.post(
            'http://localhost:50000/api/transactions/',
            data.encode(),
            headers={'Content-Type': 'application/json; charset=utf-8'})
