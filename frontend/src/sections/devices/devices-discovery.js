import { useEffect, useState } from 'react';
import {
  Card,
  CardContent,
  SvgIcon,
  IconButton,
  List,
  ListItem,
  ListItemText,
} from '@mui/material';
import CircularProgress from '@mui/material/CircularProgress';
import PlusCircleIcon from '@heroicons/react/24/outline/PlusCircleIcon';
import Pencil from "@heroicons/react/24/outline/PencilIcon"
import ArrowUturnLeftIcon from '@heroicons/react/24/outline/ArrowUturnLeftIcon'
import { useRouter } from 'next/router';

const AccessType = {
    ReadOnly: 0,
    ReadWrite: 1,
    WriteOnly: 2,
}

function ShowParamsWithValues({x, deviceParametersValue}) {
    let paths = x.supported_obj_path.split(".")
    if(paths[paths.length -2] == "{i}"){
        return Object.keys(deviceParametersValue).map((paramKey, h)=>{
            return (
            <List dense={true} key={h}>
                <ListItem
                    divider={true}
                    sx={{
                        boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;',
                        pl: 4,
                    }}
                >
                <ListItemText
                primary={<b>{paramKey}</b>}
                sx={{fontWeight:'bold'}}
                />
                </ListItem>
            {deviceParametersValue[paramKey].map((param, i) => {
                return (
                <List 
                component="div" 
                disablePadding 
                dense={true}
                key={i}
                >
                    <ListItem
                        key={i}
                        divider={true}
                        sx={{
                            boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;',
                            pl: 4 
                        }}
                        secondaryAction={
                            <div>
                                {Object.values(param)[0].value}
                                {Object.values(param)[0].access > 0 && <IconButton>
                                <SvgIcon sx={{width:'20px'}}>
                                
                                    <Pencil></Pencil>
                                
                                </SvgIcon>
                                </IconButton>}
                            </div>
                        }
                    >
                        <ListItemText
                            primary={Object.keys(param)[0]}
                        />
                    </ListItem>
                </List>
                )
            })}
            </List>
            )
        })
    }else{
        return x.supported_params.map((y, index)=>{
            return (
                <List 
                    component="div" 
                    disablePadding 
                    dense={true}
                    key={y.param_name}
                    >
                    <ListItem
                        key={index}
                        divider={true}
                        sx={{
                            boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;',
                            pl: 4 
                        }}
                        secondaryAction={
                            <div>
                                {deviceParametersValue[y.param_name].value}
                                {deviceParametersValue[y.param_name].access > 0 && <IconButton>
                                <SvgIcon sx={{width:'20px'}}>
                                 
                                    <Pencil></Pencil>
                                
                                </SvgIcon>
                                </IconButton>}
                            </div>
                        }
                    >
                        <ListItemText
                            primary={y.param_name}
                        />
                    </ListItem>
                </List>)
            })
    }
}

export const DevicesDiscovery = () => {

const router = useRouter()

const [deviceParameters, setDeviceParameters] = useState(null)
const [deviceParametersValue, setDeviceParametersValue] = useState({})

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

const getDeviceParameterInstances = async (raw) =>{
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    var requestOptions = {
        method: 'PUT',
        headers: myHeaders,
        redirect: 'follow',
        body: raw
    };

    let result = await (await fetch(`${process.env.NEXT_PUBLIC_REST_ENPOINT}/device/${router.query.id[0]}/instances`, requestOptions))
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
  // Multi instance not used, found better way to get values
  const updateDeviceParametersMultiInstance = async (param) =>{
    console.log("UpdateDeviceParametersMultiInstance => param = ", param)

    let raw = JSON.stringify({
        "obj_paths": [param],
        "first_level_only" : true,
        "return_commands" : true,
        "return_events" : true,
        "return_params" : true 
    })

    let response = await getDeviceParameterInstances(raw)
    console.log("response:", response)

    let instancesToGet = []
    if (response.req_path_results[0].curr_insts) {
        let supportedParams = response.req_path_results[0].curr_insts
        let instances = () => {
            for (let i =0; i < supportedParams.length ;i++){
                instancesToGet.push(supportedParams[i].instantiated_obj_path)
            }
        }
        instances()
    }else{
        instancesToGet.push(response.req_path_results[0].requested_path)
    }

    let rawInP = JSON.stringify({
        "obj_paths": instancesToGet,
        "first_level_only" : true,
        "return_commands" : true,
        "return_events" : true,
        "return_params" : true 
    })

    let resultParams = await getDeviceParameters(rawInP)
    console.log("result params:", resultParams)
    setDeviceParameters(resultParams)


    let paramsToFetch = []

    // console.log("parameters to fetch: ", paramsToFetch)

    //     let rawV = JSON.stringify({
    //         "param_paths": paramsToFetch,
    //         "max_depth": 1
    //     })

    //     let resultValues = await getDeviceParametersValue(rawV)
    //     console.log("result values:", resultValues)


    //     let rawP = JSON.stringify({
    //         "obj_paths": paramsToFetch,
    //         "first_level_only" : true,
    //         "return_commands" : true,
    //         "return_events" : true,
    //         "return_params" : true 
    //     })

    //     let resultParams = await getDeviceParameters(rawP)
    //     console.log("result params:", resultParams)

    //     let values = {}
    //     let setvalues = () => {resultValues.req_path_results.map((x)=>{
    //         // let path = x.requested_path.split(".")
    //         // let param = path[path.length -1]
    //         if (!x.resolved_path_results){
    //             return
    //         }
    //         x.resolved_path_results.map((y)=> {
                
    //         })
    //         // Object.keys(x.resolved_path_results[0].result_params).forEach((key, index) =>{
    //         //     values[key] = x.resolved_path_results[0].result_params[key]
    //         // })
    //         return values
    //     })}
    //     setvalues()
    //     console.log("values:",values)

    //     setDeviceParameters(resultParams)
    //     setDeviceParametersValue(values)
  }

  const updateDeviceParameters = async (param) => {
    console.log("UpdateDeviceParameters => param = ", param)
    let raw = JSON.stringify({
            "obj_paths": [param],
            "first_level_only" : true,
            "return_commands" : true,
            "return_events" : true,
            "return_params" : true 
    })

    let content = await getDeviceParameters(raw)

    console.log("content:",content)

    let values = {}
    let multiInstanceParamsInfo = {}

    let supportedParams = content.req_obj_results[0].supported_objs[0].supported_params
    let parametersToFetch = () => {
        let paramsToFetch = []
        for (let i =0; i < supportedParams.length ;i++){
            
            let supported_obj_path = content.req_obj_results[0].supported_objs[0].supported_obj_path.replaceAll("{i}","*")
            let param = supportedParams[i]
            
            paramsToFetch.push(supported_obj_path+param.param_name)

            let paths = supported_obj_path.split(".")
            if (paths[paths.length -2] !== "*"){
                values[param.param_name] = {
                    "value_change":param["value_change"],
                    "value_type":param["value_type"],
                    "access": param["access"],
                    "value": "-",
                }
            }else{
                multiInstanceParamsInfo[param.param_name] = {
                    "value_change":param["value_change"],
                    "value_type":param["value_type"],
                    "access": param["access"],
                    "value":"-",
                }
            }
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

        let setvalues = () => {result.req_path_results.map((x)=>{
            if (!x.resolved_path_results){
                return
            }

            let paths = x.requested_path.split(".")
            if(paths[paths.length -2] == "*"){
                x.resolved_path_results.map(y=>{
                    Object.keys(y.result_params).forEach((key, index) =>{
                        if (!values[y.resolved_path]){
                            values[y.resolved_path] = []
                        }
                        if (y.result_params[key] == ""){
                            y.result_params[key] = "\"\""
                        }
                        multiInstanceParamsInfo[key].value = y.result_params[key]
                        let obj = {};
                        obj[key] = multiInstanceParamsInfo[key]
                        values[y.resolved_path].push(obj)
                    })
                })
            }else{
                Object.keys(x.resolved_path_results[0].result_params).forEach((key, index) =>{
                    if (x.resolved_path_results[0].result_params[key] != ""){
                        values[key].value = x.resolved_path_results[0].result_params[key]
                    }else{
                        values[key].value = "\"\""
                    }
                })
            }

            return values
        })}
        console.log("VALUES:",values)
        setvalues()
        console.log(values)
        setDeviceParametersValue(values)
        setDeviceParameters(content)
    }else{
        setDeviceParameters(content)
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
    return deviceParameters.req_obj_results.map((a,b)=>{
        return a.supported_objs.map((x,i)=> {

            let supported_obj_path = x.supported_obj_path.split(".")
            let supportedObjPath = ""

            supported_obj_path.map((x,i)=>{
                if(i !== supported_obj_path.length -2){
                    supportedObjPath = supportedObjPath + x + "."
                }
            })

            let req_obj_path = a.req_obj_path.split(".")
            let reqObjPath = ""

            req_obj_path.map((x,i)=>{
                if(i !== req_obj_path.length -2){
                    reqObjPath = reqObjPath + x + "."
                }
            })

            // console.log("reqObjPath:", reqObjPath)
            // console.log("supportedObjPath:", supportedObjPath)

            let paramName = x.supported_obj_path
            if (supportedObjPath != "Device.."){
                if (supportedObjPath == reqObjPath){
                    paramName = a.req_obj_path
                }
            }

            return (
            <List dense={true} key={x.supported_obj_path}>
                <ListItem
                    key={x.supported_obj_path}
                    divider={true}
                    secondaryAction={
                        i == 0 && x.supported_obj_path != "Device." ?
                        <IconButton onClick={()=>
                            {   
                                let supported_obj_path = x.supported_obj_path.replaceAll("{i}.","*.")
                                let paths = supported_obj_path.split(".")
                                console.log("paths:",paths)

                                let pathsToJump = 2
                                if (paths[paths.length -2] == "*"){
                                    pathsToJump = 3
                                }
                                
                                paths.splice(paths.length - pathsToJump, pathsToJump)
                                let pathToFetch = paths.join(".")
                                
                                updateDeviceParameters(pathToFetch)
                            }
                        }>
                        <SvgIcon>
                            <ArrowUturnLeftIcon></ArrowUturnLeftIcon>
                        </SvgIcon>
                        </IconButton>
                        :
                        <IconButton onClick={()=>{
                            console.log("x.supported_obj_path:",x.supported_obj_path)
                            let supported_obj_path = x.supported_obj_path.replaceAll("{i}.","*.")
                            updateDeviceParameters(supported_obj_path)
                        }}>
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
                        primary={<b>{paramName}</b>}
                        sx={{fontWeight:'bold'}}
                    />
                </ListItem>
                { x.supported_params && 
                    <ShowParamsWithValues x={x} deviceParametersValue={deviceParametersValue}/>
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
