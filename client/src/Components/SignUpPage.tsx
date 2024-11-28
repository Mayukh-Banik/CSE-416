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
  const [isDialogOpen, setIsDialogOpen] = useState(false); // State for dialog open/close

  const navigate = useNavigate();

  const handleLogin = () => {
    navigate("/login");
  };

  // API call to backend to generate wallet
  // SignUpPage.tsx
  const handleSignup = async () => {
    console.log("Generate Wallet 버튼 클릭됨"); // 로그 추가
    try {
      const response = await fetch("http://localhost:8080/api/auth/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ passphrase: "CSE416" }), // 필요 시 추가 데이터 전송
      });

      console.log("서버 응답 상태:", response.status); // 응답 상태 로그

      if (response.ok) {
        const data = await response.json();
        console.log("서버 응답 데이터:", data); // 응답 데이터 로그
        setWalletAddress(data.user.address); // 수정된 부분
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

  const handleOpenDialog = () => {
    setIsDialogOpen(true); // Open the dialog
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false); // Close the dialog
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
            onClick={handleSignup}
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
            //   onClick={handleContinueToAccount}
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

export default SignUpPage;
