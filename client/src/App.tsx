import React, { useState } from "react";
import {
  ThemeProvider,
  createTheme,
  CssBaseline,
} from "@mui/material";
import { HashRouter as Router, Routes, Route, Outlet, Navigate } from 'react-router-dom';


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
import GlobalTransactions from "./Components/Transactions";

import SearchPage from "./Components/SearchPage";
import { FileMetadata } from "./models/fileMetadata";
import InitPage from "./Components/InitPage";

const isUserLoggedIn = true; 

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
      mode: 'dark',
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
      <CssBaseline />
      <Router>
        <Routes>
          <Route path="/" element={<InitPage />} />
          <Route path="/init" element={<InitPage />} />

          <Route path="/signup" element={<SignupPage />} />

          <Route path="/proxy" element={<ProxyPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/market" element={<MarketPage />} />
          <Route
            path="/files"
            element={
              <FilesPage
                uploadedFiles={uploadedFiles}
                setUploadedFiles={setUploadedFiles}
                initialFetch={initialFetch}
                setInitialFetch={setInitialFetch}
              />
            }
          />
          <Route path="/global-transactions" element={<GlobalTransactions />} />
          <Route
            path="/settings"
            element={<SettingPage darkMode={darkMode} toggleTheme={toggleTheme} />}
          />
          {/* <Route path='/transaction' element={<TransactionPage />} /> */}
          <Route path="/mining" element={<MiningPage />} />

          {/* Routes protected by PrivateRoute */}
          <Route element={<PrivateRoute isAuthenticated={isUserLoggedIn} />}>
            <Route path="/account" element={<AccountViewPage />} />
            <Route path="/account/:address" element={<AccountViewPage />} />
          </Route>
          
          <Route path="/fileview" element={<FileViewPage />} />
          <Route path="/fileview/:fileId" element={<FileViewPage />} />
          <Route path="/search-page" element={<SearchPage />} />

          {/* Fallback route */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
};

export default App;
