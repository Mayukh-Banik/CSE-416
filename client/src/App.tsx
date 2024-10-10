import React, { useState } from "react";
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
// Remove the BrowserRouter import and replace it with HashRouter





import GeneralTheme from "./Stylesheets/GeneralTheme";
import WelcomePage from "./Components/WelcomePage";
import RegisterPage from "./Components/RegisterPage";
import LoginPage from "./Components/LoginPage";
import LoginPage2 from "./Components/LoginPage2";
import SignupPage from "./Components/SignUpPage";
import SettingPage from "./Components/SettingPage";
import WalletPage from "./Components/WalletPage";
import FilesPage from "./Components/FilesPage";
import MiningPage from "./Components/MiningPage";

const isUserLoggedIn = true; // should add the actual login state logic here.

// Fake data for your wallet (temporary data)
const walletAddress = "0x1234567890abcdef";
const balance = 100;
const transactions = [
  {
    id: "tx001",
    sender: "0xsender001",
    receiver: "0xreceiver001",
    amount: 10,
    timestamp: "2023-10-01T10:00:00",
    status: "completed",
  },
  {
    id: "tx009",
    sender: "0xsender009",
    receiver: "0xreceiver009",
    amount: 90,
    timestamp: "2023-10-09T17:20:00",
    status: "completed",
  },
  {
    id: "tx002",
    sender: "0xsender002",
    receiver: "0xreceiver002",
    amount: 20,
    timestamp: "2023-10-02T12:00:00",
    status: "pending",
  },
  {
    id: "tx005",
    sender: "0xsender005",
    receiver: "0xreceiver005",
    amount: 50,
    timestamp: "2023-10-05T16:00:00",
    status: "completed",
  },
  {
    id: "tx011",
    sender: "0xsender011",
    receiver: "0xreceiver011",
    amount: 110,
    timestamp: "2023-10-11T14:50:00",
    status: "completed",
  },
  {
    id: "tx003",
    sender: "0xsender003",
    receiver: "0xreceiver003",
    amount: 30,
    timestamp: "2023-10-03T14:00:00",
    status: "completed",
  },
  {
    id: "tx004",
    sender: "0xsender004",
    receiver: "0xreceiver004",
    amount: 40,
    timestamp: "2023-10-04T09:00:00",
    status: "failed",
  },

  {
    id: "tx006",
    sender: "0xsender006",
    receiver: "0xreceiver006",
    amount: 60,
    timestamp: "2023-10-06T11:00:00",
    status: "pending",
  },
  {
    id: "tx007",
    sender: "0xsender007",
    receiver: "0xreceiver007",
    amount: 70,
    timestamp: "2023-10-07T13:30:00",
    status: "completed",
  },
  {
    id: "tx013",
    sender: "0xsender013",
    receiver: "0xreceiver013",
    amount: 130,
    timestamp: "2023-10-13T18:00:00",
    status: "failed",
  },
  {
    id: "tx008",
    sender: "0xsender008",
    receiver: "0xreceiver008",
    amount: 80,
    timestamp: "2023-10-08T15:45:00",
    status: "failed",
  },

  {
    id: "tx010",
    sender: "0xsender010",
    receiver: "0xreceiver010",
    amount: 100,
    timestamp: "2023-10-10T10:10:00",
    status: "pending",
  },

  {
    id: "tx012",
    sender: "0xsender012",
    receiver: "0xreceiver012",
    amount: 120,
    timestamp: "2023-10-12T09:40:00",
    status: "pending",
  },
];

const publicKey = "publicKeyExample";
const privateKey = "privateKeyExample";

interface PrivateRouteProps {
  isAuthenticated: boolean;
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ isAuthenticated }) => {
  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
};

const App: React.FC = () => {
  const [darkMode, setDarkMode] = useState(false);

  const lightTheme = createTheme({
    palette: {
      mode: "light",
      background: {
        default: "#f4f4f4", //white
      },
      primary:{ //blue background
        main:'#1876d2'
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

  return (
    <ThemeProvider theme={darkMode ? darkTheme : lightTheme}>
      <CssBaseline /> {/* This resets CSS to a consistent baseline */}
      <Router>
        <Routes>
          <Route path="/" element={<WelcomePage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/login2" element={<LoginPage2 />} />
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
          <Route element={<PrivateRoute isAuthenticated={isUserLoggedIn} />}>
            <Route
              path="/wallet"
              element={
                <WalletPage
                  walletAddress={walletAddress}
                  balance={balance}
                  transactions={transactions}
                />
              }
            />
          </Route>
        </Routes>
      </Router>
    </ThemeProvider>
  );
};

export default App;
