import { useCallback, useState, useEffect } from 'react';
import Head from 'next/head';
import MagnifyingGlassIcon from '@heroicons/react/24/solid/MagnifyingGlassIcon';
import PlusIcon from '@heroicons/react/24/solid/PlusIcon';
import { Box, Button, Container, Stack, SvgIcon, Tooltip, Typography, IconButton,
DialogActions,
Dialog,
DialogTitle,
DialogContent,
TextField,
InputAdornment,
Input, Backdrop, CircularProgress,
InputLabel, FormControl,
OutlinedInput,
Card } from '@mui/material';
import EyeIcon from '@heroicons/react/24/outline/EyeIcon';
import EyeSlashIcon from '@heroicons/react/24/outline/EyeSlashIcon';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { CredentialsTable } from 'src/sections/credentials/credentials-table';
import InformactionCircleIcon from '@heroicons/react/24/outline/InformationCircleIcon';
import { useAuth } from 'src/hooks/use-auth';
import { useRouter } from 'next/router';

const Page = () => {
  const auth = useAuth();
  const router = useRouter();

  const [page, setPage] = useState(0);
  const [devices, setDevices] = useState({});
  const [addDeviceDialogOpen, setAddDeviceDialogOpen] = useState(false);
  const [newDeviceData, setNewDeviceData] = useState({});
  const [showPassword, setShowPassword] = useState(false); 
  const [isUsernameEmpty, setIsUsernameEmpty] = useState(false);
  const [isUsernameExistent, setIsUsernameExistent] = useState(false);
  const [creatingNewCredential, setCreatingNewCredential] = useState(false);
  const [loading, setLoading] = useState(true);
  const [credentialNotFound, setCredentialNotFound] = useState(false);

  const deleteCredential = (id) => {
    console.log("request to delete credentials: ", id)

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var requestOptions = {
      method: 'DELETE',
      headers: myHeaders,
      redirect: 'follow'
    }

    return fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device/auth?id=${id}`, requestOptions)
      .then(response => {
        if (response.status === 401) {
          router.push("/auth/login")
        } else if (response.status === 403) {
          return router.push("/403")
        }
        let copiedDevices = {...devices}
        delete copiedDevices[id];
        setDevices(device => ({
            ...copiedDevices
        }))
      })
      .catch(error => {
        return console.error('Error:', error)
      });
  }

  const createCredential = async (data) => {
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var raw = JSON.stringify(data);

    var requestOptions = {
      method: 'POST',
      headers: myHeaders,
      body: raw,
      redirect: 'follow'
    };

    let result = await fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device/auth`, requestOptions)

    if (result.status == 200) {
      console.log("user created: deu boa raça !!")
    }else if (result.status == 403) {
      console.log("num tenx permissão, seu boca de sandália")
      setCreatingNewCredential(false)
      return router.push("/403")
    }else if (result.status == 401){
      console.log("taix nem autenticado, sai fora oh")
      setCreatingNewCredential(false)
      return router.push("/auth/login")
    }else if (result.status == 409){
      console.log("usuário já existe, seu boca de bagre")
      setIsUsernameExistent(true)
      setCreatingNewCredential(false)
      return
    }else if (result.status == 400){
      console.log("faltou mandar dados jow")
      setAddDeviceDialogOpen(false)
      setNewDeviceData({})
      setIsUsernameEmpty(false)
      setIsUsernameExistent(false)
      setCreatingNewCredential(false)
      return
    }else {
      console.log("agora quebrasse ux córno mô quiridu")
      const content = await result.json()
      setCreatingNewCredential(false)
      throw new Error(content);
    }
    setAddDeviceDialogOpen(false)
    let newData = {} 
    newData[data.id] = data.password
    setDevices(prevState =>({...prevState, ...newData}))
    setNewDeviceData({})
    setIsUsernameEmpty(false)
    setIsUsernameExistent(false)
    setCreatingNewCredential(false)
  }

  const fetchCredentials = async (id) => {
    console.log("fetching credentials data...")
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var requestOptions = {
      method: 'GET',
      headers: myHeaders,
      redirect: 'follow'
    }

    let url = `${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device/auth`
    if (id !== undefined && id !== "") {
      url += "?id="+id
    }

    return fetch(url, requestOptions)
      .then(response => {
        if (response.status === 401) {
          router.push("/auth/login")
        } else if (response.status === 403) {
          return router.push("/403")
        }else if (response.status === 404) {
          setLoading(false)
          setCredentialNotFound(true)
          console.log("credential not found: ", credentialNotFound)
          return 
        }
        return response.json()
      })
      .then(json => {
        if (json === undefined) {
          return
        }
        console.log("devices credentials: ", json)
        setDevices(json)
        setLoading(false)
        setCredentialNotFound(false)
      })
      .catch(error => {
        setLoading(false)
        setCredentialNotFound(false)
        return console.error('Error:', error)
      });
  }

  useEffect(() => {
    fetchCredentials()
  }, []);

  return (
    <>
      <Head>
        <title>
          Devices Credentials | Oktopus
        </title>
      </Head>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          py: 8
        }}
      >
        <Container maxWidth="xl">
          <Stack spacing={3}>
            <Stack
              direction="row"
              justifyContent="space-between"
              spacing={4}
            >
              <Stack spacing={1} direction="row" alignItems={'center'}>
                <Typography variant="h4">
                  Devices Credentials 
                </Typography>
                <Tooltip title="Defines username and password for devices authentication, this must be enabled through environment vars." placement="top">
                    <IconButton>
                        <SvgIcon>
                            <InformactionCircleIcon />
                        </SvgIcon>
                    </IconButton>
                </Tooltip>
              </Stack>
              <div>
                <Button
                  startIcon={(
                    <SvgIcon fontSize="small">
                      <PlusIcon />
                    </SvgIcon>
                  )}
                  onClick={() => setAddDeviceDialogOpen(true)}
                  variant="contained"
                >
                  Add
                </Button>
              </div>
            </Stack>
            <Card sx={{ p: 2 }}>
            <OutlinedInput
            defaultValue=""
            fullWidth
            placeholder="Search credentials by username"
            onKeyDownCapture={(e) => {
                if (e.key === 'Enter') {
                  console.log("Fetch credentials per username: ", e.target.value)
                  fetchCredentials(e.target.value)
                }
            }}
            startAdornment={(
                <InputAdornment position="start">
                <SvgIcon
                    color="action"
                    fontSize="small"
                >
                    <MagnifyingGlassIcon />
                </SvgIcon>
                </InputAdornment>
            )}
            sx={{ maxWidth: 500 }}
            />
        </Card>
            {!loading ? (credentialNotFound ? 
                <Box
                sx={{
                    display: 'flex',
                    justifyContent: 'center'
                }}
                >
                <p>Credential Not Found</p>
                </Box>:
            <CredentialsTable
              items={devices}
              deleteCredential={deleteCredential}
            //   onPageChange={handlePageChange}
              page={page}
            />):
                <CircularProgress/>
            }
          </Stack>
        </Container>
      </Box>
      <Dialog
      open={addDeviceDialogOpen}
      onClose={() => {
        setAddDeviceDialogOpen(false)
        setIsUsernameEmpty(false)
        setIsUsernameExistent(false)
        setNewDeviceData({})
      }}
      >
        <DialogTitle>Create New Credentials</DialogTitle>
        <DialogContent>
          <Stack
            alignItems="center"
            direction="row"
            spacing={2}
          >
                        <FormControl sx={{ m: 1, width: '25ch' }} variant="standard">
            <TextField
              inputProps={{
                form: {
                  autocomplete: 'new-password',
                },
              }}
              // focused={isUsernameEmpty}
              // color={isUsernameEmpty ? "error" : "primary"}
              error={isUsernameEmpty || isUsernameExistent}
              helperText={isUsernameEmpty ? "Username invalid": (isUsernameExistent ? "Username already exists" : "")}
              autoFocus
              required
              margin="dense"
              id="username"
              name="username"
              label="Username"
              type="username"
              fullWidth
              onChange={
                (event) => {
                  setNewDeviceData({...newDeviceData, id: event.target.value})
                }
              }
              variant="standard">
            </TextField>
            </FormControl>
            <FormControl sx={{ m: 1, width: '25ch' }} variant="standard">
                <InputLabel htmlFor="standard-adornment-password">Password</InputLabel>
                <Input
                id="standard-adornment-password"
                type={showPassword ? 'text' : 'password'}
                label="Password"
                endAdornment={
                    <InputAdornment position="end">
                    <IconButton
                        aria-label="toggle password visibility"
                        onClick={()=>{
                        setShowPassword(!showPassword)
                        }}
                    >
                        <SvgIcon>
                        {showPassword ? <EyeSlashIcon /> : <EyeIcon />}
                        </SvgIcon>
                    </IconButton>
                    </InputAdornment>
                }
                onChange={
                    (event) => {
                    setNewDeviceData({...newDeviceData, password: event.target.value})
                    }
                }
                />
            </FormControl>
          </Stack>
          <Stack
            alignItems="center"
            direction="row"
            spacing={2}
          >
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={
            () => {
              setAddDeviceDialogOpen(false)
              setIsUsernameEmpty(false)
              setIsUsernameExistent(false)
              setNewDeviceData({})
            }
          }>Cancel</Button>
          <Button onClick={()=>{
            console.log("new user data: ", newDeviceData)
            if (newDeviceData.id === undefined || newDeviceData.id === "") {
              setIsUsernameEmpty(true)
              return
            }else{
              setIsUsernameEmpty(false)
            }
            setCreatingNewCredential(true)
            createCredential(newDeviceData)
          }}>Confirm</Button>
        </DialogActions>
        {
        <Backdrop
        sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
        open={creatingNewCredential}
        >
        <CircularProgress color="inherit" />
        </Backdrop>
      }
      </Dialog>
    </>
  );
};

Page.getLayout = (page) => (
  <DashboardLayout>
    {page}
  </DashboardLayout>
);

export default Page;