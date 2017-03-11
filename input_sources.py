#
# THIS FILE DOESN'T WORK CURRENTLY
#


# -*- coding: utf-8 -*- 
import sys
import os
import json
import argparse
import datetime
import re
import mysql.connector



def parsePurchase_smbc(log_file):
    printLog("parse smbc")
    purchase_list = []
    for line in log_file:
        entry = line[:-1].split(',')
        # check whether line is purchase entry; Is first item Date ?
        if re.match('H\d\d.\d\d.\d\d$', entry[0]) == None :
            continue

        year = int(entry[0][1:3]) + 1988
        date = re.sub('^H\d\d', str(year), entry[0]) # 和暦を西暦に 
        desc = entry[3]
        desc = desc[1:-1] # eliminate double-quotes from head / tail
        desc = re.sub(r'\u3000+', ' ', desc) # 全角スペースを半角に
        # check income or outgo
        if entry[1] != '':
            price = entry[1]
        else:
            price = '-' + entry[2]
        purchase_list.append((date, desc, price))
    return purchase_list

def parsePurchase_pitapa(log_file):
    printLog("parse pitapa")
    purchase_list = []
    for line in log_file:
        entry = line[:-1].split(',')
        # check whether 'entry' is purchase entry; Is first item Date ?
        if re.match('\d\d\d\d/\d\d/\d\d$', entry[0]) == None :
            continue

        date = entry[0]
        desc = entry[3]
        desc = re.sub(r'\u3000+', ' ', desc) # 全角スペースを半角に
        price = entry[4]
        purchase_list.append((date, desc, price))
    return purchase_list


def parsePurchase_jpbank(log_file):
    printLog("parse jpbank")
    purchase_list = []
    for line in log_file:
        entry = line[:-1].split(',')
        # check whether 'entry' is purchase entry; Is first item Date ?
        if re.match('\d\d\d\d/\d\d/\d\d$', entry[0]) == None :
            continue

        date = entry[0]
        desc = entry[1]
        desc = re.sub(r'\u3000+', ' ', desc) # 全角スペースを半角に
        if entry[6] != '' :
            desc += ' ' + entry[6]
        price = entry[2]
        purchase_list.append((date, desc, price))
    return purchase_list


def parsePurchase_visa(log_file):
    printLog("parse visa")
    purchase_list = []
    for line in log_file:
        entry = line[:-1].split(',')
        # check whether 'entry' is purchase entry; Is first item Date ?
        if re.match('\d\d\d\d/\d\d/\d\d$', entry[0]) == None :
            continue

        date = entry[0]
        desc = entry[1]
        desc = re.sub(r'\u3000+', ' ', desc) # 全角スペースを半角に
        if entry[6] != '' :
            desc += ' ' + entry[6]
        price = entry[2]
        purchase_list.append((date, desc, price))
    return purchase_list


def printLog(message):
    if conf.isLog: print("log: " + message)

class Config(object):
    def loadparam(self, conf_dict, param_name):
        if param_name in conf_dict:
            self.__dict__[param_name] = conf_dict[param_name]
            return True
        else:
            return False

script_path = os.path.abspath(os.path.dirname(__file__))
conf_name = script_path + '/conf/WalletLog.conf'
conf = Config()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Parse a purchase log file.')
    parser.add_argument('filename', metavar='file', type=str, 
                               help='file name that is parsed')
    args = parser.parse_args()

    # load and check config
    with open(conf_name, 'r') as conf_file:
        conf_dict = json.loads(conf_file.read())
        if conf.loadparam(conf_dict, 'isLog') and \
           conf.loadparam(conf_dict, 'db_host') and \
           conf.loadparam(conf_dict, 'db_user') and \
           conf.loadparam(conf_dict, 'db_pass') and \
           conf.loadparam(conf_dict, 'db_name') and \
           conf.loadparam(conf_dict, 'db_charset'): pass
        else:
            print("error: load config file failed; lack of parameter.")
            quit(1)


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

    # connect to database
    connect = mysql.connector.connect(
            user=conf.db_user, password=conf.db_pass,
            host=conf.db_host, database=conf.db_name,
            charset=conf.db_charset, buffered=True)
    cursor = connect.cursor()
    for item in purchase_items:
        try:
            cursor.execute(
                'insert into WalletLog (date, description, source, price) \
                value (%s, %s, %s, %s)', (item[0], item[1], source, item[2]))
            print("insert : {0}, {1}, {2}, {3}".format( \
                item[0], item[1], source, item[2]))
        except mysql.connector.errors.IntegrityError:
            print('duplicate')

    connect.commit()
    cursor.close()
    connect.close()
