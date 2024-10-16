import React, { useState, useEffect } from "react";
// Replace the BrowserRouter import with HashRouter-02-01
// import {
//   BrowserRouter as Router,
//   Routes,
//   Route,
//   Outlet,
//   Navigate,
// } from "react-router-dom";
import {
  Container,
  Typography,
  Button,
  ThemeProvider,
  createTheme,
  CssBaseline,
} from "@mui/material";
import { HashRouter as Router, Routes, Route, Outlet, Navigate } from 'react-router-dom';
import axios from "axios";
// Remove the BrowserRouter import and replace it with HashRouter

import GeneralTheme from "./Stylesheets/GeneralTheme";
import WelcomePage from "./Components/WelcomePage";
import RegisterPage from "./Components/RegisterPage";
import LoginPage from "./Components/LoginPage";
import SignupPage from "./Components/SignUpPage";
import SettingPage from "./Components/SettingPage";
import WalletPage from "./Components/WalletPage";
import FilesPage from "./Components/FilesPage";
import MiningPage from "./Components/MiningPage";

const jwtDecode = require("jwt-decode");

interface PrivateRouteProps {
  isAuthenticated: boolean | null;
  children: React.ReactNode;
}

const getTokenFromCookies = (): string | null => {
  const match = document.cookie.match(new RegExp("(^| )token=([^;]+)"));
  return match ? match[2] : null;
};

const PrivateRoute: React.FC<PrivateRouteProps> = ({ isAuthenticated, children }) => {
  if (isAuthenticated === null) {
    return <div>Loading...</div>; // 로그인 여부 확인 중
  }
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
};


const App: React.FC = () => {
  const [darkMode, setDarkMode] = useState(false);
  const [isUserLoggedIn, setIsUserLoggedIn] = useState<boolean | null>(null);
  const [user, setUser] = useState<any | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const lightTheme = createTheme({
    palette: {
      mode: "light",
      background: {
        default: "#f4f4f4", //white
      },
      primary: { //blue background
        main: '#1876d2'
      },
      secondary: {
        main: "#121212", // text color
      },
    },
  });

  const darkTheme = createTheme({
    palette: {
      mode: "dark",
      background: {
        default: "#202d45", //darker gray
      },
      primary: {
        main: "#202d45", // lighter gray
      },
      secondary: {
        main: "#f4f4f4", // text color
      },
    },
  });

  const toggleTheme = () => {
    setDarkMode((prevMode) => !prevMode);
  };

  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        const response = await axios.get("http://localhost:5000/api/auth/status", {
          withCredentials: true,
        });
        console.log("Auth status response:", response); // 응답 전체 출력
        if (response.data.isAuthenticated) {
          console.log("User is authenticated:", response.data.user); // 인증된 사용자 정보 출력
          setUser(response.data.user);
          setIsUserLoggedIn(true);
        } else {
          console.log("User is not authenticated");
          setIsUserLoggedIn(false);
        }
      } catch (error) {
        console.error("Error checking auth status:", error);
        setIsUserLoggedIn(false);
      }
      setIsLoading(false);
    };

    checkAuthStatus();
  }, []);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <ThemeProvider theme={darkMode ? darkTheme : lightTheme}>
      <CssBaseline /> {/* This resets CSS to a consistent baseline */}
      <Router>
        <Routes>
          <Route path="/" element={<WelcomePage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/signup" element={<SignupPage />} />
          <Route path="/files" element={<FilesPage />} />
          <Route
            path="/settings"
            element={
              <SettingPage darkMode={darkMode} toggleTheme={toggleTheme} />
            }
          />
          {/* <Route path='/transaction' element={<TransactionPage />} /> */}
          <Route path="/mining" element={<MiningPage />} />

          {/* <Route path='/register' element={<RegisterPage />} /> */}

          {/* Routes protected by PrivateRoute */}
          <Route
            path="/wallet"
            element={
              <PrivateRoute isAuthenticated={isUserLoggedIn}>
                <WalletPage />
              </PrivateRoute>
            }
          />
        </Routes>
      </Router>
    </ThemeProvider>
  );
};

export default App;
