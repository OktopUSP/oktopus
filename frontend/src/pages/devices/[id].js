import Head from 'next/head';
import { Box, Stack, Typography, Container, Unstable_Grid2 as Grid } from '@mui/material';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { useRouter } from 'next/router';
import { DevicesRPC } from 'src/sections/devices/devices-rpc';

const Page = () => {
    const router = useRouter()
    const { id } = router.query
  

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
                <Typography variant="h4">
                    RPC
                </Typography>
                {/*<SettingsNotifications />*/}
                < DevicesRPC />
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