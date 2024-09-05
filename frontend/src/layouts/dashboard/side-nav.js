import NextLink from 'next/link';
import { usePathname } from 'next/navigation';
import PropTypes from 'prop-types';
import Link from 'next/link'
import {
  Box,
  Button,
  Collapse,
  Divider,
  Drawer,
  List,
  ListItemButton,
  ListItemText,
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

    if (currentPath.includes(itemPath) && itemPath !== '/' && itemPath !== '/mass-actions') {
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
              <img src={`${process.env.NEXT_PUBLIC_REST_ENDPOINT || ""}/images/logo.png`} 
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
                  children={item?.children}
                  padleft={2}
                  tooltip={item.tooltip}
                />
              );
            })}
            <Collapse in={open} timeout="auto" unmountOnExit> 
              <Box
                component="span"
                sx={{
                  color: 'neutral.400',
                  flexGrow: 1,
                  fontFamily: (theme) => theme.typography.fontFamily,
                  fontSize: 14,
                  fontWeight: 600,
                  lineHeight: '24px',
                  whiteSpace: 'nowrap',
                  ...(true && {
                    color: 'common.white'
                  }),
                  ...(false && {
                    color: 'neutral.500'
                  })
                }}
              >
                {""}
              </Box>
            </Collapse>
            {/* <List>
              <Collapse in={true}>
                <List disablePadding>
                <ListItemButton sx={{ pl: 4 }}>
                  <ListItemText primary="Starred" />
                  oi
                </ListItemButton>
                </List>
              </Collapse>
            </List> */}
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
            background: `linear-gradient(120deg, ${theme.palette.neutral["800"]} 0%, ${theme.palette.primary.dark} 90%);`,
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
          background: `linear-gradient(120deg, ${theme.palette.neutral["800"]} 0%, ${theme.palette.primary.dark}  90%);`,
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
