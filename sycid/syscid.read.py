#!/usr/bin/env python
"""This script prompts reads in an endless loop the sycid rfid reader and pushes the
read id to an webapp api
usb-Sycreader_RFID_Technology_Co.__Ltd_SYC_ID_IC_USB_Reader_08FF20140315-event-kbd
"""

import httplib
import logging
import logging.handlers
import time

from select import select

from evdev import InputDevice

KEYS = "X^1234567890XXXXqwertzuiopXXXXasdfghjklXXXXXyxcvbnmXXXXXXXXXXXXXXXXXXXXXXX"
DEV = InputDevice("/dev/input/by-id/usb-Sycreader_RFID_Technology_Co.__Ltd_SYC_ID_IC_USB_Reader_08FF20140315-event-kbd")

LOG_FILENAME = "/var/log/nfc.log"
LOG_LEVEL = logging.DEBUG

SERVER = 'localhost'
URI = '/api/cardreader?id='
	
LOGGER = logging.getLogger(__name__)
LOGGER.setLevel(LOG_LEVEL)
LOGHANDLER = logging.handlers.TimedRotatingFileHandler(LOG_FILENAME, when='midnight', backupCount=3)
LOGFORMATTER = logging.Formatter('%(asctime)s %(levelname)-8s %(message)s')
LOGHANDLER.setFormatter(LOGFORMATTER)
LOGGER.addHandler(LOGHANDLER)

LOGGER.info("RFID Start")

def read_data(dev):
    '''read from chip, waits for data.... ideling
    0000802843
    000802843
    '''
    LOGGER.info("__ wait for data")

    rfid = ""
    done = False
    try:
		        r, w, x = select([dev], [], [])

        while not done: 
            for event in dev.read():
                LOGGER.info("__ rfid.len: %d, type %d, value %d, code %s", 
                len(rfid), event.type, event.value, event.code)
                if event.type == 0 and event.value == 0 and event.code == 0:
                    LOGGER.info("__ noop") 
                    continue
                if event.type != 1 or event.value != 1:
                    continue
                if event.code == 28:
                    done = True
                    break

                LOGGER.info("__ adding %s to rfid %s", KEYS[event.code], rfid) 
                rfid += KEYS[event.code]
     
    except Exception:
        LOGGER.error('__ generic exception: ' + traceback.format_exc())

    LOGGER.info("__ returning " + rfid)

    return rfid

def send_to_server(rfid):
    '''send result to server'''
    conn = httplib.HTTPConnection(SERVER, 80, timeout=5)
    conn.request("GET", URI + rfid)
    res = conn.getresponse()
    if res.status == 200:
        ret = res.getheader('X-Return')
        LOGGER.info("http response: " + ret)
    conn.close()

RFID_DATA = ''
DEV.grab()

while True:
    try:
        RFID_DATA = read_data(DEV)
       # print "Read code from RFID reader: " + RFID_DATA
        if RFID_DATA:
            LOGGER.info("Chip read: " + RFID_DATA)
            # send_to_server(RFID_DATA)
            time.sleep(2)
    except Exception:
        import traceback
        LOGGER.error('generic exception: ' + traceback.format_exc())

DEV.ungrab()
