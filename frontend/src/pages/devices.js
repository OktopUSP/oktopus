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
  Pagination,
  CircularProgress
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
  const [deviceFound, setDeviceFound] = useState(true)
  const [pages, setPages] = useState(0);
  const [page, setPage] = useState(null);
  const [Loading, setLoading] = useState(true);


  useEffect(() => {
    setLoading(true)
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

    fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device`, requestOptions)
      .then(response => {
        if (response.status === 401)
          router.push("/auth/login")
        return response.json()
      })
      .then(json => {
        setPages(json.pages + 1)
        setPage(json.page +1)
        setDevices(json.devices)
        setLoading(false)
        return setDeviceFound(true)
      })
      .catch(error => {
        return console.error('Error:', error)
      });
  }, [auth.user]);

  const handleChangePage = (event, value) => {
    console.log("new page: ", value)
    setPage(value)
    fetchDevicePerPage(value)
  }

  const fetchDevicePerPage = async (p) => {
    setLoading(true)

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var requestOptions = {
      method: 'GET',
      headers: myHeaders,
      redirect: 'follow'
    }

    p = p - 1
    p = p.toString()

    fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device?page_number=+${p}`, requestOptions)
    .then(response => {
      if (response.status === 401)
        router.push("/auth/login")
      return response.json()
    })
    .then(json => {
      setDevices(json.devices)
      setLoading(false)
      return
    })
    .catch(error => {
      return console.error('Error:', error)
    });
  }

  const fetchDevicePerId = async (id) => {
    setLoading(true)
    setDeviceFound(true)
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var requestOptions = {
      method: 'GET',
      headers: myHeaders,
      redirect: 'follow'
    }

    if (id == ""){
      return fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device`, requestOptions)
      .then(response => {
        if (response.status === 401)
          router.push("/auth/login")
        return response.json()
      })
      .then(json => {
        setPages(json.pages + 1)
        setPage(json.page)
        setDevices(json.devices)
        setLoading(false)
        return setDeviceFound(true)
      })
      .catch(error => {
        return console.error('Error:', error)
      });
    }

    let response = await fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device?id=${id}`, requestOptions)
    if (response.status === 401)
      router.push("/auth/login")
    let json = await response.json()
    if (json.SN != undefined){
      setDevices([json])
      setDeviceFound(true)
      setLoading(false)
      setPages(1)
      setPage(1)
    }else{
      setDeviceFound(false)
      setDevices([])
      setPages(1)
      setPage(1)
      setLoading(false)
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
          ( !Loading ?
          <OverviewLatestOrders
            orders={devices}
            sx={{ height: '100%' }}
          /> : <CircularProgress></CircularProgress>
          )
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
