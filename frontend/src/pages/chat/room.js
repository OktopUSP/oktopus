import React, { useEffect, useState, useContext } from "react";
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { 
    Card,
    Box,
    CardContent,
    Container,
    SvgIcon,
    CircularProgress,
    Avatar,
    Tooltip
} from "@mui/material";
import { WsContext } from "src/contexts/socketio-context";
import { useRouter } from "next/router";

const Page = (props) => {

    const { callUser, callAccepted, myVideo, userVideo, callEnded, stream, call, setStream } = 
    useContext(WsContext);
    const router = useRouter()

    const stopCamera = () => {
        // stream.getTracks().forEach(function(track) {
        //     track.stop();
        //   });
        //console.log(stream)
        window.location.reload() //TODO: find better way to stop recording user
    }

    useEffect(()=>{
        callUser(router.query.user)
        navigator.mediaDevices
        .getUserMedia({ video: true, audio: true })
        .then((currentStream) => {
          setStream(currentStream);
          if (myVideo.current) {
            myVideo.current.srcObject = currentStream;
          }
        })
        .catch((err)=>{
            console.log('You cannot place/ receive a call without granting video and audio permissions! Please change your settings to use Oktopus calls.')
            console.log(err)
        })

        return(stopCamera)
    },[])

    return (
        <Card>
            <CardContent>
            { myVideo && 
            <video 
            className="userVideo" 
            playsInline 
            muted 
            ref={myVideo} 
            autoPlay />
            }
            </CardContent>
        </Card>
    )
}

Page.getLayout = (page) => (
    <DashboardLayout>
        {page}
    </DashboardLayout>
);

export default Page;