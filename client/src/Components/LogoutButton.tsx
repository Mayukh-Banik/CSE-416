import React from 'react';
import { Button } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const LogoutButton: React.FC = () => {
    const navigate = useNavigate();

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
                backgroundColor:'#f4f4f4',
                color:'#1876d2',
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
