import React, { useState, useEffect } from 'react';
import Head from 'next/head';
import { Box, Container, Unstable_Grid2 as Grid } from '@mui/material';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { OverviewLatestOrders } from 'src/sections/overview/overview-latest-orders';
import { useAuth } from 'src/hooks/use-auth';

const Page = () => {
  const auth = useAuth();
  const [devices, setDevices] = useState([]);

  useEffect(() => {

    if (auth.user.token) {
      console.log("auth.user.token =", auth.user.token)
    }else{
      auth.user.token = localStorage.getItem("token")
    }

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var requestOptions = {
      method: 'GET',
      headers: myHeaders,
      redirect: 'follow'
    }

    fetch(process.env.NEXT_PUBLIC_REST_ENPOINT+'/device', requestOptions)
      .then(response => response.json())
      .then(json => setDevices(json))
      .catch(error => console.error('Error:', error));
  }, []);

  return (devices &&
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
          py: 8,
        }}
      >
        <Container maxWidth="xl" >
          <Grid
            container
            spacing={3}
          >
          </Grid>
          <OverviewLatestOrders
            orders={devices}
            sx={{ height: '100%' }}
          />
        </Container>
      </Box>
    </>
  )
}
Page.getLayout = (page) => (
    <DashboardLayout>
        {page}
    </DashboardLayout>
);

export default Page;
