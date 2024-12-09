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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  TextField,
  CircularProgress,
  Backdrop,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import useRegisterPageStyles from "../Stylesheets/RegisterPageStyles";
import Header from "./Header";

const SignUpPage: React.FC = () => {
  const classes = useRegisterPageStyles();
  const [walletAddress, setWalletAddress] = useState<string | null>(null);
  const [privateKey, setPrivateKey] = useState<string | null>(null);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [isDialogOpen, setIsDialogOpen] = useState(false); 
  const [passphrase, setPassphrase] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const navigate = useNavigate();

  const handleLogin = () => {
    navigate("/login");
  };

  const checkWalletExistence = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/wallet/check", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ passphrase }),
      });

      if (response.ok) {
        const data = await response.json();
        return data.exists;
      } else {
        console.error("Failed to check wallet existence.");
        return false;
      }
    } catch (error) {
      console.error("Error checking wallet existence:", error);
      return false;
    }
  };

  // button click handler for "Generate Wallet" button
  const handleGenerateWalletClick = () => {
    setIsDialogOpen(true);
  };

  // hadler for "Create Wallet" button click
  const handlePassphraseSubmit = async () => {
    if (!passphrase) {
      setError("Passphrase is required.");
      return;
    }

    setIsLoading(true); // show loading spinner

    try {
      const walletExists = await checkWalletExistence();

      if (walletExists) {
        alert("이미 지갑이 존재합니다. 로그인 페이지로 이동합니다.");
        setIsDialogOpen(false);
        navigate("/login"); // send user to login page
        return;
      }

      const response = await fetch("http://localhost:8080/api/auth/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ passphrase }), // send passphrase to server
      });

      console.log("서버 응답 상태:", response.status); // response status log

      if (response.ok) {
        const data = await response.json();
        console.log("서버 응답 데이터:", data); // response data log

        const { address, private_key, message } = data;

        if (message === "Wallet successfully created.") {
          setWalletAddress(address);
          setPrivateKey(private_key);
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
    } catch (error) {
      setError("Error during signup");
      alert("Error during signup: " + error);
    } finally {
      setIsLoading(false); // hide loading spinner
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
      document.body.removeChild(element); // Cleanup
    }
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false); // close dialog
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
          height: "75vh",
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
            <img src="/images/loading.gif" alt="Loading" />
          </div>
        )}

      </Container>
    </>
  );
};

export default SignUpPage;
