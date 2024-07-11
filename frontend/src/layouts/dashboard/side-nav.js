import NextLink from 'next/link';
import { usePathname } from 'next/navigation';
import PropTypes from 'prop-types';
import Link from 'next/link'
import {
  Box,
  Button,
  Divider,
  Drawer,
  Stack,
  SvgIcon,
  Typography,
  useMediaQuery
} from '@mui/material';
import { Logo } from 'src/components/logo';
import { Scrollbar } from 'src/components/scrollbar';
import { items } from './config';
import { SideNavItem } from './side-nav-item';
import { useTheme } from '@mui/material';

export const SideNav = (props) => {
  const { open, onClose } = props;
  const pathname = usePathname();
  const lgUp = useMediaQuery((theme) => theme.breakpoints.up('lg'));

  const theme = useTheme();

  const isItemActive = (currentPath, itemPath) => {
    if (currentPath === itemPath) {
      return true;
    }

    if (currentPath.includes(itemPath) && itemPath !== '/') {
      return true;
    }

    return false;
  }

  const content = (
    <Scrollbar
      sx={{
        height: '100%',
        '& .simplebar-content': {
          height: '100%'
        },
        '& .simplebar-scrollbar:before': {
          background: 'neutral.400'
        }
      }}
    >
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          height: '100%'
        }}
      >
        <Box sx={{ p: 3 }}>
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
          <Box
            sx={{
              alignItems: 'center',
              backgroundColor: 'rgba(255, 255, 255, 0.04)',
              borderRadius: 1,
              cursor: 'pointer',
              display: 'flex',
              justifyContent: 'space-between',
              mt: 2,
              p: '12px'
            }}
          >
         <Link href="http://localhost/companylink" target="_blank">
            <div style={{display:'flex',justifyContent:'center'}}>
              <img src="/images/logo.png" 
              width={'60%'}
              />
            </div>
          </Link>
            <SvgIcon
              fontSize="small"
              sx={{ color: 'neutral.500' }}
            >
            </SvgIcon>
          </Box>
        </Box>
        <Divider sx={{ borderColor: 'neutral.700' }} />
        <Box
          component="nav"
          sx={{
            flexGrow: 1,
            px: 2,
            py: 3
          }}
        >
          <Stack
            component="ul"
            spacing={0.5}
            sx={{
              listStyle: 'none',
              p: 0,
              m: 0
            }}
          >
            {items.map((item) => {
              if (item.title == "Map" && process.env.NEXT_PUBLIC_ENTERPRISE_VERSION != "true"){
                return
              }
              const active = isItemActive(pathname, item.path);

              return (
                <SideNavItem
                  active={active}
                  disabled={item.disabled}
                  external={item.external}
                  icon={item.icon}
                  key={item.title}
                  path={item.path}
                  title={item.title}
                />
              );
            })}
          </Stack>
        </Box>
        <Stack style={{position:"absolute", bottom:"2px", left:"2px"}} direction={"row"} spacing={"1"} zIndex={9999}>  
        <Typography
          align="center"
          color="primary.contrastText"
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
        width={80}
        />
      </a>
      </Box>
    </Scrollbar>
  );

  if (lgUp) {
    return (
      <Drawer
        anchor="left"
        open
        PaperProps={{
          sx: {
            background: `linear-gradient(0deg, ${theme.palette.neutral["800"]} 0%, ${theme.palette.primary.dark} 90%);`,
            color: 'common.white',
            width: 280
          }
        }}
        variant="permanent"
      >
        {content}
      </Drawer>
    );
  }

  return (
    <Drawer
      anchor="left"
      onClose={onClose}
      open={open}
      PaperProps={{
        sx: {
          background: `linear-gradient(0deg, ${theme.palette.primary.main} 0%, ${theme.palette.primary.dark} 90%);`,
          color: 'common.white',
          width: 280
        }
      }}
      sx={{ zIndex: (theme) => theme.zIndex.appBar + 100 }}
      variant="temporary"
    >
      {content}
    </Drawer>
  );
};

SideNav.propTypes = {
  onClose: PropTypes.func,
  open: PropTypes.bool
};
