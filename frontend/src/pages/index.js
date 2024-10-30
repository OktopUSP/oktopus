import Head from 'next/head';
import React, { useEffect, useState } from 'react';
import { subDays, subHours } from 'date-fns';
import { 
  Box, 
  Container,
  CircularProgress, 
  Unstable_Grid2 as Grid } from '@mui/material';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { OverviewTasksProgress } from 'src/sections/overview/overview-tasks-progress';
import { OverviewTotalCustomers } from 'src/sections/overview/overview-total-customers';
import { OverviewTraffic } from 'src/sections/overview/overview-traffic';
import { useRouter } from 'next/router';

const now = new Date();

const Page = () => {

  const router = useRouter()

  const [generalInfo, setGeneralInfo] = useState(null)
  const [devicesStatus, setDevicesStatus] = useState([0,0])
  const [devicesCount, setDevicesCount] = useState(0)
  const [productClassLabels, setProductClassLabels] = useState(['-'])
  const [productClassValues, setProductClassValues] = useState(['0'])
  const [vendorLabels, setVendorLabels] = useState(['-'])
  const [vendorValues, setVendorValues] = useState([0])

  const fetchGeneralInfo = async () => {

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    var requestOptions = {
        method: 'GET',
        headers: myHeaders,
        redirect: 'follow',
    };

    let result = await (await fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/info/general`, requestOptions))
    if (result.status === 401){
    router.push("/auth/login")
    }else if (result.status != 200){
      console.log("Status:", result.status)
      let content = await result.json()
      console.log("Message:", content)
    }else{
      let content = await result.json()
      console.log("general info result:", content)
      let totalDevices = content.StatusCount.Offline + content.StatusCount.Online
      setDevicesCount(totalDevices)

      let onlinePercentage = ((content.StatusCount.Online * 100)/totalDevices)
      console.log("ONLINE AND OFFLINE:",onlinePercentage,100 - onlinePercentage)
      
      if(Number.isInteger(onlinePercentage)){
        setDevicesStatus([onlinePercentage, 100 - onlinePercentage])
      }else{
        onlinePercentage = Number(onlinePercentage.toFixed(1))
        let offlinePercentage = 100 - onlinePercentage
        setDevicesStatus([onlinePercentage, Number(offlinePercentage.toFixed(1))])
      }

      let prodClassLabels = []
      let prodClassValues = []
      let prodClassValue = 0
      
      content.ProductClassCount?.map((p)=>{
        if (p.productClass === ""){
          prodClassLabels.push("unknown")
        }else{
          prodClassLabels.push(p.productClass)
        }
        prodClassValue += p.count
      })

      content.ProductClassCount?.map((p)=>{
        let percentageValue = p.count * 100 / prodClassValue
        if (Number.isInteger(percentageValue)){
          prodClassValues.push(percentageValue)
        }else{
          prodClassValues.push(Number(percentageValue.toFixed(1)))
        }
      })
      
      setProductClassLabels(prodClassLabels)
      setProductClassValues(prodClassValues)
      console.log("productClassLabels:", prodClassLabels)
      console.log("productClassValues:", productClassValues)

      let vLabels = []
      let vValues = []
      let vValue = 0
      content.VendorsCount?.map((p)=>{
        if (p.vendor === ""){
          vLabels.push("unknown")
        }else{
          vLabels.push(p.vendor)
        }
        vValue = vValue + p.count
      })

      content.VendorsCount?.map((p)=>{
        let percentageValue = p.count * 100 / vValue
        if (Number.isInteger(percentageValue)){
          vValues.push(percentageValue)
        }else{
          vValues.push(Number(percentageValue.toFixed(1)))
        }
      })

      setVendorLabels(vLabels)
      setVendorValues(vValues)

      console.log("vendorLabels:", vLabels)
      console.log("vendorValues:", vValues)

      setGeneralInfo(content)
    }

  }

  useEffect(()=>{
    fetchGeneralInfo()
  },[])
  
  return(generalInfo ?
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
        py: 8
      }}
    >
      <Container maxWidth="xl">
        <Grid
          container
          spacing={3}
        >
          <Grid
            xs={12}
            sm={6}
            lg={3}
          >
            <OverviewTotalCustomers
              //difference={16}
              positive={false}
              sx={{ height: '100%' }}
              value={devicesCount.toString()}
            />
          </Grid>
          <Grid
            xs={12}
            sm={6}
            lg={3}
          >
            <OverviewTasksProgress
              sx={{ height: '100%' }}
              mtp={"STOMP Connection"}
              type={"stomp"}
              value={generalInfo.StompRtt}
            />
          </Grid>
          <Grid
            xs={12}
            sm={6}
            lg={3}
          >
            <OverviewTasksProgress
              sx={{ height: '100%' }}
              mtp={"MQTT Connection"}
              type={"mqtt"}
              value={generalInfo.MqttRtt}
            />
          </Grid>
          <Grid
            xs={12}
            sm={6}
            lg={3}
          >
          <OverviewTasksProgress
              sx={{ height: '100%' }}
              mtp={"Websockets Connection"}
              type={"websocket"}
              value={generalInfo.WebsocketsRtt}
            />
          </Grid>
          <Grid
            xs={12}
            lg={4}
          >
          <OverviewTraffic
          chartSeries={vendorValues}
          labels={vendorLabels}
          sx={{ height: '100%' }}
          title={'Vendors'}
          />
          </Grid>
          <Grid
            xs={12}
            lg={4}
          >
          <OverviewTraffic
          chartSeries={devicesStatus}
          labels={['Online', 'Offline']}
          sx={{ height: '100%' }}
          title={'Status'}
          />
          </Grid>
          <Grid
            xs={12}
            //md={6}
            lg={4}
          >
            <OverviewTraffic
              chartSeries={productClassValues}
              labels={productClassLabels}
              sx={{ height: '100%' }}
              title={'Devices Type'}
            />
          </Grid>
          <Grid
            xs={12}
            md={6}
            lg={4}
          >
          </Grid>
        </Grid>
      </Container>
    </Box>
  </>:    <Box sx={{display:'flex',justifyContent:'center'}}>
        <CircularProgress color="inherit" />
    </Box>)
};

Page.getLayout = (page) => (
  <DashboardLayout>
    {page}
  </DashboardLayout>
);

export default Page;

/*
            <OverviewSales
              chartSeries={[
                {
                  name: 'This year',
                  data: [18, 16, 5, 8, 3, 14, 14, 16, 17, 19, 18, 20]
                },
                {
                  name: 'Last year',
                  data: [12, 11, 4, 6, 2, 9, 9, 10, 11, 12, 13, 13]
                }
              ]}
              sx={{ height: '100%' }}
            />
                        <OverviewLatestProducts
              products={[
                {
                  id: '5ece2c077e39da27658aa8a9',
                  image: '/assets/products/product-1.png',
                  name: 'Healthcare Erbology',
                  updatedAt: subHours(now, 6).getTime()
                },
                {
                  id: '5ece2c0d16f70bff2cf86cd8',
                  image: '/assets/products/product-2.png',
                  name: 'Makeup Lancome Rouge',
                  updatedAt: subDays(subHours(now, 8), 2).getTime()
                },
                {
                  id: 'b393ce1b09c1254c3a92c827',
                  image: '/assets/products/product-5.png',
                  name: 'Skincare Soja CO',
                  updatedAt: subDays(subHours(now, 1), 1).getTime()
                },
                {
                  id: 'a6ede15670da63f49f752c89',
                  image: '/assets/products/product-6.png',
                  name: 'Makeup Lipstick',
                  updatedAt: subDays(subHours(now, 3), 3).getTime()
                },
                {
                  id: 'bcad5524fe3a2f8f8620ceda',
                  image: '/assets/products/product-7.png',
                  name: 'Healthcare Ritual',
                  updatedAt: subDays(subHours(now, 5), 6).getTime()
                }
              ]}
              sx={{ height: '100%' }}
            />
          </Grid>
          <Grid
            xs={12}
            md={12}
            lg={8}
          >
            <OverviewLatestOrders
              orders={[
                {
                  id: 'f69f88012978187a6c12897f',
                  ref: 'DEV1049',
                  amount: 30.5,
                  customer: {
                    name: 'Ekaterina Tankova'
                  },
                  createdAt: 1555016400000,
                  status: 'pending'
                },
                {
                  id: '9eaa1c7dd4433f413c308ce2',
                  ref: 'DEV1048',
                  amount: 25.1,
                  customer: {
                    name: 'Cao Yu'
                  },
                  createdAt: 1555016400000,
                  status: 'delivered'
                },
                {
                  id: '01a5230c811bd04996ce7c13',
                  ref: 'DEV1047',
                  amount: 10.99,
                  customer: {
                    name: 'Alexa Richardson'
                  },
                  createdAt: 1554930000000,
                  status: 'refunded'
                },
                {
                  id: '1f4e1bd0a87cea23cdb83d18',
                  ref: 'DEV1046',
                  amount: 96.43,
                  customer: {
                    name: 'Anje Keizer'
                  },
                  createdAt: 1554757200000,
                  status: 'pending'
                },
                {
                  id: '9f974f239d29ede969367103',
                  ref: 'DEV1045',
                  amount: 32.54,
                  customer: {
                    name: 'Clarke Gillebert'
                  },
                  createdAt: 1554670800000,
                  status: 'delivered'
                },
                {
                  id: 'ffc83c1560ec2f66a1c05596',
                  ref: 'DEV1044',
                  amount: 16.76,
                  customer: {
                    name: 'Adam Denisov'
                  },
                  createdAt: 1554670800000,
                  status: 'delivered'
                }
              ]}
              sx={{ height: '100%' }}
            />
*/ 
