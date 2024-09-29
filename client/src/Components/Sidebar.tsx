
import React from 'react';
import { Link } from 'react-router-dom';
import { Box,CssBaseline,Drawer,IconButton,List,ListItem,ListItemText
  , ListItemIcon,Toolbar,Typography,Button,Collapse,TextField
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
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import NotificationsIcon from '@mui/icons-material/Notifications';
import DarkModeIcon from '@mui/icons-material/DarkMode';

const drawerWidth = 275;


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
  const [dashboardOpen,setDashBoardOpen] = React.useState(false);
  

  const toggleDrawer = (newOpen:boolean) =>
  {
    setOpen(newOpen);
  };

  const toggleDashboard = () =>
  {
    setDashBoardOpen(!dashboardOpen);
  }


  const drawer = (
    <div>
      <List>
        {/*
        Not sure if this is the intended logic 
        */}
       <ListItem component = {Button} onClick={()=>toggleDashboard()}>
          <ListItemIcon><DashboardIcon /></ListItemIcon>
          <ListItemText primary={
      <Typography style={{ textTransform: 'none' }}>
        Dashboard
      </Typography>
    }  />
          {dashboardOpen ? <ExpandLessIcon /> : <ExpandMoreIcon />} {/* Icon to indicate collapse/expand */}
        </ListItem>

        <Collapse in = {dashboardOpen} timeout = "auto" unmountOnExit>
          <List component = "div" disablePadding>
            <ListItem component = {Link} to = "/overview" sx={{ pl: 4 }} >
              <ListItemIcon></ListItemIcon>
              <ListItemText primary = "Overview"></ListItemText>
            </ListItem>

            <ListItem component = {Link} to = "/notifications" sx={{ pl: 4 }} >
              <ListItemIcon></ListItemIcon>
              <ListItemText primary = "Notifications"></ListItemText>
            </ListItem>

            <ListItem component = {Link} to = "/trade-history" sx={{ pl: 4 }} >
              <ListItemIcon></ListItemIcon>
              <ListItemText primary = "Trade History"></ListItemText>
            </ListItem>
          </List>
        </Collapse>


        {/*
        Might need to implement file drop down functionality here too 
        */}

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
          <ListItemText primary = "Settings"></ListItemText>
        </ListItem>

        <Divider />

      </List>
      </div>
  )


  return (
    <Box sx = {{display:'flex'}}>
      <CssBaseline/>
      <AppBar position = "fixed" open = {open} sx = {{backgroundColor: 'white'}} >
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
            <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1,
            }}>
            Squid coin
          </Typography>
          <IconButton color = "inherit" sx = {{ml:2}}>
            <DarkModeIcon/>
          </IconButton>

          <IconButton color = "inherit" component = {Link} to = "/notifications">
            <NotificationsIcon/>
          </IconButton>

          <Button 
            color = "inherit"
            sx = {{
              border: '2px solid #808080',
              borderRadius: '4px',
              padding: '6px 12px',
              ml: 2,
              textDecoration: 'none',
              color: 'black',
            }}>
            <Typography variant = "h6" component = {Link} to = "/login" 
              style={{ 
                textTransform: 'none',
                textDecoration: 'none',
                color: 'black', }}>
              Log in
            </Typography>
          </Button>

          <Button 
            color = "inherit"
            sx = {{
              border: '2px solid #808080',
              borderRadius: '4px',
              padding: '6px 12px',
              ml: 2,
              
            }}>
            <Typography variant = "h6" component = {Link} to = "/register" 
              style={{ 
                textTransform: 'none',
                textDecoration: 'none',
                color: 'black', }}>
              Sign up
            </Typography>
          </Button>

        </Toolbar>
        <Divider/>
        <Toolbar>
        <Button 
            color = "inherit"
            sx = {{
              padding: '6px 12px',
              ml: 10,
              
            }}>
            <Typography variant = "h6" component = {Link} to = "/dashboard" 
              style={{ 
                textTransform: 'none',
                textDecoration: 'none',
                color: 'black', }}>
              Home
            </Typography>
          </Button>

          <Button 
            color = "inherit"
            sx = {{
              padding: '6px 12px',
              ml: 4,
              
            }}>
            <Typography variant = "h6" component = {Link} to = "/about" 
              style={{ 
                textTransform: 'none',
                textDecoration: 'none',
                color: 'black', }}>
              About/Info Page
            </Typography>
          </Button>

          <Button 
            color = "inherit"
            sx = {{
              padding: '6px 12px',
              ml: 4,
              
            }}>
            <Typography variant = "h6" component = {Link} to = "/transactionsr" 
              style={{ 
                textTransform: 'none',
                textDecoration: 'none',
                color: 'black', }}>
              Transactions
            </Typography>
          </Button>

          <Button 
            color = "inherit"
            sx = {{
              padding: '6px 12px',
              ml: 4,
              
            }}>
            <Typography variant = "h6" component = {Link} to = "/admin" 
              style={{ 
                textTransform: 'none',
                textDecoration: 'none',
                color: 'black', }}>
              Admin page
            </Typography>
          </Button>

          <Button 
            color = "inherit"
            sx = {{
              padding: '6px 12px',
              ml: 4,
              
            }}>
            <Typography variant = "h6" component = {Link} to = "/faq" 
              style={{ 
                textTransform: 'none',
                textDecoration: 'none',
                color: 'black', }}>
              FAQ
            </Typography>
          </Button>

          <Button 
            color = "inherit"
            sx = {{
              padding: '6px 12px',
              ml: 4,
              
            }}>
            <Typography variant = "h6" component = {Link} to = "/trading" 
              style={{ 
                textTransform: 'none',
                textDecoration: 'none',
                color: 'black', }}>
              Trading
            </Typography>
          </Button>

          <Box sx={{ flexGrow: 1 }} /> {/* Pushes the search bar to the right */}
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

                '& fieldset':
                {
                  borderColor:'grey',
                },

                '&:hover fieldset':{
                  borderColor:'darkgrey'
                }
              },
            }}
          />
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