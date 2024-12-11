import React from "react";
import Sidebar from "./Sidebar";
import { Box, Container, IconButton, Typography, Button } from '@mui/material';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';
import DeleteIcon from '@mui/icons-material/Delete';
import { useTheme } from '@mui/material/styles';
import { useNavigate } from "react-router-dom";

type SettingPageProps = {
  darkMode: boolean;
  toggleTheme: () => void;
};

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const SettingPage: React.FC<SettingPageProps> = ({ darkMode, toggleTheme }) => {
  const theme = useTheme();
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/auth/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        const errorMessage = await response.text();
        console.error("Logout failed:", errorMessage);
        alert("Logout failed. Please try again.");
        return;
      }

      const data = await response.json();
      console.log("Logout response:", data.message);
      alert("Logout successful!");
      navigate("/signup"); // Redirect to login page after logout
    } catch (err) {
      console.error("Error during logout:", err);
      alert("An error occurred. Please try again.");
    }
  };

  // Delete Account handler
  const handleDeleteAccount = async () => {
    const confirmDelete = window.confirm(
      "Are you sure you want to delete your account? This action cannot be undone."
    );

    if (!confirmDelete) return;

    try {
      const response = await fetch("http://localhost:8080/api/auth/delete", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        const errorMessage = await response.text();
        console.error("Delete account failed:", errorMessage);
        alert("Account deletion failed. Please try again.");
        return;
      }

      const data = await response.json();
      console.log("Delete account response:", data.message);
      alert("Account deleted successfully!");
      navigate("/signup"); // Redirect to signup page after account deletion
    } catch (err) {
      console.error("Error during account deletion:", err);
      alert("An error occurred. Please try again.");
    }
  };

  return (
    <Box
      sx={{
        padding: 2,
        marginTop: '70px',
        marginLeft: `${drawerWidth}px`, // Default expanded margin
        transition: 'margin-left 0.3s ease', // Smooth transition
        [theme.breakpoints.down('sm')]: {
          marginLeft: `${collapsedDrawerWidth}px`, // Adjust left margin for small screens
        },
      }}
    >
      <Sidebar />
      <Container>
        {/* Top-right Theme Toggle */}
        <Box sx={{ display: 'flex', justifyContent: 'flex-end', mb: 3 }}>
          <IconButton onClick={toggleTheme}>
            {darkMode ? <Brightness7Icon /> : <Brightness4Icon />}
          </IconButton>
        </Box>

        {/* Page Title */}
        <Typography variant="h4" sx={{ marginBottom: '30px', textAlign: 'center' }}>
          Settings
        </Typography>

        {/* Logout Button */}
        <Box sx={{ marginTop: '20px', textAlign: 'center' }}>
          <Button
            variant="contained"
            color="primary"
            onClick={handleLogout}
            sx={{ width: '50%' }}
          >
            Logout
          </Button>
        </Box>

        {/* Delete Account Button */}
        <Box sx={{ marginTop: '20px', textAlign: 'center' }}>
          <Button
            variant="outlined"
            color="error"
            startIcon={<DeleteIcon />}
            onClick={handleDeleteAccount}
            sx={{ width: '50%' }}
          >
            Delete Account
          </Button>
        </Box>
      </Container>
    </Box>
  );
};

export default SettingPage;
