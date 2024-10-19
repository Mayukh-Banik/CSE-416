import React, { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import {
  Box, CssBaseline, Drawer, List, ListItem, ListItemText,
  ListItemIcon, Toolbar, Typography, TextField, Divider
} from '@mui/material';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import SettingsIcon from '@mui/icons-material/Settings';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import { styled, useTheme } from '@mui/material/styles';
import AccountBalanceWalletIcon from '@mui/icons-material/AccountBalanceWallet';
import SecurityIcon from '@mui/icons-material/Security';
import StoreIcon from '@mui/icons-material/Store';

const drawerWidth = 275;
const collapsedDrawerWidth = 80;

const Main = styled('main')(({ theme }) => ({
  flexGrow: 1,
  padding: theme.spacing(3),
  transition: theme.transitions.create(['margin'], {
    easing: theme.transitions.easing.easeInOut,
    duration: theme.transitions.duration.standard,
  }),
  marginLeft: `${drawerWidth}px`, 
  [theme.breakpoints.down('sm')]: {
    marginLeft: `${collapsedDrawerWidth}px`,
  },
}));

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar)<AppBarProps>(({ theme }) => ({
  transition: theme.transitions.create(['margin', 'width'], {
    easing: theme.transitions.easing.easeInOut,
    duration: theme.transitions.duration.standard,
  }),
  width: `calc(100% - ${drawerWidth}px)`, 
  marginLeft: `${drawerWidth}px`,
  [theme.breakpoints.down('sm')]: {
    width: `calc(100% - ${collapsedDrawerWidth}px)`,
    marginLeft: `${collapsedDrawerWidth}px`,
  },
}));

const DrawerHeader = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  padding: theme.spacing(0, 1),
  ...theme.mixins.toolbar,
  justifyContent: 'center', 
}));

// Reusable MenuItem component to avoid repetition
const MenuItem = ({ text, icon, onClick, active }: { text: string, icon: React.ReactNode, onClick: () => void, active: boolean }) => {
  return (
    <ListItem 
      onClick={onClick} 
      sx={{ 
        cursor: "pointer", 
        backgroundColor: active ? 'rgba(0, 0, 0, 0.12)' : 'transparent',
        "&:hover": { backgroundColor: 'rgba(0, 0, 0, 0.08)' },
      }}
    >
      <ListItemIcon>{icon}</ListItemIcon>
      <ListItemText primary={text} sx={{ display: { xs: 'none', sm: 'block' } }} />
    </ListItem>
  );
};

const Sidebar: React.FC = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const location = useLocation(); // Get the current route

  const isActive = (path: string) => location.pathname === path;

  const menuItems = [
    { text: 'Market', icon: <StoreIcon />, onClick: () => navigate('/market'), path: '/market' },
    { text: 'View/Upload Files', icon: <FileCopyIcon />, onClick: () => navigate('/files'), path: '/files' },
    { text: 'Wallet', icon: <AccountBalanceWalletIcon />, onClick: () => navigate('/wallet'), path: '/wallet' },
    { text: 'Mining', icon: <img src={`${process.env.PUBLIC_URL}/pickaxe.png`} 
      alt="Pickaxe Icon" 
      style={{ width: '24px', height: '24px', filter: 'invert' }} />, 
      onClick: () => navigate('/mining'), 
      path: '/mining' },
    { text: 'Account', icon: <AccountCircleIcon />, onClick: () => navigate('/account'), path: '/account' },
    { text: 'Proxy', icon: <SecurityIcon />, onClick: () => navigate('/proxy'), path: '/proxy' },
    { text: 'Settings', icon: <SettingsIcon />, onClick: () => navigate('/settings'), path: '/settings' },
  ];

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar position="fixed" sx={{ backgroundColor: 'primary.main' }}>
        <Toolbar>
          <Box sx={{ flexGrow: 1 }} />
          <TextField
            variant="outlined"
            placeholder="Search…"
            size="small"
            sx={{
              width: '250px',
              ml: 4,
              '& .MuiOutlinedInput-root': {
                borderRadius: '4px',
                borderColor: 'grey',
                backgroundColor: 'background.default',
                color: 'secondary.main',
                '& fieldset': { borderColor: 'white' },
                '&:hover fieldset': { borderColor: 'darkgrey' }
              },
            }}
          />
        </Toolbar>
      </AppBar>

      <Drawer
        sx={{
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: drawerWidth,
            boxSizing: 'border-box',
            backgroundColor: 'background.default',
            color: 'secondary.main',
            transition: theme.transitions.create('width', {
              easing: theme.transitions.easing.easeInOut,
              duration: theme.transitions.duration.standard,
            }),
            [theme.breakpoints.down('sm')]: {
              width: collapsedDrawerWidth,
            },
          },
        }}
        variant="permanent"
        anchor="left"
      >
        <DrawerHeader>
          <img src={`${process.env.PUBLIC_URL}/squidcoin.png`} alt="Squid Icon" width="30" />
          <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1, margin: 1, display: { xs: 'none', sm: 'block' } }}>
            SquidNet
          </Typography>
        </DrawerHeader>
        <Divider />
        <List>
          {menuItems.map((item, index) => (
            <MenuItem 
              key={index}
              text={item.text}
              icon={item.icon}
              onClick={item.onClick}
              active={isActive(item.path)} // Highlight the active page
            />
          ))}
        </List>
        <Divider />
      </Drawer>
    </Box>
  );
};

export default Sidebar;
