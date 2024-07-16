import PropTypes from 'prop-types';
import NextLink from 'next/link';
import Link from 'next/link'
import { Box, Typography, Unstable_Grid2 as Grid, Stack } from '@mui/material';
import { Logo } from 'src/components/logo';
import { useTheme, useMediaQuery } from '@mui/material'

export const Layout = (props) => {
  const { children } = props;
  const lgUp = useMediaQuery((theme) => theme.breakpoints.up('lg'));
  const theme = useTheme();

  console.log("logUp", lgUp)

  return (
    <Box
      component="main"
      sx={{
        display: 'flex',
        flex: '1 1 auto'
      }}
    >
      <Grid
        container
        sx={{ flex: '1 1 auto' }}
      >
        <Grid
          xs={12}
          lg={6}
          sx={{
            backgroundColor: 'background.paper',
            display: 'flex',
            flexDirection: 'column',
            position: 'relative'
          }}
        >
          <Box
            component="header"
            sx={{
              left: 0,
              p: 3,
              position: 'fixed',
              top: 0,
              width: '100%'
            }}
          >
            <Box
              component={NextLink}
              href="/"
              sx={{
                display: 'inline-flex',
                height: 32,
                width: 32
              }}
            >
              <Logo />
            </Box>
          </Box>
          {children}
        </Grid>
        <Grid
          xs={12}
          lg={6}
          sx={{
            alignItems: 'center',
            background: `radial-gradient(50% 50% at 50% 50%, ${theme.palette.primary.main} 0%, ${theme.palette.primary.dark } 100%)`,
            color: 'white',
            display: 'flex',
            justifyContent: 'center',
            '& img': {
              maxWidth: '100%'
            }
          }}
        >
          <Box sx={{ p: 3 }}>
            <Link href="http://localhost/companylink" target="_blank">
              <img
                alt=""
                src="/images/logo.png"
              />
            </Link>
          </Box>
        </Grid>
      </Grid>
      <Stack style={{position:"absolute", bottom:"2px", left:"2px"}} direction={"row"} spacing={"1"}>  
        <Typography
          align="center"
          color={lgUp ? 'neutral[900]' : 'primary.contrastText'}
          component="footer"
          variant="body2"
          sx={{ p: 2 }}
        >
          Powered by
        </Typography>
      </Stack>
      <a href='https://oktopus.app.br' style={{position:"absolute", bottom:"10px", left:"100px"}} target='_blank'>
      <img 
        src="/assets/logo.png" 
        alt="Oktopus logo image"
        width={80}/>
      </a>
    </Box>
  );
};

Layout.prototypes = {
  children: PropTypes.node
};