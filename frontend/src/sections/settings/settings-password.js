import { useCallback, useState } from 'react';
import {
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Divider,
  Stack,
  TextField
} from '@mui/material';
import { useBackendContext } from 'src/contexts/backend-context';
import { useAlertContext } from 'src/contexts/error-context';

export const SettingsPassword = () => {

  let {httpRequest} = useBackendContext();
  let {setAlert} = useAlertContext();

  const [values, setValues] = useState({
    password: '',
    confirm: ''
  });

  const handleChange = useCallback(
    (event) => {
      setValues((prevState) => ({
        ...prevState,
        [event.target.name]: event.target.value
      }));
    },
    []
  );

  return (
    <form>
      <Card>
        <CardHeader
          subheader="Update password"
          title="Password"
        />
        <Divider />
        <CardContent>
          <Stack
            spacing={3}
            sx={{ maxWidth: 400 }}
          >
            <TextField
              fullWidth
              label="Password"
              name="password"
              onChange={handleChange}
              type="password"
              value={values.password}
            />
            <TextField
              fullWidth
              label="Password (Confirm)"
              name="confirm"
              onChange={handleChange}
              type="password"
              value={values.confirm}
            />
          </Stack>
        </CardContent>
        <Divider />
        <CardActions sx={{ justifyContent: 'flex-end' }}>
          <Button variant="contained"
            onClick={async ()=>{
              if (values.password !== values.confirm) {
                console.log("Passwords do not match")
                setAlert({
                  severity: 'error',
                  message: 'Passwords do not match'
                });
                return
              }
              let {status} = await httpRequest('/api/auth/password', 'PUT', JSON.stringify({"password": values.password}))
              if (status === 204) {
                console.log("Password updated")
                setAlert({
                  severity: 'success',
                  message: 'Password updated'
                });
              }
            }}
          >
            Update
          </Button>
        </CardActions>
      </Card>
    </form>
  );
};
