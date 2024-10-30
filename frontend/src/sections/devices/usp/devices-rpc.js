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
  Grid
} from '@mui/material';
import XMarkIcon from '@heroicons/react/24/outline/XMarkIcon';
import PaperAirplane from '@heroicons/react/24/solid/PaperAirplaneIcon';
import CircularProgress from '@mui/material/CircularProgress';
import Backdrop from '@mui/material/Backdrop';
import { useRouter } from 'next/router';
import { useBackendContext } from 'src/contexts/backend-context';
import DocumentArrowDown from '@heroicons/react/24/outline/DocumentArrowDownIcon';
import TrashIcon from '@heroicons/react/24/outline/TrashIcon';
import PlusCircleIcon from '@heroicons/react/24/outline/PlusCircleIcon';
import EnvelopeIcon from '@heroicons/react/24/outline/EnvelopeIcon';
import CheckIcon from '@heroicons/react/24/outline/CheckIcon';


export const DevicesRPC = () => {

const router = useRouter()
let { httpRequest } = useBackendContext()

const [open, setOpen] = useState(false);
const [scroll, setScroll] = useState('paper');
const [answer, setAnswer] = useState(false)
const [content, setContent] = useState('')
const [age, setAge] = useState(6);
const [newMessage, setNewMessage] = useState(false)
const [message, setMessage] = useState(null)
const [currentMsg, setCurrentMsg] = useState(0)
const [newMsgName, setNewMsgName] = useState("")
const [value, setValue] = useState()
const [saveChanges, setSaveChanges] = useState(false)
const [loadingSaveMsg, setLoadingSaveMsg] = useState(false)
const possibleMsgs = [
  `{
    "header": {
        "msg_id": "b7dc38ea-aefb-4761-aa55-edaa97adb2f0",
        "msg_type": 4
    },
    "body": {
        "request": {
            "set": {
                "allow_partial":true,
                "update_objs":[
                    {
                        "obj_path":"Device.IP.Interface.1.",
                        "param_settings":[
                            {
                            "param":"Alias",
                            "value":"test",
                            "required":true
                            }
                        ]
                    }
                ]
            }
        }
    }
}`,
`{
    "header": {
        "msg_id": "b7dc38ea-aefb-4761-aa55-edaa97adb2f0",
        "msg_type": 10
    },
    "body": {
        "request": {
            "delete": {
                "allow_partial": true,
                "obj_paths": [
                    "Device.IP.Interface.[Alias==test]."
                ]
            }
        }
    }
}`,
`
{
    "header": {
        "msg_id": "b7dc38ea-aefb-4761-aa55-edaa97adb2f0",
        "msg_type": 8
    },
    "body": {
        "request": {
            "add": {
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
            }
        }
    }
}`,
`{
    "header": {
        "msg_id": "b7dc38ea-aefb-4761-aa55-edaa97adb2f0",
        "msg_type": 6
    },
    "body": {
        "request": {
            "operate": {
                "command": "Device.Reboot()",
                "send_resp": true
            }
        }
    }
}`,
`{
    "header": {
        "msg_id": "b7dc38ea-aefb-4761-aa55-edaa97adb2f0",
        "msg_type": 1
    },
    "body": {
        "request": {
            "get": {
                "paramPaths": [
                    "Device.WiFi.SSID.[Name==wlan0].",
                    "Device.IP.Interface.*.Alias",
                    "Device.DeviceInfo.FirmwareImage.*.Alias",
                    "Device.IP.Interface.1.IPv4Address.1.IPAddress"
                ],
                "maxDepth": 2
            }
        }
    }
}`,`{
    "header": {
        "msg_id": "b7dc38ea-aefb-4761-aa55-edaa97adb2f0",
        "msg_type": 12
    },
    "body": {
        "request": {
            "get_supported_dm": {
                "obj_paths" : [
                    "Device."
                ],
                "first_level_only" : false,
                "return_commands" : false,
                "return_events" : false,
                "return_params" : true 
            }
        }
    }
}`,
`{
    "header": {
        "msg_id": "b7dc38ea-aefb-4761-aa55-edaa97adb2f0",
        "msg_type": 14
    },
    "body": {
        "request": {
            "get_instances": {
                "obj_paths" : ["Device.DeviceInfo."],
                "first_level_only" : false
            }
        }
    }
}`]
const [newMsgValue, setNewMsgValue] = useState(possibleMsgs[age-1])
const [loading, setLoading] = useState(false);

const handleNewMessageValue = (event) => {
  setNewMsgValue(event.target.value)
}

const handleClose = () => {
  setOpen(false);
};

const handleCancelNewMsgTemplate = () => {
  setNewMessage(false)
  setNewMsgName("")
  setNewMsgValue(possibleMsgs[age-1])
  // setValue(possibleMsgs[age-1])
}

const saveMsg = async () => {
  let {status} = await httpRequest(
    `/api/device/message?name=`+message[currentMsg].name,
    "PUT", 
    value,
    null,
  )
  if ( status === 204){
    setSaveChanges(false)
    setMessage(message.map((msg, index) => {
      if (index === currentMsg) {
        return {...msg, value: value}
      }else{
        return msg
      }
    }))
  }
}

const createNewMsg = async () => {
  setLoading(true)
  let {status} = await httpRequest(
    `/api/device/message/usp?name=`+newMsgName,
    "POST", 
    newMsgValue,
    null,
  )
  if ( status === 204){
    setNewMessage(false)
    setNewMsgName("")
    let result = await fetchMessages()
    if (result) {
      setCurrentMsg(result.length-1)
    }
    setValue(newMsgValue)
    setNewMsgValue(possibleMsgs[age-1])
  }
  setLoading(false)
}

const handleChangeMessage = (event) => {
  setSaveChanges(false)
  setCurrentMsg(event.target.value)
  setValue(message[event.target.value].value)
}

const handleDeleteMessage = async () => {
  let {status} = await httpRequest(
    `/api/device/message?name=`+message[currentMsg].name.replace(" ", '+'),
    "DELETE", 
  )
  if ( status === 204){
    fetchMessages()
    setCurrentMsg(0)
    setValue("")
  }
}

const handleOpen = async () => {
  setOpen(true);

  let {result, status} = await httpRequest(
    `/api/device/${router.query.id[0]}/any/generic`,
    "PUT", 
    value, 
    null,
  )
  if (status === 200){
    setAnswer(true)
    console.log("result:",result)
    let answer = JSON.stringify(result, null, 2)
    if (answer == "null"){
      answer = result
    }
    console.log(answer)
    setContent(answer)
  }
  setOpen(false)

}

const fetchMessages = async () => {
  let {result, status} = await httpRequest(
    `/api/device/message?type=usp`,
    "GET", 
    null, 
    null,
  )
  if ( status === 200){
    setMessage(result)
    setValue(result ? result[0].value : "")
    return result
  }
}

  const handleChangeRPC = (event) => {
    setAge(event.target.value);
    setNewMsgValue(possibleMsgs[event.target.value-1])
  };

  const handleEditMessage = (event) => {
    if (message) {
      setSaveChanges(true)
    }
    setValue(event.target.value)
  }

  const handleSubmit = useCallback(
    (event) => {
      event.preventDefault();
    },
    []
  );

  useEffect(() => {
    fetchMessages();
  },[]);

  return (
    <form onSubmit={handleSubmit}>
      <Card>
        <CardHeader sx={{ justifyContent: 'flex-end'}} 
          avatar={<SvgIcon>< EnvelopeIcon/></SvgIcon>}
          title="Custom Message" 
          action={ 
          <Stack direction={"row"} spacing={1} width={"100%"} justifyContent={"flex-end"}>
            <Button sx={{ backgroundColor: "rgba(48, 109, 111, 0.04)" }}
            endIcon={<SvgIcon><PlusCircleIcon /></SvgIcon>}
            onClick={()=>{setNewMessage(true)}}
            >
              <Stack direction={"row"} spacing={1}>
                New Message
              </Stack>
            </Button>
          </Stack>}
        >
        </CardHeader>
        <Divider />
        <CardContent>
          <Stack pb={4} spacing={5} direction={"row"}>
            <FormControl sx={{display:"flex", width: "15%"}} variant="standard" >
            <InputLabel>Message</InputLabel>
              <Select
                  value={currentMsg}
                  onChange={(event)=>{handleChangeMessage(event)}}
              > 
                {message && message.map((msg, index) => {
                  return  <MenuItem value={index}>{msg.name}</MenuItem>
                })}
              </Select>
            </FormControl>
          </Stack>
          <Stack
            spacing={3}
            alignItems={'stretch'}
          >
            {!loadingSaveMsg ? <TextField
              id="outlined-multiline-static"
              size="large"
              multiline="true"
              // label="Payload"
              name="password"
              onChange={handleEditMessage}
              value={value}
              variant="filled"
              fullWidth
              rows="15"
            />:<CircularProgress />}
          </Stack>
        </CardContent>
        {/* <Divider /> */}
        <CardActions>
          <Stack direction={"row"} spacing={1} width={"100%"} justifyContent={"flex-start"}>
            <Button 
            variant="contained" 
            endIcon={<SvgIcon><TrashIcon /></SvgIcon>} 
            onClick={handleDeleteMessage}
            disabled={!message}
            >
              Delete
            </Button>
            {!loadingSaveMsg ? <Button 
            variant="contained" 
            endIcon={<SvgIcon><DocumentArrowDown /></SvgIcon>} 
            onClick={saveMsg}
            disabled={!saveChanges}
            >
              Save
            </Button>: <CircularProgress />}
          </Stack>
          <Stack direction={"row"} spacing={1} width={"100%"} justifyContent={"flex-end"}>
            <Button 
            variant="contained" 
            endIcon={<SvgIcon><PaperAirplane /></SvgIcon>} 
            onClick={handleOpen}
            >
              Send
            </Button>
          </Stack>
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
          <pre style={{color: 'black'}}>
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
      <Dialog open={newMessage} maxWidth={"800px"}>
          <DialogTitle>
           <SvgIcon><EnvelopeIcon/></SvgIcon>
          </DialogTitle>
          <DialogContent>
          <Stack
              direction={"row"}
              container
              pb={3}
              // direction="column"
              // alignItems="center"
              // justifyContent="center"
              pt={1}
              spacing={3}
            >
              <TextField
                  variant='standard'
                  fullWidth
                  value={newMsgName}
                  onChange={(event)=>{setNewMsgName(event.target.value)}}
                  label="Name"
                  sx={{maxWidth: "30%", justifyContent:"center"}}
              />
                <FormControl sx={{display:"flex", width: "30%"}} variant="standard" >
                  <InputLabel>Template</InputLabel>
                  <Select
                      value={age}
                      label="Action"
                      name='action'
                      onChange={(event)=>{handleChangeRPC(event)}}
                  >
                      <MenuItem value={1}>Set</MenuItem>
                      <MenuItem value={2}>Delete</MenuItem>
                      <MenuItem value={3}>Add</MenuItem>
                      <MenuItem value={4}>Operate</MenuItem>
                      <MenuItem value={5}>Get</MenuItem>
                      <MenuItem value={6}>Get Supported DM</MenuItem>
                      <MenuItem value={7}>Get Instances</MenuItem>
                  </Select>
                </FormControl>
            </Stack>
            <Stack
              spacing={3}
              alignItems={'stretch'}
              width={"600px"}
            >
              <TextField
                id="outlined-multiline-static"
                size="large"
                multiline="true"
                label="Payload"
                name="password"
                onChange={handleNewMessageValue}
                value={newMsgValue}
                // fullWidth
                // rows="15"
              />
            </Stack>
          </DialogContent>
          {/* <Divider/> */}
          <Stack direction={"row"} spacing={1} width={"100%"} justifyContent={"flex-end"} p={2}>
            <Button 
            variant="contained" 
            onClick={handleCancelNewMsgTemplate}
            >
              Cancel
            </Button>
            {!loading ?
            <Button 
            variant="contained" 
            endIcon={<SvgIcon><CheckIcon /></SvgIcon>} 
            onClick={createNewMsg}
            >
              Save
            </Button>:<CircularProgress />}
          </Stack>
      </Dialog>
      </Card>
    </form>
  );
};
