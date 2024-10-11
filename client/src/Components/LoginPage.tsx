import React, { useState } from 'react';
import axios from 'axios';
import { TextField, Button, Container, Typography, Link, Box } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import useRegisterPageStyles from '../Stylesheets/RegisterPageStyles';
import Header from './Header';
import { useNavigate } from 'react-router-dom';

const PORT = 5000;

const LoginPage: React.FC = () => {
    const classes = useRegisterPageStyles();
    const navigate = useNavigate();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState<string | null>(null);

    const validateForm = () => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

        if (!email || !password) {
            setError('All fields are required.');
            return false;
        }
        if (!emailRegex.test(email)) {
            setError('Invalid email format.');
            return false;
        }
        if (password.length < 8) {
            setError('Password must be at least 8 characters long.');
            return false;
        }
        setError(null);
        return true;
    };

    const handleLogin = async () => {
        if (!validateForm()) {
            return;
        }

        const userData = {
            email,
            password,
        };


        try {
            const response = await axios.post(`http://localhost:${PORT}/api/users/login`, userData, {
              withCredentials: true
            });
            if (response.data && response.data.error) {
              setError(response.data.error);
            } else {
              navigate('/wallet'); // ToDo: navigate to dashboard after implemented
            }
        } catch (error: any) {
            console.error('Error during login:', error);
            if (error.response && error.response.data && error.response.data.message) {
                setError(error.response.data.message);
            } else {
                setError('Failed to login. Please try again later.');
            }
        }
    };

    const handleSignup = () => {
        navigate('/signup');
    };

    return (
        <>
            <Header></Header>
            <Container className={classes.container}>
                {/* Icon */}
                <LockOutlinedIcon className={classes.icon} />

                {/* Form Title */}
                <Typography variant="h4" gutterBottom>
                    Log In
                </Typography>

                {/* Error Message */}
                {error && (
                    <Typography color="error" style={{ marginBottom: '1rem' }}>
                        {error}
                    </Typography>
                )}

                {/* Form */}
                <Box component="form" className={classes.form}>
                    <TextField
                        label="Email Address"
                        type="email"
                        variant="outlined"
                        className={classes.inputField}
                        required
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                    />
                    <TextField
                        label="Password"
                        type="password"
                        variant="outlined"
                        className={classes.inputField}
                        required
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <Button
                        variant="contained"
                        className={classes.button}
                        onClick={handleLogin}
                    >
                        Log In
                    </Button>
                </Box>

                {/* Don't have an account link */}
                <Typography>
                    <Link
                        onClick={handleSignup}
                        className={classes.link}
                    >
                        Don't have an account? Sign Up
                    </Link>
                    <div>tempUser@example.com</div>
                    <div>password123</div>
                </Typography>
            </Container>
        </>
    );
};

export default LoginPage;
