import React, { useState } from 'react';
import { TextField, Button, Container, Typography, Box, Link } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import useRegisterPageStyles from '../Stylesheets/RegisterPageStyles';
import Header from './Header';
import { useNavigate } from 'react-router-dom';

const LoginPage: React.FC = () => {
  const [walletAddress, setWalletAddress] = useState('');
  const classes = useRegisterPageStyles();
  const navigate = useNavigate();

  const handleLogin = () => {
    // Handle login logic here
    console.log('Logging in with wallet:', walletAddress);
  };

  const handleSignUpRedirect = () => {
    navigate('/signup');
  };

  return (
    <>
      <Header />
      <Container className={classes.container}>
        {/* Icon */}
        <LockOutlinedIcon className={classes.icon} />

        {/* Form Title */}
        <Typography variant="h4" gutterBottom>
          Log In
        </Typography>

        {/* Form */}
        <Box component="form" className={classes.form} sx={{ width: '70%' }}>
          <TextField
            label="Wallet Address"
            variant="outlined"
            fullWidth
            value={walletAddress}
            onChange={(e) => setWalletAddress(e.target.value)}
            className={classes.inputField}
            required
          />
          <Button
            variant="contained"
            fullWidth
            className={classes.button}
            onClick={handleLogin}
          >
            Log In
          </Button>
        </Box>

        {/* Highlighted text for Sign Up */}
        <Typography sx={{ mt: 2 }}>
          <Link onClick={handleSignUpRedirect} className={classes.link}>
            Don't have an account? Sign Up!
          </Link>
        </Typography>
      </Container>
    </>
  );
};

export default LoginPage;
