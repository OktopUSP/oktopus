import { useCallback, useMemo, useState, useEffect } from 'react';
import Head from 'next/head';
import { subDays, subHours } from 'date-fns';
import ArrowDownOnSquareIcon from '@heroicons/react/24/solid/ArrowDownOnSquareIcon';
import ArrowUpOnSquareIcon from '@heroicons/react/24/solid/ArrowUpOnSquareIcon';
import PlusIcon from '@heroicons/react/24/solid/PlusIcon';
import { Box, Button, CircularProgress, Container, Dialog, DialogContent, DialogTitle, Stack, SvgIcon, Typography,
  DialogActions,
  TextField,
  Backdrop,
} from '@mui/material';
import { useSelection } from 'src/hooks/use-selection';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { CustomersTable } from 'src/sections/customer/customers-table';
import { CustomersSearch } from 'src/sections/customer/customers-search';
import { applyPagination } from 'src/utils/apply-pagination';
import { useAuth } from 'src/hooks/use-auth';
import { useRouter } from 'next/router';
import { is } from 'date-fns/locale';
import { set } from 'nprogress';

const Page = () => {

  const auth = useAuth();
  const router = useRouter();

  const validateEmail = (email) => {
    return email.match(
      /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
    );
  };  

  //const [page, setPage] = useState(0);
  //const [rowsPerPage, setRowsPerPage] = useState(5);
  const [loading, setLoading] = useState(true);
  const [creatingNewUser, setCreatingNewUser] = useState(false);
  const [users, setUsers] = useState([]);
  const [selected, setSelected] = useState([]);
  const [addDeviceDialogOpen, setAddDeviceDialogOpen] = useState(false);
  const [newUserData, setNewUserData] = useState({});
  const [isPasswordEmpty, setIsPasswordEmpty] = useState(false);
  const [isEmailEmpty, setIsEmailEmpty] = useState(false);
  const [isEmailExistent, setIsEmailExistent] = useState(false);

  const deleteUser = (id) => {
    console.log("request to delete user: ", id)

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var requestOptions = {
      method: 'DELETE',
      headers: myHeaders,
      redirect: 'follow'
    }

    return fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/auth/delete/${id}`, requestOptions)
      .then(response => {
        if (response.status === 401) {
          router.push("/auth/login")
        } else if (response.status === 403) {
          return router.push("/403")
        }
        setUsers(users.filter(user => user.email !== id))
      })
      .catch(error => {
        return console.error('Error:', error)
      });
  }


  const fetchUsers = async () => {
    console.log("fetching users data...")
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var requestOptions = {
      method: 'GET',
      headers: myHeaders,
      redirect: 'follow'
    }

    return fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/users`, requestOptions)
      .then(response => {
        if (response.status === 401) {
          router.push("/auth/login")
        } else if (response.status === 403) {
          return router.push("/403")
        }
        return response.json()
      })
      .then(json => {
        console.log("users: ", json)
        setUsers(json)
        // setPages(json.pages + 1)
        // setPage(json.page +1)
        // setDevices(json.devices)
        setLoading(false)
      })
      .catch(error => {
        return console.error('Error:', error)
      });
  }

  useEffect(() => {
    // if (auth.user.token) {
    //   console.log("auth.user.token =", auth.user.token)
    // }else{
    //   auth.user.token = localStorage.getItem("token")
    // }
    //console.log("auth.user.token =", auth.user.token)
    fetchUsers()
  }, []);

  // const handlePageChange = useCallback(
  //   (event, value) => {
  //     setPage(value);
  //   },
  //   []
  // );

  // const handleRowsPerPageChange = useCallback(
  //   (event) => {
  //     setRowsPerPage(event.target.value);
  //   },
  //   []
  // );

  const createUser = async (data) => {
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

    let result = await fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/auth/register`, requestOptions)

    if (result.status == 200) {
      console.log("user created: deu boa raça !!")
    }else if (result.status == 403) {
      console.log("num tenx permissão, seu boca de sandália")
      setCreatingNewUser(false)
      return router.push("/403")
    }else if (result.status == 401){
      console.log("taix nem autenticado, sai fora oh")
      setCreatingNewUser(false)
      return router.push("/auth/login")
    }else if (result.status == 409){
      console.log("usuário já existe, seu boca de bagre")
      setIsEmailExistent(true)
      setCreatingNewUser(false)
      return
    }else if (result.status == 400){
      console.log("faltou mandar dados jow")
      setAddDeviceDialogOpen(false)
      setNewUserData({})
      setIsPasswordEmpty(false)
      setIsEmailEmpty(false)
      setIsEmailExistent(false)
      setCreatingNewUser(false)
      return
    }else {
      console.log("agora quebrasse ux córno mô quiridu")
      const content = await result.json()
      setCreatingNewUser(false)
      throw new Error(content);
    }
    setAddDeviceDialogOpen(false)
    data["_id"] = data.email
    data["createdAt"] = new Date().toLocaleDateString('es-pa')
    data["level"] = 0

    setUsers([...users, data])
    setNewUserData({})
    setIsPasswordEmpty(false)
    setIsEmailEmpty(false)
    setIsEmailExistent(false)
    setCreatingNewUser(false)
  }



  return (
    <>
      <Head>
        <title>
          Oktopus | Users
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
              <Stack spacing={1}>
                <Typography variant="h4">
                  Users
                </Typography>
                <Stack
                  alignItems="center"
                  direction="row"
                  spacing={1}
                >
                  {/* <Button
                    color="inherit"
                    startIcon={(
                      <SvgIcon fontSize="small">
                        <ArrowUpOnSquareIcon />
                      </SvgIcon>
                    )}
                  >
                    Import
                  </Button> */}
                  {/* <Button
                    color="inherit"
                    startIcon={(
                      <SvgIcon fontSize="small">
                        <ArrowDownOnSquareIcon />
                      </SvgIcon>
                    )}
                  >
                    Export
                  </Button> */}
                </Stack>
              </Stack>
              <div>
                <Button
                  startIcon={(
                    <SvgIcon fontSize="small">
                      <PlusIcon />
                    </SvgIcon>
                  )}
                  variant="contained"
                  onClick={() => {
                    setAddDeviceDialogOpen(true)
                  }}
                >
                  Add
                </Button>
              </div>
            </Stack>
            {/* <CustomersSearch /> */}
            {users && !loading ?
              <CustomersTable
                count={users.length}
                items={users}
                //onDeselectAll={customersSelection.handleDeselectAll}
                onDeselectOne={(id) => {
                  setSelected(selected.filter((item) => item !== id))
                }}
                //onPageChange={handlePageChange}
                //onRowsPerPageChange={handleRowsPerPageChange}
                //onSelectAll={customersSelection.handleSelectAll}
                onSelectOne={(id) => {
                  setSelected([...selected, id])
                  console.log("added user " + id + " to selected array")
                }}
                //page={page}
                //rowsPerPage={rowsPerPage}
                deleteUser={deleteUser}
                selected={selected}
              /> :
              <CircularProgress></CircularProgress>
            }
          </Stack>
        </Container>
      </Box>
      <Dialog
      open={addDeviceDialogOpen}
      onClose={() => {
        setAddDeviceDialogOpen(false)
        setIsEmailEmpty(false)
        setIsEmailExistent(false)
        setIsPasswordEmpty(false)
        setNewUserData({})
      }}
      >
        <DialogTitle>Create User</DialogTitle>
        <DialogContent>
          <Stack
            alignItems="center"
            direction="row"
            spacing={2}
          >
            <TextField
              inputProps={{
                form: {
                  autocomplete: 'new-password',
                },
              }}
              // focused={isEmailEmpty}
              // color={isEmailEmpty ? "error" : "primary"}
              helperText={isEmailEmpty ? "Email error" : (isEmailExistent ? "Email already exists" : "")}
              autoFocus
              required
              margin="dense"
              id="email"
              name="email"
              label="Email Address"
              type="email"
              fullWidth
              onChange={
                (event) => {
                  setNewUserData({...newUserData, email: event.target.value})
                }
              }
              variant="standard">
            </TextField>
            <TextField 
              // focused={isPasswordEmpty}
              //color={isPasswordEmpty ? "error" : "primary"}
              helperText={isPasswordEmpty ? "Password cannot be empty" : ""}
              autoFocus
              required
              margin="dense"
              id="password"
              name="password"
              label="Password"
              type="password"
              autoComplete='new-password'
              fullWidth
              onChange={
                (event) => {
                  setNewUserData({...newUserData, password: event.target.value})
                }
              }
              variant="standard">
            </TextField>
          </Stack>
          <Stack
            alignItems="center"
            direction="row"
            spacing={2}
          >
            <TextField 
              inputProps={{
                form: {
                  autocomplete: 'off',
                },
              }}
              autoFocus
              margin="dense"
              id="name"
              name="name"
              label="Full Name"
              type="name"
              fullWidth
              variant="standard"
              onChange={
                (event) => {
                  setNewUserData({...newUserData, name: event.target.value})
                }
              }
              >
            </TextField>
            <TextField 
              inputProps={{
                form: {
                  autocomplete: 'off',
                },
              }}
              autoComplete="off"
              autoFocus
              margin="dense"
              id="phone"
              name="phone"
              label="Phone Number"
              type="phone"
              fullWidth
              variant="standard"
              onChange={
                (event) => {
                  setNewUserData({...newUserData, phone: event.target.value})
                }
              }
              >
            </TextField>
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={
            () => {
              setAddDeviceDialogOpen(false)
              setIsEmailEmpty(false)
              setIsEmailExistent(false)
              setIsPasswordEmpty(false)
              setNewUserData({})
            }
          }>Cancel</Button>
          <Button onClick={()=>{
            console.log("new user data: ", newUserData)
            if (newUserData.password === undefined || newUserData.password === "") {
              setIsPasswordEmpty(true)
              return
            } else{
              setIsPasswordEmpty(false)
            }
            if (newUserData.email === undefined || newUserData.email === "") {
              setIsEmailEmpty(true)
              return
            } else if(!validateEmail(newUserData.email)){
              setIsEmailEmpty(true)
              return
            }else{
              setIsEmailEmpty(false)
            }
            setIsEmailExistent(false)
            setCreatingNewUser(true)
            createUser(newUserData)
          }}>Confirm</Button>
        </DialogActions>
        {
        <Backdrop
        sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
        open={creatingNewUser}
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
