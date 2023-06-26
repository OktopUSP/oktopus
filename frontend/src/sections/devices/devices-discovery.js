import { useEffect, useState } from 'react';
import {
  Card,
  CardContent,
  SvgIcon,
  IconButton,
  List,
  ListItem,
  ListItemText,
  Collapse,
  Box,
  Tabs,
  Tab
} from '@mui/material';
import CircularProgress from '@mui/material/CircularProgress';
import { useRouter } from 'next/router';


export const DevicesDiscovery = () => {

const router = useRouter()

const [deviceParameters, setDeviceParameters] = useState(null)

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
        "first_level_only" : false,
        "return_commands" : false,
        "return_events" : false,
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

//   const updateDeviceParameters = async (param,) => {
    
//     let raw = JSON.stringify({
//             "obj_paths": [param],
//             "first_level_only" : true,
//             "return_commands" : false,
//             "return_events" : false,
//             "return_params" : true 
//     })

//     let content = await getDeviceParameters(raw)

//     console.log(content)

//     setDeviceParameters(prevState => {
//         return {...prevState, req_obj_results: [
//                 {
//                 supported_objs:[ ...prevState.req_obj_results[0].
//                 supported_objs][index]
//                 .supported_params:[]
//                 }
//             ]
//         }
//     })
//   }

  const showParameters = () => {
    console.log(deviceParameters)
    return deviceParameters.req_obj_results[0].supported_objs.map((x,i)=> {
        return (
        <List dense={true} key={x.supported_obj_path}>
            <ListItem
                key={x.supported_obj_path}
                divider={true}
                // secondaryAction={
                //     <IconButton /*onClick={()=>updateDeviceParameters(x.supported_obj_path, index)}*/>
                //         <SvgIcon>
                //             <PlusCircleIcon></PlusCircleIcon>
                //         </SvgIcon>
                //     </IconButton>
                // }
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
                x.supported_params.map((y)=>{
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
                    >
                        <ListItemText
                            primary={y.param_name}
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
