import React, { useEffect, useState } from "react";

import { Card, CardActions, CardContent, CardHeader, CircularProgress, Divider, Grid, InputLabel, MenuItem, Select, SvgIcon, Tooltip, Typography } from "@mui/material";
import CpuChipIcon from "@heroicons/react/24/solid/CpuChipIcon";
import { Stack } from "@mui/system";
import { useTheme } from "@emotion/react";
import { useRouter } from "next/router";
import { set } from "nprogress";

export const ConnectedDevices = () => {

    const theme = useTheme();
    const router = useRouter();

    const [content, setContent] = useState(null);
    const [interfaces, setInterfaces] = useState([]);
    const [interfaceValue, setInterfaceValue] = useState(null);

    const getConnectionState = (rssi) => {
        let connectionStatus = "Signal "
        if (rssi > -30) {
            return connectionStatus + "Excellent"
        } else if (rssi > -60) {
            return connectionStatus + "Good"
        } else if (rssi > -70) {
            return connectionStatus + "Bad"
        } else {
            return connectionStatus + "Awful"
        }
    }

    const fetchConnectedDevicesData = async () => {

        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", localStorage.getItem("token"));

        var requestOptions = {
            method: 'GET',
            headers: myHeaders,
            redirect: 'follow'
        };

        fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device/${router.query.id[0]}/connecteddevices`, requestOptions)
            .then(response => {
                if (response.status === 401) {
                    router.push("/auth/login")
                }
                return response.json()
            })
            .then(result => {
                console.log("connecteddevices content", result)
                let interfaces = Object.keys(result)
                setInterfaces(interfaces)
                setInterfaceValue(interfaces[0])
                setContent(result)
            })
            .catch(error => console.log('error', error));
    };

    useEffect(() => {
        fetchConnectedDevicesData();
    }, [])

    return (
        <Stack
            justifyContent="center"
            alignItems={(!content || interfaces.length == 0) && "center"}
        >
            {content && interfaces.length > 0 ?
                <Card>
                    <CardContent>
                        <Grid mb={3}>
                            <InputLabel> Interface </InputLabel>
                            <Select label="interface" variant="standard" value={interfaceValue} onChange={(e) => setInterfaceValue(e.target.value)}>
                                {(
                                    interfaces.map((item, index) => (
                                        <MenuItem key={index} value={item}>
                                            {item}
                                        </MenuItem>
                                    ))
                                )}
                            </Select>
                        </Grid>
                        {
                            content[interfaceValue].map((property, index) => (
                                <Card key={index}>
                                    <CardContent>
                                        <Grid container justifyContent={"center"}>
                                            <Stack direction="row" spacing={5}>
                                                <Stack justifyItems={"center"} direction={"row"} mt={2}>
                                                    <Tooltip title={property.active ? "Online" : "Offline"}>
                                                        <SvgIcon>
                                                            <CpuChipIcon color={property.active ? theme.palette.success.main : theme.palette.error.main}></CpuChipIcon>
                                                        </SvgIcon>
                                                    </Tooltip>
                                                    <Typography ml={"10px"}>
                                                        {property.hostname}
                                                    </Typography>
                                                </Stack>
                                                <Divider orientation="vertical" />
                                                <Stack spacing={2}>
                                                    <Typography>
                                                        IP address: {property.ip_adress}
                                                    </Typography>
                                                    <Typography>
                                                        MAC: {property.mac}
                                                    </Typography>
                                                </Stack>
                                                <Stack spacing={2}>
                                                    <Typography>
                                                        Source: {property.adress_source}
                                                    </Typography>
                                                    <Tooltip title={getConnectionState(property.rssi)}>
                                                        <Typography display={"flex"} color={() => {
                                                            let rssi = property.rssi
                                                            if(rssi == 0){
                                                                return theme.palette.neutral[900]
                                                            } else if (rssi > -30) {
                                                                return theme.palette.success.main
                                                            } else if (rssi > -60) {
                                                                return theme.palette.success.main
                                                            } else if (rssi > -70) {
                                                                return theme.palette.warning.main
                                                            } else {
                                                                return theme.palette.error.main
                                                            }
                                                        }}>
                                                            <Typography color={theme.palette.neutral[900]} sx={{pr:"5px"}}>
                                                                RSSI:
                                                            </Typography> 
                                                            {property.rssi} dbm
                                                        </Typography>
                                                    </Tooltip>
                                                </Stack>
                                            </Stack>
                                        </Grid>
                                    </CardContent>
                                </Card>
                            ))
                        }
                    </CardContent>
                </Card> : (
                    content ? <Typography> No connected devices found </Typography> : <CircularProgress />
                )}
        </Stack>
    )
}