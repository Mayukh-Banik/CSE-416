import React from "react";
import { AppBar, Toolbar, Typography } from "@mui/material";


const Header: React.FC = () => {
  return (
    <AppBar position="static" color="primary">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          Project Squid
        </Typography>
        {/* Logo */}
      </Toolbar>
    </AppBar>
  );
};

export default Header;
