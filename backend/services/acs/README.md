# Moses ACS [![Build Status](https://travis-ci.org/lucacervasio/mosesacs.svg?branch=master)](https://travis-ci.org/lucacervasio/mosesacs)

An ACS in Go for provisioning CPEs, suitable for test purposes or production deployment.

## Getting started

Install the package:

    go get oktopUSP/backend/services/acs

Run daemon:

    mosesacs -d

Connect to it and get a cli:

    mosesacs

Congratulations, you've connected to the daemon via websocket. Now you can issue commands via CLI or browse the embedded webserver at http://localhost:9292/www

## Compatibility on ARM

Moses is built on purpose only with dependencies in pure GO. So it runs on ARM processors with no issues. We tested it on QNAP devices and Raspberry for remote control.

## CLI commands

### 1. `list`: list CPEs

 example:

```
 moses@localhost:9292/> list
 cpe list
 CPE A54FD with OUI 006754
```

### 2. `readMib SERIAL LEAF/SUBTREE`: read a specific leaf or a subtree

 example:

```
 moses@localhost:9292/> readMib A54FD Device.
 Received an Inform from [::1]:58582 (3191 bytes) with SerialNumber A54FD and EventCodes 6 CONNECTION REQUEST
 InternetGatewayDevice.Time.NTPServer1 : pool.ntp.org
 InternetGatewayDevice.Time.CurrentLocalTime : 2014-07-11T09:08:25
 InternetGatewayDevice.Time.LocalTimeZone : +00:00
 InternetGatewayDevice.Time.LocalTimeZoneName : Greenwich Mean Time : Dublin
 InternetGatewayDevice.Time.DaylightSavingsUsed : 0
```

### 3. `writeMib SERIAL LEAF VALUE`: issue a SetParameterValues and write a value into a leaf

 example:

```
 moses@localhost:9292/> writeMib A54FD InternetGatewayDevice.Time.Enable false
 Received an Inform from [::1]:58582 (3191 bytes) with SerialNumber A54FD and EventCodes 6 CONNECTION REQUEST
```

### 4. `GetParameterNames SERIAL LEAF/SUBTREE`: issue a GetParameterNames and get all leaves/objects at first level

 example:

```
moses@localhost:9292/> GetParameterNames A54FD InternetGatewayDevice.
Received an Inform from [::1]:55385 (3119 bytes) with SerialNumber A54FD and EventCodes 6 CONNECTION REQUEST
InternetGatewayDevice.LANDeviceNumberOfEntries : 0
InternetGatewayDevice.WANDeviceNumberOfEntries : 0
InternetGatewayDevice.DeviceInfo. : 0
InternetGatewayDevice.ManagementServer. : 0
InternetGatewayDevice.Time. : 0
InternetGatewayDevice.Layer3Forwarding. : 0
InternetGatewayDevice.LANDevice. : 0
InternetGatewayDevice.WANDevice. : 0
InternetGatewayDevice.X_00507F_InternetAcc. : 0
InternetGatewayDevice.X_00507F_LAN. : 0
InternetGatewayDevice.X_00507F_NAT. : 0
InternetGatewayDevice.X_00507F_VLAN. : 0
InternetGatewayDevice.X_00507F_Firewall. : 0
InternetGatewayDevice.X_00507F_Applications. : 0
InternetGatewayDevice.X_00507F_System. : 0
InternetGatewayDevice.X_00507F_Status. : 0
InternetGatewayDevice.X_00507F_Diagnostics. : 0
```




## Services exposed

Moses exposes three services:

 - http://localhost:9292/acs is the endpoint for the CPEs to connect
 - http://localhost:9292/www is the embedded webserver to control your CPEs
 - ws://localhost:9292/ws is the websocket endpoint used by the cli to issue commands. Read about the API specification if you want to build a custom frontend which interacts with mosesacs daemon.


