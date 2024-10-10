import React from 'react';
import { Button } from '@mui/material';
import { styled, useTheme } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';

const LogoutButton: React.FC = () => {
    const navigate = useNavigate();
    const theme = useTheme();

    // work on this later
    const handleLogout = async () => {
        // bring to welcome page
        navigate('/');
    };

    return (
        <Button
            variant="contained"
            color="inherit"
            onClick={handleLogout}
            style={{
                backgroundColor:'background.default',
                color:'secondary.main',
                position: 'absolute',
                top: '16px',
                right: '16px',
                transition: 'left 0.3s ease, right 0.3s ease',
            }}
        >
        Logout
        </Button>
    );
};

export default LogoutButton;
