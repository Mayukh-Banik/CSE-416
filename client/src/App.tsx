import React, { useState } from "react";
import {
  Container,
  Typography,
  Button,
  ThemeProvider,
  createTheme,
  CssBaseline,
} from "@mui/material";
import { HashRouter as Router, Routes, Route, Outlet, Navigate } from 'react-router-dom';

import GeneralTheme from "./Stylesheets/GeneralTheme";
import WelcomePage from "./Components/WelcomePage";
import RegisterPage from "./Components/RegisterPage";
import LoginPage from "./Components/LoginPage";
import SignupPage from "./Components/SignUpPage";
import SettingPage from "./Components/SettingPage";
import FilesPage from "./Components/FilesPage";
import MiningPage from "./Components/MiningPage";
import MarketPage from "./Components/MarketPage";
import FileViewPage from "./Components/FileViewPage";
import AccountViewPage from "./Components/AccountViewPage";
import ProxyPage from "./Components/ProxyPage";
// import GlobalTransactions from "./Components/GlobalTransactions";
import GlobalTransactions from "./Components/Transactions";

import SearchPage from "./Components/SearchPage";
import { FileMetadata } from "./models/fileMetadata";

const isUserLoggedIn = true; // should add the actual login state logic here.

// Fake data for your wallet (temporary data)
const walletAddress = "0x1234567890abcdef";
const balance = 100;
const transactions= [
  {
    id: "tx001",
    sender: "0xsender001",
    receiver: "0xreceiver001",
    amount: 10,
    timestamp: "2023-10-01T10:00:00",
    status: "completed",
    file: "file001.pdf"
  },
  {
    id: "tx009",
    sender: "0xsender009",
    receiver: "0xreceiver009",
    amount: 90,
    timestamp: "2023-10-09T17:20:00",
    status: "completed",
    file: "file009.docx"
  },
  // Add other transactions...
];


interface PrivateRouteProps {
  isAuthenticated: boolean;
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ isAuthenticated }) => {
  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
};

const App: React.FC = () => {
  const [darkMode, setDarkMode] = useState(false);
  const [uploadedFiles, setUploadedFiles] = useState<FileMetadata[]>([])
  const [initialFetch, setInitialFetch] = useState(false);
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
          <Route path="/" element={<MarketPage />} />   
          <Route path="/proxy" element={<ProxyPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/login" element={<LoginPage />} />
          {/* <Route path="/signup" element={<SignupPage />} /> */}
          <Route path="/market" element={<MarketPage />} />
          <Route path="/files" element={<FilesPage 
              uploadedFiles={uploadedFiles} 
              setUploadedFiles={setUploadedFiles} 
              initialFetch={initialFetch} 
              setInitialFetch={setInitialFetch} 
            />} />
          <Route path="/global-transactions" element={<GlobalTransactions />} />
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
              path="/account"
            />
          </Route>
          <Route path="/fileview" element={<FileViewPage />} />
          <Route path="/account" element={<AccountViewPage />} />
          <Route path="/account/:address" element={<AccountViewPage />} />
          <Route path="/fileview/:fileId" element={<FileViewPage />} />
          <Route path="/search-page" element={<SearchPage />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
};

export default App;