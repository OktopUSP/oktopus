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
import CubeTransparentIcon from '@heroicons/react/24/outline/CubeTransparentIcon';
import MapPin from '@heroicons/react/24/outline/MapPinIcon';
import ArrowTrendingUp from '@heroicons/react/24/outline/ArrowTrendingUpIcon';

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
                router.replace(`/devices/usp/${deviceID}/discovery`)
        }
    }
  
    return(
    <>
        <Head>
            <title>
                Oktopus | Controller
            </title>
        </Head>
        <Box
            component="main"
            sx={{
                flexGrow: 1,
                py: 0,
            }}
        >
            <Container maxWidth="xg" >
                <Stack spacing={3} mb={3}>
                    <Breadcrumbs separator="›" aria-label="breadcrumb"ml={10}>
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
                }}>
                    <Tabs value={router.query.id[1]}  aria-label="icon label tabs example" variant='scrollable'>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><WifiIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Wi-Fi" 
                        style={{cursor:"default", opacity: 0.5}}
                        value={"wifi"}/>
                        </Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                            icon={<SvgIcon><SignalIcon/></SvgIcon>} 
                            iconPosition={"end"} 
                            label="Site Survey" 
                            style={{cursor:"default", opacity: 0.5}} 
                            value={"site-survey"} 
                        />
                        </Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><DevicePhoneMobile/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Connected Devices" 
                        style={{cursor:"default", opacity: 0.5}}
                        value={"connected-devices"} 
                        />
                        </Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><WrenchScrewDriverIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Diagnostic" 
                        style={{cursor:"default", opacity: 0.5}}
                        value={"diagnostic"} /></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><ServerStackIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Ports" 
                        style={{cursor:"default", opacity: 0.5}}
                        value={"ports"} /></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><ArrowTrendingUp/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Historic" 
                        style={{cursor:"default", opacity: 0.5}}
                        value={"historic"} /></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><CubeTransparentIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="LCM" 
                        style={{cursor:"default", opacity: 0.5}}
                        value={"lcm"} /></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><MapPin/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Location" 
                        style={{cursor:"default", opacity: 0.5}}
                        value={"location"} /></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><CommandLineIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Actions" 
                        style={{cursor:"default", opacity: 0.5}}
                        value={"actions"} /></Tooltip>
                        <Tab 
                        value={"discovery"} 
                        onClick={()=>{router.push(`/devices/usp/${deviceID}/discovery`)}}
                        icon={<SvgIcon><MagnifyingGlassIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Parameters" />
                        <Tab 
                        value={"msg"} 
                        onClick={()=>{router.push(`/devices/usp/${deviceID}/msg`)}}
                        icon={<SvgIcon><EnvelopeIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Messages" />
                    </Tabs>
                </Box>
                </Stack>
            </Container>
            <Container maxWidth="lg">
                <Stack spacing={3}>
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