import React from "react";
import Sidebar from "./Sidebar";
import { Box, Container, IconButton, Typography } from '@mui/material';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';
import DeleteIcon from '@mui/icons-material/Delete';
import LogoutButton from "./LogoutButton";
import { useTheme } from '@mui/material/styles';

type SettingPageProps = {
    darkMode: boolean;
    toggleTheme: () => void;
  };

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const SettingPage: React.FC<SettingPageProps> = ({ darkMode, toggleTheme }) => {
    const theme = useTheme();

    return (
        <Box
            sx={{
            padding: 2,
            marginTop: '70px',
            marginLeft: `${drawerWidth}px`, // Default expanded margin
            transition: 'margin-left 0.3s ease', // Smooth transition
            [theme.breakpoints.down('sm')]: {
                marginLeft: `${collapsedDrawerWidth}px`, // Adjust left margin for small screens
            },
            }}
        >
            <Sidebar />
            <Container> {/* Add top margin here */}
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
        </Box>
    );
}

export default SettingPage;
