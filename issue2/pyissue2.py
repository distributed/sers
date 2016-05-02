import serial
from time import sleep

head = serial.Serial("/dev/ttyAMA0", 9600, timeout=0.01)

#while True:
#    head.write('\x55')
#    sleep(0.010)

