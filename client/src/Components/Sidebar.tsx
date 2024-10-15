
import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Link, useNavigate } from 'react-router-dom';
import { Box,CssBaseline,Drawer,IconButton,List,ListItem,ListItemText
  , ListItemIcon,Toolbar,Typography,Button,Collapse,TextField
} from '@mui/material';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import SettingsIcon from '@mui/icons-material/Settings';
import MenuIcon from '@mui/icons-material/Menu';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { styled, useTheme } from '@mui/material/styles';
import Divider from '@mui/material/Divider';
import AccountBalanceWalletIcon from '@mui/icons-material/AccountBalanceWallet';
const drawerWidth = 275;


const Main = styled('main', { shouldForwardProp: (prop) => prop !== 'open' })<{
  open?: boolean;
}>(({ theme, open }) => ({
  flexGrow: 1,
  padding: theme.spacing(3),
  transition: theme.transitions.create('margin', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  marginLeft: `-${drawerWidth}px`,
  ...(open && {
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
    marginLeft: 0,
  }),
}));


interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
  transition: theme.transitions.create(['margin', 'width'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  ...(open && {
    width: `calc(100% - ${drawerWidth}px)`,
    marginLeft: `${drawerWidth}px`,
    transition: theme.transitions.create(['margin', 'width'], {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
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
  const [open, setOpen]= React.useState(false);
  const navigate = useNavigate();

  const handleFiles = async () => {
      navigate('/files');
  };

  const handleWallet = async () => {
    navigate('/wallet');
  }; 

  const handleMining = async () => {
    navigate('/mining');
  };  

  const handleAccount = async () => {
    navigate('/account');
  };

  const handleSettings = async () => {
    navigate('/settings');
  };

  const handleMarket = async () => {
    navigate('/market')
  }
  // const [dashboardOpen,setDashBoardOpen] = React.useState(false);
  

  const toggleDrawer = (newOpen:boolean) =>
  {
    setOpen(newOpen);
  };

  // const toggleDashboard = () =>
  // {
  //   setDashBoardOpen(!dashboardOpen);
  // }

  const drawer = (
    <div>
      <List>

        <ListItem onClick={handleMarket} sx={{cursor:"pointer", "&:hover": {backgroundColor: 'rgba(0, 0, 0, 0.08)', },}}>
          <ListItemIcon><FileCopyIcon/></ListItemIcon>
          <ListItemText primary = "Market"></ListItemText>
        </ListItem>

        <ListItem onClick={handleFiles} sx={{cursor:"pointer", "&:hover": {backgroundColor: 'rgba(0, 0, 0, 0.08)', },}}>
          <ListItemIcon><FileCopyIcon/></ListItemIcon>
          <ListItemText primary = "Files"></ListItemText>
        </ListItem>

        <ListItem onClick={handleWallet} sx={{cursor:"pointer", "&:hover": {backgroundColor: 'rgba(0, 0, 0, 0.08)', },}}>
          <ListItemIcon><AccountBalanceWalletIcon/></ListItemIcon>
          <ListItemText primary = "Wallet"></ListItemText>
        </ListItem>

        <ListItem onClick={handleMining} sx={{cursor:"pointer", "&:hover": {backgroundColor: 'rgba(0, 0, 0, 0.08)', },}}>
          <ListItemIcon>
            <img 
              src={`${process.env.PUBLIC_URL}/pickaxe.png`} 
              alt="Pickaxe Icon" 
              style={{ width: '24px', height: '24px', filter: 'invert'}} 
            />
          </ListItemIcon>
          <ListItemText primary = "Mining"></ListItemText>
        </ListItem>

        <ListItem onClick={handleAccount} sx={{cursor:"pointer", "&:hover": {backgroundColor: 'rgba(0, 0, 0, 0.08)', },}}>
          <ListItemIcon><AccountCircleIcon/></ListItemIcon>
          <ListItemText primary = "Account"></ListItemText>
        </ListItem>

        <ListItem onClick={handleSettings} sx={{cursor:"pointer", "&:hover": {backgroundColor: 'rgba(0, 0, 0, 0.08)', },}}>
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
      <AppBar position = "fixed" open = {open} sx={{backgroundColor:'primary.main'}}>
      {/* <AppBar position = "fixed" open = {open} sx={{backgroundColor:'#000000'}}> */}

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

          <Toolbar>
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
                backgroundColor: 'background.default',
                color: 'secondary.main',

                '& fieldset':
                {
                  borderColor:'white',
                },

                '&:hover fieldset':{
                  borderColor:'darkgrey'
                }
              },
            }}
          />
          </Toolbar>
        </Toolbar>
        <Divider/>
        
      </AppBar>

      <Drawer 
        sx={{
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: drawerWidth,
            boxSizing: 'border-box',
            backgroundColor:'background.default',
            color:'secondary.main'
          },
        }}
        variant = "persistent"
        anchor = "left"
        open = {open}
        
      >
        <DrawerHeader>
          <img src="/images/squidcoin.png" alt="Squid Icon" width="30" />
          <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1, margin: 1}}>
            Squid Coin
          </Typography>
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