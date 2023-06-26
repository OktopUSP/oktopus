import { createContext, useContext, useEffect, useState, useRef } from 'react';
import PropTypes from 'prop-types';
import io from 'socket.io-client';
import { useAuth } from 'src/hooks/use-auth';
import Peer from "simple-peer";


// The role of this context is to propagate socketio io state through app tree
export const WsContext = createContext({ undefined });

export const WsProvider = (props) => {
  const { children } = props;
  const [users, setUsers] = useState(null)
  const [callAccepted, setCallAccepted] = useState(false);
	const [callEnded, setCallEnded] = useState(false);
	const [stream, setStream] = useState();
	const [name, setName] = useState("");
	const [call, setCall] = useState({});

	const myVideo = useRef();
	const userVideo = useRef();
	const connectionRef = useRef();
  const auth = useAuth()
  const socket = io(process.env.NEXT_PUBLIC_WS_ENPOINT)

  const initialize = async () => {
    // Prevent from calling twice in development mode with React.StrictMode enable

    socket.on('connect', () => {
        console.log('[IO] Connect => A new connection has been established')

        socket.on("users", (data) => {
            setUsers(data)
            console.log("data received from users event: ", data)
        })

        socket.emit("newuser",{
            id: socket.id,
            name: window.sessionStorage.getItem("email")
        })

      socket.on("callUser", ({ from, name: callerName, signal }) => {
        console.log("you're receiving call brow")
        setCall({ isReceivingCall: true, from, name: callerName, signal });
      });

      socket.on('disconnect', function(){
        
    });

    })
  };

  const answerCall = () => {
		setCallAccepted(true);

		const peer = new Peer({ 
      initiator: false, 
      trickle: false, 
      stream: stream,
      config: {
        iceServers: [
            {url:'stun:stun.l.google.com:19302'},
            {url:'stun:stun1.l.google.com:19302'},
        ]
    },
    });

		peer.on("signal", (data) => {
			socket.emit("answerCall", { signal: data, to: call.from });
		});

		peer.on("stream", (currentStream) => {
			userVideo.current.srcObject = currentStream;
		});

		peer.signal(call.signal);

		connectionRef.current = peer;
	};


	const callUser = (id) => {

    console.log("calling user ",id)
		const peer = new Peer({ initiator: true, trickle: false, stream:stream, config: {
      iceServers: [
          {url:'stun:stun.l.google.com:19302'},
          {url:'stun:stun1.l.google.com:19302'},
          {url:'stun:stun2.l.google.com:19302'},
          {url:'stun:stun3.l.google.com:19302'},
          {url:'stun:stun4.l.google.com:19302'},
      ]
  }, });

		peer.on("signal", (data) => {
			socket.emit("callUser", {
				userToCall: id,
				signalData: data,
				from: window.sessionStorage.getItem("email"),
			});
		});

		peer.on("stream", (currentStream) => {
			userVideo.current.srcObject = currentStream;
		});

		socket.on("callAccepted", (signal) => {
			setCallAccepted(true);

			peer.signal(signal);
		});

		connectionRef.current = peer;
	};

  const leaveCall = () => {
		setCallEnded(true);

		connectionRef.current.destroy();

		window.location.reload();
	};

  useEffect(
    () => {
        if(auth.isAuthenticated){
            initialize();
        }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [auth.isAuthenticated]
  );

  return (
    <WsContext.Provider
      value={{
        users,
        call,
				callAccepted,
				myVideo,
				userVideo,
				stream,
				callEnded,
				callUser,
				leaveCall,
				answerCall,
        setStream
      }}
    >
      {children}
    </WsContext.Provider>
  );
};

WsProvider.propTypes = {
  children: PropTypes.node
};

export const WsConsumer = WsContext.Consumer;

export const useWsContext = () => useContext(WsContext);
