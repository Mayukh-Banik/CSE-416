import React, { useState } from 'react';
import {
  Box, Table, TableBody, TableCell, TableHead, TableRow,
  TableSortLabel, Typography, Button
} from '@mui/material';
import Sidebar from './Sidebar'; // Adjust the path as necessary
import { styled } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';

const Main = styled('main')(({ theme }) => ({
  flexGrow: 1,
  padding: theme.spacing(3),
  transition: theme.transitions.create(['margin'], {
    easing: theme.transitions.easing.easeInOut,
    duration: theme.transitions.duration.standard,
  }),
  marginLeft: `275px`, // Default expanded
  [theme.breakpoints.down('sm')]: {
    marginLeft: '80px', // Collapsed on small screens
  },
}));

interface Transaction {
  date: string;
  time: string;
  fileName: string;
  sender: string;
  senderAddress: string; // New field for sender's address
  receiver: string;
  receiverAddress: string; // New field for receiver's address
  amount: string;
  status: string;
}

const initialData: Transaction[] = [
  { date: '2024-10-18', time: '12:30 PM', fileName: 'file1.txt', sender: 'john_doe', senderAddress: 'john_doe', receiver: 'john_doe', receiverAddress: 'john_doe', amount: '$100', status: 'Completed' },
  { date: '2024-10-17', time: '11:00 AM', fileName: 'file2.txt', sender: 'john_doe', senderAddress: 'john_doe', receiver: 'john_doe', receiverAddress: 'john_doe', amount: '$200', status: 'Pending' },
  // Add more dummy data here
];

const GlobalTransactions = () => {
  const [transactionData, setTransactionData] = useState<Transaction[]>(initialData);
  const [sortConfig, setSortConfig] = useState<{ key: keyof Transaction; direction: 'asc' | 'desc' }>({
    key: 'date',
    direction: 'asc',
  });

  const navigate = useNavigate();

  const handleSort = (column: keyof Transaction) => {
    const newDirection = sortConfig.key === column && sortConfig.direction === 'asc' ? 'desc' : 'asc';
    const sortedData = [...transactionData].sort((a, b) => {
      if (a[column] < b[column]) return newDirection === 'asc' ? -1 : 1;
      if (a[column] > b[column]) return newDirection === 'asc' ? 1 : -1;
      return 0;
    });
    setSortConfig({ key: column, direction: newDirection });
    setTransactionData(sortedData);
  };

  const handleFileClick = (fileName: string) => {
    navigate(`/fileview/${fileName}`); // Navigate to the file's page
  };

  const handleAccountClick = (address: string) => {
    navigate(`/account/${address}`); // Navigate to the account page
  };

  return (
    <Box sx={{ display: 'flex' }}>
      <Sidebar />
      <Main>
        {/* Header for Global Transactions */}
        <Typography variant="h4" component="h1" gutterBottom sx={{ textAlign: 'center', marginBottom: 4 }}>
          Global Transactions
        </Typography>

        <Table>
          <TableHead>
            <TableRow>
              <TableCell>
                <TableSortLabel
                  active={sortConfig.key === 'date'}
                  direction={sortConfig.direction}
                  onClick={() => handleSort('date')}
                >
                  Date
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={sortConfig.key === 'time'}
                  direction={sortConfig.direction}
                  onClick={() => handleSort('time')}
                >
                  Time
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={sortConfig.key === 'fileName'}
                  direction={sortConfig.direction}
                  onClick={() => handleSort('fileName')}
                >
                  File Name
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={sortConfig.key === 'sender'}
                  direction={sortConfig.direction}
                  onClick={() => handleSort('sender')}
                >
                  Sender
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={sortConfig.key === 'receiver'}
                  direction={sortConfig.direction}
                  onClick={() => handleSort('receiver')}
                >
                  Receiver
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={sortConfig.key === 'amount'}
                  direction={sortConfig.direction}
                  onClick={() => handleSort('amount')}
                >
                  Transaction Amount
                </TableSortLabel>
              </TableCell>
              <TableCell>
                <TableSortLabel
                  active={sortConfig.key === 'status'}
                  direction={sortConfig.direction}
                  onClick={() => handleSort('status')}
                >
                  Status
                </TableSortLabel>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {transactionData.map((transaction, index) => (
              <TableRow key={index}>
                <TableCell>{transaction.date}</TableCell>
                <TableCell>{transaction.time}</TableCell>
                <TableCell>
                  <Button
                    onClick={() => handleFileClick(transaction.fileName)}
                    color="primary"
                    variant="outlined"
                  >
                    {transaction.fileName}
                  </Button>
                </TableCell>
                <TableCell>
                  <Button onClick={() => handleAccountClick(transaction.senderAddress)} color="primary" variant="outlined">
                    {transaction.sender}
                  </Button>
                </TableCell>
                <TableCell>
                  <Button onClick={() => handleAccountClick(transaction.receiverAddress)} color="primary" variant="outlined">
                    {transaction.receiver}
                  </Button>
                </TableCell>
                <TableCell>{transaction.amount}</TableCell>
                <TableCell>{transaction.status}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Main>
    </Box>
  );
};

export default GlobalTransactions;
