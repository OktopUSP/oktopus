import { format } from 'date-fns';
import PropTypes from 'prop-types';
import ArrowRightIcon from '@heroicons/react/24/solid/ArrowRightIcon';
import ArrowTopRightOnSquareIcon from '@heroicons/react/24/solid/ArrowTopRightOnSquareIcon';
import {
  Box,
  Button,
  Card,
  CardActions,
  CardHeader,
  Divider,
  SvgIcon,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow
} from '@mui/material';
import { Scrollbar } from 'src/components/scrollbar';
import { SeverityPill } from 'src/components/severity-pill';
import { useRouter } from 'next/router';

const statusMap = {
  1: 'warning',
  2: 'success',
  0: 'error'
};

const status = (s)=>{
  if (s == 0){
    return "Offline"
  } else if (s == 1){
    return "Associating"
  }else if (s==2){
    return "Online"
  }else {
    return "Unknown"
  }
}

export const OverviewLatestOrders = (props) => {
  const { orders = [], sx } = props;

  const router = useRouter()

  return (
    <Card sx={sx}>
      <CardHeader title="Devices" />
      <Scrollbar sx={{ flexGrow: 1 }}>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell align="center">
                  Serial Number
                </TableCell>
                <TableCell>
                  Model
                </TableCell>
                <TableCell>
                  Vendor
                </TableCell>
                <TableCell>
                  Version
                </TableCell>
                <TableCell>
                  Status
                </TableCell>
                <TableCell>
                  Access
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {orders && orders.map((order) => {

                return (
                  <TableRow
                    hover
                    key={order.SN}
                  >
                    <TableCell TableCell align="center">
                      {order.SN}
                    </TableCell>
                    <TableCell>
                      {order.Model}
                    </TableCell>
                    <TableCell>
                      {order.Vendor}
                    </TableCell>
                    <TableCell>
                      {order.Version}
                    </TableCell>
                    <TableCell>
                    <SeverityPill color={statusMap[order.Status]}>
                        {status(order.Status)}
                    </SeverityPill>
                    </TableCell>
                    <TableCell>
                    <SvgIcon 
                      fontSize="small" 
                      sx={{cursor: order.Status == 2 && 'pointer'}} 
                      onClick={()=>{
                          if (order.Status == 2){
                            router.push("devices/"+order.SN+"/discovery")
                          }
                        }
                      }
                    >
                      <ArrowTopRightOnSquareIcon />
                    </SvgIcon>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </Box>
      </Scrollbar>
      {/*<Divider />
        <CardActions sx={{ justifyContent: 'flex-end' }}>
          <Button
            color="inherit"
            endIcon={(
              <SvgIcon fontSize="small">
                <ArrowRightIcon />
              </SvgIcon>
            )}
            size="small"
            variant="text"
          >
            View all
          </Button>
            </CardActions>*/}
    </Card>
  );
};

OverviewLatestOrders.prototype = {
  orders: PropTypes.array,
  sx: PropTypes.object
};
