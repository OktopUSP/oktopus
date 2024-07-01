import { use, useCallback, useEffect, useState } from 'react';
import {
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Stack,
  TextField,
  SvgIcon,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions, 
  Box,
  IconButton,
  Input,
  Typography,
  DialogContentText
} from '@mui/material';
import XMarkIcon from '@heroicons/react/24/outline/XMarkIcon';
import Check from '@heroicons/react/24/outline/CheckIcon';
import CircularProgress from '@mui/material/CircularProgress';
import Backdrop from '@mui/material/Backdrop';
import { useRouter } from 'next/router';
import ArrowsUpDownIcon from '@heroicons/react/24/solid/ArrowsUpDownIcon';
import PaperAirplane from '@heroicons/react/24/solid/PaperAirplaneIcon';

export const DevicesDiagnostic = () => {

    const router = useRouter()

    const [content, setContent] = useState(null)
    const [applyPing, setApplyPing] = useState(false)
    const [pingResponse, setPingResponse] = useState(null)
    const [progress, setProgress] = useState(0);

    //TODO: fixme
    // useEffect(()=>{
    //     let timeout = content?.number_of_repetitions.value * content?.timeout.value
    //     if (timeout <= 0) return;

    //     const increment = 100 / timeout ;// Calculate increment based on the timeout

    //     const interval = setInterval(() => {
    //         setProgress((prevProgress) => (
    //             prevProgress >= 100 ? 0 : prevProgress + increment
    //         ));
    //     }, 1000);

    //     return () => {
    //         clearInterval(interval);
    //       };
    // },[content])

    const fetchPingData = async () => {

        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", localStorage.getItem("token"));
      
        var requestOptions = {
          method: 'GET',
          headers: myHeaders,
          redirect: 'follow'
        };

        fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT}/device/${router.query.id[0]}/ping`, requestOptions)
        .then(response => {
            if (response.status === 401) {
                router.push("/auth/login")
            }
            return response.json()
        })
        .then(result => {
            console.log("ping content", result)
            setContent(result)
        })
        .catch(error => console.log('error', error));
    };

    useEffect(()=>{
        fetchPingData()
    },[])

    const handlePing = async () => {
        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", localStorage.getItem("token"));
      
        var requestOptions = {
          method: 'PUT',
          headers: myHeaders,
          redirect: 'follow',
          body: JSON.stringify(content)
        };

        fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT}/device/${router.query.id[0]}/ping`, requestOptions)
        .then(response => {
            if (response.status === 401) {
                router.push("/auth/login")
            }
            return response.json()
        })
        .then(result => {
            console.log("ping content", result)
            setProgress(100)
            setApplyPing(false)
            setPingResponse(result)
        })
        .catch(error => console.log('error', error));
    }


    return ( content &&
        <div>
        <Stack 
        direction="row" 
        spacing={2}   
        justifyContent="center" 
        alignItems="center"
        >
            <Card sx={{minWidth:"500px"}}>
                <CardHeader
                title="Ping"
                avatar={
                    <SvgIcon>
                        <ArrowsUpDownIcon/>
                    </SvgIcon>
                }
                />
                <CardContent>
                    <Stack spacing={3}>
                        <TextField
                            fullWidth
                            label="Host"
                            name="host"
                            type="text"
                            value={content.host.value}
                            onChange={(e) => setContent({...content, host: {value: e.target.value}})}
                        />
                        <TextField
                            fullWidth
                            label="Count"
                            name="count"
                            type="number"
                            value={content.number_of_repetitions.value}
                            onChange={(e) => setContent({...content, number_of_repetitions: {value: e.target.valueAsNumber}})}
                        />
                        <TextField
                            fullWidth
                            label="Timeout"
                            name="timeout"
                            type="number"
                            value={content.timeout.value}
                            onChange={(e) => setContent({...content, timeout: {value: e.target.valueAsNumber}})}
                        />  
                    </Stack>
                </CardContent>
                <CardActions  sx={{ justifyContent: 'flex-end' }}>
                    <Button
                        variant="contained"
                        color="primary"
                        endIcon={
                            <SvgIcon>
                                <PaperAirplane/>
                            </SvgIcon>
                        }
                        onClick={()=>{
                            setApplyPing(true)
                            handlePing()
                        }}
                    >
                        Ping
                    </Button>
                </CardActions>
            </Card>
        </Stack>
        {applyPing &&
        <Backdrop
        sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
        open={applyPing}
        >
        <CircularProgress /*variant="determinate" value={progress}*/ color="inherit"/>
        </Backdrop>
        }{ pingResponse &&
        <Dialog open={!applyPing && pingResponse}>
            <DialogTitle>
            <Box display="flex" alignItems="center">
                <Box flexGrow={1} >Ping Result</Box>
                <Box>
                    <IconButton >
                            <SvgIcon 
                            onClick={()=>{
                            setPingResponse(null)
                            }}
                            >
                            <XMarkIcon />
                        </SvgIcon>
                    </IconButton>
                </Box>
            </Box>
            </DialogTitle>
            <DialogContent>
                <DialogContentText>
                    <Stack spacing={2}>
                        {!pingResponse.failure_count && !pingResponse.success_count ?
                        <Typography sx={{display:"flex", justifyContent:"center"}}>
                            Error: {pingResponse}
                        </Typography>:<div>
                        <Stack spacing={2}>
                        <Typography sx={{display:"flex", justifyContent:"center", fontWeight:'fontWeightMedium'}}>
                            Ping Statistics for {content.host.value}
                        </Typography>
                        <Typography sx={{display:"flex", justifyContent:"center"}}>
                            Failure Count: {pingResponse.failure_count} | Success Count: {pingResponse.success_count}
                        </Typography>
                        <Typography>
                            Average Time: {pingResponse.average_rtt}s | Minimum Time: {pingResponse.minimum_rtt}s | Maximum Time: {pingResponse.maximum_rtt}s
                        </Typography></Stack></div>
                        }
                    </Stack>
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button
                    variant="contained"
                    color="primary"
                    onClick={()=>{
                        setPingResponse(null)
                    }}
                >
                    OK
                </Button>
            </DialogActions>
        </Dialog>}
        </div>
    )
};
