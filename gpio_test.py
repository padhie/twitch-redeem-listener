from gpiozero import LED
from time import sleep

led17 = LED(17)
led18 = LED(18)

while True:
    led17.on()
    sleep(1)
    led17.off()

    led18.on()
    sleep(1)
    led18.off()
