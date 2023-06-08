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
  IconButton
} from '@mui/material';
import XMarkIcon from '@heroicons/react/24/outline/XMarkIcon';
import PaperAirplane from '@heroicons/react/24/solid/PaperAirplaneIcon';
import CircularProgress from '@mui/material/CircularProgress';
import Backdrop from '@mui/material/Backdrop';
import { useRouter } from 'next/router';


export const DevicesRPC = () => {

const router = useRouter()

const [open, setOpen] = useState(false);
const [scroll, setScroll] = useState('paper');
const [answer, setAnswer] = useState(false)
const [content, setContent] = useState('')
const [age, setAge] = useState(2);

const [value, setValue] = useState(`{
  "param_paths": [
      "Device.WiFi.SSID.[Name==wlan0].",
      "Device.IP.Interface.*.Alias",
      "Device.DeviceInfo.FirmwareImage.*.Alias",
      "Device.IP.Interface.1.IPv4Address.1.IPAddress"
  ],
  "max_depth": 2
}`)

const handleClose = () => {
  setOpen(false);
};
const handleOpen = () => {
  setOpen(true);
  var myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "<token>");

 var raw = value

  var requestOptions = {
    method: 'PUT',
    headers: myHeaders,
    body: raw,
    redirect: 'follow'
  };

  var method;

  switch(age) {
    case 1:
      method="add"
      break;
    case 2:
      method="get"
      break;
    case 3:
      method="set"
      break;
    case 4:
      method="del"
    break;
  }

 
  fetch(`${process.env.NEXT_PUBLIC_REST_ENPOINT}/device/${router.query.id}/${method}`, requestOptions)
    .then(response => response.text())
    .then(result => {
      setOpen(false)
      setAnswer(true)
      let teste =  JSON.stringify(JSON.parse(result), null, 2)
      console.log(teste)
      setContent(teste)
    })
    .catch(error => console.log('error', error));
  };

  const handleChangeRPC = (event) => {
    setAge(event.target.value);
    switch(event.target.value) {
      case 1:
        setValue(`{
          "allow_partial": true,
          "create_objs": [
              {
                  "obj_path": "Device.IP.Interface.",
                  "param_settings": [
                      {
                          "param": "Alias",
                          "value": "test",
                          "required": true
                      }
                  ]
              }
          ]
      }`)
        break;
      case 2:
        setValue(`{
          "param_paths": [
              "Device.WiFi.SSID.[Name==wlan0].",
              "Device.IP.Interface.*.Alias",
              "Device.DeviceInfo.FirmwareImage.*.Alias",
              "Device.IP.Interface.1.IPv4Address.1.IPAddress"
          ],
          "max_depth": 2
      }`)
        break;
      case 3:
        setValue(`
        {
          "allow_partial":true,
          "update_objs":[
              {
                  "obj_path":"Device.IP.Interface.[Alias==pamonha].",
                  "param_settings":[
                      {
                      "param":"Alias",
                      "value":"goiaba",
                      "required":true
                      }
                  ]
              }
          ]
      }`)
        break;
      case 4:
        setValue(`{
          "allow_partial": true,
          "obj_paths": [
              "Device.IP.Interface.3."
          ]
      }`)
        break;
      default:
        // code block
    }
  };

  const handleChange = (event) => {
    setValue(event.target.value);
  };

  const handleSubmit = useCallback(
    (event) => {
      event.preventDefault();
    },
    []
  );

  return (
    <form onSubmit={handleSubmit}>
      <Card>
        <CardActions sx={{ justifyContent: 'flex-end'}}>
            <FormControl sx={{width:'100px'}}>
                <Select
                    labelId="demo-simple-select-standard-label"
                    id="demo-simple-select-standard"
                    value={age}
                    label="Action"
                    onChange={(event)=>{handleChangeRPC(event)}}
                    variant='standard'
                >
                    <MenuItem value={1}>Create</MenuItem>
                    <MenuItem value={2}>Read</MenuItem>
                    <MenuItem value={3}>Update</MenuItem>
                    <MenuItem value={4}>Delete</MenuItem>
                </Select>
            </FormControl>
        </CardActions>
        <Divider />
        <CardContent>
          <Stack
            spacing={3}
            alignItems={'stretch'}
          >
            <TextField
              id="outlined-multiline-static"
              size="large"
              multiline="true"
              label="Mensagem"
              name="password"
              onChange={handleChange}
              value={value}
              fullWidth
              rows="9"
            />
          </Stack>
        </CardContent>
        <Divider />
        <CardActions sx={{ justifyContent: 'flex-end' }}>
          <Button 
          variant="contained" 
          endIcon={<SvgIcon><PaperAirplane /></SvgIcon>} 
          onClick={handleOpen}
          >
            Send
          </Button>
        </CardActions>
        <Backdrop
            sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
            open={open}
            onClick={handleClose}
            >
            <CircularProgress color="inherit" />
        </Backdrop>
        <Dialog
        fullWidth={ true } 
        maxWidth={"md"}
        open={answer}
        scroll={scroll}
        aria-labelledby="scroll-dialog-title"
        aria-describedby="scroll-dialog-description"
      >
        <DialogTitle id="scroll-dialog-title">
        <Box display="flex" alignItems="center">
              <Box flexGrow={1} >Response</Box>
              <Box>
                  <IconButton >
                        <SvgIcon 
                          onClick={()=>{
                          setAnswer(false);
                          handleClose;
                          //setContent("");
                          }}
                        >
                          <XMarkIcon />
                      </SvgIcon>
                  </IconButton>
              </Box>
        </Box>
        </DialogTitle>
        <DialogContent dividers={scroll === 'paper'}>
          <DialogContentText
            id="scroll-dialog-description"
            //ref={descriptionElementRef}
            tabIndex={-1}
          >
          <pre>
            {content}
          </pre>
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={()=>{
            setAnswer(false);
            handleClose;
            //setContent("");
          }}>Ok</Button>
        </DialogActions>
      </Dialog>
      </Card>
    </form>
  );
};
