## Development (on raspberry pi)
### prepare system
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y golang git make
go get -u github.com/stianeikeland/go-rpio/v4
sudo apt install -y build-essential
```


## Hardware
- Raspberry Pi 4
### LED Lib
activate GPIO pins
```bash
sudo raspi-config
```
navigate to "Interface Options" -> "SPI" -> "Enable"
navigate to "Interface Options" -> "I2C" -> "Enable"

## Twitch
### generate twitch token
url: [twitchtokengenerator](https://twitchtokengenerator.com/)   
relevant scopes:
- channel:read:redemptions


## Autostart
### Service
copy service file to system
```bash
cp ./service/twitch-redeem.service /etc/systemd/system/twitch-redeem.service
```

enable autostart
```bash
sudo systemctl enable twitch-redeem.service
sudo systemctl start twitch-redeem.service
```


## Debug
### LEDs (GPIO)
Redeem detected (1x blink) -> GPIO 17 (Pin 11)   
Tapo Request (2x blink)    -> GPIO 18 (Pin 12)

### Logfile
Primary Log: `/var/log/twitch-redeem.log`   
Fallback:    `Syslog`

### Web-UI
Set `.env`-varaible `ENABLE_WEB_INTERFACE=true`

