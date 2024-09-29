import React from "react";
import { Button, Typography, Box, Container } from "@mui/material";
import { useNavigate } from "react-router-dom";
import Header from "./Header";

const WelcomePage: React.FC = () => {
  const navigate = useNavigate();

  const handleLogin = () => {
    navigate("/login");
  };

  const handleSignup = () => {
    navigate("/register");
  };

  const handleGuest = () => {
    navigate("/dashboard");
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
          height: "80vh", // Vertically center the content
          textAlign: "center",
        }}
      >
        {/* Increased font size and weight for the Welcome message */}
        <Typography variant="h2" sx={{ fontWeight: 700, mb: 2 }}>
          Welcome to Project Squid
        </Typography>

        {/* Supporting text - slightly smaller */}
        <Typography
          variant="body1"
          sx={{ mb: 4, fontSize: "1.2rem", color: "#666" }}
        >
          Your go-to solution for secure file sharing.
        </Typography>

        {/* Buttons styled to match blue theme */}
        <Box sx={{ display: "flex", flexDirection: "column", width: "100%" }}>
          <Button
            onClick={handleLogin}
            variant="contained"
            color="primary"
            sx={{ mb: 2, width: "100%", padding: "15px 0", fontSize: "1.2rem" }}
          >
            Log In
          </Button>

          <Button
            onClick={handleSignup}
            variant="contained"
            color="primary"
            sx={{ mb: 2, width: "100%", padding: "15px 0", fontSize: "1.2rem" }}
          >
            Sign Up
          </Button>

          <Button
            onClick={handleGuest}
            variant="contained"
            color="primary"
            sx={{ width: "100%", padding: "15px 0", fontSize: "1.2rem" }}
          >
            Continue as Guest
          </Button>
        </Box>
      </Container>
    </>
  );
};

export default WelcomePage;
