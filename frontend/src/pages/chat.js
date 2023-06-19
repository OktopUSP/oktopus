import React, { useEffect, useState, useContext } from "react";
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import PhoneIcon from "@heroicons/react/24/solid/PhoneIcon";
import PhoneXMarkIcon from "@heroicons/react/24/solid/PhoneXMarkIcon"
import { 
    Card,
    Box,
    CardContent,
    Container,
    SvgIcon,
    CircularProgress,
    Avatar,
    Backdrop,
} from "@mui/material";
import { WsContext } from "src/contexts/socketio-context";

const Page = () => {

    //const [isConnected, setIsConnected] = useState(socket.connected);
    const [users, setUsers] = useState([]) 
    //const [onlineUsers, setOnlineUsers] = useState([])

    const ws = useContext(WsContext)

    useEffect(()=>{
        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", localStorage.getItem("token"));
        
        var requestOptions = {
            method: 'GET',
            headers: myHeaders,
            redirect: 'follow'
        };

        fetch(`${process.env.NEXT_PUBLIC_REST_ENPOINT}/users`,requestOptions)
        .then(response => response.json())
        .then(result => {
            // let teste =  JSON.stringify(JSON.parse(result), null, 2)
            setUsers(result)
        })
        .catch(error => console.log('error', error));
    },[])

    const renderUsers = () => {
        console.log("users: ", users)
        console.log("wsUsers: ", ws.users)
        if(users.length == 0){
            console.log("users is empty")
            return (
                <div style={{display:'flex', justifyContent:'center'}} height={'100%'} >
                    <CircularProgress color="inherit" width='100%'/>
                </div>
            )
        }else {
            return (
                <Card sx={{
                    display: 'flex',
                    justifyContent:'center',
                }}>
                    <CardContent>
                        
                        <Container sx={{display:'flex',justifyContent:'center'}}>
                        {users.map((x)=> {

                            let color = "#CB1E02"
                            let status = "offline"

                            if (ws.users.findIndex(y => y.name === x.email) >= 0){
                                console.log("user: "+x.email+" is online")
                                //color = "#11ADFB"
                                color = "#17A000"
                                status = "online"
                            }

                            if (x.email !== window.sessionStorage.getItem("email")){
                                return (
                                    <Box sx={{margin:"30px",textAlign:'center'}}>
                                        <Avatar
                                        sx={{
                                            height: 150,
                                            width: 150,
                                            border: '3px solid '+color
                                        }}
                                        src={"/assets/avatars/default-avatar.png"}
                                        />
                                        <div style={{marginTop:'10px'}}>
                                        </div>
                                        <SvgIcon
                                        sx={{cursor:'pointer'}}
                                        >
                                            {status === "online" ?
                                                <PhoneIcon 
                                                color={color}
                                                onClick={()=>{
                                                    console.log("call", x.email)
                                                }}
                                                title={"call"}
                                                />
                                            :
                                            <PhoneXMarkIcon 
                                            color={color}
                                            onClick={()=>{
                                                console.log("call", x.email)
                                            }}
                                            title={"offline"}
                                            />
                                            }
                                        </SvgIcon>
                                        <p style={{marginTop:'-2.5px'}}>{x.email}</p>
                                    </Box> 
                                )   
                            }
                        })}
                        </Container>
                    </CardContent>
                </Card>
            )
        }
    }

    return(ws.users ?
        <Box
        component="main"
        sx={{
                flexGrow: 1,
                py: 10,
                alignItems: 'center',
                flexDirection: 'column',
            }}
        >
        <Container maxWidth="md">
            {renderUsers()}
        </Container>
        </Box>
        :
        <CircularProgress color="inherit" />
    )
}

Page.getLayout = (page) => (
    <DashboardLayout>
        {page}
    </DashboardLayout>
);

export default Page;