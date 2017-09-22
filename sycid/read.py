#!/usr/bin/env python

import sys

while(1):

    #sys.stdin = open('/dev/tty0', 'r')
    sys.stdin = open('/dev/hidraw0', 'r')
    rfid_data = raw_input()
    print "Read code from RFID reader: {0}".format(rfid_data)
