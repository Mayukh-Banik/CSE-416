import React from 'react';
import { Button } from '@mui/material';
import { useNavigate } from "react-router-dom";

const LogoutButton: React.FC<React.ComponentProps<typeof Button>> = (props) => {
    const navigate = useNavigate();

    const handleLogout = () => {
        navigate('/');
        console.log('User logged out'); // Placeholder for logout logic
    };

    return (
        <Button
            variant="contained"
            color="primary"
            onClick={handleLogout}
            {...props} // Spread the additional props like style
        >
            Logout
        </Button>
    );
};

export default LogoutButton;
