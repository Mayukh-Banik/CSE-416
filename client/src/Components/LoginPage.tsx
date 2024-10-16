import React, { useState } from "react";
import {
  TextField,
  Button,
  Container,
  Typography,
  Box,
  Link,
} from "@mui/material";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import useRegisterPageStyles from "../Stylesheets/RegisterPageStyles";
import Header from "./Header";
import { useNavigate } from "react-router-dom";

const LoginPage: React.FC = () => {
  const [walletAddress, setWalletAddress] = useState(""); // Wallet ID (Public Key)
  const [walletError, setWalletError] = useState(""); // Error for wallet address
  const classes = useRegisterPageStyles();
  const navigate = useNavigate();

  // Function to handle the "Log In" button click
  const handleLogin = async () => {
    if (walletAddress.trim() === "") {
      setWalletError("Wallet address cannot be empty."); // Set error if the input is empty
      return;
    }

    // Clear any previous errors
    setWalletError("");

    try {
      // Send the walletId (public key) to the backend for login
      const response = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ walletId: walletAddress.trim() }), // Using walletAddress as walletId
      });

      if (response.ok) {
        const data = await response.json();
        console.log("Login successful:", data);
        // Navigate to the wallet page upon successful login
        navigate("/wallet"); // Adjust the path as per your routing
      } else {
        const error = await response.text();
        console.error("Login failed:", error);
        setWalletError("Invalid wallet address."); // Show error if login fails
      }
    } catch (error) {
      console.error("Error during login:", error);
      setWalletError("An error occurred. Please try again.");
    }
  };

  const handleSignUpRedirect = () => {
    navigate("/signup");
  };

  return (
    <>
      <Header />
      <Container className={classes.container}>
        <LockOutlinedIcon className={classes.icon} />
        <Typography variant="h4" gutterBottom>
          Log In
        </Typography>

        <Box component="form" className={classes.form} sx={{ width: "70%" }}>
          <TextField
            label="Wallet Address"
            variant="outlined"
            fullWidth
            value={walletAddress}
            onChange={(e) => setWalletAddress(e.target.value)}
            className={classes.inputField}
            required
            error={!!walletError} // Show error state if walletError exists
            helperText={walletError} // Display error message
          />

          <Button
            variant="contained"
            fullWidth
            className={classes.button}
            onClick={handleLogin} // Log in with walletId
          >
            Log In
          </Button>
        </Box>

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
