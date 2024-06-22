import Head from 'next/head';
import { Box, Stack, Typography, Container, Unstable_Grid2 as Grid,
Tab, 
Tabs,
SvgIcon,
Breadcrumbs,
Link } from '@mui/material';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { useRouter } from 'next/router';
import { DevicesRPC } from 'src/sections/devices/cwmp/devices-rpc';
import EnvelopeIcon from '@heroicons/react/24/outline/EnvelopeIcon';
import MagnifyingGlassIcon from '@heroicons/react/24/solid/MagnifyingGlassIcon';
import WifiIcon from '@heroicons/react/24/solid/WifiIcon';
import WrenchScrewDriverIcon from '@heroicons/react/24/solid/WrenchScrewdriverIcon';
import SignalIcon from '@heroicons/react/24/solid/SignalIcon';
import DevicePhoneMobile from '@heroicons/react/24/solid/DevicePhoneMobileIcon';
import { useEffect, useState } from 'react';
import { DevicesWiFi } from 'src/sections/devices/cwmp/devices-wifi';
import { DevicesDiagnostic } from 'src/sections/devices/cwmp/devices-diagnostic';
import { SiteSurvey } from 'src/sections/devices/cwmp/site-survey';
import { ConnectedDevices } from 'src/sections/devices/cwmp/connecteddevices';


const Page = () => {
    const router = useRouter()

    const deviceID = router.query.id[0]
    const section = router.query.id[1]

    const sectionHandler = () => {
        switch(section){
            case "msg":
                return <DevicesRPC/>
            case "wifi":
                return <DevicesWiFi/>
            case "diagnostic":
                    return <DevicesDiagnostic/>
            case "connected-devices":
                return <ConnectedDevices/>
            case "site-survey":
                return <SiteSurvey/>
            default:
                return <p>Hello World</p>
        }
    }

    useEffect(()=>{
        console.log("deviceid:",deviceID)
    })
  
    return(
    <>
        <Head>
            <title>
                Oktopus | TR-369
            </title>
        </Head>
        <Box
            component="main"
            sx={{
                flexGrow: 1,
                py: 0,
            }}
        >
            <Container maxWidth="lg" >
            <Stack spacing={3} >
                <Breadcrumbs separator="›" aria-label="breadcrumb" sx={{md: 40, mr: 20}}>
                {[<Link underline="hover" key="1" color="inherit" href="/devices">
                    Devices
                </Link>,
                <Link
                underline="none"
                key="2"
                color="inherit"
                hre={`/devices/${deviceID}`}
                >
                {deviceID}
                </Link>]}
                </Breadcrumbs>
                <Box sx={{
                    display:'flex',
                    justifyContent:'center',
                }}>
                    <Tabs value={router.query.id[1]}  aria-label="icon label tabs example">
                        <Tab icon={<SvgIcon><WifiIcon/></SvgIcon>} iconPosition={"end"} label="Wi-Fi" onClick={()=>{router.push(`/devices/cwmp/${deviceID}/wifi`)}} value={"wifi"}/>
                        <Tab icon={<SvgIcon><SignalIcon/></SvgIcon>} iconPosition={"end"} label="Site Survey" onClick={()=>{router.push(`/devices/cwmp/${deviceID}/site-survey`)}} value={"site-survey"}/>
                        <Tab icon={<SvgIcon><DevicePhoneMobile/></SvgIcon>} iconPosition={"end"} label="Connected Devices" onClick={()=>{router.push(`/devices/cwmp/${deviceID}/connected-devices`)}} value={"connected-devices"}/>
                        {/* <Tab value={"discovery"} onClick={()=>{router.push(`/devices/cwmp/${deviceID}/discovery`)}} icon={<SvgIcon><MagnifyingGlassIcon/></SvgIcon>} iconPosition={"end"} label="Discover Parameters" /> */}
                        <Tab icon={<SvgIcon><WrenchScrewDriverIcon/></SvgIcon>} iconPosition={"end"} label="Diagnostic" onClick={()=>{router.push(`/devices/cwmp/${deviceID}/diagnostic`)}} value={"diagnostic"}/>
                        <Tab value={"msg"} onClick={()=>{router.push(`/devices/cwmp/${deviceID}/msg`)}} icon={<SvgIcon><EnvelopeIcon/></SvgIcon>} iconPosition={"end"} label="Remote Messages" />
                    </Tabs>
                </Box>
                {
                   sectionHandler()
                }
                </Stack>
            </Container>
        </Box>
    </>
    );
};

Page.getLayout = (page) => (
    <DashboardLayout>
        {page}
    </DashboardLayout>
);

export default Page;