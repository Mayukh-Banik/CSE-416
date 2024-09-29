import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles'; // Import ThemeProvider
import GeneralTheme from './Stylesheets/GeneralTheme';
import WelcomePage from './Components/WelcomePage';
import RegisterPage from './Components/RegisterPage';
import LoginPage from './Components/LoginPage';
import LoginPage2 from './Components/LoginPage2';
import SignupPage from './Components/SignUpPage';

const App: React.FC = () => {
  return (
    <ThemeProvider theme={GeneralTheme}>  {/* Wrap your app in ThemeProvider */}
      <Router>
        <Routes>
          <Route path="/" element={<WelcomePage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path='/login' element={<LoginPage />} />
          <Route path='/login2' element={<LoginPage2 />} />
          <Route path='/signup' element={<SignupPage />} />
          {/* <Route path='/register' element={<RegisterPage />} /> */}
        </Routes>
      </Router>
    </ThemeProvider>
  );
};

export default App;

  

