import React from 'react';
import { Link } from 'react-router-dom';
import { Box,CssBaseline,Drawer,IconButton,List,ListItem,ListItemText
  , ListItemIcon,Toolbar,Typography,Button
} from '@mui/material';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import HomeIcon from '@mui/icons-material/Home';
import DashboardIcon from '@mui/icons-material/Dashboard';
import SettingsIcon from '@mui/icons-material/Settings';
import MenuIcon from '@mui/icons-material/Menu';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import SupportAgentIcon from '@mui/icons-material/SupportAgent';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { FileCopy } from '@mui/icons-material';
import { styled, useTheme } from '@mui/material/styles';
import Divider from '@mui/material/Divider';

const drawerWidth = 240;


const Main = styled('main', { shouldForwardProp: (prop) => prop !== 'open' })<{
  open?: boolean;
}>(({ theme }) => ({
  flexGrow: 1,
  padding: theme.spacing(3),
  transition: theme.transitions.create('margin', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  marginLeft: `-${drawerWidth}px`,
  variants: [
    {
      props: ({ open }) => open,
      style: {
        transition: theme.transitions.create('margin', {
          easing: theme.transitions.easing.easeOut,
          duration: theme.transitions.duration.enteringScreen,
        }),
        marginLeft: 0,
      },
    },
  ],
}));

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme }) => ({
  transition: theme.transitions.create(['margin', 'width'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  variants: [
    {
      props: ({ open }) => open,
      style: {
        width: `calc(100% - ${drawerWidth}px)`,
        marginLeft: `${drawerWidth}px`,
        transition: theme.transitions.create(['margin', 'width'], {
          easing: theme.transitions.easing.easeOut,
          duration: theme.transitions.duration.enteringScreen,
        }),
      },
    },
  ],
}));

const DrawerHeader = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  padding: theme.spacing(0, 1),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar,
  justifyContent: 'flex-end',
}));

const Sidebar: React.FC = () => 
{
  const theme = useTheme();
  const [open,setOpen]= React.useState(false);

  const toggleDrawer = (newOpen:boolean) =>
  {
    setOpen(newOpen);
  };

  const drawer = (
    <div>
      <List>
        <ListItem component = {Link} to = "/">
          <ListItemIcon><HomeIcon/></ListItemIcon>
          <ListItemText primary = "home"/>
        </ListItem>
        <ListItem component = {Link} to ="/dashboard">
        <ListItemIcon><DashboardIcon/></ListItemIcon>
        <ListItemText primary = "Dashboard"/>
        </ListItem>

        <ListItem component = {Link} to = "/files">
          <ListItemIcon><FileCopyIcon/></ListItemIcon>
          <ListItemText primary = "Files"></ListItemText>
        </ListItem>

        <ListItem component = {Link} to = "/account">
          <ListItemIcon><AccountCircleIcon/></ListItemIcon>
          <ListItemText primary = "Account"></ListItemText>
        </ListItem>

        <ListItem component = {Link} to = "/support">
          <ListItemIcon><SupportAgentIcon/></ListItemIcon>
          <ListItemText primary = "Support"></ListItemText>
        </ListItem>

        <ListItem component = {Link} to = "/settings">
          <ListItemIcon><SettingsIcon/></ListItemIcon>
          <ListItemText primary = "Files"></ListItemText>
        </ListItem>
      </List>
      </div>
  )


  return (
    <Box sx = {{display:'flex'}}>
      <CssBaseline/>
      <AppBar position = "fixed" open = {open}>
        <Toolbar>
          <IconButton
            color = "inherit"
            aria-label = "open-drawer"
            onClick = {()=>toggleDrawer(true)}
            edge = "start"
            sx={[
              {
                mr: 2,
              },
              open && { display: 'none' },
            ]}
            >
              <MenuIcon />
            </IconButton>
        </Toolbar>
      </AppBar>

      <Drawer 
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: drawerWidth,
            boxSizing: 'border-box',
          },
        }}
        variant = "persistent"
        anchor = "left"
        open = {open}
      >
        <DrawerHeader>
          <IconButton onClick ={()=>toggleDrawer(false)}>
            {theme.direction === 'ltr'? <ChevronLeftIcon />:<ChevronRightIcon />}
          </IconButton>
        </DrawerHeader>
        <Divider />
        {drawer}
      </Drawer>
    </Box>
  );
}

export default Sidebar;