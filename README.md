# Gree HVAC MQTT bridge GO

Bridge service for communicating with Gree air conditioners using MQTT broadcasts. It can also be used as a [Hass.io](https://home-assistant.io/) addon.

## Requirements

- GO (>=18.0)
- An MQTT broker and Gree smart HVAC device on the same network

## Running locally

Make sure you have GO (>=18.0) installed and run the following (adjust the arguments to match your setup):

```shell
app -DIR 192.168.1.255 \
    -MBU 192.168.1.1 \
    -MBP 1883 \
    -MTP home/greehvac \
    -MU admin \
    -MP admin \
    -MR false
```

## Supported commands

MQTT topic scheme:

- `MQTT_TOPIC_PREFIX/COMMAND/get` Get value
- `MQTT_TOPIC_PREFIX/COMMAND/set` Set value

Note: _boolean_ values are set using 0 or 1

| Command         | Values                                                                                                                                                            | Description                                                                                                          |
|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------|
| **temperature** | any integer                                                                                                                                                       | In degrees Celsius by default                                                                                        |
| **mode**        | _off_, _auto_, _cool_, _heat_, _dry_, _fan_only_                                                                                                                  | Operation mode                                                                                                       |
| **fanspeed**    | _auto_, _low_, _mediumLow_, _medium_, _mediumHigh_, _high_                                                                                                        | Fan speed                                                                                                            |
| **swinghor**    | _default_, _full_, _fixedLeft_, _fixedMidLeft_, _fixedMid_, _fixedMidRight_, _fixedRight_                                                                         | Horizontal Swing                                                                                                     |
| **swingvert**   | _default_, _full_, _fixedTop_, _fixedMidTop_, _fixedMid_, _fixedMidBottom_, _fixedBottom_, _swingBottom_, _swingMidBottom_, _swingMid_, _swingMidTop_, _swingTop_ | Vetical swing                                                                                                        |
| **power**       | _0_, _1_                                                                                                                                                          | Turn device on/off                                                                                                   |
| **health**      | _0_, _1_                                                                                                                                                          | Health ("Cold plasma") mode, only for devices equipped with "anion generator", which absorbs dust and kills bacteria |
| **powersave**   | _0_, _1_                                                                                                                                                          | Power Saving mode                                                                                                    |
| **lights**      | _0_, _1_                                                                                                                                                          | Turn on/off device lights                                                                                            |
| **quiet**       | _0_, _1_, _2_, _3_                                                                                                                                                | Quiet modes                                                                                                          |
| **blow**        | _0_, _1_                                                                                                                                                          | Keeps the fan running for a while after shutting down (also called "X-Fan", only usable in Dry and Cool mode)        |
| **air**         | _off_, _inside_, _outside_, _mode3_                                                                                                                               | Fresh air valve                                                                                                      |
| **sleep**       | _0_, _1_                                                                                                                                                          | Sleep mode                                                                                                           |
| **turbo**       | _0_, _1_                                                                                                                                                          | Turbo mode                                                                                                           |

## Hass.io addon

The service can be used as a 3rd party addon for the Hass.io [MQTT climate platform](https://home-assistant.io/components/climate.mqtt/), although not all commands are supported.

1. [Install](https://home-assistant.io/hassio/installing_third_party_addons/) the addon
2. Customize addon options (HVAC host, MQTT broker URL, MQTT topic prefix)
3. Add the following to your `configuration.yaml`

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

### How to power on/off

Hass.io doesn't supply separate on/off switch. Use the dedicated mode for that.


## Configuring HVAC WiFi

1. Make sure your HVAC is running in AP mode. You can reset the WiFi config by pressing MODE +WIFI (or MODE + TURBO) on the AC remote for 5s.
2. Connect with the AP wifi network (the SSID name should be a 8-character alfanumeric, e.g. "u34k5l166").
3. Run the following in your UNIX terminal:

```shell
echo -n "{\"psw\": \"YOUR_WIFI_PASSWORD\",\"ssid\": \"YOUR_WIFI_SSID\",\"t\": \"wlan\"}" | nc -cu 192.168.1.1 7000
````

Note: This command may vary depending on your OS (e.g. Linux, macOS, CygWin). If facing problems, please consult the appropriate netcat manual.

## Changelog

[1.0.0]

First release

## License

This project is licensed under the GNU GPLv3 - see the [LICENSE](LICENSE) file for details

## Acknowledgments

- [arthurkrupa](https://github.com/arthurkrupa) for gree-hvac-mqtt-bridge node.js project
- [tomikaa87](https://github.com/tomikaa87) for reverse-engineering the Gree protocol
- [oroce](https://github.com/oroce) for inspiration
- [arthurkrupa](https://https://github.com/arthurkrupa) for the actual service
- [bkbilly](https://github.com/bkbilly) for service improvements to MQTT
- [aaronsb](https://github.com/aaronsb) for sweeping the Node floor
