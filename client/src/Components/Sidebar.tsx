// // Sidebar.tsx
// import React from 'react';
// import { Link } from 'react-router-dom';
// import { Drawer, List, ListItem, ListItemText, ListItemIcon } from '@mui/material';
// import HomeIcon from '@mui/icons-material/Home';
// import DashboardIcon from '@mui/icons-material/Dashboard';
// import SettingsIcon from '@mui/icons-material/Settings';

// const Sidebar: React.FC = () => {
//   return (
//     <Drawer variant="permanent" anchor="left">
//       <List>
//         <ListItem  component={Link} to="/">
//           <ListItemIcon><HomeIcon /></ListItemIcon>
//           <ListItemText primary="Home" />
//         </ListItem>
//         <ListItem  component={Link} to="/dashboard">
//           <ListItemIcon><DashboardIcon /></ListItemIcon>
//           <ListItemText primary="Dashboard" />
//         </ListItem>
//         <ListItem  component={Link} to="/transaction">
//           <ListItemIcon><DashboardIcon /></ListItemIcon>
//           <ListItemText primary="Transaction" />
//         </ListItem>
//         <ListItem  component={Link} to="/settings">
//           <ListItemIcon><SettingsIcon /></ListItemIcon>
//           <ListItemText primary="Settings" />
//         </ListItem>
//       </List>
//     </Drawer>
//   );
// }

// export default Sidebar;
import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Drawer, List, ListItem, ListItemIcon, ListItemText } from '@mui/material';
import HomeIcon from '@mui/icons-material/Home';
import DashboardIcon from '@mui/icons-material/Dashboard';
import SettingsIcon from '@mui/icons-material/Settings';
import ReceiptIcon from '@mui/icons-material/Receipt';
import { makeStyles } from '@mui/styles';

const useStyles = makeStyles({
  drawer: {
    width: 240, // Sidebar width
    flexShrink: 0,
  },
  drawerPaper: {
    width: 240, // Sidebar width
  },
  active: {
    backgroundColor: '#f4f4f4',
  },
});

const Sidebar: React.FC = () => {
  const classes = useStyles();
  const location = useLocation(); // Get the current path

  return (
    <Drawer
      className={classes.drawer}
      variant="permanent"
      anchor="left"
      classes={{ paper: classes.drawerPaper }}
    >
      <List>
        <ListItem
          component={Link}
          to="/"
          className={location.pathname === '/' ? classes.active : ''}
        >
          <ListItemIcon><HomeIcon /></ListItemIcon>
          <ListItemText primary="Home" />
        </ListItem>
        <ListItem
          component={Link}
          to="/dashboard"
          className={location.pathname === '/dashboard' ? classes.active : ''}
        >
          <ListItemIcon><DashboardIcon /></ListItemIcon>
          <ListItemText primary="Dashboard" />
        </ListItem>
        <ListItem
          component={Link}
          to="/transaction"
          className={location.pathname === '/transaction' ? classes.active : ''}
        >
          <ListItemIcon><ReceiptIcon /></ListItemIcon>
          <ListItemText primary="Transactions" />
        </ListItem>
        <ListItem
          component={Link}
          to="/settings"
          className={location.pathname === '/settings' ? classes.active : ''}
        >
          <ListItemIcon><SettingsIcon /></ListItemIcon>
          <ListItemText primary="Settings" />
        </ListItem>
      </List>
    </Drawer>
  );
};

export default Sidebar;
