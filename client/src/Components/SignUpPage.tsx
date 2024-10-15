import React, { useState } from "react";
import {
  Button,
  Typography,
  Box,
  Container,
  Link,
  Paper,
  IconButton,
  Tooltip,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import useRegisterPageStyles from "../Stylesheets/RegisterPageStyles";
import Header from "./Header";

const SignUpPage: React.FC = () => {
  const classes = useRegisterPageStyles();
  const [walletAddress, setWalletAddress] = useState<string | null>(null);
  const [privateKey, setPrivateKey] = useState<string | null>(null);
  const [isSubmitted, setIsSubmitted] = useState(false);

  const navigate = useNavigate();

  const handleLogin = () => {
    navigate("/login");
  };

  // API call to backend to generate wallet
  const handleSignup = async () => {
    try {
      const response = await fetch("http://localhost:8080/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json();
        setWalletAddress(data.public_key);
        setPrivateKey(data.private_key);
        setIsSubmitted(true);
      } else {
        console.error("Failed to signup");
      }
    } catch (error) {
      console.error("Error during signup:", error);
    }
  };

  const copyToClipboard = (text: string | null) => {
    if (text) {
      navigator.clipboard.writeText(text);
      alert("Copied to clipboard!");
    }
  };

  const downloadPrivateKey = () => {
    if (privateKey) {
      const element = document.createElement("a");
      const file = new Blob([privateKey], { type: "text/plain" });
      element.href = URL.createObjectURL(file);
      element.download = "privateKey.txt";
      document.body.appendChild(element); // Required for this to work in FireFox
      element.click();
    }
  };

  return (
    <>
      <Header />
      <Container className={classes.container}>
        <Box sx={{ marginTop: "6rem", textAlign: "center" }}>
          <Typography variant="h4" gutterBottom>
            Sign Up
          </Typography>
          <LockOutlinedIcon
            className={classes.icon}
            sx={{ fontSize: "4rem", marginBottom: "1rem" }}
          />
        </Box>

        {!isSubmitted ? (
          <Box component="form" className={classes.form}>
            <Button
              variant="contained"
              className={classes.button}
              onClick={handleSignup} // Backend signup
            >
              Generate Wallet
            </Button>
          </Box>
        ) : (
          <Box sx={{ mt: 4, width: "70%", margin: "0 auto" }}>
            <Paper
              elevation={2}
              sx={{
                padding: 3,
                textAlign: "center",
                width: "100%",
              }}
            >
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  marginBottom: "1rem",
                }}
              >
                <Typography variant="body1" sx={{ marginRight: "1rem" }}>
                  <strong>Your wallet address (public key):</strong>
                </Typography>
                <Tooltip title="Copy to Clipboard">
                  <IconButton onClick={() => copyToClipboard(walletAddress)}>
                    <ContentCopyIcon />
                  </IconButton>
                </Tooltip>
              </Box>
              <Typography
                variant="body2"
                sx={{
                  whiteSpace: "nowrap",
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                  maxWidth: "100%",
                  marginBottom: "1rem",
                  fontFamily: "monospace",
                  wordBreak: "break-all",
                }}
              >
                {walletAddress}
              </Typography>

              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  marginBottom: "1rem",
                }}
              >
                <Typography variant="body1" sx={{ marginRight: "1rem" }}>
                  <strong>Your private key:</strong>
                </Typography>
                <Tooltip title="Copy to Clipboard">
                  <IconButton onClick={() => copyToClipboard(privateKey)}>
                    <ContentCopyIcon />
                  </IconButton>
                </Tooltip>
              </Box>
              <Typography
                variant="body2"
                sx={{
                  whiteSpace: "nowrap",
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                  maxWidth: "100%",
                  marginBottom: "1rem",
                  fontFamily: "monospace",
                  color: "red",
                  wordBreak: "break-all",
                }}
              >
                {privateKey}
              </Typography>

              <Typography
                variant="body2"
                sx={{
                  marginBottom: "1rem",
                  fontWeight: "bold",
                  color: "orange",
                }}
              >
                Important: Keep your private key secure and do not share it with
                anyone.
              </Typography>

              <Button
                variant="contained"
                className={classes.button}
                onClick={downloadPrivateKey}
              >
                Download Private Key
              </Button>
            </Paper>
          </Box>
        )}

        <Typography sx={{ marginTop: "2rem", textAlign: "center" }}>
          <Link onClick={handleLogin} className={classes.link}>
            Already have an account? Log In
          </Link>
        </Typography>
      </Container>
    </>
  );
};

export default SignUpPage;
