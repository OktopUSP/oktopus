import { createContext, useContext, } from 'react';
import { useRouter } from 'next/router';
import { useAlertContext } from './error-context';

export const BackendContext = createContext({ undefined });

export const BackendProvider = (props) => {
    const { children } = props;

    const { setAlert } = useAlertContext();
    
    const router = useRouter();

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    const httpRequest = async (path, method, body, headers, encoding) => {
        
        var requestOptions = {
            method: method,
            redirect: 'follow',
        };

        if (body) {
            requestOptions.body = body
        }

        console.log("headers:", headers)
        if (headers) {
            requestOptions.headers = headers
        }else {
            requestOptions.headers = myHeaders
        }

        const response = await fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}${path}`, requestOptions)
        console.log("status:", response.status)
        if (response.status != 200) {
            if (response.status == 401) {
                router.push("/auth/login")
            }
            if (response.status == 204 || response.status == 201) {
                return {status : response.status, result: null}
            }
            const data = await response.text()
            // setAlert({
            //     severity: "error",
            //     title: "Error",
            //     message: `Status: ${response.status}, Content: ${data}`,
            // })
            setAlert({
                severity: "error",
                // title: "Error",
                message: `${data}`,
            })
            return {status : response.status, result: null}
        }
        if (encoding) {
            console.log("encoding:", encoding)
            if (encoding == "text") {
                const data = await response.text()
                return {status: response.status, result: data}
            }else if (encoding == "blob") {
                const data = await response.blob()
                return {status: response.status, result: data}
            }else if (encoding == "json") {
                const data = await response.json()
                return {status: response.status, result: data}
            }else{
                return {status: response.status, result: response}
            }
        }
        const data = await response.json()
        return {status: response.status, result: data}
    }
    
    return (
        <BackendContext.Provider
        value={{
            httpRequest,
            setAlert,
        }}
        >
        {children}
        </BackendContext.Provider>
    );
};

export const BackendConsumer = BackendContext.Consumer;

export const useBackendContext = () => useContext(BackendContext);