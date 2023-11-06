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
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  TextField,
  Button,
  Backdrop,
} from '@mui/material';
import ArrowRightIcon from '@heroicons/react/24/solid/ArrowRightIcon';
import CircularProgress from '@mui/material/CircularProgress';
import PlayCircle from '@heroicons/react/24/outline/PlayCircleIcon'
import PlusCircleIcon from '@heroicons/react/24/outline/PlusCircleIcon';
import Pencil from "@heroicons/react/24/outline/PencilIcon"
import ArrowUturnLeftIcon from '@heroicons/react/24/outline/ArrowUturnLeftIcon'
import XMarkIcon from '@heroicons/react/24/outline/XMarkIcon';

import { useRouter } from 'next/router';
import TrashIcon from '@heroicons/react/24/outline/TrashIcon';
import PlayCircleIcon from '@heroicons/react/24/outline/PlayCircleIcon';

/*
    OBJ_READ_ONLY (0)
    OBJ_ADD_DELETE (1)
    OBJ_ADD_ONLY (2)
    OBJ_DELETE_ONLY (3)

*/

const ObjAccessType = {
    ReadOnly: 0,
    AddDelete: 1,
    AddOnly: 2,
    DeleteOnly:3,
}


const ParamAccessType = {
    ReadOnly: 0,
    ReadWrite: 1,
    WriteOnly: 2,
}

const ParamValueType = {
    Unknown: 0,
    Base64: 1,
    Boolean: 2,
    DateTime: 3,
    Decimal: 4, 
    HexBinary: 5,
    Int: 6,
    Long: 7,
    String: 8,
    UnisgnedInt: 9,
    UnsignedLong: 10,
}
//TODO: refact all of this mess
const addDeviceObj = async(obj, setShowLoading, router, updateDeviceParameters) => {
    console.log("AddDeviceObj => obj = ", obj)
    let raw = JSON.stringify(
        {
            "allow_partial": true,
            "create_objs": [
                {
                    "obj_path": obj,
                    //TODO: create queue and create new obj with parameters
                    // "param_settings": [
                    //     {
                    //         "param": "Alias",
                    //         "value": "test",
                    //         "required": true
                    //     }
                    // ]
                }
            ]
        }
    )
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    var requestOptions = {
        method: 'PUT',
        headers: myHeaders,
        redirect: 'follow',
        body: raw
    };
    setShowLoading(true)
    let result = await (await fetch(`${process.env.NEXT_PUBLIC_REST_ENPOINT}/device/${router.query.id[0]}/add`, requestOptions))
    if (result.status != 200) {
        setShowLoading(false)
        throw new Error('Please check your email and password');
    }else{
        setShowLoading(false)
        console.log(result.json())
        updateDeviceParameters(obj+"*.")
    }
}

const deleteDeviceObj = async(obj, setShowLoading, router, updateDeviceParameters) => {
    console.log("deleteDeviceObj => obj = ", obj)
    let raw = JSON.stringify(
        {
            "allow_partial": true,
            "obj_paths": [
                obj
            ]
        }
    )
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    var requestOptions = {
        method: 'PUT',
        headers: myHeaders,
        redirect: 'follow',
        body: raw
    };
    setShowLoading(true)
    let result = await (await fetch(`${process.env.NEXT_PUBLIC_REST_ENPOINT}/device/${router.query.id[0]}/del`, requestOptions))
    if (result.status != 200) {
        setShowLoading(false)
        throw new Error('Please check your email and password');
    }else{
        setShowLoading(false)
        console.log(result.json())

        let paths = obj.split(".")
        console.log("paths:",paths)

        let pathsToJump = 2
        if (paths[paths.length -2] == "*"){
            pathsToJump = 3
        }
        
        paths.splice(paths.length - pathsToJump, pathsToJump)
        let pathToFetch = paths.join(".")
        
        updateDeviceParameters(pathToFetch)
    }
}

function ShowPath({x,updateDeviceParameters,setShowLoading, router}) {
    // console.log(x)
    // console.log("x.supported_obj_path:", x.supported_obj_path)
    // console.log("x.access:", x.access)
    if(x.supported_obj_path != "Device."){
        if (x.access === ObjAccessType.ReadOnly || x.access === undefined){
            return (
                <IconButton onClick={()=>{
                    console.log("x.supported_obj_path:",x.supported_obj_path)
                    let supported_obj_path = x.supported_obj_path.replaceAll("{i}.","*.")
                    updateDeviceParameters(supported_obj_path)
                }}>
            <SvgIcon>
                <ArrowRightIcon></ArrowRightIcon>
            </SvgIcon>
            </IconButton>)
        }else if (x.access === ObjAccessType.AddDelete){
            console.log("addDelete")
            return (
                <div style={{justifyContent:"center", display:'flex'}}>
                    <IconButton>
                    <SvgIcon onClick={()=>addDeviceObj(
                        x.supported_obj_path.replace("{i}.",""),
                        setShowLoading, router, updateDeviceParameters)}>
                    <PlusCircleIcon></PlusCircleIcon>
                    </SvgIcon>
                    </IconButton>
                    <IconButton onClick={()=>{
                    console.log("x.supported_obj_path:",x.supported_obj_path)
                    let supported_obj_path = x.supported_obj_path.replaceAll("{i}.","*.")
                    updateDeviceParameters(supported_obj_path)
                }}>
            <SvgIcon>
                <ArrowRightIcon></ArrowRightIcon>
            </SvgIcon>
            </IconButton>
                </div>
            )
        }else if (x.access === ObjAccessType.AddOnly){
            return <IconButton>
                <SvgIcon 
                onClick={()=>addDeviceObj(x.supported_obj_path.replace("{i}.",""),
                setShowLoading, router, updateDeviceParameters)}>
                <PlusCircleIcon></PlusCircleIcon>
                </SvgIcon>
                </IconButton>
        }
    }
    return <></>
}

function ShowParamsWithValues({
    x, deviceParametersValue, 
    setOpen, setParameter, 
    setParameterValue, deviceParameters, 
    setShowLoading, router,
    updateDeviceParameters, deviceCommands,
    openCommandDialog
}) {
    console.log("HEY jow:", deviceParametersValue)
    let paths = x.supported_obj_path.split(".")
    const showDialog = (param, paramvalue) => {
        setParameter(param);
        if (paramvalue == "\"\"") {
            setParameterValue("")
        }else{
            setParameterValue(paramvalue);
        }
        setOpen(true);
    }

    if(paths[paths.length -2] == "{i}"){
        return Object.keys(deviceParametersValue).map((paramKey, h)=>{
            console.log("deviceParametersValue:", deviceParametersValue)
            console.log("paramKey:", paramKey)
            console.log("deviceParameters.req_obj_results[0].supported_objs[h]?.access:", deviceParameters.req_obj_results[0].supported_objs[h]?.access)
            let obj = deviceParameters.req_obj_results[0].supported_objs[0]
            let access = obj?.access
            return (
            <List dense={true} key={h}>
                <ListItem
                    divider={true}
                    sx={{
                        boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;',
                        pl: 4,
                    }}
                    secondaryAction={
                        access > ObjAccessType.ReadOnly &&
                        <IconButton onClick={()=>{
                            deleteDeviceObj(
                                paramKey,
                                setShowLoading,
                                router,
                                updateDeviceParameters
                            )
                        }}>
                        <SvgIcon>
                            <TrashIcon></TrashIcon>
                        </SvgIcon>
                        </IconButton>
                    }
                >
                <ListItemText
                primary={<b>{paramKey}</b>}
                sx={{fontWeight:'bold'}}
                />
                </ListItem>
            {deviceParametersValue[paramKey].length > 0 ?
            deviceParametersValue[paramKey].map((param, i) => {
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
                                {Object.values(param)[0].access > ParamAccessType.ReadOnly && <IconButton>
                                <SvgIcon sx={{width:'20px'}}
                                onClick={()=>{
                                    showDialog(
                                        paramKey+Object.keys(param)[0],
                                        Object.values(param)[0].value)
                                }
                                }>
                                
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
            }):<></>}
            {
              deviceCommands && 
              Object.keys(deviceCommands).map(commando =>{
                console.log("Comando:", commando)
                return <List 
                component="div" 
                disablePadding 
                dense={true}
                key={commando}
                >
                    <ListItem
                        key={commando}
                        divider={true}
                        sx={{
                            boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;',
                            pl: 4 
                        }}
                        secondaryAction={
                            <IconButton>
                                <SvgIcon>
                                    <PlayCircleIcon>
                                    </PlayCircleIcon>
                                </SvgIcon>
                            </IconButton>
                        }
                    >
                        <ListItemText
                            primary={commando}
                        />
                    </ListItem>
                </List>
              })
            }
            </List>
            )
        })
    }else{
        return (
            <>
            {x.supported_params && x.supported_params.map((y, index)=>{
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
                                    {deviceParametersValue[y.param_name].access > ParamAccessType.ReadOnly && <IconButton>
                                    <SvgIcon sx={{width:'20px'}}
                                    onClick={()=>{
                                        showDialog(
                                            x.supported_obj_path + y.param_name,
                                            deviceParametersValue[y.param_name].value)
                                    }
                                    }>
                                    
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
                    </List>
                )
            })}
            {
              deviceCommands && 
              Object.keys(deviceCommands).map(commando =>{
                console.log("Comando:", commando)
                return <List 
                component="div" 
                disablePadding 
                dense={true}
                key={commando}
                >
                    <ListItem
                        key={commando}
                        divider={true}
                        sx={{
                            boxShadow: 'rgba(149, 157, 165, 0.2) 0px 0px 5px;',
                            pl: 4 
                        }}
                        secondaryAction={
                            <IconButton>
                                <SvgIcon>
                                    <PlayCircleIcon>
                                    </PlayCircleIcon>
                                </SvgIcon>
                            </IconButton>
                        }
                    >
                        <ListItemText
                            primary={commando}
                        />
                    </ListItem>
                </List>
              })
            }
            </>
        )
    }
}

export const DevicesDiscovery = () => {

const router = useRouter()

const [deviceCommands, setDeviceCommands] = useState({})
const [deviceParameters, setDeviceParameters] = useState(null)
const [parameter, setParameter] = useState(null)
const [parameterValue, setParameterValue] = useState(null)
const [parameterValueChange, setParameterValueChange] = useState(null)
const [deviceParametersValue, setDeviceParametersValue] = useState({})
const [open, setOpen] = useState(false)
const [errorModal, setErrorModal] = useState(false)
const [errorModalText, setErrorModalText] = useState("")
const [showLoading, setShowLoading] = useState(false)
const [openCommandDialog, setOpenCommandDialog] = useState(false)
const [deviceCommandToExecute, setDeviceCommandToExecute] = useState(null)
const [inputArgsValue, setInputArgsValue] = useState({})


// const initDeviceCommands = (content) => {
//     let supportedCommands = content?.req_obj_results[0].supported_objs[0].supported_commands
    
//     if (supportedCommands === undefined){
//         return paramsToFetch
//     }

//     let commands = {}

//     for(let i =0; i < supportedCommands.length; i++){
//         let command = supportedCommands[i]
//         commands[command.command_name] = {
//             "type":command["command_type"]
//         }
//     }

//     console.log("commands:", commands)
//     setDeviceCommands(commands)
// }   

const initialize = async (raw) => {
    let content = await getDeviceParameters(raw)
    setDeviceParameters(content)
    //initDeviceCommands(content)
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
        if (result.status === 401){
            router.push("/auth/login")
        }
        console.log('Please check your email and password');
    }else {
        return result.json()
    }
}

/*
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
    }else if (result.status === 401){
    router.push("/auth/login")
}else{
        return result.json()
    }
}*/

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
//   const updateDeviceParametersMultiInstance = async (param) =>{
//     console.log("UpdateDeviceParametersMultiInstance => param = ", param)

//     let raw = JSON.stringify({
//         "obj_paths": [param],
//         "first_level_only" : true,
//         "return_commands" : true,
//         "return_events" : true,
//         "return_params" : true 
//     })

//     let response = await getDeviceParameterInstances(raw)
//     console.log("response:", response)

//     let instancesToGet = []
//     if (response.req_path_results[0].curr_insts) {
//         let supportedParams = response.req_path_results[0].curr_insts
//         let instances = () => {
//             for (let i =0; i < supportedParams.length ;i++){
//                 instancesToGet.push(supportedParams[i].instantiated_obj_path)
//             }
//         }
//         instances()
//     }else{
//         instancesToGet.push(response.req_path_results[0].requested_path)
//     }

//     let rawInP = JSON.stringify({
//         "obj_paths": instancesToGet,
//         "first_level_only" : true,
//         "return_commands" : true,
//         "return_events" : true,
//         "return_params" : true 
//     })

//     let resultParams = await getDeviceParameters(rawInP)
//     console.log("result params:", resultParams)
//     setDeviceParameters(resultParams)


//     let paramsToFetch = []

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
//  }


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

    let paramsInfo = {}
    let commandsInfo = {}

    let supportedParams = content.req_obj_results[0].supported_objs[0].supported_params
    let supportedCommands = content.req_obj_results[0].supported_objs[0].supported_commands

    let parametersToFetch = () => {
        let paramsToFetch = []
        for (let i =0; i < supportedParams.length ;i++){
            
            let supported_obj_path = content.req_obj_results[0].supported_objs[0].supported_obj_path.replaceAll("{i}","*")
            let param = supportedParams[i]
            
            paramsToFetch.push(supported_obj_path+param.param_name)

            paramsInfo[param.param_name] = {
                "value_change":param["value_change"],
                "value_type":param["value_type"],
                "access": param["access"],
                "value": "-",
            }
        }

        if (supportedCommands === undefined){
            return paramsToFetch
        }

        for(let i =0; i < supportedCommands.length; i++){
            let command = supportedCommands[i]
            commandsInfo[command.command_name] = {
                "type":command["command_type"]
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
        console.log("/-------------------------------------------------------/")

        let values = {}
        let commands = {}

        console.log("VALUES:",values)
        result.req_path_results.map((x)=>{
            if (!x.resolved_path_results){
                values[x.requested_path] = {}
                setDeviceParametersValue(values)
                return
            }

            let paths = x.requested_path.split(".")
            if(paths[paths.length -2] == "*"){
                x.resolved_path_results.map(y=>{
                    // console.log(y.result_params)
                    // console.log(y.resolved_path)
                    let key = Object.keys(y.result_params)[0]
                    // console.log(key)
                    // console.log(paramsInfo[key].value)
                    // console.log(paramsInfo[key])
                    // console.log(y.result_params[key])
                    // console.log({[key]:paramsInfo[key]})

                    console.log("Take a look here mate: ",{...paramsInfo[key], value: y.result_params[key]})
                    if (!values[y.resolved_path]){
                        values[y.resolved_path] = []
                    }

                    if (!commands[y.resolved_path]){
                        commands[y.resolved_path] = []
                    }

                    if (y.result_params[key] == ""){
                        y.result_params[key] = "\"\""
                    }
                    
                    values[y.resolved_path].push({[key]:{...paramsInfo[key], value: y.result_params[key]}})
                })
            }else{
                Object.keys(x.resolved_path_results[0].result_params).forEach((key, index) =>{
                    if (x.resolved_path_results[0].result_params[key] != ""){
                        paramsInfo[key].value = x.resolved_path_results[0].result_params[key]
                    }else{
                        paramsInfo[key].value = "\"\""
                    }
                    values = paramsInfo
                })
            }

            console.log("values:", values)
            setDeviceParametersValue(values)
            console.log("commands:", commandsInfo)
            setDeviceCommands(commandsInfo)
        })
        
        console.log("/-------------------------------------------------------/")
        setDeviceParameters(content)
    }else{
        console.log("fixme")
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
    }else if (result.status === 401){
    router.push("/auth/login")
}else{
        return result.json()
    }

  }

  function isInteger(value) {
    return /^\d+$/.test(value);
  }

  const inputjow = () => {
    if (inputArgsValue === ""){
        return {"":""}
    }else{
        return inputArgsValue
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
                        <ShowPath x={x} setShowLoading={setShowLoading} updateDeviceParameters={updateDeviceParameters} router={router}></ShowPath>
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
                {   x.supported_params &&
                    <ShowParamsWithValues 
                    x={x} 
                    deviceParametersValue={deviceParametersValue} 
                    setOpen={setOpen} 
                    setParameter={setParameter}
                    setParameterValue={setParameterValue}
                    deviceParameters={deviceParameters}
                    setShowLoading={setShowLoading}
                    router={router}
                    updateDeviceParameters={updateDeviceParameters}
                    deviceCommands={deviceCommands}
                    openCommandDialog={openCommandDialog}
                    />
                }
                { x.supported_commands && Object.keys(deviceCommands).length == 0 &&
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
                            secondaryAction={
                                <IconButton onClick={()=> {
                                    setDeviceCommandToExecute(
                                        {
                                            [x.supported_obj_path+y.command_name]:
                                            {"input_arg_names":
                                                [
                                                    y.input_arg_names
                                                ]
                                            }
                                        }
                                    )
                                    setOpenCommandDialog(true)
                                    }}>
                                    <SvgIcon>
                                        <PlayCircleIcon >
                                        </PlayCircleIcon>
                                    </SvgIcon>
                                </IconButton>
                            }
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
                    <Dialog open={open} 
                    slotProps={{ backdrop: { style: { backgroundColor: 'rgba(255,255,255,0.5)' } } }}
                    >
                    <DialogContent>
                    <DialogContentText>
                        {parameter}
                    </DialogContentText>
                    <TextField
                        autoFocus
                        margin="dense"
                        id="parameterValue"
                        fullWidth
                        variant="standard"
                        defaultValue={parameterValue}
                        autoComplete='off'
                        onChange={(e)=>setParameterValueChange(e.target.value)}
                    />
                    </DialogContent>
                    <DialogActions>
                    <Button onClick={()=>{setOpen(false)}}>Cancel</Button>
                    <Button onClick={async ()=>{
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    let params = parameter.split(".")
    let objToChange;
    let parameterToChange;
    console.log("params",params)
    parameterToChange = params.pop()
    objToChange = params.join(".")

    var requestOptions = {
        method: 'PUT',
        headers: myHeaders,
        redirect: 'follow',
        body: JSON.stringify(
            {
                "allow_partial":true,
                "update_objs":[
                    {
                        "obj_path":objToChange,
                        "param_settings":[
                            {
                            "param":parameterToChange,
                            "value":parameterValueChange,
                            "required":true
                            }
                        ]
                    }
                ]
            }
        )
    };

    console.log(requestOptions.body)

    setOpen(false)
    setShowLoading(true)
    let result = await (await fetch(`${process.env.NEXT_PUBLIC_REST_ENPOINT}/device/${router.query.id[0]}/set`, requestOptions))
    if (result.status != 200) {
        setShowLoading(false)
        throw new Error('Please check your email and password');
    }else{
        setShowLoading(false)
        let response = await result.json()
        let feedback = JSON.stringify(response, null, 2)

        if (response.updated_obj_results[0].oper_status.OperStatus["OperSuccess"] === undefined) {
            console.log("Error to set parameter change")
            setErrorModalText(feedback)
            setErrorModal(true)
            return
        }

        //Means it has more than one instance
        if(isInteger(params[params.length -1])){
            setDeviceParametersValue((prevState) => ({
                ...prevState, [objToChange+"."]: prevState[objToChange+"."].map(el => {
                    if (el[parameterToChange] !== undefined){
                        console.log(el[parameterToChange])
                        el[parameterToChange].value = parameterValueChange
                        return el
                    }else{
                        console.log(el)
                        return el
                    }
                })
            }));
        }else{
            setDeviceParametersValue((prevState) => ({
                ...prevState, 
                [parameterToChange] : {
                    ...prevState[parameterToChange], 
                    value: parameterValueChange}
            }));
        }

        setOpen(false)
    }
                    }}>Apply</Button>
                    </DialogActions>
                </Dialog>
                <Dialog open={errorModal} 
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
                {deviceCommandToExecute && <Dialog open={openCommandDialog} 
                    slotProps={{ backdrop: { style: { backgroundColor: 'rgba(255,255,255,0.5)' } } }}
                    fullWidth={ true } 
                    maxWidth={"md"}   
                    scroll={"paper"}
                    aria-labelledby="scroll-dialog-title"
                    aria-describedby="scroll-dialog-description" 
                >
                <DialogTitle id="scroll-dialog-title">
                {Object.keys(deviceCommandToExecute)[0]}
                </DialogTitle>    
                    <DialogContent dividers={scroll === 'paper'}>
                    {deviceCommandToExecute[Object.keys(deviceCommandToExecute)[0]].input_arg_names[0]!=undefined && <DialogContentText id="scroll-dialog-description" tabIndex={-1}>
                    Input Arguments:
                    </DialogContentText>}
                    {deviceCommandToExecute[Object.keys(deviceCommandToExecute)[0]].input_arg_names[0] !=undefined && 
                    deviceCommandToExecute[Object.keys(deviceCommandToExecute)[0]].input_arg_names?.map(arg => {
                        return <TextField
                        autoFocus
                        margin="dense"
                        id={arg}
                        label={arg}
                        type="text"
                        onChange={(e)=> {setInputArgsValue(prevState=>{
                            return {...prevState, [arg] : e.target.value}
                        })}}
                        onClick={()=>{console.log(deviceCommandToExecute)}}
                        value={inputArgsValue[arg]}
                      />
                    })}
                    </DialogContent>
                    <DialogActions>
                    <Button onClick={()=>{
                        setInputArgsValue("")
                        setDeviceCommandToExecute(null)
                        setOpenCommandDialog(false)
                    }}>Cancel</Button>
                    <Button onClick={async ()=>{
                       let raw = JSON.stringify(
                        {
                            "command": Object.keys(deviceCommandToExecute)[0],
                            "input_args": inputjow(),
                            "send_resp": true
                        }
                       )
                       console.log("deviceOperate => obj = ", raw)
                       var myHeaders = new Headers();
                       myHeaders.append("Content-Type", "application/json");
                       myHeaders.append("Authorization", localStorage.getItem("token"));
                   
                       var requestOptions = {
                           method: 'PUT',
                           headers: myHeaders,
                           redirect: 'follow',
                           body: raw
                       };
                       setShowLoading(true)
                       let result = await fetch(`${process.env.NEXT_PUBLIC_REST_ENPOINT}/device/${router.query.id[0]}/operate`, requestOptions)
                       let content = await result.json()
                       if (result.status != 200) {
                           setShowLoading(false)
                           if (result.status === 401){
                            router.push("/auth/login")
                            }
                            setInputArgsValue("")
                            setDeviceCommandToExecute(null)
                            setOpenCommandDialog(false)
                            setShowLoading(false)
                       }else{
                            setInputArgsValue("")
                            setDeviceCommandToExecute(null)
                            setOpenCommandDialog(false)
                            setShowLoading(false)
                            if (content.operation_results[0].OperationResp.CmdFailure != undefined){
                                setErrorModalText(JSON.stringify(content, null, 2))
                                setErrorModal(true)
                                return
                            }
                       }
                    }}>Apply</Button>
                    </DialogActions>
                </Dialog>}
        <Backdrop
            sx={{ 
            color: '#fff', 
            zIndex: (theme) => theme.zIndex.drawer + 1, 
            overflow: 'hidden'
            }}
            open={showLoading}
            >
            <CircularProgress color="inherit" />
        </Backdrop>
    </Card> 
    :
    <Box sx={{display:'flex',justifyContent:'center'}}>
        <CircularProgress color="inherit" />
    </Box>
  )
};
