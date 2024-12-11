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

const FileTransactions : React.FC = () => {
  const theme = useTheme();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [activeTab, setActiveTab] = useState<number>(2); // Default to Pending Requests tab
  const navigate = useNavigate();

  
  useEffect(() => {
    const fetchTransactions = async () => {
      try {
        const response = await fetch("http://localhost:8081/files/getTransactions")
        if (!response.ok) {
          throw new Error(`Error fetching all transactions: ${response.statusText}`);
        }
        const data = await response.json();
        if (Array.isArray(data)) {
          console.log("all transactions: ", data)
          setTransactions(data); // Directly set the array of transactions
        } else {
          console.error("Unexpected response format:", data);
        }
      } catch (error) {
        console.error("failed to fetch all transactions: ", error)
      }
    } 
    fetchTransactions();
  }, [])


  const renderTable = (data: Transaction[]) => (
    <Table component={Paper}>
      <TableHead>
        <TableRow>
          <TableCell>Date</TableCell>
          {/* <TableCell>Time</TableCell> */}
          <TableCell>Transaction ID</TableCell>
          <TableCell>File Name</TableCell>
          <TableCell>File Hash</TableCell>
          <TableCell>Sender</TableCell>
          <TableCell>Receiver</TableCell>
          <TableCell>Total Fee</TableCell>
          <TableCell>Status</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {data.map((transaction, index) => (
          <TableRow key={index}>
            {/* <TableCell>{transaction.date}</TableCell> */}
            {/* <TableCell>{transaction.time}</TableCell> */}
            <TableCell>{transaction.CreatedAt}</TableCell>
            <TableCell>{transaction.TransactionID}</TableCell>
            <TableCell>{transaction.FileName}</TableCell>
            <TableCell>
              <Button onClick={() => navigate(`/fileview/${transaction.FileHash}`)}>{transaction.FileHash}</Button>
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
        Transaction History
      </Typography>

      {/* <Tabs value={activeTab} onChange={(_, value) => setActiveTab(value)}>
        <Tab label="Global Transactions" />
        <Tab label="My Transactions" />
        <Tab label="Pending Requests" />
      </Tabs> */}

      {/* {activeTab === 0 && renderTable(transactionData)}  Global Transactions */}
      {/* {activeTab === 1 && renderTable(transactionData.filter(tx => tx.sender === 'john_doe' || tx.receiver === 'john_doe'))} My Transactions */}
      {/* {activeTab === 2 && renderTable(transactions)}   */}
      {renderTable(transactions)}
    </Box>
  );
};

export default FileTransactions;