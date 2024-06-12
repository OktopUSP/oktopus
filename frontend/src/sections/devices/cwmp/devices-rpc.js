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

const [value, setValue] = useState(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:GetParameterValues>
      <ParameterNames>
        <string>InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.</string>
        <string>InternetGatewayDevice.LANDevice.1.WLANConfiguration.2.</string>
       </ParameterNames>
    </cwmp:GetParameterValues>
  </soap:Body>
</soap:Envelope>`)

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

const handleClose = () => {
  setOpen(false);
};
const handleOpen = () => {
  setOpen(true);
  var myHeaders = new Headers();
  myHeaders.append("Authorization", localStorage.getItem("token"));

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
      method="addObject"
      break;
    case 2:
      method="getParameterValues"
      break;
    case 3:
      method="setParameterValues"
      break;
    case 4:
      method="deleteObject"
    break;
  }

 
  fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT}/device/cwmp/${router.query.id[0]}/${method}`, requestOptions)
    .then(response => response.text())
    .then(result => {
      if (result.status === 401){
        router.push("/auth/login")
      }
      setOpen(false)
      setAnswer(true)
      let teste =  prettifyXml(result)
      console.log(teste)
      setContent(teste)
    })
    .catch(error => console.log('error', error));
  };

  const handleChangeRPC = (event) => {
    setAge(event.target.value);
    switch(event.target.value) {
      case 1:
        setValue(`<?xml version="1.0" encoding="UTF-8"?>
        <soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
          <soap:Header/>
          <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
            <cwmp:AddObject>
              <ObjectName>InternetGatewayDevice.LANDevice.</ObjectName>
              <ParameterKey></ParameterKey>
            </cwmp:AddObject>
          </soap:Body>
        </soap:Envelope>`)
        break;
      case 2:
        setValue(`<?xml version="1.0" encoding="UTF-8"?>
        <soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
          <soap:Header/>
          <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
            <cwmp:GetParameterValues>
              <ParameterNames>
                <string>InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.</string>
                <string>InternetGatewayDevice.LANDevice.1.WLANConfiguration.2.</string>
                <string>InternetGatewayDevice.LANDevice.2.WLANConfiguration.2.</string>
                <string>InternetGatewayDevice.LANDevice.2.WLANConfiguration.1.</string>
               </ParameterNames>
            </cwmp:GetParameterValues>
          </soap:Body>
        </soap:Envelope>`)
        break;
      case 3:
        setValue(`
        <?xml version="1.0" encoding="UTF-8"?>
        <soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
          <soap:Header/>
          <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
            <cwmp:SetParameterValues>
              <ParameterList soapenc:arrayType="cwmp:ParameterValueStruct[3]">
            <ParameterValueStruct>
                <Name>InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.Enable</Name>
                <Value>0</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                <Name>InternetGatewayDevice.LANDevice.1.WLANConfiguration.2.SSID</Name>
                <Value>HUAWEI_TEST-2</Value>
                </ParameterValueStruct>
              </ParameterList>
              <ParameterKey>LC1309123</ParameterKey>
            </cwmp:SetParameterValues>
          </soap:Body>
        </soap:Envelope>`)
        break;
      case 4:
        setValue(`<?xml version="1.0" encoding="UTF-8"?>
        <soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
          <soap:Header/>
          <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
            <cwmp:DeleteObject>
              <ObjectName>InternetGatewayDevice.LANDevice.3.</ObjectName>
              <ParameterKey></ParameterKey>
            </cwmp:DeleteObject>
          </soap:Body>
        </soap:Envelope>`)
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
              rows="15"
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
      </Card>
    </form>
  );
};
