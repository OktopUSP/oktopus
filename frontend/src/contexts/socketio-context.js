import { createContext, useContext, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import io from 'socket.io-client';
import { useAuth } from 'src/hooks/use-auth';

// The role of this context is to propagate socketio io state through app tree
export const WsContext = createContext({ undefined });

export const WsProvider = (props) => {
  const { children } = props;
  const [users, setUsers] = useState([])
  const auth = useAuth()
  const initialize = async () => {
    // Prevent from calling twice in development mode with React.StrictMode enable
    const socket = io(process.env.NEXT_PUBLIC_WS_ENPOINT)

    socket.on('connect', () => {
        console.log('[IO] Connect => A new connection has been established')

        socket.on("users", (data) => {
            setUsers(data)
            console.log("data received from users event: ", users)
        })

        socket.emit("newuser",{
            id:socket.id,
            name: window.sessionStorage.getItem("email")
        })
    })
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
