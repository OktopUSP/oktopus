import { useCallback, useState } from 'react';
import {
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Divider,
  Stack,
  TextField,
  InputLabel,
  MenuItem, 
  Select,
  FormControl,
  SvgIcon
} from '@mui/material';
import PaperAirplane from '@heroicons/react/24/solid/PaperAirplaneIcon';

export const DevicesRPC = () => {
  const [value, setValue] = useState(`
  {opa, 
    teste123:goiaba}`
  )

  const [age, setAge] = useState(1);

  const handleChangeRPC = (event) => {
    setAge(event.target.value);
  };

  const handleChange = (event) => {
    setValue(event.target.value);
  };

  const handleSubmit = useCallback(
    (event) => {
      event.preventDefault();
    },
    []
  );

  return (
    <form onSubmit={handleSubmit}>
      <Card>
        <CardActions sx={{ justifyContent: 'flex-end'}}>
            <FormControl sx={{width:'100px'}}>
                <Select
                    labelId="demo-simple-select-standard-label"
                    id="demo-simple-select-standard"
                    value={age}
                    label="Action"
                    onChange={(event)=>{handleChangeRPC(event)}}
                    variant='standard'
                >
                    <MenuItem value={1}>Create</MenuItem>
                    <MenuItem value={2}>Read</MenuItem>
                    <MenuItem value={3}>Update</MenuItem>
                    <MenuItem value={4}>Delete</MenuItem>
                </Select>
            </FormControl>
        </CardActions>
        <Divider />
        <CardContent>
          <Stack
            spacing={3}
            alignItems={'stretch'}
          >
            <TextField
              id="outlined-multiline-static"
              size="large"
              multiline="true"
              label="Mensagem"
              name="password"
              onChange={handleChange}
              value={value}
              fullWidth
              rows="10"
            />
          </Stack>
        </CardContent>
        <Divider />
        <CardActions sx={{ justifyContent: 'flex-end' }}>
          <Button variant="contained" endIcon={<SvgIcon><PaperAirplane /></SvgIcon>}>
            Send
          </Button>
        </CardActions>
      </Card>
    </form>
  );
};
