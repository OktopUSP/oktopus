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
            default:
                router.push(`/devices/cwmp/${deviceID}/wifi`)
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
                    justifyContent:'center',
                }}>                                        
                    <Tabs value={router.query.id[1]}  aria-label="icon label tabs example">
                        <Tab 
                        icon={<SvgIcon><WifiIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Wi-Fi" 
                        onClick={()=>{router.push(`/devices/cwmp/${deviceID}/wifi`)}} 
                        value={"wifi"}/>
                        <Tab 
                        icon={<SvgIcon><SignalIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Site Survey" 
                        onClick={()=>{router.push(`/devices/cwmp/${deviceID}/site-survey`)}} 
                        value={"site-survey"} 
                        disabled={true}/>
                        <Tab 
                        icon={<SvgIcon><DevicePhoneMobile/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Connected Devices" 
                        disabled={true}
                        onClick={()=>{router.push(`/devices/cwmp/${deviceID}/connected-devices`)}} 
                        value={"connected-devices"} 
                        />
                        <Tab 
                        icon={<SvgIcon><WrenchScrewDriverIcon/></SvgIcon>} 
                        iconPosition={"end"} 
                        label="Diagnostic" 
                        onClick={()=>{router.push(`/devices/cwmp/${deviceID}/diagnostic`)}} 
                        value={"diagnostic"} 
                        disabled={true}/>
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
                        value={"msg"} 
                        onClick={()=>{router.push(`/devices/cwmp/${deviceID}/msg`)}} 
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