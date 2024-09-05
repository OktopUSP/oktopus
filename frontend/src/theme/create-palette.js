import { common } from '@mui/material/colors';
import { alpha } from '@mui/material/styles';
import { error, indigo, info, neutral, success, warning, graphics } from './colors';

const getColorScheme =  () => {
  return JSON.stringify({
    "buttons": "#c05521",
    "sidebar_end": "#305a85",
    "sidebar_initial": "#173033",
    "tables": "#214256",
    "words_outside_sidebar": "#173033",
    "connected_mtps_color": "#c05521"
  });
}

export function createPalette() {

  let colors = getColorScheme();
  console.log("colors scheme:", colors);

  let neutralColors = neutral(colors);

  return { 
    action: {
      active: neutralColors[500],
      disabled: alpha(neutralColors[900], 0.38),
      disabledBackground: alpha(neutralColors[900], 0.12),
      focus: alpha(neutralColors[900], 0.16),
      hover: alpha(neutralColors[900], 0.04),
      selected: alpha(neutralColors[900], 0.12)
    },
    background: {
      default: common.white,
      paper: common.white
    },
    divider: '#F2F4F7',
    error,
    graphics,
    info,
    mode: 'light',
    neutral: neutralColors,
    primary: indigo(colors),
    success,
    text: {
      primary: neutralColors[900],
      secondary: neutralColors[500],
      disabled: alpha(neutralColors[900], 0.38)
    },
    warning
  };
}
