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
        public_key: "gen-public-key-123",
        private_key: "gen-private-key-456",
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

  const handleContinueToAccount = () => {
    navigate("/account/1");
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
          height: "75vh",
          textAlign: "center",
          marginTop: "2rem",
        }}
      >
        {/* Change the title based on whether wallet is generated */}
        <Typography variant="h3" sx={{ fontWeight: 600, mb: 2, marginTop:10 }}>
          {isSubmitted ? "Wallet Successfully Generated" : "Welcome to Project Squid"}
        </Typography>
        {!isSubmitted && (
          <Typography
            variant="body1"
            sx={{ mb: 4, fontSize: "1.2rem", color: "#888" }}
          >
            Your go-to solution for secure file sharing.
          </Typography>
        )}

        {!isSubmitted ? (
          <Button
            variant="contained"
            color="primary"
            sx={{
              mb: 2,
              width: "100%",
              padding: "15px 0",
              fontSize: "1.2rem",
              borderRadius: "8px",
              boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
              "&:hover": {
                backgroundColor: "#1976d2",
              },
            }}
            onClick={handleGenerateWallet}
          >
            Generate Wallet
          </Button>
        ) : (
          <Box sx={{ mt: 4, width: "100%" }}>
            <Paper
              elevation={4}
              sx={{
                padding: 4,
                textAlign: "center",
                width: "100%",
                borderRadius: "12px",
                boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
              }}
            >
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  marginBottom: "1rem",
                  marginTop: "-2rem",
                }}
              >
                <Typography variant="body1" sx={{ fontWeight: 600 }}>
                  Your wallet address (public key):
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
                <Typography variant="body1" sx={{ fontWeight: 600 }}>
                  Your private key:
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
                sx={{
                  marginBottom: "1rem",
                  fontWeight: "bold",
                  color: "#ffa726",
                }}
              >
                Important: Keep your private key secure and do not share it with
                anyone.
              </Typography>

              <Button
                variant="contained"
                color="secondary"
                sx={{
                  width: "100%",
                  padding: "10px 0",
                  fontSize: "1rem",
                  borderRadius: "8px",
                  boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
                }}
                onClick={downloadPrivateKey}
              >
                Download Private Key
              </Button>
            </Paper>

            {/* Continue to Account Button */}
            <Button
              variant="contained"
              color="success"
              sx={{
                mt: 3,
                width: "100%",
                padding: "12px 0",
                fontSize: "1.2rem",
                borderRadius: "8px",
                boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
              }}
              onClick={handleContinueToAccount}
            >
              Continue to Account
            </Button>
          </Box>
        )}

        <Typography sx={{ marginTop: "2rem" }}>
          <Link
            onClick={handleLogin}
            sx={{
              cursor: "pointer",
              fontSize: "1rem",
              textDecoration: "underline",
            }}
          >
            Already have an account? Log In
          </Link>
        </Typography>
      </Container>
    </>
  );
};

export default WelcomePage;
