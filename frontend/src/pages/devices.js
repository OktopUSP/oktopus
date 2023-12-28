import React, { useState, useEffect } from 'react';
import Head from 'next/head';
import { 
  Box, 
  Container, 
  Unstable_Grid2 as Grid,
  Card, 
  OutlinedInput,
  InputAdornment,
  SvgIcon,
  Stack,
  Pagination
} from '@mui/material';
import MagnifyingGlassIcon from '@heroicons/react/24/solid/MagnifyingGlassIcon';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { OverviewLatestOrders } from 'src/sections/overview/overview-latest-orders';
import { useAuth } from 'src/hooks/use-auth';
import { useRouter } from 'next/router';

const Page = () => {
  const router = useRouter()
  const auth = useAuth();
  const [devices, setDevices] = useState([]);
  const [deviceFound, setDeviceFound] = useState(false)
  const [pages, setPages] = useState(0);
  const [page, setPage] = useState(null);


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
      .then(response => {
        if (response.status === 401)
          router.push("/auth/login")
        return response.json()
      })
      .then(json => {
        setPages(json.pages + 1)
        setPage(json.page)
        setDevices(json.devices)
        return setDeviceFound(true)
      })
      .catch(error => {
        return console.error('Error:', error)
      });
  }, [auth.user]);

  const handleChangePage = (e) => {
    console.log("new page: ", e.target.value)
    //TODO: Handle page change
  }

  const fetchDevicePerId = async (id) => {

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var requestOptions = {
      method: 'GET',
      headers: myHeaders,
      redirect: 'follow'
    }

    if (id == ""){
      fetch(process.env.NEXT_PUBLIC_REST_ENPOINT+'/device', requestOptions)
      .then(response => {
        if (response.status === 401)
          router.push("/auth/login")
        return response.json()
      })
      .then(json => {
        setPages(json.pages + 1)
        setPage(json.page)
        setDevices(json.devices)
        return setDeviceFound(true)
      })
      .catch(error => {
        return console.error('Error:', error)
      });
    }

    let response = await fetch(process.env.NEXT_PUBLIC_REST_ENPOINT+'/device?id='+id, requestOptions)
    if (response.status === 401)
      router.push("/auth/login")
    let json = await response.json()
    if (json.device != undefined){
      setDevices({"devices":[
        json.device
      ]})
      setPages(1)
      setPage(1)
    }else{
      setDeviceFound(false)
      setDevices([])
      setPages(1)
      setPage(1)
    }

  }

  return (
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
        <Stack spacing={3}>
          <Grid 
          container
          >
          <Grid xs={8}>
          </Grid>
          <OutlinedInput
            xs={4}
            defaultValue=""
            fullWidth
            placeholder="Search Device"
            onKeyDownCapture={(e) => {
              if (e.key === 'Enter') {
                console.log("Fetch devices per id: ", e.target.value)
                fetchDevicePerId(e.target.value)
              }
            }}
            startAdornment={(
              <InputAdornment position="start">
                <SvgIcon
                  color="action"
                  fontSize="small"
                >
                  <MagnifyingGlassIcon />
                </SvgIcon>
              </InputAdornment>
            )}
            sx={{ maxWidth: 500 }}
          />
        </Grid>
        {deviceFound ?
          <OverviewLatestOrders
            orders={devices}
            sx={{ height: '100%' }}
          />
        :
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'center'
          }}
        >
        <p>Device Not Found</p>
        </Box>
        }
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'center'
          }}
        >
          {pages ? <Pagination
            count={pages}
            size="small"
            page={page}
            onChange={handleChangePage}
          />: null} 
          {/* //TODO: show loading */}
        </Box>
        </Stack>
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
