import React from "react";
import Sidebar from "./Sidebar";
import { Box, Container, IconButton, Typography } from '@mui/material';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';
import DeleteIcon from '@mui/icons-material/Delete';
import LogoutButton from "./LogoutButton";

type SettingPageProps = {
    darkMode: boolean;
    toggleTheme: () => void;
  };

const SettingPage: React.FC<SettingPageProps> = ({ darkMode, toggleTheme }) => {

    return (
        <div>
            <Sidebar />
            <Container sx={{ marginTop: '100px', marginLeft:'50px' }}> {/* Add top margin here */}
                <Typography variant="h4">Identification</Typography>
                <Typography variant="body1">Your public key: </Typography>

                <Typography variant="h6" sx={{ marginTop: '20px' }}>Theme</Typography>
                <IconButton onClick={toggleTheme}>
                    {darkMode ? <Brightness7Icon /> : <Brightness4Icon />}
                </IconButton>

                <Box sx={{ marginTop: '20px' }}>
                    <LogoutButton sx={{ display: 'block', visibility: 'visible' }} /> {/* Force visibility */}
                </Box>


                <IconButton sx={{ marginTop: '20px' }}> 
                    <DeleteIcon/>  Delete Account
                </IconButton>
            </Container>
        </div>
    );
}

export default SettingPage;
