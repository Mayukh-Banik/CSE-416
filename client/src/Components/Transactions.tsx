import React, { useState, useEffect } from 'react';
import {
  Box, Table, TableBody, TableCell, TableHead, TableRow,
  TableSortLabel, Typography, Button, Paper, Tabs, Tab,
} from '@mui/material';
import Sidebar from './Sidebar';
import { useTheme } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Transaction } from '../models/transactions';

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

// interface Transaction {
//   date: string;
//   time: string;
//   fileName: string;
//   sender: string;
//   senderAddress: string;
//   receiver: string;
//   receiverAddress: string;
//   amount: string;
//   status: string;
// }

const GlobalTransactions : React.FC = () => {
  const theme = useTheme();
  const [pendingRequests, setPendingRequests] = useState<Transaction[]>([]);
  const [activeTab, setActiveTab] = useState<number>(2); // Default to Pending Requests tab
  const navigate = useNavigate();

  useEffect(() => {
    // Fetch pending requests from the backend
    const fetchPendingRequests = async () => {
      try {
        const response = await fetch('http://localhost:8081/download/getRequests');
        if (!response.ok) {
          throw new Error(`Error fetching pending requests: ${response.statusText}`);
        }
  
        // Parse JSON response
        const data = await response.json();
        if (Array.isArray(data)) {
          console.log("pending requests: ", data)
          setPendingRequests(data); // Directly set the array of transactions
        } else {
          console.error("Unexpected response format:", data);
        }
      } catch (error) {
        console.error('Failed to fetch pending requests:', error);
      }
    };
  
    fetchPendingRequests();
  }, []);
  

  const renderTable = (data: Transaction[]) => (
    <Table component={Paper}>
      <TableHead>
        <TableRow>
          <TableCell>Date</TableCell>
          {/* <TableCell>Time</TableCell> */}
          <TableCell>File Name</TableCell>
          <TableCell>Sender</TableCell>
          <TableCell>Receiver</TableCell>
          <TableCell>Transaction Amount</TableCell>
          <TableCell>Status</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {data.map((transaction, index) => (
          <TableRow key={index}>
            {/* <TableCell>{transaction.date}</TableCell> */}
            {/* <TableCell>{transaction.time}</TableCell> */}
            <TableCell>{transaction.CreatedAt}</TableCell>
            <TableCell>
              <Button onClick={() => navigate(`/fileview/${transaction.FileName}`)}>{transaction.FileName}</Button>
            </TableCell>
            <TableCell>
              <Button onClick={() => navigate(`/account/${transaction.RequesterID}`)}>{transaction.RequesterID}</Button>
            </TableCell>
            <TableCell>
              <Button onClick={() => navigate(`/account/${transaction.TargetID}`)}>{transaction.TargetID}</Button>
            </TableCell>
            <TableCell>{transaction.Fee}</TableCell>
            <TableCell>{transaction.Status}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );

  return (
    <Box
      sx={{
        padding: 2,
        marginTop: '70px',
        marginLeft: `${drawerWidth}px`,
        [theme.breakpoints.down('sm')]: {
          marginLeft: `${collapsedDrawerWidth}px`,
        },
      }}
    >
      <Sidebar />
      <Typography variant="h4" component="h1" gutterBottom>
        Transactions
      </Typography>

      <Tabs value={activeTab} onChange={(_, value) => setActiveTab(value)}>
        <Tab label="Global Transactions" />
        <Tab label="My Transactions" />
        <Tab label="Pending Requests" />
      </Tabs>

      {/* {activeTab === 0 && renderTable(transactionData)}  Global Transactions */}
      {/* {activeTab === 1 && renderTable(transactionData.filter(tx => tx.sender === 'john_doe' || tx.receiver === 'john_doe'))} My Transactions */}
      
      {activeTab === 2 && renderTable(pendingRequests)}  {/* Pending Requests */}
    </Box>
  );
};

export default GlobalTransactions;
