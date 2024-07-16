import { common } from '@mui/material/colors';
import { alpha } from '@mui/material/styles';
import { error, indigo, info, neutral, success, warning, graphics } from './colors';

const getColorScheme = async () => {

  let result = await fetch('/custom-frontend/colors').catch((error) => {
    console.log('Error fetching colors');
    sessionStorage.setItem('colors', JSON.stringify({
      "buttons": "#306d6f",
      "sidebar_end": "#306d6f",
      "sidebar_initial": "#306d6f",
      "tables": "#306d6f",
      "words_outside_sidebar": "#30596f",
      "connected_mtps_color": "#f28950"
    }));
    location.reload();
    return
  });
  if (result.status!=200) {
    console.log('Error fetching colors');
    sessionStorage.setItem('colors', JSON.stringify({
      "buttons": "#306d6f",
      "sidebar_end": "#306d6f",
      "sidebar_initial": "#306d6f",
      "tables": "#306d6f",
      "words_outside_sidebar": "#30596f",
      "connected_mtps_color": "#f28950"
    }));
    location.reload();
    return
  }

  let response = await result.json();
  let fmtresponse = JSON.stringify(response);
  sessionStorage.setItem('colors', fmtresponse);
  location.reload();
}

export function createPalette() {

  let colors = sessionStorage.getItem('colors');

  if (colors !== null) {
    console.log('colors already fetched');
  } else {
    getColorScheme();
  }
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
