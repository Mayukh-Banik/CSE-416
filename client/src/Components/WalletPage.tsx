import React, { useState } from "react";
import { Box, Typography, Container, Table, TableBody, TableCell, TableHead, TableRow, TableSortLabel, TableContainer, Paper, Dialog, DialogTitle, DialogContent, DialogActions, Button } from "@mui/material";
import { QRCodeCanvas } from "qrcode.react";
import Sidebar from "./Sidebar";
import { useTheme } from '@mui/material/styles';
import { Link } from "react-router-dom";

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

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

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
    const theme = useTheme();
    const [order, setOrder] = useState<'asc' | 'desc'>('asc');
    const [orderBy, setOrderBy] = useState<keyof Transaction>('timestamp');
    const [selectedTransaction, setSelectedTransaction] = useState<Transaction | null>(null);
    const [isDialogOpen, setDialogOpen] = useState(false);

    const handleSort = (property: keyof Transaction) => {
        const isAsc = orderBy === property && order === 'asc';
        setOrder(isAsc ? 'desc' : 'asc');
        setOrderBy(property);
    };

    const handleTransactionClick = (transaction: Transaction) => {
        setSelectedTransaction(transaction);
        setDialogOpen(true);
    };

    const handleCloseDialog = () => {
        setDialogOpen(false);
    };

    const sortedTransactions = sortTransactions([...transactions], orderBy, order);

    return (
        <Box
            sx={{
                padding: 2,
                marginTop: '70px',
                marginLeft: `${drawerWidth}px`,
                transition: 'margin-left 0.3s ease',
                [theme.breakpoints.down('sm')]: {
                    marginLeft: `${collapsedDrawerWidth}px`,
                },
            }}
        >
            <Sidebar />
            <Container maxWidth="md">
                <Typography variant="h4" sx={{ mb: 2 }}>
                    {walletLabel || 'Wallet Details'}
                </Typography>

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
                    <Box
                        sx={{
                            display: 'flex',
                            flexDirection: 'row',
                            alignItems: 'center',
                            padding: 2,
                            borderRadius: 2,
                            boxShadow: 2,
                            width: '45%',
                            background: 'white',
                        }}
                    >
                        <img src={`${process.env.PUBLIC_URL}/squidcoin.png`} alt="Squid Icon" width="30" />
                        <Box sx={{ ml: 2, textAlign: 'left' }}>
                            <Typography variant="h6">Balance</Typography>
                            <Typography variant="body1">{balance} Coins</Typography>
                        </Box>
                    </Box>

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
                            background: 'white',
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
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {sortedTransactions.map((transaction) => (
                                <TableRow
                                    key={transaction.id}
                                    hover
                                    sx={{ cursor: 'pointer' }}
                                    onClick={() => handleTransactionClick(transaction)}
                                >
                                    <TableCell>{transaction.id}</TableCell>
                                    <TableCell>{transaction.sender}</TableCell>
                                    <TableCell>{transaction.receiver}</TableCell>
                                    <TableCell>{transaction.amount}</TableCell>
                                    <TableCell>{new Date(transaction.timestamp).toLocaleString()}</TableCell>
                                    <TableCell>{transaction.status}</TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>

                <Dialog open={isDialogOpen} onClose={handleCloseDialog}>
                    <DialogTitle>Transaction Details</DialogTitle>
                    <DialogContent>
                        {selectedTransaction && (
                            <Box>
                                <Typography variant="body1">
                                    <strong>ID:</strong> {selectedTransaction.id} <br />
                                    <strong>Sender:</strong> {selectedTransaction.sender} <br />
                                    <strong>Receiver:</strong> {selectedTransaction.receiver} <br />
                                    <strong>Amount:</strong> {selectedTransaction.amount} <br />
                                    <strong>Date:</strong>{" "}
                                    {new Date(selectedTransaction.timestamp).toLocaleString()} <br />
                                    <strong>Status:</strong> {selectedTransaction.status} <br />
                                </Typography>

                                <Box
                                    sx={{
                                        display: "flex",
                                        flexDirection: "column",
                                        alignItems: "center",
                                        marginTop: "20px"
                                    }}
                                >
                                    <Button
                                        variant="contained"
                                        component={Link}
                                        to={`/fileview/${selectedTransaction.id}`}
                                    >
                                        View File
                                    </Button>
                                </Box>
                            </Box>
                        )}
                    </DialogContent>
                    <DialogActions sx={{ marginTop: "-10px", justifyContent: "center" }}>
                        <Button onClick={handleCloseDialog} color="primary">
                            Close
                        </Button>
                    </DialogActions>
                </Dialog>
            </Container>
        </Box>
    );
};

export default WalletPage;
