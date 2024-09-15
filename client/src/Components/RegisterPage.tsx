import React from 'react';
import { TextField, Button, Container, Typography, Link, Box } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import { useNavigate } from 'react-router-dom';
import useRegisterPageStyles from '../Stylesheets/RegisterPageStyles';
import Header from './Header';

const RegisterPage: React.FC = () => {
  const classes = useRegisterPageStyles()
  const navigate = useNavigate();

  const handleRegister = () => {
    // Registration logic (placeholder)
    navigate('/dashboard');
  };
  const handleLogin = () => {
    // Registration logic (placeholder)
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

      {/* Form */}
      <Box component="form" className={classes.form}>
        <TextField
          label="Username"
          variant="outlined"
          className={classes.inputField}
          required
        />
        <TextField
          label="Email Address"
          type="email"
          variant="outlined"
          className={classes.inputField}
          required
        />
        <TextField
          label="Password"
          type="password"
          variant="outlined"
          className={classes.inputField}
          required
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

export default RegisterPage;
