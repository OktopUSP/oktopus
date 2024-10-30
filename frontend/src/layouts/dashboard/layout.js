import { useCallback, useEffect, useState } from 'react';
import { usePathname } from 'next/navigation';
import { styled } from '@mui/material/styles';
import { withAuthGuard } from 'src/hocs/with-auth-guard';
import { SideNav } from './side-nav';
import { TopNav } from './top-nav';
import { useAlertContext } from 'src/contexts/error-context';
import { Alert, AlertTitle, Snackbar } from '@mui/material';

const SIDE_NAV_WIDTH = 280;

const LayoutRoot = styled('div')(({ theme }) => ({
  display: 'flex',
  flex: '1 1 auto',
  maxWidth: '100%',
  [theme.breakpoints.up('lg')]: {
    paddingLeft: SIDE_NAV_WIDTH
  }
}));

const LayoutContainer = styled('div')({
  display: 'flex',
  flex: '1 1 auto',
  flexDirection: 'column',
  width: '100%'
});

export const Layout = withAuthGuard((props) => {
  const { children } = props;
  const pathname = usePathname();
  const [openNav, setOpenNav] = useState(false);

  const {alert, setAlert} = useAlertContext();

  const handlePathnameChange = useCallback(
    () => {
      if (openNav) {
        setOpenNav(false);
      }
    },
    [openNav]
  );

  useEffect(
    () => {
      handlePathnameChange();
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [pathname]
  );

  return (
    <>
      {pathname != "/map" && <TopNav onNavOpen={() => setOpenNav(true)} />}
      <SideNav
        onClose={() => setOpenNav(false)}
        open={openNav}
      />
      <LayoutRoot>
        <LayoutContainer>
          {children}
        </LayoutContainer>
      </LayoutRoot>
      {alert && <Snackbar 
      open={true} 
      autoHideDuration={4000} 
      anchorOrigin={{vertical:'bottom', horizontal: 'right'}}
      onClose={() => setAlert(null)}
      >
        <Alert
          severity={alert?.severity}
          variant={alert?.severity == 'success' ? 'standard' : 'filled'}
          sx={{ width: '100%' }}
        >
          {alert?.title && <AlertTitle>{alert.title}</AlertTitle>}
          {alert?.message}
        </Alert>
      </Snackbar>}
    </>
  );
});