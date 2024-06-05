import { useCallback, useEffect, useState } from 'react';
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
  SvgIcon,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions, 
  Box,
  IconButton,
  Icon,
  Checkbox,
  FormControlLabel
} from '@mui/material';
import XMarkIcon from '@heroicons/react/24/outline/XMarkIcon';
import PaperAirplane from '@heroicons/react/24/solid/PaperAirplaneIcon';
import Check from '@heroicons/react/24/outline/CheckIcon'
import CircularProgress from '@mui/material/CircularProgress';
import Backdrop from '@mui/material/Backdrop';
import { useRouter } from 'next/router';
import GlobeAltIcon from '@heroicons/react/24/outline/GlobeAltIcon';


export const DevicesWiFi = () => {

  return (
        <Stack 
        direction="row" 
        spacing={2}   
        justifyContent="center" 
        alignItems="center"
        >
            <Card>
                <CardHeader
                    title="2.4GHz"
                    avatar={
                        <SvgIcon>
                                <GlobeAltIcon/>
                        </SvgIcon>
                    }
                />
                <CardContent>
                    <Stack spacing={3}>
                        <FormControlLabel control={<Checkbox defaultChecked />} label="Enabled" />
                        <TextField
                            fullWidth
                            label="SSID"
                            value="wlan0"
                            variant="outlined"
                        />
                        <TextField
                            fullWidth
                            label="Encryption"
                            value="WPA2-PSK"
                            variant="outlined"
                        />
                        <TextField
                            fullWidth
                            label="Key"
                            value="password"
                            variant="outlined"
                        />
                    </Stack>
                </CardContent>
            </Card>
            <Card>
                <CardHeader
                    title="5GHz"
                    avatar={
                        <SvgIcon>
                                <GlobeAltIcon/>
                        </SvgIcon>
                    }
                />
                <CardContent>
                    <Stack spacing={4}>
                        <FormControlLabel control={<Checkbox defaultChecked />} label="Enabled" />
                        <TextField
                            fullWidth
                            label="SSID"
                            value="wlan0"
                            variant="outlined"
                        />
                        <FormControl variant="outlined" sx={{ m: 1, minWidth: 120 }}>
                            <InputLabel id="demo-simple-select-standard-label">Security</InputLabel>
                            <Select
                            labelId="demo-simple-select-standard-label"
                            id="demo-simple-select-standard"
                            value={"WPA2-PSK"}
                            //onChange={handleChange}
                            label="Security"
                            >
                            <MenuItem value={30}>Open</MenuItem>
                            <MenuItem value={"WPA2-PSK"}>WPA2-PSKnp</MenuItem>
                            <MenuItem value={20}>WPA3</MenuItem>
                            </Select>
                        </FormControl>
                        <TextField
                            fullWidth
                            label="Key"
                            value="password"
                            variant="outlined"
                        />
                    </Stack>
                    <CardActions sx={{display:"flex", justifyContent:"flex-end"}}>
                    <Button 
                        variant="contained" 
                        endIcon={<SvgIcon><Check /></SvgIcon>} 
                       // onClick={}
                        sx={{mt:'25px', mb:'-15px'}}
                        >
                        Apply
                    </Button>
                    </CardActions>
                </CardContent>
            </Card>
        </Stack>
      
  );
};
