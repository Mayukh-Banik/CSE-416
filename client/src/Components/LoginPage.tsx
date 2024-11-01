import React, { useState, useEffect } from 'react';
import { TextField, Button, Container, Typography, Link, Box } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import useRegisterPageStyles from '../Stylesheets/RegisterPageStyles';
import Header from './Header';
import { useNavigate } from 'react-router-dom';
import forge from 'node-forge';

const PORT = 8080;

const LoginPage: React.FC = () => {
    const classes = useRegisterPageStyles();
    const navigate = useNavigate();

    const [publicKey, setPublicKey] = useState('');
    const [privateKey, setPrivateKey] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [challenge, setChallenge] = useState<string | null>(null);

    // Automatically request challenge on page load
    useEffect(() => {
        const requestChallenge = async () => {
            try {
                const response = await fetch(`http://localhost:${PORT}/api/auth/request-challenge`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                if (!response.ok) {
                    const errorMessage = await response.text();
                    console.error("Server Error:", errorMessage);
                    throw new Error('Failed to request challenge');
                }

                const data = await response.json();
                setChallenge(data.challenge);
                console.log("Challenge received:", data.challenge);
            } catch (err) {
                console.error("Error requesting challenge:", err);
                setError('Failed to request challenge.');
            }
        };

        requestChallenge();
    }, []);

    // Signature and authentication function
    const handleLogin = async () => {
        if (!publicKey || !privateKey || !challenge) {
            setError('Public key, private key, and challenge are required.');
            return;
        }

        try {
            // Parse the private key
            const rsaPrivateKey = forge.pki.privateKeyFromPem(privateKey);
            console.log('Private key parsed successfully:', rsaPrivateKey);

            // Decode the challenge
            const challengeBytes = forge.util.decode64(challenge);
            console.log('Decoded challenge bytes:', challengeBytes);

            // Generate signature
            const md = forge.md.sha256.create();
            md.update(challengeBytes);
            const signature = rsaPrivateKey.sign(md);
            console.log('Signature generated:', signature);

            // Encode the signature in Base64
            const signatureBase64 = forge.util.encode64(signature);

            // Send signature and public key to the server
            const response = await fetch(`http://localhost:${PORT}/api/auth/verify-challenge`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ public_key: publicKey, signature: signatureBase64 }),
                credentials: 'include' // Send request with cookies
            });

            if (!response.ok) {
                const errorMessage = await response.text();
                console.error('Server Error:', errorMessage);
                throw new Error('Failed to verify signature');
            }
    
            // Process login success
            const data = await response.json();
            console.log('Verification Response:', data);

            // Attempt to access cookie
            // const tokenCookie = document.cookie
            //     .split('; ')
            //     .find(row => row.startsWith('token='));
            
            // if (tokenCookie) {
            //     console.log('Token cookie found:', tokenCookie);
            // } else {
            //     console.warn('HttpOnly token cookie is not accessible by JavaScript.');
            // }

        } catch (err) {
            console.error('Error during login:', err);
            setError('Failed to login.');
        }
    };

    const handleSignup = () => {
        navigate('/');
    };

    return (
        <>
            <Header />
            <Container maxWidth="sm" sx={{ mt: 4 }}>
                <LockOutlinedIcon sx={{ fontSize: 40, mb: 2 }} />
                <Typography variant="h4" component="h1" gutterBottom>
                    Login
                </Typography>
                {error && <Typography color="error">{error}</Typography>}
                <Box sx={{ mt: 2 }}>
                    <TextField
                        label="Public Key"
                        multiline
                        rows={6}
                        fullWidth
                        value={publicKey}
                        onChange={(e) => setPublicKey(e.target.value)}
                        sx={{ mb: 2 }}
                    />
                    <TextField
                        label="Private Key"
                        multiline
                        rows={6}
                        fullWidth
                        value={privateKey}
                        onChange={(e) => setPrivateKey(e.target.value)}
                    />
                    <Button variant="contained" fullWidth sx={{ mt: 2 }} onClick={handleLogin}>
                        Login
                    </Button>
                </Box>
                <Typography sx={{ mt: 2 }}>
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
