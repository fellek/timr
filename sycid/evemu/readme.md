# eveemu

From [Evemu](https://wiki.ubuntu.com/Multitouch/Testing/Evemu)

# Overview

To describe a device in your machine:

  `$ sudo evemu-describe /dev/input/by-id/usb-Sycreader_RFID_Technology_Co.__Ltd_SYC_ID_IC_USB_Reader_08FF20140315-event-kbd > device.prop`

To record data from your device:

  `$ sudo evemu-record /dev/input/by-id/usb-Sycreader_RFID_Technology_Co.__Ltd_SYC_ID_IC_USB_Reader_08FF20140315-event-kbd > device.event`

To remotely setup a copy of the device:

  `$ sudo evemu-device device.prop`

The device node will be printed on the terminal, and will stay valid until the program is 
terminated. To replay the input stream (beware that pointer clicks can execute commands in your window manager):

  `$ sudo evemu-play /dev/input/by-id/usb-Sycreader_RFID_Technology_Co.__Ltd_SYC_ID_IC_USB_Reader_08FF20140315-event-kbd < device.event`

More interestingly, one can first setup a program like mtview to grab the input:

  `$ sudo mtview /dev/input/by-id/usb-Sycreader_RFID_Technology_Co.__Ltd_SYC_ID_IC_USB_Reader_08FF20140315-event-kbd`

and then replay the data, which will show every individual finger on the screen.

Given the evemu lib, one can time gestures by first setting up the device, replay gestures on the device, and time the appearance of gestures in grail.