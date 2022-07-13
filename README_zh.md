# Gree HVAC MQTT bridge GO

## [Enginsh](README.md) 简体中文


使用MQTT广播与格力空调进行通信的桥接服务。它也可以作为[Hass.io](https://home-assistant.io/)插件使用。
## 要求

- 同一网络上的MQTT代理和格力智能HVAC设备

## 本地运行

```shell
app -DIR 192.168.1.255 \
    -MBU 192.168.1.1 \
    -MBP 1883 \
    -MTP home/greehvac \
    -MU admin \
    -MP admin \
    -MR false
```

## 支持的命令

MQTT主题方案:

- `MTP/COMMAND/get` 获取 值
- `MTP/COMMAND/set` 设置 值

注意：开关是用0或1来设置的。

| Command  | Values                                                                                                                                                            | Description                                  |
|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------|
| **温度**   | any integer                                                                                                                                                       | 默认以摄氏度为单位                                    |
| **模式**   | _off_, _auto_, _cool_, _heat_, _dry_, _fan_only_                                                                                                                  | 运行模式                                         |
| **风速**   | _auto_, _low_, _mediumLow_, _medium_, _mediumHigh_, _high_                                                                                                        | 风扇速度                                         |
| **左右摆风** | _default_, _full_, _fixedLeft_, _fixedMidLeft_, _fixedMid_, _fixedMidRight_, _fixedRight_                                                                         | 水平摆动                                         |
| **上下摆风** | _default_, _full_, _fixedTop_, _fixedMidTop_, _fixedMid_, _fixedMidBottom_, _fixedBottom_, _swingBottom_, _swingMidBottom_, _swingMid_, _swingMidTop_, _swingTop_ | 上下摆动                                         |
| **电源**   | _0_, _1_                                                                                                                                                          | 开启/关闭设备                                      |
| **健康**   | _0_, _1_                                                                                                                                                          | 健康（"冷等离子体"）模式，仅适用于配备 "阴离子发生器 "的设备，可吸收灰尘和杀灭细菌 |
| **省电**   | _0_, _1_                                                                                                                                                          | 省电模式                                         |
| **灯光**   | _0_, _1_                                                                                                                                                          | 开启/关闭设备灯                                     |
| **安静**   | _0_, _1_, _2_, _3_                                                                                                                                                | 静音模式                                         |
| **BLOW** | _0_, _1_                                                                                                                                                          | 关机后保持风扇运行一段时间（也叫 "X-风扇"，只在干燥和冷却模式下可用）。       |
| **空气循环** | _off_, _inside_, _outside_, _mode3_                                                                                                                               | 新鲜空气阀                                        |
| **睡眠**   | _0_, _1_                                                                                                                                                          | 睡眠模式                                         |
| **强劲**   | _0_, _1_                                                                                                                                                          | 涡轮模式                                         |

## Hass.io附加组件

该服务可以作为Hass.io [MQTT气候平台](https://home-assistant.io/components/climate.mqtt/)的第三方插件使用，尽管不是所有的命令都支持。

1. [安装](https://home-assistant.io/hassio/installing_third_party_addons/) 该插件
2. 自定义附加选项（HVAC主机、MQTT代理URL、MQTT主题前缀）。
3. 在你的`configuration.yaml`中加入以下内容：

```yaml
climate:
  - platform: mqtt

    # Change to whatever you want
    name: Gree HVAC

    # Change MQTT_TOPIC_PREFIX to what you've set in addon options
    current_temperature_topic: "MQTT_TOPIC_PREFIX/temperature/get"
    temperature_command_topic: "MQTT_TOPIC_PREFIX/temperature/set"
    temperature_state_topic: "MQTT_TOPIC_PREFIX/temperature/get"
    mode_state_topic: "MQTT_TOPIC_PREFIX/mode/get"
    mode_command_topic: "MQTT_TOPIC_PREFIX/mode/set"
    fan_mode_state_topic: "MQTT_TOPIC_PREFIX/fanspeed/get"
    fan_mode_command_topic: "MQTT_TOPIC_PREFIX/fanspeed/set"
    swing_mode_state_topic: "MQTT_TOPIC_PREFIX/swingvert/get"
    swing_mode_command_topic: "MQTT_TOPIC_PREFIX/swingvert/set"
    power_state_topic: "MQTT_TOPIC_PREFIX/power/get"
    power_command_topic: "MQTT_TOPIC_PREFIX/power/set"

    # Keep the following as is
    payload_off: 0
    payload_on: 1
    modes:
      - "off"
      - "auto"
      - "cool"
      - "heat"
      - "dry"
      - "fan_only"
    swing_modes:
      - "default"
      - "full"
      - "fixedTop"
      - "fixedMidTop"
      - "fixedMid"
      - "fixedMidBottom"
      - "fixedBottom"
      - "swingBottom"
      - "swingMidBottom"
      - "swingMid"
      - "swingMidTop"
      - "swingTop"
    fan_modes:
      - "auto"
      - "low"
      - "mediumLow"
      - "medium"
      - "mediumHigh"
      - "high"
```

### 如何开启/关闭电源

Hass.io不提供单独的开/关开关。请使用专用模式。


## 配置空调WiFi

1. 确保你的暖通空调在AP模式下运行。你可以通过按下空调遥控器上的MODE +WIFI（或 MODE + TURBO）5秒钟来重置WiFi配置。
2. 连接AP wifi网络（SSID名称应该是一个8位数字，例如 "u34k5l166"，默认密码应该是：12345678）。
3. 在你的UNIX终端运行以下程序：

```shell
echo -n "{\"psw\": \"YOUR_WIFI_PASSWORD\",\"ssid\": \"YOUR_WIFI_SSID\",\"t\": \"wlan\"}" | nc -cu 192.168.1.1 7000
````

注意：这个命令可能会根据你的操作系统（如Linux、macOS、CygWin）而有所不同。如果遇到问题，请查阅相应的netcat手册。

## 更新日志

[1.0.0]

First release

## 许可证

这个项目是根据GNU GPLv3授权的--详见[LICENSE](LICENSE)文件

## 鸣谢

- [arthurkrupa](https://github.com/arthurkrupa) for gree-hvac-mqtt-bridge node.js project
- [tomikaa87](https://github.com/tomikaa87) for reverse-engineering the Gree protocol
- [oroce](https://github.com/oroce) for inspiration
- [arthurkrupa](https://https://github.com/arthurkrupa) for the actual service
- [bkbilly](https://github.com/bkbilly) for service improvements to MQTT
- [aaronsb](https://github.com/aaronsb) for sweeping the Node floor
