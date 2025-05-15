import ChartBarIcon from '@heroicons/react/24/solid/ChartBarIcon';
import CogIcon from '@heroicons/react/24/solid/CogIcon';
import ChatBubbleLeftRightIcon from '@heroicons/react/24/solid/ChatBubbleLeftRightIcon'
import RectangleGroupIcon from '@heroicons/react/24/solid/RectangleGroupIcon'
import ArrowDownOnSquareStackIcon from '@heroicons/react/24/solid/ArrowDownOnSquareStackIcon'
import UserGroupIcon from '@heroicons/react/24/solid/UserGroupIcon'
import KeyIcon from '@heroicons/react/24/solid/KeyIcon'
import CpuChip from '@heroicons/react/24/solid/CpuChipIcon';
import { SvgIcon } from '@mui/material';
import FolderIcon from '@heroicons/react/24/solid/FolderIcon';
import ShieldCheckIcon from '@heroicons/react/24/solid/ShieldCheckIcon';
import EnvelopeIcon from '@heroicons/react/24/solid/EnvelopeIcon';
import UserIcon from '@heroicons/react/24/solid/UserIcon';
import BuildingOffice2Icon from '@heroicons/react/24/solid/BuildingOffice2Icon';
import BookOpenIcon from '@heroicons/react/24/solid/BookOpenIcon';
import CommandLineIcon from '@heroicons/react/24/solid/CommandLineIcon';

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
    disabled: true,
    children: [
      {
        title: 'Firmware Update',
        icon: (
          <SvgIcon fontSize="small">
            <ArrowDownOnSquareStackIcon color='gray'/>
          </SvgIcon>
        ),
        disabled: true
      },
      {
        title: 'Message',
        disabled: true,
        icon: (
          <SvgIcon fontSize="small">
           <EnvelopeIcon color='gray'/>
          </SvgIcon>
        )
      },
    ]
  },
  {
    title: 'Scripts',
    disabled: true,
    icon: (
      <SvgIcon fontSize="small">
        <CommandLineIcon color='gray'/>
      </SvgIcon>
    )
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
    title: 'Access Control',
    disabled: true,
    icon: (
      <SvgIcon fontSize="small">
        <UserGroupIcon color='gray'/>
      </SvgIcon>
    ),
    children: [
      {
        title: 'Tenants',
        disabled: true,
        icon: (
          <SvgIcon fontSize="small">
            <BuildingOffice2Icon color='gray'/>
          </SvgIcon>
        )
      },
      {
        title: 'Roles',
        disabled: true,
        icon: (
          <SvgIcon fontSize="small">
            <ShieldCheckIcon color='gray'/>
          </SvgIcon>
        )
      },
      {
        title: 'Users',
        path: '/access-control/users',
        icon: (
          <SvgIcon fontSize="small">
           <UserIcon/>
          </SvgIcon>
        )
      },
     ]
   },
  {
    title: 'File  Server',
    disabled: true,
    icon: (
      <SvgIcon fontSize="small">
        <FolderIcon color='gray'/>
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
  {
    title: 'Docs',
    path: 'https://docs.oktopus.app.br',
    icon: (
      <SvgIcon fontSize="small">
        <BookOpenIcon />
      </SvgIcon>
    ),
    external: true,
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