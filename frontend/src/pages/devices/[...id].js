import Head from 'next/head';
import { Box, Stack, Typography, Container, Unstable_Grid2 as Grid,
Tab, 
Tabs,
SvgIcon } from '@mui/material';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { useRouter } from 'next/router';
import { DevicesRPC } from 'src/sections/devices/devices-rpc';
import { DevicesDiscovery } from 'src/sections/devices/devices-discovery';
import EnvelopeIcon from '@heroicons/react/24/outline/EnvelopeIcon';
import MagnifyingGlassIcon from '@heroicons/react/24/solid/MagnifyingGlassIcon';
import WifiIcon from '@heroicons/react/24/solid/WifiIcon';
import { useEffect, useState } from 'react';

const Page = () => {
    const router = useRouter()

    const deviceID = router.query.id[0]
    const section = router.query.id[1]

    const sectionHandler = () => {
        switch(section){
            case "msg":
                return <DevicesRPC/>
            case "discovery":
                return <DevicesDiscovery/>
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
                <Box sx={{
                display:'flex',
                justifyContent:'center'
                }}
                mb={3}>
                    <Tabs value={router.query.id[1]}  aria-label="icon label tabs example">
                        <Tab icon={<SvgIcon><WifiIcon/></SvgIcon>} iconPosition={"end"} label="Wi-Fi" />
                        <Tab value={"discovery"} onClick={()=>{router.push(`/devices/${deviceID}/discovery`)}} icon={<SvgIcon><MagnifyingGlassIcon/></SvgIcon>} iconPosition={"end"} label="Discover Parameters" />
                        <Tab value={"msg"} onClick={()=>{router.push(`/devices/${deviceID}/msg`)}} icon={<SvgIcon><EnvelopeIcon/></SvgIcon>} iconPosition={"end"} label="Remote Messages" />
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