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
  const [walletAddress, setWalletAddress] = useState("");
  const [privateKey, setPrivateKey] = useState(""); // To hold the private key
  const [challenge, setChallenge] = useState<string | null>(null); // Challenge from the server
  const [isChallengeReceived, setIsChallengeReceived] = useState(false); // Track if challenge received
  const classes = useRegisterPageStyles();
  const navigate = useNavigate();

  // Function to handle the "Log In" button click
  const handleLogin = async () => {
    try {
      // Step 1: Request Challenge from the backend
      const response = await fetch("http://localhost:8080/login/challenge", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ user_id: walletAddress.trim() }), // Send the wallet address (userID)
      });

      if (response.ok) {
        const data = await response.json();
        setChallenge(data.challenge); // Store the challenge received from the server
        setIsChallengeReceived(true); // Indicate challenge was received
      } else {
        console.error("Failed to get challenge");
      }
    } catch (error) {
      console.error("Error during login:", error);
    }
  };

  // Function to sign the challenge with the private key
  const signChallenge = async () => {
    if (!challenge || !privateKey) {
      console.error('Challenge or private key is missing');
      return;
    }
  
    try {
      // Sanitize the private key by removing newlines, spaces, and headers/footers
      const sanitizedPrivateKey = privateKey
        .replace(/-----BEGIN RSA PRIVATE KEY-----/g, '')
        .replace(/-----END RSA PRIVATE KEY-----/g, '')
        .replace(/\s+/g, ''); // Remove any spaces or newlines
  
      console.log("Sanitized Private Key:", sanitizedPrivateKey);
  
      // Decode base64 private key
      const privateKeyBytes = Uint8Array.from(atob(sanitizedPrivateKey), c => c.charCodeAt(0));
  
      // Import the private key using SubtleCrypto API
      const privateKeyObj = await window.crypto.subtle.importKey(
        'pkcs8',
        privateKeyBytes.buffer, // Ensure we pass the ArrayBuffer here, not the Uint8Array itself
        { name: 'RSA-PSS', hash: 'SHA-256' },
        false,
        ['sign']
      );
  
      console.log("Imported Private Key Object:", privateKeyObj);
  
      // Ensure the challenge is properly encoded
      const enc = new TextEncoder();
      const challengeArrayBuffer = enc.encode(challenge);
  
      console.log("Challenge Array Buffer:", challengeArrayBuffer);
  
      // Sign the challenge
      const signatureArrayBuffer = await window.crypto.subtle.sign(
        { name: 'RSA-PSS', saltLength: 32 },
        privateKeyObj,
        challengeArrayBuffer
      );
  
      // Convert the signature to base64
      const signature = btoa(String.fromCharCode(...new Uint8Array(signatureArrayBuffer)));
  
      console.log("Generated Signature:", signature);
  
      // Step 3: Send the signature to the backend for verification
      const response = await fetch('http://localhost:8080/login/verify', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_id: walletAddress,
          challenge: challenge,
          signature: signature,  // Send the signature to the backend
        }),
      });
  
      if (response.ok) {
        const data = await response.json();
        console.log('Login successful:', data);
        // Navigate to wallet page upon successful login
        navigate('/wallet');  // Adjust the path as per your routing
      } else {
        const error = await response.text();
        console.error('Login failed:', error);
      }
    } catch (error) {
      console.error('Error signing challenge:', error);
    }
  };

  const handleSignUpRedirect = () => {
    navigate("/signup");
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
        <Box component="form" className={classes.form} sx={{ width: "70%" }}>
          <TextField
            label="Wallet Address"
            variant="outlined"
            fullWidth
            value={walletAddress}
            onChange={(e) => setWalletAddress(e.target.value)}
            className={classes.inputField}
            required
            disabled={isChallengeReceived} // Disable input after challenge is received
          />

          {!isChallengeReceived && (
            <Button
              variant="contained"
              fullWidth
              className={classes.button}
              onClick={handleLogin} // Request the challenge
            >
              Log In
            </Button>
          )}

          {isChallengeReceived && (
            <>
              <Typography variant="body2" gutterBottom>
                Challenge received. Please provide your private key to sign it.
              </Typography>

              <TextField
                label="Private Key"
                variant="outlined"
                fullWidth
                value={privateKey}
                onChange={(e) => setPrivateKey(e.target.value)}
                className={classes.inputField}
                required
              />

              <Button
                variant="contained"
                fullWidth
                className={classes.button}
                onClick={signChallenge} // Sign the challenge
              >
                Submit Signature
              </Button>
            </>
          )}
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
