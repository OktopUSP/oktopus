import PropTypes from 'prop-types';
import ExclamationTriangle from '@heroicons/react/24/solid/ExclamationTriangleIcon';
import Signal from '@heroicons/react/24/solid/SignalIcon';
import Image from 'next/image';
import {
  Avatar,
  Box,
  Card,
  CardContent,
  LinearProgress,
  Stack,
  SvgIcon,
  Typography
} from '@mui/material';


export const OverviewTasksProgress = (props) => {
  var { value, mtp, sx, type } = props;
  var valueRaw;
  if( value !== undefined) {
    valueRaw = value.substring(1);
    console.log("rtt:", valueRaw)
  }

  const formatMilliseconds = (timeString) => {
  
  if (timeString === "") {
    return "";
  }
  // Regular expression to extract value and unit
  const regex = /^(\d+(\.\d+)?)\s*([mµ]?s|s)?$/;

  // Extract value and unit
  const matches = timeString.match(regex);
  if (!matches) {
      return "Invalid time format";
  }

  let value = parseFloat(matches[1]);
  const unit = matches[3] || "ms";

  // Convert units to milliseconds
  switch (unit) {
      case "s":
          value *= 1000;
          break;
      case "µs":
          value /= 1000;
          break;
      default:
          // For "ms", do nothing
          break;
  }

  // Round the number to two decimal places
  const roundedValue = value.toFixed(2);

  return `${roundedValue}ms`;
}

const showIcon = (mtpType) => {
  if (valueRaw === "") {
    return <SvgIcon><ExclamationTriangle></ExclamationTriangle></SvgIcon>
  }
  switch (mtpType) {
    case "mqtt":
      return <SvgIcon><Signal/></SvgIcon>
    case "stomp":
      return <Image src="/assets/mtp/boot-stomp.svg" alt="STOMP" width={24} height={24} />;
    case "websocket":
      return <Image src="/assets/mtp/websocket.svg" alt="WebSocket" width={24} height={24} />;
    default:
      return <ExclamationTriangle />;
  }
}

  return (
    <Card sx={sx}>
      <CardContent>
        <Stack
          alignItems="flex-start"
          direction="row"
          justifyContent="space-between"
          spacing={3}
        >
          <Stack spacing={1}>
            <Typography
              color="sucess.main"
              gutterBottom
              variant="overline"
            >
              {mtp}
            </Typography>
            <Typography variant="h4">
              {formatMilliseconds(valueRaw)}
            </Typography>
          </Stack>
          <Avatar
            sx={{
              backgroundColor: 'primary.darkest',
              height: 56,
              width: 56
            }}
          >
            {showIcon(type)}
          </Avatar>
        </Stack>
        <Box sx={{ mt: 3 }}>
          <LinearProgress
            value={valueRaw ? 80 : 0}
            variant="determinate"
          />
        </Box>
      </CardContent>
    </Card>
  );
};

OverviewTasksProgress.propTypes = {
  value: PropTypes.string.isRequired,
  sx: PropTypes.object
};
