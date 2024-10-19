import React, { useState } from 'react';
import axios from 'axios';
import { TextField, Button, Container, Typography, Link, Box, IconButton, Tooltip, Paper } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import useRegisterPageStyles from '../Stylesheets/RegisterPageStyles';
import Header from './Header';
import { useNavigate } from 'react-router-dom';

const PORT = 5000;

const LoginPage: React.FC = () => {
    const classes = useRegisterPageStyles();
    const navigate = useNavigate();

    const [publicKey, setPublicKey] = useState('');
    const [error, setError] = useState<string | null>(null);

    const validateForm = () => {
        if (!publicKey) {
            setError('Public key is required.');
            return false;
        }
        setError(null);
        return true;
    };

    const handleLogin = async () => {
        if (validateForm()) {
            navigate("/account/1"); // Simply navigate to the account page if public key is not empty
        }
    };

    const handleSignup = () => {
        navigate('/');
    };

    return (
        <>
            <Header />
            <Container
              maxWidth="sm"
              sx={{
                display: "flex",
                flexDirection: "column",
                justifyContent: "center",
                alignItems: "center",
                height: "80vh", 
                textAlign: "center",
                marginTop: "2rem",
              }}
            >
                {/* Icon centered and blue */}
                <LockOutlinedIcon sx={{ color: "primary.main", fontSize: "3rem", marginBottom: "0.5rem" }} />

                {/* Log In text, no bold */}
                <Typography variant="h4" sx={{ fontWeight: "normal", mb: 2 }}>
                    Log In
                </Typography>

                {/* Error Message */}
                {error && (
                    <Typography color="error" sx={{ marginBottom: '1rem' }}>
                        {error}
                    </Typography>
                )}

                {/* Form */}
                <Box
                  component="form"
                  sx={{
                    display: "flex",
                    flexDirection: "column",
                    gap: 2,
                    width: "100%",
                    maxWidth: "450px", // Made the form wider
                    padding: "2rem",
                    borderRadius: "12px",
                    boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
                  }}
                >
                    <TextField
                        label="Public Key"
                        type="text"
                        variant="outlined"
                        fullWidth
                        required
                        value={publicKey}
                        onChange={(e) => setPublicKey(e.target.value)}
                    />
                    <Button
                        variant="contained"
                        color="primary"
                        sx={{
                          padding: "12px 0",
                          fontSize: "1.2rem",
                          borderRadius: "8px",
                          boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
                          "&:hover": {
                            backgroundColor: "#1976d2",
                          },
                        }}
                        onClick={handleLogin}
                    >
                        Authenticate with Public Key
                    </Button>
                </Box>

                {/* Don't have an account link */}
                <Typography sx={{ marginTop: "1rem" }}>
                    <Link
                        onClick={handleSignup}
                        sx={{ cursor: "pointer", textDecoration: "underline", fontSize: "1rem", color: "#1976d2" }}
                    >
                        Don't have a wallet? Generate One
                    </Link>
                </Typography>
            </Container>
        </>
    );
};

export default LoginPage;
