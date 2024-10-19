import React, { useState } from "react";
import {
  Button,
  Typography,
  Box,
  Container,
  Link,
  IconButton,
  Tooltip,
  Paper,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import useRegisterPageStyles from "../Stylesheets/RegisterPageStyles";
import Header from "./Header";

const WelcomePage: React.FC = () => {
  const classes = useRegisterPageStyles();
  const [walletAddress, setWalletAddress] = useState<string | null>(null);
  const [privateKey, setPrivateKey] = useState<string | null>(null);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const navigate = useNavigate();

  const handleLogin = () => {
    navigate("/login");
  };

  // Placeholder function for wallet generation
  const handleGenerateWallet = async () => {
    try {
      // Replace this with actual backend call when ready
      const fakeResponse = {
        public_key: "fake-public-key-123",
        private_key: "fake-private-key-456",
      };
      setWalletAddress(fakeResponse.public_key);
      setPrivateKey(fakeResponse.private_key);
      setIsSubmitted(true);
    } catch (error) {
      console.error("Error generating wallet:", error);
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
      document.body.appendChild(element); // Required for this to work in Firefox
      element.click();
    }
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
        }}
      >
        <Typography variant="h2" sx={{ fontWeight: 700, mb: 2 }}>
          Welcome to Project Squid
        </Typography>
        <Typography
          variant="body1"
          sx={{ mb: 4, fontSize: "1.2rem", color: "#666" }}
        >
          Your go-to solution for secure file sharing.
        </Typography>

        {!isSubmitted ? (
          <Button
            variant="contained"
            color="primary"
            sx={{ mb: 2, width: "100%", padding: "15px 0", fontSize: "1.2rem" }}
            onClick={handleGenerateWallet}
          >
            Generate Wallet
          </Button>
        ) : (
          <Box sx={{ mt: 4, width: "100%" }}>
            <Paper elevation={2} sx={{ padding: 3, textAlign: "center" }}>
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
                }}
              >
                {privateKey}
              </Typography>

              <Typography
                variant="body2"
                sx={{ marginBottom: "1rem", fontWeight: "bold", color: "orange" }}
              >
                Important: Keep your private key secure and do not share it with
                anyone.
              </Typography>

              <Button
                variant="contained"
                color="secondary"
                onClick={downloadPrivateKey}
              >
                Download Private Key
              </Button>
            </Paper>
          </Box>
        )}

        <Typography sx={{ marginTop: "2rem" }}>
          <Link onClick={handleLogin} sx={{ cursor: "pointer" }}>
            Already have an account? Log In
          </Link>
        </Typography>
      </Container>
    </>
  );
};

export default WelcomePage;
