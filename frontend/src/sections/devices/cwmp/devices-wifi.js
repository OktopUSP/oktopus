import { useCallback, useEffect, useState } from 'react';
import {
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Divider,
  Stack,
  TextField,
  InputLabel,
  MenuItem, 
  Select,
  FormControl,
  SvgIcon,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions, 
  Box,
  IconButton,
  Icon,
  SnackbarContent,
  Snackbar,
  Checkbox,
  FormControlLabel,
  useTheme,
} from '@mui/material';
import XMarkIcon from '@heroicons/react/24/outline/XMarkIcon';
import Check from '@heroicons/react/24/outline/CheckIcon';
//import ExclamationTriangleIcon from '@heroicons/react/24/solid/ExclamationTriangleIcon';
import CircularProgress from '@mui/material/CircularProgress';
import Backdrop from '@mui/material/Backdrop';
import { useRouter } from 'next/router';
import GlobeAltIcon from '@heroicons/react/24/outline/GlobeAltIcon';

export const DevicesWiFi = () => {

    const theme = useTheme();
    const router = useRouter()

    const [content, setContent] = useState([])
    const [applyContent, setApplyContent] = useState([])
    const [apply, setApply] = useState(false)

    const [errorModal, setErrorModal] = useState(false)
    const [errorModalText, setErrorModalText] = useState("")

    const fetchWifiData = async () => {

        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", localStorage.getItem("token"));
      
        var requestOptions = {
          method: 'GET',
          headers: myHeaders,
          redirect: 'follow'
        };

        fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT}/device/${router.query.id[0]}/wifi`, requestOptions)
        .then(response => {
            if (response.status === 401) {
                router.push("/auth/login")
            }
            return response.json()
        })
        .then(result => {
            console.log("wifi content", result)
            result.map((item) => {
                let contentToApply = {
                    hasChanges: false,
                    path: item.path,
                }
                setApplyContent(oldValue => [...oldValue, contentToApply])
            })
            setContent(result)
        })
        .catch(error => console.log('error', error));
    };

    useEffect(()=>{
        fetchWifiData()
    },[])

    return (<div>
        <Stack 
        direction="row" 
        spacing={2}   
        justifyContent="center" 
        alignItems="center"
        >
            {content.length > 1 ?
                (content.map((item, index) => {
                    return (
                        <Card key={index}>
                            <CardHeader
                                title={item.name.value}
                                avatar={
                                    <SvgIcon>
                                            <GlobeAltIcon/>
                                    </SvgIcon>
                                }
                            />
                            <CardContent>
                                <Stack spacing={3}>
                                    { item.enable.value != null &&
                                    <FormControlLabel control={<Checkbox defaultChecked={item.enable.value == 1 ? true : false} 
                                    onChange={(e) => {
                                        let enable = e.target.value == 1 ? "1" : "0"
                                        applyContent[index].hasChanges = true
                                        applyContent[index].enable = {
                                            value : enable
                                        }
                                        setApplyContent([...applyContent])
                                        item.enable.value = enable
                                    }}/>}
                                    label="Enabled" />}
                                    {item.ssid.value != null && <TextField
                                        fullWidth
                                        label="SSID"
                                        value={item.ssid.value}
                                        disabled={!item.ssid.writable}
                                        onChange={(e) => {
                                            applyContent[index].hasChanges = true
                                            applyContent[index].ssid = {
                                                value : e.target.value
                                            }
                                            setApplyContent([...applyContent])
                                            item.ssid.value = e.target.value
                                        }}
                                    />}
                                    {item.securityCapabilities &&
                                    <TextField
                                        fullWidth
                                        label="Encryption"
                                        value={""}
                                    />}
                                    {item.password.value != null &&
                                    <TextField
                                        fullWidth
                                        label="Password"
                                        disabled={!item.password.writable}
                                        value={item.password.value}
                                    />}
                                    {item.standard.value != null &&
                                    <TextField
                                        fullWidth
                                        label="Standard"
                                        disabled={!item.standard.writable}
                                        value={item.standard.value}
                                    />}
                                </Stack>
                                <CardActions sx={{display:"flex", justifyContent:"flex-end"}}>
                                    <Button 
                                        variant="contained" 
                                        disabled={!applyContent[index].hasChanges}
                                        endIcon={<SvgIcon><Check /></SvgIcon>} 
                                        onClick={
                                            ()=>{
                                                setApply(true)
                                                var myHeaders = new Headers();
                                                myHeaders.append("Content-Type", "application/json");
                                                myHeaders.append("Authorization", localStorage.getItem("token"));
                                                
                                                delete applyContent[index].hasChanges
                                                let contentToApply = [applyContent[index]]
                                                console.log("contentToApply: ", contentToApply)
                                                var data = JSON.stringify(contentToApply);
                                                var requestOptions = {
                                                    method: 'PUT',
                                                    headers: myHeaders,
                                                    body: data,
                                                    redirect: 'follow'
                                                };
                                                fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT}/device/${router.query.id[0]}/wifi`, requestOptions)
                                                .then(response => {
                                                    if (response.status === 401) {
                                                        router.push("/auth/login")
                                                    }
                                                    if (response.status == 500) {
                                                        setErrorModal(true)
                                                    }
                                                    return response.json()
                                                })
                                                .then(result => {
                                                    if (errorModal) {
                                                        setErrorModalText(result)
                                                    }
                                                    setApply(false)
                                                    if (result == 1) {
                                                        setErrorModalText("This change could not be applied, or It's gonna be applied later on")
                                                        setErrorModal(true)
                                                        //TODO: fetch wifi data again
                                                    }
                                                })
                                                .catch(error => console.log('error', error));
                                            }
                                        }
                                        sx={{mt:'25px', mb:'-15px'}}
                                        >
                                        Apply
                                    </Button>
                                </CardActions>
                            </CardContent>
                        </Card>
                    )
                })):
                <CircularProgress />
            }
            {/* <Card>
                <CardHeader
                    title="2.4GHz"
                    avatar={
                        <SvgIcon>
                                <GlobeAltIcon/>
                        </SvgIcon>
                    }
                />
                <CardContent>
                    <Stack spacing={3}>
                        <FormControlLabel control={<Checkbox defaultChecked />} label="Enabled" />
                        <TextField
                            fullWidth
                            label="SSID"
                            value="wlan0"
                        />
                        <TextField
                            fullWidth
                            label="Encryption"
                            value="WPA2-PSK"
                        />
                        <TextField
                            fullWidth
                            label="Key"
                            value="password"
                        />
                    </Stack>
                </CardContent>
            </Card>
            <Card>
                <CardHeader
                    title="5GHz"
                    avatar={
                        <SvgIcon>
                                <GlobeAltIcon/>
                        </SvgIcon>
                    }
                />
                <CardContent>
                    <Stack spacing={4}>
                        <FormControlLabel control={<Checkbox defaultChecked />} label="Enabled" />
                        <TextField
                            fullWidth
                            label="SSID"
                            value="wlan0"
                        />
                            <InputLabel id="demo-simple-select-standard-label">Security</InputLabel>
                            <Select
                            labelId="demo-simple-select-standard-label"
                            id="demo-simple-select-standard"
                            value={"WPA2-PSK"}
                            //onChange={handleChange}
                            label="Security"
                            >
                            <MenuItem value={30}>Open</MenuItem>
                            <MenuItem value={"WPA2-PSK"}>WPA2-PSKnp</MenuItem>
                            <MenuItem value={20}>WPA3</MenuItem>
                            </Select>
                        </FormControl>
                        <TextField
                            fullWidth
                            label="Key"
                            value="password"
                        />
                    </Stack>
                    <CardActions sx={{display:"flex", justifyContent:"flex-end"}}>
                    <Button 
                        variant="contained" 
                        endIcon={<SvgIcon><Check /></SvgIcon>} 
                       // onClick={}
                        sx={{mt:'25px', mb:'-15px'}}
                        >
                        Apply
                    </Button>
                    </CardActions>
                </CardContent>
            </Card> */}
            {/* <Snackbar
                open={errorToApplyContent}
                TransitionComponent={"Slide"}
                color="red"
                onClose={() => setErrorToApplyContent(false)}
                autoHideDuration={1200} 
            >
                <SnackbarContent style={{
                    backgroundColor:theme.palette.warning.main,
                    }}
                    message={
                        <div style={{display:"flex"}}>
                            <SvgIcon>
                                <ExclamationTriangleIcon />
                            </SvgIcon>
                            <div style={{margin: "5px"}}></div>
                            <span id="client-snackbar">
                                No changes to apply
                            </span>
                        </div>
                    }
                />
            </Snackbar> */}
        </Stack>
        <Backdrop
            sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
            open={apply}
        >
            <CircularProgress color="inherit" />
        </Backdrop>
        <Dialog open={errorModal && errorModalText != ""} 
                    slotProps={{ backdrop: { style: { backgroundColor: 'rgba(255,255,255,0.5)' } } }}
                    fullWidth={ true } 
                    maxWidth={"md"}   
                    scroll={"paper"}
                    aria-labelledby="scroll-dialog-title"
                    aria-describedby="scroll-dialog-description" 
                >
                <DialogTitle id="scroll-dialog-title">
                <Box display="flex" alignItems="center">
                    <Box flexGrow={1} >Response</Box>
                    <Box>
                        <IconButton onClick={()=>{
                                    setErrorModalText("")
                                    setErrorModal(false)
                                }}>
                                <SvgIcon 
                                >
                                < XMarkIcon/>
                            </SvgIcon>
                        </IconButton>
                    </Box>
                </Box>
                </DialogTitle>    
                    <DialogContent dividers={scroll === 'paper'}>
                    <DialogContentText id="scroll-dialog-description"tabIndex={-1}>
                    <pre style={{color: 'black'}}>
                        {errorModalText}
                    </pre>
                    </DialogContentText>
                    </DialogContent>
                    <DialogActions>
                    <Button onClick={()=>{
                        setErrorModalText("")
                        setErrorModal(false)
                    }}>OK</Button>
                    </DialogActions>
                </Dialog>
        </div>
  );
};
