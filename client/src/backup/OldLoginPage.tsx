// import React from 'react';
// import { TextField, Button, Container, Typography, Link, Box } from '@mui/material';
// import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
// import { useNavigate } from 'react-router-dom';
// import useRegisterPageStyles from '../Stylesheets/RegisterPageStyles';
// import Header from './Header';

// const LoginPage: React.FC = () => {
//   const classes = useRegisterPageStyles()
//   const navigate = useNavigate();

//   const handleSignup = () => {
//     // Registration logic (placeholder)
//     navigate('/register');
//   };

//   const handleLogin = () => {
//     // Registration logic (placeholder)
//     navigate('/dashboard');
//   }

//   return (
//     <>
//     <Header></Header>
//     <Container className={classes.container}>
//       {/* Icon */}
//       <LockOutlinedIcon className={classes.icon} />

//       {/* Form Title */}
//       <Typography variant="h4" gutterBottom>
//         Login
//       </Typography>

//       {/* Form */}
//       <Box component="form" className={classes.form}>
//         <TextField
//           label="Username/Email"
//           variant="outlined"
//           className={classes.inputField}
//           required
//         />
//         <TextField
//           label="Password"
//           type="password"
//           variant="outlined"
//           className={classes.inputField}
//           required
//         />
//         <Button
//           variant="contained"
//           className={classes.button}
//           onClick={handleLogin}
//         >
//           Login
//         </Button>
//       </Box>

//       {/* Don't have an account have an account? */}
//       <Typography>
//         <Link
//           onClick={handleSignup}
//           className={classes.link}
//         >
//           Don't have an account? Sign up
//         </Link>
//       </Typography>
//     </Container>
//     </>
//   );
// };

// export default LoginPage;
