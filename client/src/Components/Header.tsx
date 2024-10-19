import React from "react";
import { AppBar, Toolbar, Typography, Box } from "@mui/material";
import { useNavigate } from "react-router-dom";

const Header: React.FC = () => {
  const navigate = useNavigate();

  const handleWelcome = () => {
    navigate("/");
  };

  return (
    <AppBar position="static" color="primary">
      <Toolbar>
        <Box
          onClick={handleWelcome}
          sx={{
            display: "flex",
            alignItems: "center",
            flexGrow: 1,
            cursor: "pointer",
          }}
        >
          <img
            src={`${process.env.PUBLIC_URL}/squidcoin.png`}
            alt="Squid Icon"
            width="30"
            height="30"
            style={{ marginRight: "8px" }} // Adds spacing between the icon and text
          />
          <Typography variant="h6" component="div">
            Squid Coin
          </Typography>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
