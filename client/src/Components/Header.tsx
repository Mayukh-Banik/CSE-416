import React from "react";
import { AppBar, Toolbar, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";


const Header: React.FC = () => {
  const navigate = useNavigate();

  const handleWelcome = () => {
    navigate("/");
  };

  return (
    <AppBar position="static" color="primary">
      <Toolbar>
        <Typography variant="h6" component="div" onClick={handleWelcome} sx={{ flexGrow: 1, cursor: "pointer"}}>
          Project Squid
        </Typography>
        {/* Logo */}
      </Toolbar>
    </AppBar>
  );
};

export default Header;
