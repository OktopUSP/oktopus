import Head from 'next/head';
import { Box, Stack, Typography, Container, Unstable_Grid2 as Grid,
Tab, 
Tabs,
SvgIcon,
Breadcrumbs,
Link, 
Tooltip,
IconButton} from '@mui/material';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { useRouter } from 'next/router';
import { DevicesRPC } from 'src/sections/devices/cwmp/devices-rpc';
import EnvelopeIcon from '@heroicons/react/24/outline/EnvelopeIcon';
import MagnifyingGlassIcon from '@heroicons/react/24/solid/MagnifyingGlassIcon';
import WifiIcon from '@heroicons/react/24/solid/WifiIcon';
import ServerStackIcon from '@heroicons/react/24/outline/ServerStackIcon';
import SignalIcon from '@heroicons/react/24/solid/SignalIcon';
import DevicePhoneMobile from '@heroicons/react/24/solid/DevicePhoneMobileIcon';
import WrenchScrewDriverIcon from '@heroicons/react/24/outline/WrenchScrewdriverIcon';
import CommandLineIcon from '@heroicons/react/24/outline/CommandLineIcon';
import { DevicesWiFi } from 'src/sections/devices/cwmp/devices-wifi';
import ArrowTrendingUpIcon from '@heroicons/react/24/outline/ArrowTrendingUpIcon';
import DocumentTextIcon from '@heroicons/react/24/outline/DocumentTextIcon';
import MapPinIcon from '@heroicons/react/24/outline/MapPinIcon';


const Page = () => {
    const router = useRouter()

    const deviceID = router.query.id[0]
    const section = router.query.id[1]

    const sectionHandler = () => {
        switch(section){
            case "msg":
                return <DevicesRPC/>
            /* case "wifi":
                return <DevicesWiFi/> */
            default:
                router.replace(`/devices/cwmp/${deviceID}/msg`)
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
            <Container maxWidth="xg">
            <Stack spacing={3} mb={3}>
                <Breadcrumbs separator="›" aria-label="breadcrumb" ml={10}>
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
                    <Tabs value={router.query.id[1]}  aria-label="icon label tabs example" variant='scrollable'>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><WifiIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Wi-Fi" 
                        value={"wifi"}
                        style={{opacity:"0.5", cursor:"default"}}/></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><SignalIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Site Survey" 
                        value={"site-survey"} 
                        style={{opacity:"0.5", cursor:"default"}}/></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><DevicePhoneMobile/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Connected Devices" 
                        style={{opacity:"0.5", cursor:"default"}}
                        value={"connected-devices"} 
                        /></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><WrenchScrewDriverIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Diagnostic" 
                        value={"diagnostic"} 
                        style={{opacity:"0.5", cursor:"default"}}/></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><ServerStackIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Ports" 
                        style={{opacity:"0.5", cursor:"default"}} 
                        value={"ports"} /></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><ArrowTrendingUpIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Historic" 
                        value={"historic"} 
                        style={{opacity:"0.5", cursor:"default"}}/></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><CommandLineIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Actions" 
                        style={{opacity:"0.5", cursor:"default"}} 
                        value={"actions"} /></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><DocumentTextIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Logs" 
                        style={{opacity:"0.5", cursor:"default"}} 
                        value={"logs"} /></Tooltip>
                        <Tooltip placement="bottom">
                        <Tab 
                        icon={<SvgIcon><MapPinIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Location" 
                        style={{opacity:"0.5", cursor:"default"}} 
                        value={"location"} /></Tooltip>
                        <Tab 
                        value={"msg"} 
                        onClick={()=>{router.push(`/devices/cwmp/${deviceID}/msg`)}} 
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