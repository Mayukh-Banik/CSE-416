import React, { useState } from 'react';
import axios from 'axios';
import { TextField, Button, Container, Typography, Link, Box } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import { useNavigate } from 'react-router-dom';
import useRegisterPageStyles from '../Stylesheets/RegisterPageStyles';
import Header from './Header';

const PORT = 5000;

const SignupPage: React.FC = () => {
  const classes = useRegisterPageStyles();
  const navigate = useNavigate();

  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState<string | null>(null);

  const validateForm = () => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    if (!name || !email || !password || !confirmPassword) {
      setError('All fields are required.');
      return false;
    }
    if (!emailRegex.test(email)) {
      setError('Invalid email format.');
      return false;
    }
    if (password !== confirmPassword) {
      setError('Passwords do not match.');
      return false;
    }
    if (password.length < 8) {
      setError('Password must be at least 8 characters long.');
      return false;
    }
    setError(null);
    return true;
  };

  const handleRegister = async () => {
    if (!validateForm()) {
      return;
    }

    const userData = {
      name,
      email,
      password,
    };

    try {
      const response = await axios.post(`http://localhost:${PORT}/api/users/signup`, userData);
      if (response.data && response.data.error) {
        setError(response.data.error);
      } else {
        navigate('/dashboard');
      }
    } catch (error: any) {
      console.error('Error during registration:', error);
      if (error.response && error.response.data && error.response.data.message) {
        setError(error.response.data.message);
      } else {
        setError('Failed to register. Please try again later.');
      }
    }
  };

  const handleLogin = () => {
    navigate('/login');
  };

  return (
    <>
      <Header></Header>
      <Container className={classes.container}>
        {/* Icon */}
        <LockOutlinedIcon className={classes.icon} />

        {/* Form Title */}
        <Typography variant="h4" gutterBottom>
          Sign up
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
            label="Name"
            variant="outlined"
            className={classes.inputField}
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
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
          <TextField
            label="Confirm Password"
            type="password"
            variant="outlined"
            className={classes.inputField}
            required
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
          />
          <Button
            variant="contained"
            className={classes.button}
            onClick={handleRegister}
          >
            Sign Up
          </Button>
        </Box>

        {/* Already have an account link */}
        <Typography>
          <Link
            onClick={handleLogin}
            className={classes.link}
          >
            Already have an account? Log In
          </Link>
        </Typography>
      </Container>
    </>
  );
};

export default SignupPage;
