import ChartBarIcon from '@heroicons/react/24/solid/ChartBarIcon';
import CogIcon from '@heroicons/react/24/solid/CogIcon';
import ChatBubbleLeftRightIcon from '@heroicons/react/24/solid/ChatBubbleLeftRightIcon'
import MapIcon from '@heroicons/react/24/solid/MapIcon'
import UserGroupIcon from '@heroicons/react/24/solid/UserGroupIcon'
import KeyIcon from '@heroicons/react/24/solid/KeyIcon'
import CpuChip from '@heroicons/react/24/solid/CpuChipIcon';
import { SvgIcon } from '@mui/material';

export const items = [
  {
    title: 'Overview',
    path: '/',
    icon: (
      <SvgIcon fontSize="small">
        <ChartBarIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Devices',
    path: '/devices',
    icon: (
      <SvgIcon fontSize="small">
       <CpuChip/>
      </SvgIcon>
    )
  },
  // {
  //   title: 'Chat (beta)',
  //   path: '/chat',
  //   icon: (
  //     <SvgIcon fontSize="small">
  //      <ChatBubbleLeftRightIcon/>
  //     </SvgIcon>
  //   )
  // },
   {
     title: 'Map',
     path: '/map',
     icon: (
       <SvgIcon fontSize="small">
        <MapIcon/>
       </SvgIcon>
     ),
   },
  {
    title: 'Credentials',
    path: '/credentials',
    icon: (
      <SvgIcon fontSize="small">
      <KeyIcon/>
      </SvgIcon>
    )
  },
  {
    title: 'Users',
    path: '/users',
    icon: (
      <SvgIcon fontSize="small">
      <UserGroupIcon/>
      </SvgIcon>
    )
   },
  {
    title: 'Settings',
    path: '/settings',
    icon: (
      <SvgIcon fontSize="small">
        <CogIcon />
      </SvgIcon>
    )
  },
];

/*
  {
    title: 'Customers',
    path: '/customers',
    icon: (
      <SvgIcon fontSize="small">
        <UsersIcon />
      </SvgIcon>
    )
  },
    {
    title: 'Account',
    path: '/account',
    icon: (
      <SvgIcon fontSize="small">
        <UserIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Register',
    path: '/auth/register',
    icon: (
      <SvgIcon fontSize="small">
        <UserPlusIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Login',
    path: '/auth/login',
    icon: (
      <SvgIcon fontSize="small">
        <LockClosedIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Companies',
    path: '/companies',
    icon: (
      <SvgIcon fontSize="small">
        <ShoppingBagIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Error',
    path: '/404',
    icon: (
      <SvgIcon fontSize="small">
        <XCircleIcon />
      </SvgIcon>
    )
  }
*/ 