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
  Button
} from '@mui/material';
import { Scrollbar } from 'src/components/scrollbar';
import { getInitials } from 'src/utils/get-initials';
import TrashIcon from '@heroicons/react/24/outline/TrashIcon';
import { useState } from 'react';

export const CustomersTable = (props) => {
  const {
    count = 0,
    items = [],
    onDeselectAll,
    onDeselectOne,
    onPageChange = () => {},
    onRowsPerPageChange,
    onSelectAll,
    onSelectOne,
    deleteUser,
    page = 0,
    rowsPerPage = 0,
    selected = []
  } = props;

  // const selectedSome = (selected.length > 0) && (selected.length < items.length);
  // const selectedAll = (items.length > 0) && (selected.length === items.length);

  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [userToDelete, setUserToDelete] = useState("")
  
  return (
    <Card>
      <Scrollbar>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                {/* <TableCell padding="checkbox"> */}
                  {/* <Checkbox
                    checked={selectedAll}
                    indeterminate={selectedSome}
                    onChange={(event) => {
                      if (event.target.checked) {
                        onSelectAll?.();
                      } else {
                        onDeselectAll?.();
                      }
                    }}
                  /> */}
                {/* </TableCell> */}
                <TableCell sx={{marginLeft:"30px"}}>
                  Name
                </TableCell>
                <TableCell>
                  Email
                </TableCell>
                {/* <TableCell>
                  Location
                </TableCell> */}
                <TableCell>
                  Phone
                </TableCell>
                <TableCell>
                  Created At
                </TableCell>
                <TableCell>
                  Level
                </TableCell>
                <TableCell>
                  Actions
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {items.map((customer) => {
                const isSelected = selected.includes(customer._id);
                return (
                  <TableRow
                    hover
                    key={customer._id}
                    selected={isSelected}
                  >
                    {/* <TableCell padding="checkbox"> */}
                      {/* <Checkbox
                        checked={isSelected}
                        onChange={(event) => {
                          if (event.target.checked) {
                            console.log(customer._id+" is selected");
                            onSelectOne(customer._id);
                          } else {
                            onDeselectOne(customer._id);
                          }
                        }}
                      /> */}
                    {/* </TableCell> */}
                    <TableCell align="center" sx={{margin: 'auto', textAlign: 'center'}}>
                      <Stack
                        alignItems="center"
                        direction="row"
                        spacing={2}
                      >
                        <Avatar src={customer.avatar ? customer.avatar : "/assets/avatars/default-avatar.png"}>
                          {getInitials(customer.name)}
                        </Avatar>
                        <Typography variant="subtitle2">
                          {customer.name}
                        </Typography>
                      </Stack>
                    </TableCell>
                    <TableCell>
                      {customer.email}
                    </TableCell>
                    {/* <TableCell>
                      {customer.address}
                    </TableCell> */}
                    <TableCell>
                      {customer.phone}
                    </TableCell>
                    <TableCell>
                      {customer.createdAt}
                    </TableCell>
                    <TableCell>
                    {customer.level == 1 ? "Admin" : "User"}
                    </TableCell>
                    <TableCell>
                      { customer.level == 0 ? <Button
                        onClick={() => {
                          console.log("delete user: ", customer._id)
                          setUserToDelete(customer.email);
                          setShowDeleteDialog(true);
                        }}
                      ><SvgIcon
                        color="action"
                        fontSize="small"
                        sx={{ cursor: 'pointer'}}
                      >
                        <TrashIcon
                        ></TrashIcon>
                      </SvgIcon></Button>: <span></span>}
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </Box>
      </Scrollbar>
      {/* <TablePagination
        component="div"
        count={count}
        //onPageChange={onPageChange}
        //onRowsPerPageChange={onRowsPerPageChange}
        //page={page}
        //rowsPerPage={rowsPerPage}
        //rowsPerPageOptions={[5, 10, 25]}
      /> */}
      <Dialog
      open={showDeleteDialog}
      onClose={() => setShowDeleteDialog(false)}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">{"Delete User"}</DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          Are you sure you want to delete this user?
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => {
          setShowDeleteDialog(false)
          setUserToDelete("")
        }} color="primary">
          Cancel
        </Button>
        <Button onClick={() => {
          deleteUser(userToDelete);
          setShowDeleteDialog(false);
          setUserToDelete("")
        }} color="primary" autoFocus>
          Delete
        </Button>
      </DialogActions>
    </Dialog>
    </Card>
  );
};

CustomersTable.propTypes = {
  count: PropTypes.number,
  items: PropTypes.array,
  onDeselectAll: PropTypes.func,
  onDeselectOne: PropTypes.func,
  onPageChange: PropTypes.func,
  //onRowsPerPageChange: PropTypes.func,
  onSelectAll: PropTypes.func,
  onSelectOne: PropTypes.func,
  deleteUser: PropTypes.func,
  //page: PropTypes.number,
  //rowsPerPage: PropTypes.number,
  selected: PropTypes.array
};
