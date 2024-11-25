// import React, { useEffect, useState } from 'react';
// import {
//   Box, Table, TableBody, TableCell, TableHead, TableRow,
//   TableSortLabel, Typography, Button, Paper, Tabs, Tab
// } from '@mui/material';
// import Sidebar from './Sidebar'; // Adjust the path as necessary
// import { styled, useTheme } from '@mui/material/styles';
// import { useNavigate } from 'react-router-dom';
// import {  PendingRequest } from '../models/transactions';
// import axios from 'axios';

// const drawerWidth = 300;
// const collapsedDrawerWidth = 100;

// interface Transaction {
//   CreatedAt?: string;
//   FileName: string;
//   FileHash: string;
//   RequesterID: string | "";
//   TargetID: string;
//   Status?: 'pending' | 'accepted' | 'declined' | 'completed';
//   Fee?:  number;
// }

// // const initialData: SampleTransaction[] = [
// //   { date: '2024-10-18', time: '12:30 PM', fileName: 'file1.txt', sender: 'john_doe', senderAddress: 'john_doe', receiver: 'john_doe', receiverAddress: 'john_doe', amount: '1.2 ORC', status: 'Completed' },
// //   { date: '2024-10-17', time: '11:00 AM', fileName: 'file2.txt', sender: 'john_doe', senderAddress: 'john_doe', receiver: 'john_doe', receiverAddress: 'john_doe', amount: '0.4 ORC', status: 'Pending' },
// // ];

// // const pendingRequestsData: SampleTransaction[] = [
// //   { date: '2024-10-20', time: '2:15 PM', fileName: 'file3.txt', sender: 'alice_smith', senderAddress: 'alice_smith', receiver: 'john_doe', receiverAddress: 'john_doe', amount: '0.5 ORC', status: 'Pending' },
// // ];

// const GlobalTransactions = () => {
//   const theme = useTheme();
//   const [transactionData, setTransactionData] = useState<Transaction[]>([]);
//   const [pendingRequests, setPendingRequests] = useState<Transaction[]>([]);
//   const [activeTab, setActiveTab] = useState<number>(0); // Track active tab
//   const [sortConfig, setSortConfig] = useState<{ key: keyof SampleTransaction; direction: 'asc' | 'desc' }>({
//     key: 'date',
//     direction: 'asc',
//   });

//   const navigate = useNavigate();

//   useEffect(() => {
//     // Fetch pending requests from the backend
//     const fetchPendingRequests = async () => {
//       try {
//         const response = await axios.get<Transaction[]>('http://localhost:8080/api/pending-requests');
//         setPendingRequests(Object.values(response.data)); // Convert the map to an array
//       } catch (error) {
//         console.error('Failed to fetch pending requests:', error);
//       }
//     };

//     fetchPendingRequests();
//   }, []);
  
//   const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
//     setActiveTab(newValue);
//   };

//   const handleSort = (column: keyof Transaction) => {
//     const newDirection = sortConfig.key === column && sortConfig.direction === 'asc' ? 'desc' : 'asc';
//     const sortedData = [...transactionData].sort((a, b) => {
//       if (a[column] < b[column]) return newDirection === 'asc' ? -1 : 1;
//       if (a[column] > b[column]) return newDirection === 'asc' ? 1 : -1;
//       return 0;
//     });
//     setSortConfig({ key: column, direction: newDirection });
//     setTransactionData(sortedData);
//   };

//   const handleFileClick = (fileName: string) => {
//     navigate(`/fileview/${fileName}`);
//   };

//   const handleAccountClick = (address: string) => {
//     navigate(`/account/${address}`);
//   };

//   const renderTable = (data: Transaction[]) => (
//     <Table component={Paper}>
//       <TableHead>
//         <TableRow>
//           <TableCell>
//             <TableSortLabel
//               active={sortConfig.key === 'date'}
//               direction={sortConfig.direction}
//               onClick={() => handleSort('date')}
//             >
//               Date
//             </TableSortLabel>
//           </TableCell>
//           <TableCell>
//             <TableSortLabel
//               active={sortConfig.key === 'time'}
//               direction={sortConfig.direction}
//               onClick={() => handleSort('time')}
//             >
//               Time
//             </TableSortLabel>
//           </TableCell>
//           <TableCell>
//             <TableSortLabel
//               active={sortConfig.key === 'fileName'}
//               direction={sortConfig.direction}
//               onClick={() => handleSort('fileName')}
//             >
//               File Name
//             </TableSortLabel>
//           </TableCell>
//           <TableCell>
//             <TableSortLabel
//               active={sortConfig.key === 'sender'}
//               direction={sortConfig.direction}
//               onClick={() => handleSort('sender')}
//             >
//               Sender
//             </TableSortLabel>
//           </TableCell>
//           <TableCell>
//             <TableSortLabel
//               active={sortConfig.key === 'receiver'}
//               direction={sortConfig.direction}
//               onClick={() => handleSort('receiver')}
//             >
//               Receiver
//             </TableSortLabel>
//           </TableCell>
//           <TableCell>
//             <TableSortLabel
//               active={sortConfig.key === 'amount'}
//               direction={sortConfig.direction}
//               onClick={() => handleSort('amount')}
//             >
//               Transaction Amount
//             </TableSortLabel>
//           </TableCell>
//           <TableCell>
//             <TableSortLabel
//               active={sortConfig.key === 'status'}
//               direction={sortConfig.direction}
//               onClick={() => handleSort('status')}
//             >
//               Status
//             </TableSortLabel>
//           </TableCell>
//         </TableRow>
//       </TableHead>
//       <TableBody>
//         {data.map((transaction, index) => (
//           <TableRow key={index}>
//             <TableCell>{transaction.CreatedAt}</TableCell>
//             {/* <TableCell>{transaction.time}</TableCell> */}
//             <TableCell>
//               <Button onClick={() => handleFileClick(transaction.FileName)} color="primary" variant="outlined">
//                 {transaction.FileName}
//               </Button>
//             </TableCell>
//             <TableCell>
//               <Button onClick={() => handleAccountClick(transaction.RequesterID)} color="primary" variant="outlined">
//                 {transaction.RequesterID}
//               </Button>
//             </TableCell>
//             <TableCell>
//               <Button onClick={() => handleAccountClick(transaction.TargetID)} color="primary" variant="outlined">
//                 {transaction.TargetID}
//               </Button>
//             </TableCell>
//             <TableCell>{transaction.Fee}</TableCell>
//             <TableCell>{transaction.Status}</TableCell>
//           </TableRow>
//         ))}
//       </TableBody>
//     </Table>
//   );

//   return (
//     <Box
//       sx={{
//         padding: 2,
//         marginTop: '70px',
//         marginLeft: `${drawerWidth}px`, 
//         transition: 'margin-left 0.3s ease',
//         [theme.breakpoints.down('sm')]: {
//           marginLeft: `${collapsedDrawerWidth}px`,
//         },
//       }}
//     >
//       <Sidebar />
//       <Typography variant="h4" component="h1" gutterBottom sx={{ textAlign: 'left', marginBottom: 4 }}>
//         Transactions
//       </Typography>

//       <Tabs value={activeTab} onChange={handleTabChange} aria-label="transaction tabs">
//         <Tab label="Global Transactions" />
//         <Tab label="My Transactions" />
//         <Tab label="Pending Requests" />
//       </Tabs>

//       {activeTab === 0 && renderTable(transactionData)}  {/* Global Transactions */}
//       {activeTab === 1 && renderTable(transactionData.filter(tx => tx.sender === 'john_doe' || tx.receiver === 'john_doe'))} {/* My Transactions */}
//       {activeTab === 2 && renderTable(pendingRequests)}  {/* Pending Requests */}
//     </Box>
//   );
// };

// export default GlobalTransactions;

export {}