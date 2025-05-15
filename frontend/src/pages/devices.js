import React, { useState, useEffect } from 'react';
import Head from 'next/head';

import {
  Box,
  Container,
  Unstable_Grid2 as Grid,
  Card,
  OutlinedInput,
  InputAdornment,
  SvgIcon,
  Stack,
  Pagination,
  CircularProgress,
  Button,
  Menu,
  MenuItem,
  Checkbox,
  ListItemText,
  SpeedDial,
  SpeedDialAction,
  SpeedDialIcon,
  CardHeader,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  FormControl,
  Tooltip,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  InputLabel,
  Input,
  TextField,
  Select,
  DialogContentText,
  TableContainer,
  TablePagination,
  Typography
} from '@mui/material';

import ViewColumnsIcon from '@heroicons/react/24/outline/ViewColumnsIcon';
import ArrowTopRightOnSquareIcon from '@heroicons/react/24/solid/ArrowTopRightOnSquareIcon';
import FunnelIcon from "@heroicons/react/24/outline/FunnelIcon";
import PencilIcon from '@heroicons/react/24/outline/PencilIcon';
import TagIcon from '@heroicons/react/24/outline/TagIcon';
import ShareIcon from '@heroicons/react/24/outline/ShareIcon';
import CommandLineIcon from '@heroicons/react/24/outline/CommandLineIcon';
import ChevronUpIcon from '@heroicons/react/24/outline/ChevronUpIcon';
import ChevronDownIcon from '@heroicons/react/24/outline/ChevronDownIcon';
import MagnifyingGlassIcon from '@heroicons/react/24/solid/MagnifyingGlassIcon';
import TrashIcon from '@heroicons/react/24/outline/TrashIcon';

import { Scrollbar } from 'src/components/scrollbar';
import { Layout as DashboardLayout } from 'src/layouts/dashboard/layout';
import { useAuth } from 'src/hooks/use-auth';
import { useRouter } from 'next/router';
import { useTheme } from '@emotion/react';

const Page = () => {

  const theme = useTheme();
  const router = useRouter()
  const auth = useAuth();

  const [devices, setDevices] = useState([]);
  const [total, setTotal] = useState(null);
  const [deviceFound, setDeviceFound] = useState(true)
  const [pages, setPages] = useState(0);
  const [page, setPage] = useState(0);
  const [Loading, setLoading] = useState(true);
  const [anchorEl, setAnchorEl] = useState(null);
  const [statusOrder, setStatusOrder] = useState("desc");

  const [filterOptions, setFilterOptions] = useState(null);
  const defaultFiltersList = {
    alias: "",
    model: "",
    vendor: "",
    version: "",
    status: "",
    type: "",
  }
  const [filtersList, setFiltersList] = useState(defaultFiltersList);
  const [newFiltersList, setNewFiltersList] = useState(defaultFiltersList);

  const cleanFilters = () => {
    setFiltersList(defaultFiltersList)
  }

  const rowsPerPageOptions = [20,30,40];
  const [rowsPerPage, setRowsPerPage] = useState(20);

  const [showSetDeviceAlias, setShowSetDeviceAlias] = useState(false);
  const [showSetDeviceToBeRemoved, setShowSetDeviceToBeRemoved] = useState(false);
  const [deviceAlias, setDeviceAlias] = useState(null);
  const [deviceToBeChanged, setDeviceToBeChanged] = useState(null);
  const [deviceToBeRemoved, setDeviceToBeRemoved] = useState(null);
  const [showFilter, setShowFilter] = useState(false);
  const [selected, setSelected] = useState([]);
  const [selectAll, setSelectAll] = useState(false);


  const [columns, setColumns] = useState({
    version: true,
    sn: true,
    alias: true,
    model: true,
    vendor: true,
    status: true,
    actions: true,
    label: false
  });

  const [showSpeedDial, setShowSpeedDial] = useState(false);

  const getColumns = () => {
    localStorage.getItem("columns") ? setColumns(JSON.parse(localStorage.getItem("columns"))) : setColumns({
      version: true,
      sn: true,
      alias: false,
      model: true,
      vendor: true,
      status: true,
      actions: true,
      label: false
    })
  }

  const changeColumn = (column) => {
    console.log("columns old:", columns)
    setColumns({ ...columns, [column]: !columns[column] })
    localStorage.setItem("columns", JSON.stringify({ ...columns, [column]: !columns[column] }))
  }

  function objsEqual(obj1,obj2){
    return JSON.stringify(obj1)===JSON.stringify(obj2);
 }

  useEffect(() => {
    if (selected.length > 0) {
      let speedDial = false
      selected.map((s) => {
        if (s == true) {
          speedDial = true
        } else if (s == false) {
          setSelectAll(false)
        }
      })
      setShowSpeedDial(speedDial)
      return
    }
  }, [selected])

  const statusMap = {
    1: 'warning',
    2: 'success',
    0: 'error'
  };

  const status = (s) => {
    if (s == 0) {
      return "Offline"
    } else if (s == 1) {
      return "Associating"
    } else if (s == 2) {
      return "Online"
    } else {
      return "Unknown"
    }
  }

  const getDeviceProtocol = (order) => {
    console.log("order:", order)
    if (order.Cwmp == 2) {
      return "cwmp"
    } else {
      return "usp"
    }
  }

  const open = Boolean(anchorEl);
  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };

  const actions = [
    /*{ icon: <SvgIcon><PencilIcon /></SvgIcon>, name: 'Alias', onClickEvent: () => {
      console.log("edit device alias")
      // setDeviceToBeChanged(selected.indexOf(true))
      // setDeviceAlias(orders[selected.indexOf(true)].Alias)
      // setShowSetDeviceAlias(true)
    } },*/
    //{ icon: <SvgIcon><TagIcon /></SvgIcon>, name: 'Label' },
    //{ icon: <SvgIcon><ShareIcon /></SvgIcon>, name: 'Share' },
    { icon: <SvgIcon><TrashIcon /></SvgIcon>, name: 'Remove' },
    { icon: <SvgIcon><CommandLineIcon /></SvgIcon>, name: 'Action' },
  ];

  useEffect(() => {
    getColumns()
    setLoading(true)

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    var requestOptions = {
      method: 'GET',
      headers: myHeaders,
      redirect: 'follow'
    }

    let status;

    fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device?statusOrder=${statusOrder}&page_number=${page}&page_size=${rowsPerPage}&vendor=${filtersList["vendor"]}&version=${filtersList["version"]}&alias=${filtersList["alias"]}&type=${filtersList["type"]}&status=${filtersList["status"]}&model=${filtersList["model"]}`, requestOptions)
      .then(response => {
        if (response.status === 401){
          router.push("/auth/login")
        }
        status = response.status
        return response.json()
      })
      .then(json => {
        if (status == 404) {
          console.log("device not found")
          setLoading(false)
          setDeviceFound(false)
          return
        }
        console.log("Status:", status)
        setPages(json.pages + 1)
        setPage(json.page + 1)
        setTotal(json.total)
        setDevices(json.devices)
        setSelected(new Array(json.devices.length).fill(false))
        setLoading(false)
        return setDeviceFound(true)
      })
      .catch(error => {
        return console.error('Error:', error)
      });

    fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device/filterOptions`, requestOptions)
      .then(response => {
        if (response.status === 401)
          router.push("/auth/login")
        return response.json()
      })
      .then(json => {
        return setFilterOptions(json)
      })
      .catch(error => {
        return console.error('Error:', error)
      });

  }, [auth.user]);

  const removeDevice = async (sn) => {
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    var requestOptions = {
      method: 'DELETE',
      headers: myHeaders,
      redirect: 'follow'
    };

    let result = await fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device?id=${sn}`, requestOptions)
    console.log("result:", result)
    if (result.status === 401) {
      router.push("/auth/login")
    } else if (result.status != 200) {
      console.log("Status:", result.status)
      let content = await result.json()
      console.log("Message:", content)
      setShowSetDeviceToBeRemoved(false)
      setDeviceToBeRemoved(null)
    } else {
      let content = await result.json()
      console.log("remove device result:", content)
      setShowSetDeviceToBeRemoved(false)
      setDeviceToBeRemoved(null)
      devices.splice(deviceToBeRemoved, 1)
      setDevices([...devices])
    }

  }

  const setNewDeviceAlias = async (alias, sn) => {
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", localStorage.getItem("token"));

    var requestOptions = {
      method: 'PUT',
      headers: myHeaders,
      body: alias,
      redirect: 'follow'
    };

    let result = await fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device/alias?id=${sn}`, requestOptions)
    console.log("result:", result)
    if (result.status === 401) {
      router.push("/auth/login")
    } else if (result.status != 200) {
      console.log("Status:", result.status)
      let content = await result.json()
      console.log("Message:", content)
      setShowSetDeviceAlias(false)
      setDeviceAlias(null)
      setDeviceToBeChanged(null)
    } else {
      let content = await result.json()
      console.log("set alias result:", content)
      setShowSetDeviceAlias(false)
      setDeviceAlias(null)
      devices[deviceToBeChanged].Alias = alias
      setDeviceToBeChanged(null)
      setDevices([...devices])
    }
    // .then(response => {
    //   if (response.status === 401) {
    //     router.push("/auth/login")
    //   }
    //   return response.json()
    // })
    // .then(result => {
    //   console.log("alias result:", result)
    //   setShowSetDeviceAlias(false)
    //   setDeviceAlias(null)
    // })
    // .catch(error => {
    //   console.log('error:', error)
    //   setShowSetDeviceAlias(false)
    //   setDeviceAlias(null)
    // })
  }

  const fetchDevicePerPage = async (p, s, localFilterList, page_size) => {
    //setLoading(true)

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var requestOptions = {
      method: 'GET',
      headers: myHeaders,
      redirect: 'follow'
    }

    if (localFilterList == undefined) {
      localFilterList = filtersList
    }

    if (page_size == undefined) {
      page_size = rowsPerPage
    }

    p = p - 1
    p = p.toString()

    fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device?page_number=${p}&page_size=${page_size}&statusOrder=${s}&vendor=${localFilterList["vendor"]}&version=${localFilterList["version"]}&alias=${localFilterList["alias"]}&type=${localFilterList["type"]}&status=${localFilterList["status"]}&model=${localFilterList["model"]}`, requestOptions)
      .then(response => {
        if (response.status === 401)
          router.push("/auth/login")
        return response.json()
      })
      .then(json => {
        setTotal(json.total)
        setDevices(json.devices)
        setPages(json.pages + 1)
        setPage(json.page + 1)
        //setLoading(false)
        return
      })
      .catch(error => {
        return console.error('Error:', error)
      });
  }

  const fetchDevicePerId = async (id) => {
    setLoading(true)
    setDeviceFound(true)
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", auth.user.token);

    var requestOptions = {
      method: 'GET',
      headers: myHeaders,
      redirect: 'follow'
    }

    if (id == "") {
      return fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device?vendor=${filtersList["vendor"]}&version=${filtersList["version"]}&alias=${filtersList["alias"]}&type=${filtersList["type"]}&status=${filtersList["status"]}&model=${filtersList["model"]}`, requestOptions)
        .then(response => {
          if (response.status === 401)
            router.push("/auth/login")
          return response.json()
        })
        .then(json => {
          setPages(json.pages + 1)
          setPage(json.page + 1)
          setTotal(json.total)
          setDevices(json.devices)
          setLoading(false)
          return setDeviceFound(true)
        })
        .catch(error => {
          return console.error('Error:', error)
        });
    }

    let response = await fetch(`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/api/device?id=${id}`, requestOptions)
    if (response.status === 401)
      router.push("/auth/login")
    let json = await response.json()
    if (json.SN != undefined) {
      setDevices([json])
      setTotal(1)
      setDeviceFound(true)
      setLoading(false)
      setPages(1)
      setPage(1)
    } else {
      setDeviceFound(false)
      setDevices(null)
      setTotal(null)
      setPages(null)
      setPage(null)
    }

  }

  return (
    <>
      <Head>
        <title>
          Oktopus | Controller
        </title>
      </Head>

      <Box
        component="main"
        sx={{
          flexGrow: 1,
        }}
      >
        <Container maxWidth="xl" >
          <Stack spacing={1}>
            <Grid
              container
              spacing={1}
              py={1}
            >
              <Grid xs={4} item>
                  <OutlinedInput
                    xs={4}
                    defaultValue=""
                    fullWidth
                    placeholder="Search Device"
                    onKeyDownCapture={(e) => {
                      if (e.key === 'Enter') {
                        console.log("Fetch devices per id: ", e.target.value)
                        fetchDevicePerId(e.target.value)
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
              </Grid>
              <Grid xs={4} item>
                <Button
                    sx={{ backgroundColor: "rgba(48, 109, 111, 0.04)" }}
                    onClick={() => { setShowFilter(true) }}
                  >
                    <SvgIcon>
                      <FunnelIcon />
                    </SvgIcon>
                </Button>
                {Object.keys(filtersList).map((key) => (
                  (filtersList[key] && 
                  <Chip label={`${key} : ${filtersList[key]}`} sx={{ml:1, mt:1}} onDelete={()=>{
                    setFiltersList({ ...filtersList, [key]: "" })
                    setNewFiltersList({ ...newFiltersList, [key]: ""})
                    fetchDevicePerPage(1, statusOrder, { ...filtersList, [key]: "" })
                  }}/>)
                ))}
              </Grid>
            </Grid>
            <div style={{ display: "flex", justifyContent: "flex-end" }}>
              <Button
                id="basic-button"
                aria-controls={open ? 'basic-menu' : undefined}
                aria-haspopup="true"
                aria-expanded={open ? 'true' : undefined}
                onClick={handleClick}
              >
                <SvgIcon>
                  <ViewColumnsIcon />
                </SvgIcon>
              </Button>
              <Menu
                id="basic-menu"
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}
                MenuListProps={{
                  'aria-labelledby': 'basic-button',
                }}
                anchorOrigin={{
                  vertical: "bottom",
                  horizontal: "center"
                }}
                transformOrigin={{
                  vertical: "top",
                  horizontal: "center"
                }}
              >
                <MenuItem dense onClick={() => changeColumn("sn")}><Checkbox checked={columns["sn"]} /*onChange={()=>changeColumn("sn")}*/ /><ListItemText primary="Serial Number" /></MenuItem>
                <MenuItem dense onClick={() => changeColumn("alias")}><Checkbox checked={columns["alias"]} /*onChange={()=>changeColumn("alias")}*/ /><ListItemText primary="Alias" /></MenuItem>
                <MenuItem dense onClick={() => changeColumn("model")}><Checkbox checked={columns["model"]} /*onChange={() => changeColumn("model")}*/ /><ListItemText primary="Model" /></MenuItem>
                <MenuItem dense onClick={() => changeColumn("vendor")}><Checkbox checked={columns["vendor"]} /*onChange={() => changeColumn("vendor")}*/ /><ListItemText primary="Vendor" /></MenuItem>
                <MenuItem dense onClick={() => changeColumn("version")}><Checkbox checked={columns["version"]} /*onChange={() => changeColumn("version")}*/ /><ListItemText primary="Version" /></MenuItem>
                <MenuItem dense onClick={() => changeColumn("status")}><Checkbox checked={columns["status"]} /*onChange={() => changeColumn("status")}*/ /><ListItemText primary="Status" /></MenuItem>
                <MenuItem dense onClick={() => changeColumn("actions")}><Checkbox checked={columns["actions"]} /*onChange={() => changeColumn("actions")}*/ /><ListItemText primary="Actions" /></MenuItem>
                {/* <MenuItem dense onClick={() => changeColumn("label")}><Checkbox checked={columns["label"]} /><ListItemText primary="Labels" /></MenuItem> */}
              </Menu>
            </div>
                <div>
                  <Card sx={{ height: "100%" }}>
                    <CardHeader title="Devices" />
                    <Scrollbar sx={{ flexGrow: 1 }}>
                      <Box sx={{ minWidth: 800 }}>
                        <TableContainer sx={{ maxHeight: 600 }}>
                          <Table stickyHeader>
                            <TableHead>
                              <TableRow>
                                {/* <TableCell align="center">
                                  <Checkbox
                                    size="small"
                                    style={{ margin: 0, padding: 0, color: theme.palette.primary.lightest }}
                                    onChange={(e) => {
                                      setSelected(new Array(devices.length).fill(e.target.checked))
                                      //console.log("selected:", selected)
                                      setSelectAll(e.target.checked)
                                    }}
                                    checked={selectAll}
                                  />
                                </TableCell> */}
                                {columns["sn"] && <TableCell align="center">
                                  Serial Number
                                </TableCell>}
                                {columns["alias"] && <TableCell>
                                  Alias
                                </TableCell>}
                                {
                                  columns["label"] && <TableCell >
                                    Labels
                                  </TableCell>
                                }
                                {columns["model"] && <TableCell>
                                  Model
                                </TableCell>}
                                {columns["vendor"] && <TableCell>
                                  Vendor
                                </TableCell>}
                                {columns["version"] && <TableCell>
                                  Version
                                </TableCell>}
                                {columns["status"] &&
                                  <TableCell>
                                    {/*//TODO: create function to fetch devices by status order*/}
                                    <Tooltip title="Change status display order" placement="top">
                                      <span style={{ cursor: "pointer" }} onClick={() => {
                                        if (statusOrder == "asc") {
                                          setStatusOrder("desc")
                                          fetchDevicePerPage(page, "desc")
                                        } else {
                                          setStatusOrder("asc")
                                          fetchDevicePerPage(page, "asc")
                                        }
                                      }}>Status ↑↓</span>
                                    </Tooltip>
                                    {/* <Box >
                                    <Tooltip title="Change status display order" placement="top">
                                      <SvgIcon fontSize='small' style={{
                                        marginLeft: "10px",
                                        cursor: "pointer"
                                      }}>
                                        <ArrowsUpDownIcon />
                                      </SvgIcon>
                                    </Tooltip>
                                  </Box> */}
                                  </TableCell>}
                                {columns["actions"] && <TableCell align="center">
                                  Actions
                                </TableCell>}
                              </TableRow>
                            </TableHead>
                            {!Loading ? <TableBody>
                              {devices && devices.map((order, index) => {
                                return (
                                  <TableRow
                                    hover
                                    key={order.SN}
                                  >
                                    {/* <TableCell align="center">
                                      <FormControl>
                                        <Checkbox
                                          size="small"
                                          onChange={(e) => {
                                            let newData = [...selected]
                                            newData.splice(index, 1, e.target.checked);
                                            console.log("newData:", newData)
                                            setSelected(newData)
                                          }}
                                          checked={selected[index]}
                                        />
                                      </FormControl>
                                    </TableCell> */}
                                    {columns["sn"] && <TableCell align="center">
                                      {order.SN}
                                    </TableCell>}
                                    {columns["alias"] && <TableCell>
                                      {order.Alias}
                                    </TableCell>}
                                    {
                                      columns["label"] && <TableCell>
                                        <Chip label="Teste1" />
                                        <Chip label="Teste2" />
                                        <Chip label="Teste3" />
                                      </TableCell>
                                    }
                                    {columns["model"] && <TableCell>
                                      {order.Model || order.ProductClass}
                                    </TableCell>}
                                    {columns["vendor"] && <TableCell>
                                      {order.Vendor}
                                    </TableCell>}
                                    {columns["version"] && <TableCell>
                                      {order.Version}
                                    </TableCell>}
                                    {columns["status"] && <TableCell>
                                      {/* <SeverityPill color={statusMap[order.Status]}>
                                        {status(order.Status)}
                                      </SeverityPill> */}
                                      <Chip label={status(order.Status)} color={statusMap[order.Status]} />
                                    </TableCell>}
                                    {columns["actions"] && <TableCell align="center">
                                      {order.Status == 2 &&
                                        <Tooltip title="Access the device">
                                          <Button
                                            onClick={() => {
                                              router.push("devices/" + getDeviceProtocol(order) + "/" + order.SN)
                                            }}
                                          >
                                            <SvgIcon
                                              fontSize="small"
                                              sx={{ cursor: 'pointer' }}
                                            >
                                              <ArrowTopRightOnSquareIcon />
                                            </SvgIcon>
                                          </Button>
                                        </Tooltip>}
                                      <Tooltip title="Edit the device alias">
                                    <Button
                                      onClick={()=>{
                                        setDeviceToBeChanged(index)
                                        setDeviceAlias(order.Alias)
                                        setShowSetDeviceAlias(true)
                                      }}
                                    >
                                      <SvgIcon 
                                        fontSize="small" 
                                        sx={{cursor: 'pointer'}} 
                                      >
                                        <PencilIcon />
                                      </SvgIcon>
                                    </Button>
                                  </Tooltip>
                                  <Tooltip title="Delete device">
                                    <Button
                                      onClick={()=>{
                                        setDeviceToBeRemoved(index)
                                        setShowSetDeviceToBeRemoved(true)
                                      }}
                                    >
                                      <SvgIcon 
                                        fontSize="small" 
                                        sx={{cursor: 'pointer'}} 
                                      >
                                        <TrashIcon />
                                      </SvgIcon>
                                    </Button>
                                  </Tooltip>
                                  {/* <Tooltip title="Edit device labels">
                                    <Button
                                      onClick={()=>{
                                        setDeviceToBeChanged(index)
                                      }}
                                    >
                                      <SvgIcon 
                                        fontSize="small" 
                                        sx={{cursor: 'pointer'}} 
                                      >
                                        <TagIcon />
                                      </SvgIcon>
                                    </Button>
                                  </Tooltip> */}
                                    </TableCell>}
                                  </TableRow>
                                );
                              })}
                            </TableBody>: 
                              <TableBody>
                                <TableRow>
                                  <TableCell colSpan={7} align="center">
                                    {
                                      deviceFound ? <CircularProgress/> : "No device found"
                                    }
                                  </TableCell>
                                </TableRow>
                              </TableBody>
                            }
                          </Table>
                        </TableContainer>
                        {(pages > 0) && total && <TablePagination 
                          rowsPerPageOptions={rowsPerPageOptions}
                          component="div"
                          count={total}
                          rowsPerPage={rowsPerPage}
                          page={page-1}
                          onPageChange={(e, newPage)=>{
                            setPage(newPage+1)
                            fetchDevicePerPage(newPage+1, statusOrder, filtersList, rowsPerPage)
                          }}
                          onRowsPerPageChange={(e)=>{
                            setRowsPerPage(e.target.value)
                            setPage(1)
                            fetchDevicePerPage(1, statusOrder, filtersList, e.target.value)
                          }}
                        />}
                      </Box>
                    </Scrollbar>
                  </Card>
                  {showSetDeviceAlias &&
                    <Dialog open={showSetDeviceAlias}>
                      <DialogContent>
                        <InputLabel>Device Alias</InputLabel>
                        <Input value={deviceAlias} onChange={(e) => { setDeviceAlias(e.target.value) }}
                          onKeyUp={e => {
                            if (e.key === 'Enter') {
                              setNewDeviceAlias(deviceAlias, devices[deviceToBeChanged].SN)
                            }
                          }}>
                        </Input>
                      </DialogContent>
                      <DialogActions>
                        <Button onClick={() => {
                          setShowSetDeviceAlias(false)
                          setDeviceAlias(null)
                          setDeviceToBeChanged(null)
                        }}>Cancel</Button>
                        <Button onClick={() => {
                          setNewDeviceAlias(deviceAlias, devices[deviceToBeChanged].SN)
                        }}>Save</Button>
                      </DialogActions>
                    </Dialog>}
                    {showSetDeviceToBeRemoved &&
                    <Dialog open={showSetDeviceToBeRemoved}>
                      <DialogContent>
                        <DialogContentText>Are you sure you want to remove <b>{devices[deviceToBeRemoved].SN}</b> device?</DialogContentText>
                      </DialogContent>
                      <DialogActions>
                        <Button onClick={() => {
                          setShowSetDeviceToBeRemoved(false)
                          setDeviceToBeRemoved(null)
                        }}>Cancel</Button>
                        <Button 
                          endIcon={<SvgIcon><TrashIcon /></SvgIcon>}
                          onClick={() => {
                          removeDevice(devices[deviceToBeRemoved].SN)
                        }}>Apply</Button>
                      </DialogActions>
                    </Dialog>}
                </div>
            {/* / :
            //   <Box
            //     sx={{
            //       display: 'flex',
            //       justifyContent: 'center'
            //     }}
            //   >
            //     <p>Device Not Found</p>
            //   </Box>
            // }
            // <Box
            //   sx={{
            //     display: 'flex',
            //     justifyContent: 'center'
            //   }}
            // > */}
              {/* {pages ? <Pagination
                count={pages}
                size="small"
                page={page}
                onChange={handleChangePage}
              /> : null} */}
              {/* //TODO: show loading */}
             {/* </Box> */}
          </Stack>
        </Container>
      </Box>
      <SpeedDial
        ariaLabel="SpeedDial basic example"
        hidden={!showSpeedDial}
        //FabProps={{size: 'small'}}
        sx={{ position: 'fixed', bottom: 16, right: 16, }}
        icon={<SpeedDialIcon icon={<SvgIcon ><ChevronDownIcon /></SvgIcon>}
          openIcon={<SvgIcon><ChevronUpIcon /></SvgIcon>} />}
      >
        {actions.map((action) => (
          <SpeedDialAction
            key={action.name}
            icon={action.icon}
            tooltipTitle={action.name}
            onClick={action.onClickEvent}
          />
        ))}
      </SpeedDial>
      <Dialog open={showFilter} fullWidth>
        <DialogTitle>
          <SvgIcon style={{ marginRight: "10px", marginBottom: "-5px" }}>
            <FunnelIcon />
          </SvgIcon>
          Filter
        </DialogTitle>
        {filterOptions && <DialogContent>
          <Stack spacing={2} marginTop={1} minWidth={400}>
            <Stack
              spacing={2}
              direction={'row'}
            >
              <TextField label="Alias" variant="filled" sx={{minWidth:"48%"}}
                value={newFiltersList["alias"]}
                onChange={(e) => setNewFiltersList({ ...newFiltersList, "alias": e.target.value })}
              />
              <FormControl variant="filled" sx={{ minWidth: "48%" }}>
                <InputLabel>Type</InputLabel>
                <Select
                  value={newFiltersList["type"]}
                  onChange={(e) => setNewFiltersList({ ...newFiltersList, "type": e.target.value })}
                  fullWidth
                >
                  {
                    filterOptions.productClasses.map((v) => {
                      return <MenuItem value={v}>{v}</MenuItem>
                    })
                  }
                </Select>
              </FormControl>
              {/* <TextField label="Serial Number" variant="filled" /> */}
            </Stack>
            <Stack
              spacing={2}
              direction={'row'}
            >
              <FormControl variant="filled" sx={{ minWidth: "48%" }}>
                <InputLabel>Vendor</InputLabel>
                <Select
                  value={newFiltersList["vendor"]}
                  onChange={(e) => setNewFiltersList({ ...newFiltersList, "vendor": e.target.value })}
                  fullWidth
                >
                  {
                    filterOptions.vendors.map((v) => {
                      return <MenuItem value={v}>{v}</MenuItem>
                    })
                  }
                </Select>
              </FormControl>
              <FormControl variant="filled" sx={{ minWidth: "48%" }}>
                <InputLabel>Version</InputLabel>
                <Select
                  value={newFiltersList["version"]}
                  onChange={(e) => setNewFiltersList({ ...newFiltersList, "version": e.target.value })}
                  fullWidth
                >
                  {
                    filterOptions.versions.map((v) => {
                      return <MenuItem value={v}>{v}</MenuItem>
                    })
                  }
                </Select>
              </FormControl>
            </Stack>
            <Stack
              spacing={2}
              direction={'row'}
            >
              <FormControl variant="filled" sx={{ minWidth: "48%" }}>
                <InputLabel>Status</InputLabel>
                <Select
                  value={
                    newFiltersList["status"]
                  //   ()=>{
                  //   if (newFiltersList["status"] == 2) {
                  //     return "online"
                  //   }else if (newFiltersList["status"] == 0) {
                  //     return "offline"
                  //   }else {
                  //     return ""
                  //   }
                  // }
                  }
                  onChange={(e) => {
                    // let value = 0
                    // if (e.target.value == "online") {
                    //   value = 2
                    // }else if (e.target.value == "offline") {
                    //   value = 0
                    // }else {
                    //   return
                    // }
                    setNewFiltersList({ ...newFiltersList, "status": e.target.value })
                  }}
                  fullWidth
                >
                  <MenuItem value={"2"}>Online</MenuItem>
                  <MenuItem value={"0"}>Offline</MenuItem>
                </Select>
              </FormControl>
              {/* <FormControl variant="filled" sx={{ minWidth: "48%" }}>
                  <InputLabel>Label</InputLabel>
                  <Select
                    value={labelFilter}
                    onChange={(e)=> setLabelFilter(e.target.value)}
                    fullWidth
                  >
                    {
                      filterOptions.labels.map((v) => {
                        return <MenuItem value={v}>{v}</MenuItem>
                      })
                    }
                  </Select>
              </FormControl> */}
              <FormControl variant="filled" sx={{ minWidth: "48%" }}>
                <InputLabel>Model</InputLabel>
                <Select
                  value={newFiltersList["model"]}
                  onChange={(e) => setNewFiltersList({ ...newFiltersList, "model": e.target.value })}
                  fullWidth
                >
                  {
                    filterOptions.models.map((v) => {
                      return <MenuItem value={v}>{v}</MenuItem>
                    })
                  }
                </Select>
              </FormControl>
            </Stack>
            {/* <Stack
              spacing={2}
              direction={'row'}
            >
              <FormControl variant="filled" sx={{ minWidth: "48%" }}>
                <InputLabel>Type</InputLabel>
                <Select
                  value={typeFilter}
                  onChange={(e) => setTypeFilter(e.target.value)}
                  fullWidth
                >
                  {
                    filterOptions.productClasses.map((v) => {
                      return <MenuItem value={v}>{v}</MenuItem>
                    })
                  }
                </Select>
              </FormControl>
            </Stack> */}
          </Stack>
        </DialogContent>}
        <DialogActions>
          <Button onClick={() => {
            
            setNewFiltersList(filtersList)
            //cleanFilters()
            setShowFilter(false)
            //if (!objsEqual(filtersList,defaultFiltersList)) {
/*               fetchDevicePerPage(1, statusOrder, defaultFiltersList)
            }
            if (!objsEqual(filtersList, newFiltersList)){

            } */

            //setFiltersList(defaultFiltersList)
          }}>Cancel</Button>
          <Button onClick={() => { 
            setFiltersList(newFiltersList)
            setShowFilter(false)
            console.log("filters list:", filtersList)
            fetchDevicePerPage(1, statusOrder, newFiltersList)
          }}>Apply</Button>
        </DialogActions>
      </Dialog>
    </>
  )
}
Page.getLayout = (page) => (
  <DashboardLayout>
    {page}
  </DashboardLayout>
);

export default Page;
