import { alpha } from '@mui/material/styles';
  
const withAlphas = (color) => {
  return {
    ...color,
    alpha4: alpha(color.main, 0.04),
    alpha8: alpha(color.main, 0.08),
    alpha12: alpha(color.main, 0.12),
    alpha30: alpha(color.main, 0.30),
    alpha50: alpha(color.main, 0.50)
  };
};

export const neutral = (colors) => {
  console.log("neutral colors:", colors);
  let parsedColors = JSON.parse(colors);

  let tableColor = parsedColors["tables"]
  let sidebarColorInitial = parsedColors["sidebar_initial"]
  let wordsOutsideSidebarColor = parsedColors["words_outside_sidebar"]

  return {
  50: tableColor,
  100: '#F3F4F6',
  200: '#E5E7EB',
  300: '#D2D6DB',
  400: '#FFFFFF',
  500: '#6C737F',
  600: '#4D5761',
  700: '#FFFFFF',
  800: sidebarColorInitial,
  900: wordsOutsideSidebarColor,
}
};

export const indigo = (colors) => {

  console.log("indigo colors:", colors);
  let parsedColors = JSON.parse(colors);

  let buttonColor = parsedColors["buttons"]
  let sidebarColorEnd = parsedColors["sidebar_end"]
  let mtpsColor = parsedColors["connected_mtps_color"]

  return withAlphas({
    lightest: '#FFFFFF',
    light: '#ff3383',
    main: buttonColor,
    dark: sidebarColorEnd,
    darkest: mtpsColor,
    contrastText: '#FFFFFF'
  });
}

export const success = withAlphas({
  lightest: '#F0FDF9',
  light: '#3FC79A',
  main: '#10B981',
  dark: '#0B815A',
  darkest: '#134E48',
  contrastText: '#FFFFFF'
});

export const info = withAlphas({
  lightest: '#ECFDFF',
  light: '#CFF9FE',
  main: '#06AED4',
  dark: '#0E7090',
  darkest: '#164C63',
  contrastText: '#FFFFFF'
});

export const warning = withAlphas({
  lightest: '#FFFAEB',
  light: '#FEF0C7',
  main: '#F79009',
  dark: '#B54708',
  darkest: '#7A2E0E',
  contrastText: '#FFFFFF'
});

export const error = withAlphas({
  lightest: '#FEF3F2',
  light: '#FEE4E2',
  main: '#F04438',
  dark: '#B42318',
  darkest: '#7A271A',
  contrastText: '#FFFFFF'
});

export const graphics = withAlphas({
  lightest: '#9EC8B9',
  light: '#706233',
  main: '#1B4242',
  dark: '#FFC5C5',
  darkest: '#7071E8'
});