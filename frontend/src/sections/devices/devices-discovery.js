import { useEffect, useState } from 'react';
import {
  Card,
  CardContent,
  SvgIcon,
  IconButton,
  List,
  ListItem,
  ListItemText,
  Box,
} from '@mui/material';
import CircularProgress from '@mui/material/CircularProgress';
import PlusCircleIcon from '@heroicons/react/24/outline/PlusCircleIcon';
import ArrowUturnLeftIcon from '@heroicons/react/24/outline/ArrowUturnLeftIcon'
import { useRouter } from 'next/router';


export const DevicesDiscovery = () => {

const router = useRouter()

const [deviceParameters, setDeviceParameters] = useState(null)
const [deviceParametersValue, setDeviceParametersValue] = useState([])

const initialize = async (raw) => {
    let content = await getDeviceParameters(raw)
    setDeviceParameters(content)
}

const getDeviceParameters = async (raw) =>{
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    var requestOptions = {
        method: 'PUT',
        headers: myHeaders,
        redirect: 'follow',
        body: raw
    };

    let result = await (await fetch(`${process.env.NEXT_PUBLIC_REST_ENPOINT}/device/${router.query.id[0]}/parameters`, requestOptions))
    if (result.status != 200) {
        throw new Error('Please check your email and password');
    }else{
        return result.json()
    }
}

  useEffect(()=> {

    initialize(
    JSON.stringify({
        "obj_paths": ["Device."],
        "first_level_only" : true,
        "return_commands" : true,
        "return_events" : true,
        "return_params" : true 
        })
    );
  },[])

//Together with showParameters, this function renders all the device parameters the device supports
//but you must set req with first_level_only property to false
//   const showPathParameters = (pathParamsList) => {
//     return pathParamsList.map((x,i)=>{
//         return(
//         <List component="div" disablePadding dense={true}>
//             <ListItem
//             key={i}
//             divider={true}
//                 sx={{
//                     boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;',
//                     pl: 4 
//                 }}
//             >
//             <ListItemText
//                 primary={x.param_name}
//             />
//             </ListItem>
//     </List>
//         )
//     })
//   }

  const updateDeviceParameters = async (param) => {
    
    let raw = JSON.stringify({
            "obj_paths": [param],
            "first_level_only" : true,
            "return_commands" : true,
            "return_events" : true,
            "return_params" : true 
    })

    let content = await getDeviceParameters(raw)

    console.log("content:",content)

    setDeviceParameters(content)

    let supportedParams = content.req_obj_results[0].supported_objs[0].supported_params
    let parametersToFetch = () => {
        let paramsToFetch = []
        for (let i =0; i < supportedParams.length ;i++){
            paramsToFetch.push(content.req_obj_results[0].supported_objs[0].supported_obj_path+supportedParams[i].param_name)
        }
        return paramsToFetch
    }

    if (supportedParams !== undefined) {
        const fetchparameters = parametersToFetch()
        console.log("parameters to fetch: ", fetchparameters)

        raw = JSON.stringify({
            "param_paths": fetchparameters,
            "max_depth": 1
        })

        let result = await getDeviceParametersValue(raw)
        console.log("result:", result)

        let values = []
        let setValues = result.req_path_results.map((x)=>{
            let path = x.requested_path.split(".")
            let param = path[path.length -1]
            return values.push(x.resolved_path_results[0].result_params[param])
        })
        console.log(values)
        setDeviceParametersValue(values)
    }
  }

  const getDeviceParametersValue = async (raw) => {

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    var requestOptions = {
        method: 'PUT',
        headers: myHeaders,
        redirect: 'follow',
        body: raw
    };

    let result = await (await fetch(`${process.env.NEXT_PUBLIC_REST_ENPOINT}/device/${router.query.id[0]}/get`, requestOptions))
    if (result.status != 200) {
        throw new Error('Please check your email and password');
    }else{
        return result.json()
    }

  }

  const showParameters = () => {
    return deviceParameters.req_obj_results[0].supported_objs.map((x,i)=> {
        return (
        <List dense={true} key={x.supported_obj_path}>
            <ListItem
                key={x.supported_obj_path}
                divider={true}
                secondaryAction={
                    i == 0 && x.supported_obj_path != "Device." ?
                    <IconButton onClick={()=>
                        {
                            let paths = x.supported_obj_path.split(".")
                            console.log(paths)
                            updateDeviceParameters(paths[paths.length -3]+".")
                        }
                    }>
                    <SvgIcon>
                        <ArrowUturnLeftIcon></ArrowUturnLeftIcon>
                    </SvgIcon>
                    </IconButton>
                    :
                    <IconButton onClick={()=>updateDeviceParameters(x.supported_obj_path)}>
                        <SvgIcon>
                            {
                            x.supported_obj_path != "Device." && 
                            <PlusCircleIcon></PlusCircleIcon>
                            }
                        </SvgIcon>
                    </IconButton>
                }
                sx={{
                    boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;'
                }}
            >
                <ListItemText
                    primary={<b>{x.supported_obj_path}</b>}
                    sx={{fontWeight:'bold'}}
                />
            </ListItem>
            { x.supported_params &&
                x.supported_params.map((y, index)=>{
                    return <List 
                    component="div" 
                    disablePadding 
                    dense={true}
                    key={y.param_name}
                    >
                    <ListItem
                        key={i}
                        divider={true}
                        sx={{
                            boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;',
                            pl: 4 
                        }}
                        secondaryAction={
                            <div>{deviceParametersValue[index]}</div>
                        }
                    >
                        <ListItemText
                            primary={y.param_name}
                        />
                    </ListItem>
                </List>
                })
            }
            { x.supported_commands &&
                x.supported_commands.map((y)=>{
                    return <List 
                    component="div" 
                    disablePadding 
                    dense={true}
                    key={y.command_name}
                    >
                    <ListItem
                        key={i}
                        divider={true}
                        sx={{
                            boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;',
                            pl: 4 
                        }}
                    >
                        <ListItemText
                            primary={y.command_name}
                        />
                    </ListItem>
                </List>
                })
            }
            { x.supported_events &&
                x.supported_events.map((y)=>{
                    return <List 
                    component="div" 
                    disablePadding 
                    dense={true}
                    key={y.event_name}
                    >
                    <ListItem
                        key={i}
                        divider={true}
                        sx={{
                            boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;',
                            pl: 4 
                        }}
                    >
                        <ListItemText
                            primary={y.event_name}
                        />
                    </ListItem>
                </List>
                })
            }
        </List>)
    })
  }
  
  return ( deviceParameters ?
    <Card>
        <CardContent>
            {showParameters()}
        </CardContent>
    </Card> 
    :
    <Box sx={{display:'flex',justifyContent:'center'}}>
        <CircularProgress color="inherit" />
    </Box>
  )
};
