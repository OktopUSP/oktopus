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

var prettifyXml = function(sourceXml)
{
    var xmlDoc = new DOMParser().parseFromString(sourceXml, 'application/xml');
    var xsltDoc = new DOMParser().parseFromString([
        // describes how we want to modify the XML - indent everything
        '<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform">',
        '  <xsl:strip-space elements="*"/>',
        '  <xsl:template match="para[content-style][not(text())]">', // change to just text() to strip space in text nodes
        '    <xsl:value-of select="normalize-space(.)"/>',
        '  </xsl:template>',
        '  <xsl:template match="node()|@*">',
        '    <xsl:copy><xsl:apply-templates select="node()|@*"/></xsl:copy>',
        '  </xsl:template>',
        '  <xsl:output indent="yes"/>',
        '</xsl:stylesheet>',
    ].join('\n'), 'application/xml');

    var xsltProcessor = new XSLTProcessor();    
    xsltProcessor.importStylesheet(xsltDoc);
    var resultDoc = xsltProcessor.transformToDocument(xmlDoc);
    var resultXml = new XMLSerializer().serializeToString(resultDoc);
    return resultXml;
};

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
  `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:SetParameterValues>
      <ParameterList soapenc:arrayType="cwmp:ParameterValueStruct[4]">
                <ParameterValueStruct>
                    <Name>InternetGatewayDevice.TraceRouteDiagnostics.Host</Name>
                    <Value>192.168.60.4</Value>
                </ParameterValueStruct>
      </ParameterList>
      <ParameterKey></ParameterKey>
    </cwmp:SetParameterValues>
  </soap:Body>
</soap:Envelope>`,
`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:DeleteObject>
      <ObjectName>InternetGatewayDevice.LANDevice.1.WLANConfiguration.2.</ObjectName>
      <ParameterKey></ParameterKey>
    </cwmp:DeleteObject>
  </soap:Body>
</soap:Envelope>`,
`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:AddObject>
      <ObjectName>InternetGatewayDevice.LANDevice.1.WLANConfiguration.</ObjectName>
      <ParameterKey></ParameterKey>
    </cwmp:AddObject>
  </soap:Body>
</soap:Envelope>`,
`<?xml version="1.0" encoding="UTF-8"?>
	<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	  <soap:Header/>
	  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
		<cwmp:Reboot>
			<CommandKey>
				007
			</CommandKey>
		</cwmp:Reboot>
	  </soap:Body>
	</soap:Envelope>`,
`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..schemaswt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:GetParameterValues>
      <ParameterNames>
        <string>InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.</string>
       </ParameterNames>
    </cwmp:GetParameterValues>
  </soap:Body>
</soap:Envelope>`,`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:GetParameterNames>
      <ParameterPath>InternetGatewayDevice.</ParameterPath>
      <NextLevel>1</NextLevel>
    </cwmp:GetParameterNames>
  </soap:Body>
</soap:Envelope>`,
`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:GetParameterAttributes>
      <ParameterNames>
        <string>InternetGatewayDevice.LANDevice.1.WLANConfiguration.</string>
       </ParameterNames>
    </cwmp:GetParameterAttributes>
  </soap:Body>
</soap:Envelope>`]
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
    `/api/device/message/cwmp?name=`+newMsgName,
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
    `/api/device/cwmp/${router.query.id[0]}/generic`,
    "PUT", 
    value, 
    null,
    "text",
  )
  if (status === 200){
    setAnswer(true)
    console.log("result:",result)
    let answer = prettifyXml(result)
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
    `/api/device/message?type=cwmp`,
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
    setSaveChanges(true)
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
                      <MenuItem value={1}>SetParameterValues</MenuItem>
                      <MenuItem value={2}>DeleteObject</MenuItem>
                      <MenuItem value={3}>AddObject</MenuItem>
                      <MenuItem value={4}>Reboot</MenuItem>
                      <MenuItem value={5}>GetParameterValues</MenuItem>
                      <MenuItem value={6}>GetParameterNames</MenuItem>
                      <MenuItem value={7}>GetParameterAttributes</MenuItem>
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
