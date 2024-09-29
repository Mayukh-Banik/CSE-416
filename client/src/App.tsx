import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material'; 
import GeneralTheme from './Stylesheets/GeneralTheme';
import WelcomePage from './Components/WelcomePage';
import RegisterPage from './Components/RegisterPage';
import LoginPage from './Components/LoginPage';
import LoginPage2 from './Components/LoginPage2';
import SignupPage from './Components/SignUpPage';

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
    <ThemeProvider theme={GeneralTheme}>  {/* Wrap your app in ThemeProvider */}
      <Router>
        <Routes>
          <Route path="/" element={<WelcomePage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path='/login' element={<LoginPage />} />
          <Route path='/login2' element={<LoginPage2 />} />
          <Route path='/signup' element={<SignupPage />} />
          <Route path="/settings" element={<SettingPage darkMode={darkMode} toggleTheme={toggleTheme} />} />
          {/* <Route path='/register' element={<RegisterPage />} /> */}
        </Routes>
      </Router>
    </ThemeProvider>
  );
};

export default App;

  

