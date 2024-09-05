import Head from 'next/head';
import { Box, Stack, Typography, Container, Unstable_Grid2 as Grid,
Tab, 
Tabs,
SvgIcon,
Breadcrumbs,
Link, 
CircularProgress,
Tooltip} from '@mui/material';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { useRouter } from 'next/router';
import { DevicesRPC } from 'src/sections/devices/usp/devices-rpc';
import { DevicesDiscovery } from 'src/sections/devices/usp/devices-discovery';
import EnvelopeIcon from '@heroicons/react/24/outline/EnvelopeIcon';
import MagnifyingGlassIcon from '@heroicons/react/24/solid/MagnifyingGlassIcon';
import WifiIcon from '@heroicons/react/24/outline/WifiIcon';
import ServerStackIcon from '@heroicons/react/24/outline/ServerStackIcon';
import { useState } from 'react';
import SignalIcon from '@heroicons/react/24/solid/SignalIcon';
import DevicePhoneMobile from '@heroicons/react/24/solid/DevicePhoneMobileIcon';
import WrenchScrewDriverIcon from '@heroicons/react/24/outline/WrenchScrewdriverIcon';
import CommandLineIcon from '@heroicons/react/24/outline/CommandLineIcon';

const Page = () => {
    const router = useRouter()

    const deviceID = router.query.id[0]
    const section = router.query.id[1]
    
    const [loading, setLoading] = useState(true)

    const sectionHandler = () => {
        switch(section){
            case "msg":
                return <DevicesRPC/>
            case "discovery":
                return <DevicesDiscovery/>
            default:
                router.push(`/devices/usp/${deviceID}/discovery`)
        }
    }
  
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
                justifyContent:'center'
                }}
                mb={3}>
                    <Tabs value={router.query.id[1]}  aria-label="icon label tabs example">
                        <Tab 
                        icon={<SvgIcon><WifiIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Wi-Fi" 
                        onClick={()=>{router.push(`/devices/usp/${deviceID}/wifi`)}}
                        disabled={true} 
                        value={"wifi"}/>
                        <Tab 
                            icon={<SvgIcon><SignalIcon/></SvgIcon>} 
                            iconPosition={"end"} 
                            label="Site Survey" 
                            style={{opacity:"0.5"}}
                            onClick={()=>{router.push(`/devices/usp/${deviceID}/site-survey`)}}
                            disabled={true} 
                            value={"site-survey"} 
                        />
                        <Tab 
                        icon={<SvgIcon><DevicePhoneMobile/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Connected Devices" 
                        onClick={()=>{router.push(`/devices/usp/${deviceID}/connected-devices`)}}
                        disabled={true} 
                        value={"connected-devices"} 
                        />
                        <Tab 
                        icon={<SvgIcon><WrenchScrewDriverIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Diagnostic" 
                        onClick={()=>{router.push(`/devices/usp/${deviceID}/diagnostic`)}}
                        disabled={true} 
                        value={"diagnostic"} />
                        <Tab 
                        icon={<SvgIcon><ServerStackIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Ports" 
                        onClick={()=>{router.push(`/devices/usp/${deviceID}/ports`)}}
                        disabled={true} 
                        value={"ports"} />
                        <Tab 
                        icon={<SvgIcon><CommandLineIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Actions" 
                        onClick={()=>{router.push(`/devices/usp/${deviceID}/actions`)}}
                        disabled={true} 
                        value={"actions"} />
                        <Tab 
                        value={"discovery"} 
                        onClick={()=>{router.push(`/devices/usp/${deviceID}/discovery`)}}
                        icon={<SvgIcon><MagnifyingGlassIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Discover Parameters" />
                        <Tab 
                        value={"msg"} 
                        onClick={()=>{router.push(`/devices/usp/${deviceID}/msg`)}}
                        icon={<SvgIcon><EnvelopeIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Remote Messages" />
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