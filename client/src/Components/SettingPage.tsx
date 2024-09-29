import React, { useState }  from "react";
import Sidebar from "./Sidebar";
import Header from "./Header";
import { Container, IconButton, Typography } from '@mui/material';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';
import DeleteIcon from '@mui/icons-material/Delete';

type SettingPageProps = {
    darkMode: boolean;
    toggleTheme: () => void;
  };

const SettingPage: React.FC<SettingPageProps> = ({ darkMode, toggleTheme }) => {

    return (
        <div>
            <Sidebar />
            <Header />
            <h1></h1>
            <h2>Identification</h2>
            Your public key: 
            <div>
            <Typography variant="h6">Theme</Typography>
            <IconButton onClick={toggleTheme}>
            {darkMode ? <Brightness7Icon /> : <Brightness4Icon />}
            </IconButton>
            <br></br>
            <IconButton> <DeleteIcon/>  Delete Account</IconButton>
        </div>

        </div>
    );
}

export default SettingPage