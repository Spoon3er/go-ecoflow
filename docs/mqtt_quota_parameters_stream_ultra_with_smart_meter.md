# MQTT Quota Topic Parameters

This document describes the parameters received from the MQTT quota topic: `/open/{certificateAccount}/{sn}/quota`

The payload is a flat JSON object with incremental updates. Parameters are sent only when they change.

## Power Parameters

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| powGetSysGrid | float | Power from grid | W |
| powGetSysLoad | float | Total system load power | W |
| powGetSysLoadFromGrid | float | Load power from grid | W |
| powGetSysLoadFromBp | float | Load power from battery pack | W |
| powGetSysLoadFromPv | float | Load power from PV/solar | W |
| powGetBpCms | float | Battery pack CMS power | W |
| powGetPv | float | PV input power | W |
| powGetPv2 | float | PV2 input power | W |
| powGetPv3 | float | PV3 input power | W |
| powGetPv4 | float | PV4 input power | W |
| powGetPvSum | float | Total PV power | W |
| powGetSchuko1 | float | Schuko socket 1 power | W |
| powGetSchuko2 | float | Schuko socket 2 power | W |
| powSysAcInMax | int | Maximum AC input power | W |
| powSysAcOutMax | int | Maximum AC output power | W |
| socketMeasurePower | float | Socket measured power | W |
| gridConnectionPower | float | Grid connection power | W |
| sysGridConnectionPower | float | System grid connection power | W |
| powConsumptionMeasurement | int | Power consumption measurement mode | - |

## Battery Management System (BMS) Parameters

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| bmsBattSoc | float | Battery state of charge | % |
| bmsBattSoh | float | Battery state of health | % |
| bmsDesignCap | int | Battery design capacity | Wh |
| bmsBattHeating | bool | Battery heating status | - |
| bmsChgDsgState | int | Battery charge/discharge state | - |
| bmsChgRemTime | int | Battery charging remaining time | seconds |
| bmsDsgRemTime | int | Battery discharging remaining time | seconds |
| bmsMaxCellTemp | float | Maximum battery cell temperature | °C |
| bmsMinCellTemp | float | Minimum battery cell temperature | °C |
| bmsMaxMosTemp | float | Maximum MOSFET temperature | °C |
| bmsMinMosTemp | float | Minimum MOSFET temperature | °C |

## CMS (Central Management System) Parameters

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| cmsBattSoc | float | CMS battery state of charge | % |
| cmsBattSoh | float | CMS battery state of health | % |
| cmsBattFullEnergy | int | CMS battery full energy | Wh |
| cmsBattPowInMax | int | CMS battery max input power | W |
| cmsBattPowOutMax | int | CMS battery max output power | W |
| cmsBmsRunState | int | CMS BMS run state | - |
| cmsChgDsgState | int | CMS charge/discharge state | - |
| cmsChgRemTime | int | CMS charging remaining time | seconds |
| cmsDsgRemTime | int | CMS discharging remaining time | seconds |
| cmsMaxChgSoc | int | CMS maximum charge SOC | % |
| cmsMinDsgSoc | int | CMS minimum discharge SOC | % |

## Grid Connection Parameters

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| gridConnectionVol | float | Grid connection voltage | V |
| gridConnectionFreq | float | Grid connection frequency | Hz |
| gridConnectionSta | string | Grid connection status | - |
| gridCodeSelection | string | Grid code selection | - |
| gridCodeVersion | int | Grid code version | - |
| sysGridInPwrLimit | int | System grid input power limit | W |
| gridSysDeviceCnt | int | Grid system device count | - |
| sysOffgrid | bool | System off-grid status | - |

## PV (Photovoltaic) Input Parameters

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| plugInInfoPvVol | float | PV voltage | V |
| plugInInfoPvAmp | float | PV current | A |
| plugInInfoPvFlag | bool | PV input flag | - |
| plugInInfoPv2Vol | float | PV2 voltage | V |
| plugInInfoPv2Amp | float | PV2 current | A |
| plugInInfoPv2Flag | bool | PV2 input flag | - |
| plugInInfoPv3Vol | float | PV3 voltage | V |
| plugInInfoPv3Amp | float | PV3 current | A |
| plugInInfoPv3Flag | bool | PV3 input flag | - |
| plugInInfoPv4Vol | float | PV4 voltage | V |
| plugInInfoPv4Amp | float | PV4 current | A |
| plugInInfoPv4Flag | bool | PV4 input flag | - |

## Inverter/Power Limits

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| maxInvInput | int | Maximum inverter input | W |
| maxInvOutput | int | Maximum inverter output | W |
| maxBpInput | int | Maximum battery pack input | W |
| maxBpOutput | int | Maximum battery pack output | W |
| busbarPowLimit | int | Busbar power limit | W |

## Feed to Grid Parameters

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| feedGridMode | int | Feed to grid mode | - |
| feedGridModePowLimit | int | Feed to grid mode power limit | W |
| feedGridModePowMax | int | Feed to grid mode max power | W |
| feedGridSafetyPowMax | int | Feed to grid safety max power | W |

## Relay Status

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| relay1Onoff | bool | Relay 1 on/off status | - |
| relay2Onoff | bool | Relay 2 on/off status | - |
| relay3Onoff | bool | Relay 3 on/off status | - |
| relay4Onoff | bool | Relay 4 on/off status | - |

## Socket/Outlet Parameters

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| scoket1BindDeviceSn | string | Socket 1 bind device serial number | - |
| scoket2BindDeviceSn | string | Socket 2 bind device serial number | - |
| socket1PwrInUnbind | bool | Socket 1 power input unbind | - |
| socket2PwrInUnbind | bool | Socket 2 power input unbind | - |

## System Configuration

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| brightness | int | Display brightness | % |
| moduleWifiRssi | float | WiFi signal strength | dBm |
| utcTimezone | int | UTC timezone offset | minutes |
| utcTimezoneId | string | Timezone ID | - |
| utcSetMode | bool | UTC set mode | - |
| townCode | int | Town code | - |
| updateBanFlag | int | Update ban flag | - |
| devCtrlStatus | int | Device control status | - |

## Energy Strategy

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| energyBackupState | int | Energy backup state | - |
| backupReverseSoc | int | Backup reverse SOC | % |

## Series/Distributed System

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| seriesConnectDeviceId | int | Series connect device ID | - |
| seriesConnectDeviceStatus | string | Series connect device status | - |
| distributedDeviceStatus | string | Distributed device status | - |
| systemGroupId | float | System group ID | - |
| systemMeshId | int | System mesh ID | - |

## Smart Meter

| Parameter | Type | Description | Unit |
|-----------|------|-------------|------|
| useLanMeter | bool | Use LAN meter | - |

## Complex Objects

### cloudMetter
```json
{
  "hasMeter": bool,
  "model": string,
  "phaseAPower": float,
  "phaseBPower": float,
  "phaseCPower": float,
  "sn": string
}
```

### gridConnectionPortBind
```json
{
  "err": int,
  "portNum": int,
  "sn": string
}
```

### energyStrategyOperateMode
```json
{
  "operateIntelligentScheduleModeOpen": bool,
  "operateScheduledOpen": bool,
  "operateSelfPoweredOpen": bool,
  "operateTouModeOpen": bool
}
```

### allTimerTask
```json
{
  "timeTask": [
    {
      "taskIndex": int,
      "role": int,
      "timeParam": int,
      "chgTask": {
        "chgMode": int,
        "chgSource": int,
        "devTargetSoc": [
          {
            "targetSoc": int,
            "chgFromGridPowerLimited": int,
            "sn": string
          }
        ]
      },
      "isCfg": int,
      "isEffect": bool,
      "timeMode": int,
      "timeTable": [int]
    }
  ]
}
```

### timezoneChangeList
```json
{
  "timeZoneChangeItem": [
    {
      "utcTime": float,
      "utcTimezone": int
    }
  ]
}
```

### powerSocket
```json
{
  "powerSocketCfg": []
}
```

### controllableLoadList
```json
{
  "list": []
}
```

### dayResidentLoadList
```json
{
  "load": []
}
```

### devErrcodeList
```json
{
  "devErrcode": []
}
```

### wifiApMeshId
```json
{
  "idList": []
}
```

## Notes

1. **Incremental Updates**: The device only sends parameters that have changed, not the full state every time.
2. **Data Types**: Most numeric values are `float64` when parsed from JSON, even if they represent integers.
3. **Boolean Values**: Boolean fields use Go bool type (true/false).
4. **Empty Strings**: Some string fields may be empty ("") when not configured.
5. **Arrays**: Some parameters contain arrays or nested objects for complex data structures.
