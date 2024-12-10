import React, { useState } from 'react';
import { TextField, Button, Container, Typography, Link, Box } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import useRegisterPageStyles from '../Stylesheets/RegisterPageStyles';
import Header from './Header';
import { useNavigate } from 'react-router-dom';

const PORT = 8080;

const LoginPage: React.FC = () => {
    const classes = useRegisterPageStyles();
    const navigate = useNavigate();

    const [walletAddress, setWalletAddress] = useState('');
    const [passphrase, setPassphrase] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [successMessage, setSuccessMessage] = useState<string | null>(null);

    const handleLogin = async () => {
        setError(null);
        setSuccessMessage(null);

        // Validate inputs
        if (!walletAddress || !passphrase) {
            setError('Wallet address and passphrase are required.');
            return;
        }

        try {
            // Send login request to the backend
            const response = await fetch(`http://localhost:${PORT}/api/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ wallet_address: walletAddress, passphrase }),
            });

            if (!response.ok) {
                const errorMessage = await response.text();
                console.error('Login failed:', errorMessage);
                throw new Error(errorMessage);
            }

            // Process successful login response
            const data = await response.json();
            console.log('Login successful:', data);
            setSuccessMessage(`Login success: ${data.message}`);

            // Redirect or handle post-login actions
        } catch (err) {
            console.error('Error during login:', err);
            setError('Login failed. Please check your credentials.');
        }
    };

    const handleSignup = () => {
        navigate('/signup'); // Redirect to the signup page
    };

    return (
        <>
            <Header />
            <Container maxWidth="sm" sx={{ mt: 4 }}>
                <Box display="flex" flexDirection="column" alignItems="center">
                    <LockOutlinedIcon sx={{ fontSize: 40, mb: 2 }} />
                    <Typography variant="h4" component="h1" gutterBottom>
                        Login
                    </Typography>
                </Box>
                {error && <Typography color="error" sx={{ mt: 2 }}>{error}</Typography>}
                {successMessage && <Typography color="success" sx={{ mt: 2 }}>{successMessage}</Typography>}
                <Box sx={{ mt: 2 }}>
                    <TextField
                        label="Wallet Address"
                        fullWidth
                        value={walletAddress}
                        onChange={(e) => setWalletAddress(e.target.value)}
                        sx={{ mb: 2 }}
                    />
                    <TextField
                        label="Passphrase"
                        type="password"
                        fullWidth
                        value={passphrase}
                        onChange={(e) => setPassphrase(e.target.value)}
                        sx={{ mb: 2 }}
                    />
                    <Button variant="contained" fullWidth sx={{ mt: 2 }} onClick={handleLogin}>
                        Login
                    </Button>
                </Box>
                <Typography sx={{ mt: 2, textAlign: 'center' }}>
                    Don't have a wallet?{' '}
                    <Link onClick={handleSignup} sx={{ cursor: 'pointer' }}>
                        Generate One
                    </Link>
                </Typography>
            </Container>
        </>
    );
};

export default LoginPage;
