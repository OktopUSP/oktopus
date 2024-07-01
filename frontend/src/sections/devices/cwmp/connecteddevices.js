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

    const fetchConnectedDevicesData = async () => {

        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", localStorage.getItem("token"));
      
        var requestOptions = {
          method: 'GET',
          headers: myHeaders,
          redirect: 'follow'
        };

        fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT}/device/${router.query.id[0]}/connecteddevices`, requestOptions)
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
    },[])

    return (         
    <Stack 
        justifyContent="center" 
        alignItems={(!content || interfaces.length == 0) &&"center"}
        >
        {content && interfaces.length > 0 ?
        <Card>
            <CardContent>
                <Grid mb={3}>
                    <InputLabel> Interface </InputLabel>
                    <Select label="interface" variant="standard" value={interfaceValue} onChange={(e)=> setInterfaceValue(e.target.value)}>
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
                    content[interfaceValue].map((property) => (
                        <Card>
                            <CardContent>
                                <Grid container justifyContent={"center"}>
                                    <Stack direction="row" spacing={5}>
                                        <Stack justifyItems={"center"} direction={"row"} mt={2}>
                                            <Tooltip title={property.active ? "Online": "Offline"}>
                                                <SvgIcon>
                                                    <CpuChipIcon color={property.active ? theme.palette.success.main : theme.palette.error.main}></CpuChipIcon>
                                                </SvgIcon>
                                            </Tooltip>
                                            <Typography ml={"10px"}>
                                                {property.hostname}
                                            </Typography>
                                        </Stack>
                                        <Divider orientation="vertical"/>
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
                                                RSSI: {property.rssi} dbm
                                            </Typography>
                                            <Typography>
                                                Source: {property.adress_source}
                                            </Typography>
                                        </Stack>
                                    </Stack>
                                </Grid>
                            </CardContent>
                        </Card>
                    ))
                }
            </CardContent>
        </Card>: (
            content ? <Typography> No connected devices found </Typography> : <CircularProgress/>
        )}
        </Stack>
    )
}