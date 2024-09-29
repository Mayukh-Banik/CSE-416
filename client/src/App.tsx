import React, {useState} from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material'; // Import CssBaseline here
import { Container, Typography, Button } from '@mui/material';
import WelcomePage from './Components/WelcomePage';
import RegisterPage from './Components/RegisterPage';
import LoginPage from './Components/LoginPage';
import SettingPage from './Components/SettingPage';
import Dashboard from './Components/Dashboard';
import TransactionPage from './Components/TransactionPage';

const App: React.FC = () => {
  const [darkMode, setDarkMode] = useState(false);

  const lightTheme = createTheme({
    palette: {
      mode: 'light',
    },
  });

  const darkTheme = createTheme({
    palette: {
      mode: 'dark',
      background: {
        default: '#121212',
      },
      primary: {
        main: '#f48fb1', // Main color for dark theme
      },
      secondary: {
        main: '#f48fb1', // Secondary color for dark theme
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
          <Route path="/" element={<WelcomePage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path='/login' element={<LoginPage />} />
<<<<<<< HEAD
          <Route path='/transaction' element={<TransactionPage />} />
          <Route path="/settings" element={<SettingPage />} /> 
=======
          {/* <Route path='/register' element={<RegisterPage />} /> */}
          <Route path="/settings" element={<SettingPage darkMode={darkMode} toggleTheme={toggleTheme} />} />
>>>>>>> 4cec157 (save)
          <Route path="/dashboard" element={<Dashboard />} />  
        </Routes>
      </Router>
      </ThemeProvider>
    );
  };
  

export default App;