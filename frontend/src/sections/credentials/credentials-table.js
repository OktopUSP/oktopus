import PropTypes from 'prop-types';
import {
  Avatar,
  Box,
  Card,
  Checkbox,
  Icon,
  Stack,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableHead,
  //TablePagination,
  TableRow,
  Typography,
  SvgIcon,
  Dialog,
  DialogActions,
  DialogTitle,
  DialogContent,
  DialogContentText,
  Button,
  TablePagination,
  TextField,
  InputAdornment,
  IconButton,
  Input
} from '@mui/material';
import EyeIcon from '@heroicons/react/24/outline/EyeIcon';
import EyeSlashIcon from '@heroicons/react/24/outline/EyeSlashIcon';
import { Scrollbar } from 'src/components/scrollbar';
import PencilIcon from '@heroicons/react/24/outline/PencilIcon';
import TrashIcon from '@heroicons/react/24/outline/TrashIcon';
import { useEffect, useState } from 'react';

export const CredentialsTable = (props) => {
  const {
    count = 0,
    items = {},
    onDeselectAll,
    onDeselectOne,
    onPageChange = () => {},
    onRowsPerPageChange,
    onSelectAll,
    onSelectOne,
    deleteCredential,
    page = 0,
    rowsPerPage = 0,
    // selected = []
  } = props;

  const [showPassword, setShowPassword] = useState({})
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [credentialToDelete, setCredentialToDelete] = useState("")

  useEffect(()=>{
    Object.keys(items).map((key) => {
      let newData = {};
      newData[key] = false
      setShowPassword(prevState => ({
        ...prevState,
        ...newData
      }))
    })
    // console.log("showPassword: "+ showPassword)
  },[])
  
  return (
    <Card>
      <Scrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell align='center'>
                  Username
                </TableCell>
                <TableCell align='center'>
                  Password
                </TableCell>
                <TableCell align='center'>
                  Actions
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {Object.keys(items).map((key) => {
                let value = items[key];

                return (
                  <TableRow
                    hover
                    key={key}
                  >
                    <TableCell align='center'>
                      {key}
                    </TableCell>
                    <TableCell align='center'>
                    <Input
                        id="standard-adornment-password"
                        type={showPassword[key] ? 'text' : 'password'}
                        endAdornment={
                          <InputAdornment position="end">
                            <IconButton
                              aria-label="toggle password visibility"
                              onClick={()=>{
                                let newData = {};
                                newData[key] = !showPassword[key]
                                setShowPassword(previous => ({...previous, ...newData}))
                              }}
                              //onMouseDown={handleMouseDownPassword}
                            >
                              <SvgIcon>
                                {showPassword[key] ? <EyeSlashIcon /> : <EyeIcon />}
                              </SvgIcon>
                            </IconButton>
                          </InputAdornment>
                        }
                        value={value}
                      />
                    </TableCell>
                    <TableCell align='center'>
                    <Button
                        onClick={() => {
                          console.log("delete user: ", key)
                          setCredentialToDelete(key);
                          setShowDeleteDialog(true);
                        }}
                      ><SvgIcon
                        color="action"
                        fontSize="small"
                        sx={{ cursor: 'pointer'}}
                      >
                        <TrashIcon
                        ></TrashIcon>
                      </SvgIcon></Button>
                    </TableCell>
                  </TableRow>
                  );
                })
              }
            </TableBody>
          </Table>
        </Box>
      </Scrollbar>
      <Dialog
      open={showDeleteDialog}
      onClose={() => setShowDeleteDialog(false)}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">{"Delete User"}</DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          Are you sure you want to delete this credential?
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => {
          setShowDeleteDialog(false)
          setCredentialToDelete("")
        }} color="primary">
          Cancel
        </Button>
        <Button onClick={() => {
          deleteCredential(credentialToDelete);
          setShowDeleteDialog(false);
          setCredentialToDelete("")
        }} color="primary" autoFocus>
          Delete
        </Button>
      </DialogActions>
    </Dialog>
    </Card>
  );
};

CredentialsTable.propTypes = {
  count: PropTypes.number,
  items: PropTypes.object,
  //onDeselectAll: PropTypes.func,
  //onDeselectOne: PropTypes.func,
  onPageChange: PropTypes.func,
  //onRowsPerPageChange: PropTypes.func,
  //onSelectAll: PropTypes.func,
  //onSelectOne: PropTypes.func,
  deleteCredential: PropTypes.func,
  //page: PropTypes.number,
  //rowsPerPage: PropTypes.number,
  //selected: PropTypes.array
};
