import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Container, Typography, Button } from '@mui/material';
import WelcomePage from './Components/WelcomePage';
import RegisterPage from './Components/RegisterPage';
import LoginPage from './Components/LoginPage';
import SettingPage from './Components/SettingPage';
import Dashboard from './Components/Dashboard';
import TransactionPage from './Components/TransactionPage';

const App: React.FC = () => {
    return (
      <Router>
        <Routes>
          <Route path="/" element={<WelcomePage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path='/login' element={<LoginPage />} />
          <Route path='/transaction' element={<TransactionPage />} />
          <Route path="/settings" element={<SettingPage />} /> 
          <Route path="/dashboard" element={<Dashboard />} />  
        </Routes>
      </Router>
    );
  };
  

export default App;