import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
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
const collapsedDrawerWidth = 80; // Width when collapsed

const Main = styled('main')(({ theme }) => ({
  flexGrow: 1,
  padding: theme.spacing(3),
  transition: theme.transitions.create(['margin'], {
    easing: theme.transitions.easing.easeInOut,
    duration: theme.transitions.duration.standard,
  }),
  marginLeft: `${drawerWidth}px`, // Default expanded
  [theme.breakpoints.down('sm')]: {
    marginLeft: `${collapsedDrawerWidth}px`, // Collapsed on small screens
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
  width: `calc(100% - ${drawerWidth}px)`, // Default expanded
  marginLeft: `${drawerWidth}px`,
  [theme.breakpoints.down('sm')]: {
    width: `calc(100% - ${collapsedDrawerWidth}px)`, // Adjust width when collapsed
    marginLeft: `${collapsedDrawerWidth}px`,
  },
}));

const DrawerHeader = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  padding: theme.spacing(0, 1),
  ...theme.mixins.toolbar,
  justifyContent: 'center', // Center the content when collapsed
}));

const Sidebar: React.FC = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');

  const handleSearchKeyPress = (event: React.KeyboardEvent<HTMLDivElement>) => {
    if (event.key === 'Enter') {
        navigate(`/search-page?q=${searchQuery}`); // Include search query in the URL
    }
};

  const handleFiles = async () => navigate('/files');
  const handleWallet = async () => navigate('/wallet');
  const handleMining = async () => navigate('/mining');
  const handleAccount = async () => navigate('/account/1');
  const handleSettings = async () => navigate('/settings');
  const handleMarket = async () => navigate('/market');
  const handleProxy = async () => navigate('/proxy');
  const handleGlobalTransactions = async () => navigate('/global-transactions');

  const drawer = (
    <div>
      <List>
        <ListItem onClick={handleMarket} sx={{ cursor: "pointer", "&:hover": { backgroundColor: 'rgba(0, 0, 0, 0.08)', }, }}>
          <ListItemIcon><StoreIcon /></ListItemIcon>
          <ListItemText primary="Market" sx={{ display: { xs: 'none', sm: 'block' } }} />
        </ListItem>

        <ListItem onClick={handleFiles} sx={{cursor:"pointer", "&:hover": {backgroundColor: 'rgba(0, 0, 0, 0.08)', },}}>
          <ListItemIcon><FileCopyIcon/></ListItemIcon>
          <ListItemText primary="View/Upload Files" sx={{ display: { xs: 'none', sm: 'block' } }} />
        </ListItem>

        <ListItem onClick={handleMining} sx={{ cursor: "pointer", "&:hover": { backgroundColor: 'rgba(0, 0, 0, 0.08)', }, }}>
          <ListItemIcon>
            <img
              src={`${process.env.PUBLIC_URL}/pickaxe.png`}
              alt="Pickaxe Icon"
              style={{ width: '24px', height: '24px', filter: 'invert' }}
            />
          </ListItemIcon>
          <ListItemText primary="Mining" sx={{ display: { xs: 'none', sm: 'block' } }} />
        </ListItem>

        <ListItem onClick={handleAccount} sx={{ cursor: "pointer", "&:hover": { backgroundColor: 'rgba(0, 0, 0, 0.08)', }, }}>
          <ListItemIcon><AccountCircleIcon /></ListItemIcon>
          <ListItemText primary="Account" sx={{ display: { xs: 'none', sm: 'block' } }} />
        </ListItem>

        <ListItem onClick={handleGlobalTransactions} sx={{ cursor: "pointer", "&:hover": { backgroundColor: 'rgba(0, 0, 0, 0.08)', }, }}>
          <ListItemIcon><AccountCircleIcon /></ListItemIcon>
          <ListItemText primary="Global Transactions" sx={{ display: { xs: 'none', sm: 'block' } }} />
        </ListItem>

        <ListItem onClick={handleProxy} sx={{ cursor: "pointer", "&:hover": { backgroundColor: 'rgba(0, 0, 0, 0.08)', }, }}>
          <ListItemIcon>
            <SecurityIcon />
          </ListItemIcon>
          <ListItemText primary="Proxy" />
        </ListItem>

        <ListItem onClick={handleSettings} sx={{ cursor: "pointer", "&:hover": { backgroundColor: 'rgba(0, 0, 0, 0.08)', }, }}>
          <ListItemIcon><SettingsIcon /></ListItemIcon>
          <ListItemText primary="Settings" sx={{ display: { xs: 'none', sm: 'block' } }} />
        </ListItem>
        <Divider />
      </List>
    </div>
  );

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
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)} // Update search query state
            onKeyPress={handleSearchKeyPress} // Handle Enter key press
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
              width: collapsedDrawerWidth, // Collapse sidebar on small screens
            },
          },
        }}
        variant="permanent"
        anchor="left"
      >
        <DrawerHeader>
          <img src={`${process.env.PUBLIC_URL}/squidcoin.png`} alt="Squid Icon" width="30" />
          <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1, margin: 1, display: { xs: 'none', sm: 'block' }}}>
            SquidNet
          </Typography>
        </DrawerHeader>
        <Divider />
        {drawer}
      </Drawer>
    </Box>
  );
}

export default Sidebar;