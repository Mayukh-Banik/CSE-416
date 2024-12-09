// SignUpPage.tsx
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
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import useRegisterPageStyles from "../Stylesheets/RegisterPageStyles";
import Header from "./Header";

const SignUpPage: React.FC = () => {
  const classes = useRegisterPageStyles();
  const [walletAddress, setWalletAddress] = useState<string | null>(null);
  const [privateKey, setPrivateKey] = useState<string | null>(null);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [isDialogOpen, setIsDialogOpen] = useState(false); // Dialog 상태
  const [passphrase, setPassphrase] = useState<string>(""); // 패스프레이즈 상태
  const [error, setError] = useState<string | null>(null); // 에러 상태

  const navigate = useNavigate();

  const handleLogin = () => {
    navigate("/login");
  };

  // "Generate Wallet" 버튼 클릭 시 다이얼로그 열기
  const handleGenerateWalletClick = () => {
    setIsDialogOpen(true);
  };

  // 패스프레이즈 제출 핸들러
  const handlePassphraseSubmit = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/auth/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ passphrase }), // 사용자가 입력한 패스프레이즈 전송
      });

      console.log("서버 응답 상태:", response.status); // 응답 상태 로그

      if (response.ok) {
        const data = await response.json();
        if (data.message === "Wallet already exists.") {
          alert("이미 지갑이 존재합니다.");
          setIsDialogOpen(false);
          // 필요 시 추가적인 로그인 로직을 여기에 추가할 수 있습니다.
        } else if (data.message === "Wallet successfully created.") {
          setWalletAddress(data.address);
          setPrivateKey(data.private_key);
          setIsSubmitted(true);
          setIsDialogOpen(false);
          setPassphrase("");
        }
      } else {
        const errorData = await response.json();
        setError(errorData.message || "Failed to signup");
        alert("Failed to signup: " + (errorData.message || "Unknown error"));
      }
    } catch (error) {
      // console.error("Error during signup:", error);
      setError("Error during signup");
      alert("Error during signup: " + error);
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
      document.body.appendChild(element); // FireFox를 위해 필요
      element.click();
    }
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false); // 다이얼로그 닫기
    setPassphrase("");
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
            onClick={handleGenerateWalletClick} // 수정된 핸들러 사용
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
              // onClick={handleContinueToAccount}
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

        {/* 패스프레이즈 입력 다이얼로그 */}
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
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseDialog}>Cancel</Button>
            <Button onClick={handlePassphraseSubmit} disabled={!passphrase}>
              Create Wallet
            </Button>
          </DialogActions>
        </Dialog>
      </Container>
    </>
  );
};

export default SignUpPage;
