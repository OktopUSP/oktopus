import ChartBarIcon from '@heroicons/react/24/solid/ChartBarIcon';
import CogIcon from '@heroicons/react/24/solid/CogIcon';
import ChatBubbleLeftRightIcon from '@heroicons/react/24/solid/ChatBubbleLeftRightIcon'
import MapIcon from '@heroicons/react/24/solid/MapIcon'
import RectangleGroupIcon from '@heroicons/react/24/solid/RectangleGroupIcon'
import ArrowDownOnSquareStackIcon from '@heroicons/react/24/solid/ArrowDownOnSquareStackIcon'
import UserGroupIcon from '@heroicons/react/24/solid/UserGroupIcon'
import KeyIcon from '@heroicons/react/24/solid/KeyIcon'
import CpuChip from '@heroicons/react/24/solid/CpuChipIcon';
import BriefCaseIcon from '@heroicons/react/24/outline/BriefcaseIcon';
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
        <CpuChip />
      </SvgIcon>
    )
  },
  {
    title: 'Mass Actions',
    icon: (
      <SvgIcon fontSize="small">
        <RectangleGroupIcon color='gray'/>
      </SvgIcon>
    ),
    tooltip: 'Upgrade to Business Plan',
    disabled: true,
    children: [
      {
        title: 'Firmware Update',
        tooltip: 'Upgrade to Business Plan',
        icon: (
          <SvgIcon fontSize="small">
            <ArrowDownOnSquareStackIcon color='gray'/>
          </SvgIcon>
        ),
        disabled: true
      }
    ]
  },
  {
    title: 'Map',
    tooltip: 'Upgrade to Business Plan',
    icon: (
      <SvgIcon fontSize="small">
        <MapIcon color='gray'/>
      </SvgIcon>
    ),
    disabled: true
  },
  {
    title: 'Credentials',
    path: '/credentials',
    icon: (
      <SvgIcon fontSize="small">
        <KeyIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Users',
    path: '/users',
    icon: (
      <SvgIcon fontSize="small">
        <UserGroupIcon />
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