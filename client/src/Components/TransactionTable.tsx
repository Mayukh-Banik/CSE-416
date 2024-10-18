import React, { useState } from "react";
import { Box, Typography, Container, Table, TableBody, TableCell, TableHead, TableRow, TableSortLabel, TableContainer, Paper, IconButton } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { QRCodeCanvas } from "qrcode.react";
import ArrowForwardIcon from '@mui/icons-material/ArrowForward'; // Add an icon to indicate clickability
import Sidebar from "./Sidebar";

interface Transaction {
    id: string;
    sender: string;
    receiver: string;
    amount: number;
    timestamp: string;
    status: string;
}

interface WalletDetailsProps {
    walletAddress: string;
    balance: number;
    transactions: Transaction[];
    walletLabel?: string;
    fee?: number;
}

// Sorting helper function
const sortTransactions = (transactions: Transaction[], orderBy: keyof Transaction, order: 'asc' | 'desc') => {
    return transactions.sort((a, b) => {
        if (order === 'asc') {
            return a[orderBy] < b[orderBy] ? -1 : 1;
        } else {
            return a[orderBy] > b[orderBy] ? -1 : 1;
        }
    });
};

const WalletPage: React.FC<WalletDetailsProps> = ({
    walletAddress,
    balance,
    transactions,
    walletLabel,
    fee,
}) => {
    const navigate = useNavigate();

    const [order, setOrder] = useState<'asc' | 'desc'>('asc');
    const [orderBy, setOrderBy] = useState<keyof Transaction>('timestamp');

    // Handle sorting click
    const handleSort = (property: keyof Transaction) => {
        const isAsc = orderBy === property && order === 'asc';
        setOrder(isAsc ? 'desc' : 'asc');
        setOrderBy(property);
    };

    // Handle transaction click (navigating to detail page)
    const handleTransactionClick = (id: string) => {
        navigate(`/transaction/${id}`);
    };

    const sortedTransactions = sortTransactions([...transactions], orderBy, order);

    return (
        <>
            <Sidebar />
            <Container maxWidth="md" sx={{ mt: 4 }}>
                <Typography variant="h4" sx={{ mb: 2 }}>
                    {walletLabel || 'Wallet Details'}
                </Typography>
                
                {/* Balance and Wallet Address */}
                <Box
                    sx={{
                        display: 'flex',
                        flexDirection: 'row',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        padding: 2,
                        gap: 2,
                    }}
                >
                    {/* Balance Box */}
                    <Box
                        sx={{
                            display: 'flex',
                            flexDirection: 'row',
                            alignItems: 'center',
                            padding: 2,
                            borderRadius: 2,
                            boxShadow: 2,
                            width: '45%',
                        }}
                    >
                        <img src="file://absolute/path/to/images/walletBalance.png" alt="Squid Icon" width="40" />
                        <Box sx={{ ml: 2, textAlign: 'left' }}>
                            <Typography variant="h6">Balance</Typography>
                            <Typography variant="body1">{balance} Coins</Typography>
                        </Box>
                    </Box>

                    {/* Wallet Address Box */}
                    <Box
                        sx={{
                            display: 'flex',
                            flexDirection: 'row',
                            alignItems: 'center',
                            padding: 2,
                            borderRadius: 2,
                            boxShadow: 2,
                            width: '45%',
                            ml: 2,
                        }}
                    >
                        <QRCodeCanvas value={walletAddress} size={50} style={{ marginRight: '10px' }} />
                        <Box sx={{ ml: 2, textAlign: 'left' }}>
                            <Typography variant="h6">Wallet Address</Typography>
                            <Typography variant="body1">{walletAddress}</Typography>
                        </Box>
                    </Box>
                </Box>

                <Typography variant="h6" sx={{ mt: 4 }}>Transaction History</Typography>

                {/* Transaction Table */}
                <TableContainer component={Paper}>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'id'}
                                        direction={orderBy === 'id' ? order : 'asc'}
                                        onClick={() => handleSort('id')}
                                    >
                                        ID
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'sender'}
                                        direction={orderBy === 'sender' ? order : 'asc'}
                                        onClick={() => handleSort('sender')}
                                    >
                                        Sender
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'receiver'}
                                        direction={orderBy === 'receiver' ? order : 'asc'}
                                        onClick={() => handleSort('receiver')}
                                    >
                                        Receiver
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'amount'}
                                        direction={orderBy === 'amount' ? order : 'asc'}
                                        onClick={() => handleSort('amount')}
                                    >
                                        Amount
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>
                                    <TableSortLabel
                                        active={orderBy === 'timestamp'}
                                        direction={orderBy === 'timestamp' ? order : 'asc'}
                                        onClick={() => handleSort('timestamp')}
                                    >
                                        Date
                                    </TableSortLabel>
                                </TableCell>
                                <TableCell>Status</TableCell>
                                <TableCell>Action</TableCell> {/* Add action column */}
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {sortedTransactions.map((transaction) => (
                                <TableRow
                                    key={transaction.id}
                                    hover
                                    sx={{
                                        cursor: 'pointer',
                                        transition: 'background-color 0.3s ease',
                                        '&:hover': {
                                            backgroundColor: 'rgba(0, 0, 0, 0.04)',
                                        },
                                    }}
                                    onClick={() => handleTransactionClick(transaction.id)}
                                >
                                    <TableCell>{transaction.id}</TableCell>
                                    <TableCell>{transaction.sender}</TableCell>
                                    <TableCell>{transaction.receiver}</TableCell>
                                    <TableCell>{transaction.amount}</TableCell>
                                    <TableCell>{new Date(transaction.timestamp).toLocaleString()}</TableCell>
                                    <TableCell>{transaction.status}</TableCell>
                                    <TableCell>
                                        <IconButton>
                                            <ArrowForwardIcon color="primary" />
                                        </IconButton>
                                    </TableCell> {/* Add an icon to indicate action */}
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </Container>
        </>
    );
};

export default WalletPage;
