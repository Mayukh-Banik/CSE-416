import React, { useState } from "react";
import {
  Button,
  Typography,
  Box,
  Container,
  IconButton,
  Tooltip,
  Paper,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  TextField,
  CircularProgress,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import useRegisterPageStyles from "../Stylesheets/RegisterPageStyles";
import Header from "./Header";

const SignUpPage: React.FC = () => {
  const classes = useRegisterPageStyles();
  const [walletAddress, setWalletAddress] = useState<string | null>(null);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [isDialogOpen, setIsDialogOpen] = useState(false); 
  const [passphrase, setPassphrase] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const navigate = useNavigate();

  const handleLogin = () => {
    navigate("/login");
  };

  // Button click handler for "Generate Wallet" button
  const handleGenerateWalletClick = () => {
    setIsDialogOpen(true);
  };

  // Handler for "Create Wallet" button click
  const handlePassphraseSubmit = async () => {
    if (!passphrase) {
      setError("Passphrase is required.");
      return;
    }

    setIsLoading(true); // Show loading spinner

    const controller = new AbortController();
    const timeoutId = setTimeout(() => {
      controller.abort();
    }, 30000); // 30 seconds timeout

    try {
      const response = await fetch("http://localhost:8080/api/auth/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ passphrase }), // Send passphrase to server
        signal: controller.signal, // Attach the signal to the fetch request
      });

      clearTimeout(timeoutId); // Clear the timeout since the request completed

      console.log("Server response status:", response.status); // Log response status

      if (response.ok) {
        const data = await response.json();
        console.log("Server response data:", data); // Log response data

        const { address, message } = data;

        if (message === "Wallet successfully created.") {
          setWalletAddress(address);
          setIsSubmitted(true);
          setIsDialogOpen(false);
          setPassphrase("");
          setError(null);
        }
      } else {
        const errorData = await response.json();
        setError(errorData.message || "Failed to signup");
        alert("Failed to signup: " + (errorData.message || "Unknown error"));
      }
    } catch (error: any) {
      if (error.name === "AbortError") {
        setError("Request timed out. Please try again.");
        alert("Error: Request timed out. Please try again.");
      } else {
        setError("Error during signup");
        alert("Error during signup: " + error.message);
      }
    } finally {
      clearTimeout(timeoutId); // Ensure the timeout is cleared
      setIsLoading(false); // Hide loading spinner
    }
  };

  const copyToClipboard = (text: string | null) => {
    if (text) {
      navigator.clipboard.writeText(text);
      alert("Copied to clipboard!");
    }
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false); // Close dialog
    setPassphrase("");
    setError(null);
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
          minHeight: "75vh",
          textAlign: "center",
          marginTop: "2rem",
        }}
      >
        <Typography variant="h3" sx={{ fontWeight: 600, mb: 2, marginTop: 10 }}>
          {isSubmitted
            ? "Wallet Successfully Generated"
            : "Welcome to Project Squid"}
        </Typography>
        {!isSubmitted && (
          <Typography
            variant="body1"
            sx={{ mb: 4, fontSize: "1.2rem", color: "#888" }}
          >
            Your go-to solution for secure file sharing.
          </Typography>
        )}

        {!isSubmitted && (
          <>
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
              onClick={handleGenerateWalletClick}
              disabled={isLoading}
            >
              Generate Wallet
            </Button>

            <Button
              variant="outlined"
              color="secondary"
              sx={{
                mb: 2,
                width: "100%",
                padding: "15px 0",
                fontSize: "1.2rem",
                borderRadius: "8px",
                boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
              }}
              onClick={handleLogin}
              disabled={isLoading}
            >
              Login
            </Button>
          </>
        )}

        {isSubmitted && (
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
                  Your wallet address:
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

              <Typography
                variant="body2"
                sx={{
                  marginBottom: "1rem",
                  fontWeight: "bold",
                  color: "#ffa726",
                }}
              >
                Important: Keep your wallet address secure.
              </Typography>

              <Button
                variant="contained"
                color="success"
                sx={{
                  width: "100%",
                  padding: "12px 0",
                  fontSize: "1.2rem",
                  borderRadius: "8px",
                  boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
                  mt: 3,
                }}
                onClick={handleLogin}
              >
                Continue to Login
              </Button>
            </Paper>
          </Box>
        )}

        <Dialog open={isDialogOpen} onClose={handleCloseDialog}>
          <DialogTitle>Enter Passphrase</DialogTitle>
          <DialogContent>
            <DialogContentText>
              Please enter a passphrase to create your wallet. This passphrase is
              essential for recovering or accessing your wallet, so do not forget
              it.
            </DialogContentText>
            <TextField
              autoFocus
              margin="dense"
              label="Passphrase"
              type="password"
              fullWidth
              variant="standard"
              value={passphrase}
              onChange={(e) => setPassphrase(e.target.value)}
              error={!!error}
              helperText={error}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseDialog} disabled={isLoading}>
              Cancel
            </Button>
            <Button onClick={handlePassphraseSubmit} disabled={!passphrase || isLoading}>
              Create Wallet
            </Button>
          </DialogActions>
        </Dialog>

        {isLoading && (
          <div
            style={{
              position: "fixed",
              top: 0,
              left: 0,
              width: "100%",
              height: "100%",
              backgroundColor: "rgba(0, 0, 0, 0.5)",
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              zIndex: 9999,
            }}
          >
            <CircularProgress color="inherit" />
          </div>
        )}
      </Container>
    </>
  );
};

export default SignUpPage;
