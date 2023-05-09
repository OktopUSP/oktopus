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
  Associating: 'warning',
  Online: 'success',
  Offline: 'error'
};

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
                <TableCell>
                  SN
                </TableCell>
                <TableCell>
                  MODEL
                </TableCell>
                <TableCell sortDirection="desc">
                  CUSTOMER
                </TableCell>
                <TableCell>
                  VENDOR
                </TableCell>
                <TableCell>
                  VERSION
                </TableCell>
                <TableCell>
                  STATUS
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {orders.map((order) => {

                return (
                  <TableRow
                    hover
                    key={order.SN}
                  >
                    <TableCell>
                      {order.SN}
                    </TableCell>
                    <TableCell>
                      {order.Model}
                    </TableCell>
                    <TableCell>
                      {order.Customer}
                    </TableCell>
                    <TableCell>
                      {order.Vendor}
                    </TableCell>
                    <TableCell>
                      {order.Version}
                    </TableCell>
                    <TableCell>
                      {order.Status}
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
