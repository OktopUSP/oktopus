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
  SnackbarContent,
  Snackbar,
  Checkbox,
  FormControlLabel,
  useTheme,
} from '@mui/material';
import XMarkIcon from '@heroicons/react/24/outline/XMarkIcon';
import Check from '@heroicons/react/24/outline/CheckIcon';
//import ExclamationTriangleIcon from '@heroicons/react/24/solid/ExclamationTriangleIcon';
import CircularProgress from '@mui/material/CircularProgress';
import Backdrop from '@mui/material/Backdrop';
import { useRouter } from 'next/router';
import GlobeAltIcon from '@heroicons/react/24/outline/GlobeAltIcon';

export const DevicesDiagnostic = () => {
    return (
        <div>
            <p>Diagnostic Page</p>
        </div>
    )
};
