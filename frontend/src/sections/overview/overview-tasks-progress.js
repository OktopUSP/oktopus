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
  const { value, sx } = props;

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
              {value}ms
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
            value={value}
            variant="determinate"
          />
        </Box>
      </CardContent>
    </Card>
  );
};

OverviewTasksProgress.propTypes = {
  value: PropTypes.number.isRequired,
  sx: PropTypes.object
};
