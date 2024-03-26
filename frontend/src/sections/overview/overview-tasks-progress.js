import PropTypes from 'prop-types';
import ListBulletIcon from '@heroicons/react/24/solid/ListBulletIcon';
import Signal from '@heroicons/react/24/solid/SignalIcon';
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
  var { value, sx } = props;
  var valueRaw;
  if( value !== undefined) {
    valueRaw = value.substring(1);
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
              Conexão MQTT
            </Typography>
            <Typography variant="h4">
              {valueRaw}
            </Typography>
          </Stack>
          <Avatar
            sx={{
              backgroundColor: '#f28950',
              height: 56,
              width: 56
            }}
          >
            <SvgIcon>
              <Signal />
            </SvgIcon>
          </Avatar>
        </Stack>
        <Box sx={{ mt: 3 }}>
          <LinearProgress
            value={80}
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
