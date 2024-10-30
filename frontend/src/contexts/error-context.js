import { createContext, useContext, useState } from 'react';


export const AlertContext = createContext({ undefined });

export const AlertProvider = (props) => {
    const { children } = props;
    /*
    {
        severity: '', // options => error, warning, info, success
        message: '',
        title: '',
    }
    */
    // const [alert, setAlert] = useState(null);
    const [alert, setAlert] = useState();
    
    return (
        <AlertContext.Provider
        value={{
            alert,
            setAlert,
        }}
        >
        {children}
        </AlertContext.Provider>
    );
};

export const AlertConsumer = AlertContext.Consumer;

export const useAlertContext = () => useContext(AlertContext);