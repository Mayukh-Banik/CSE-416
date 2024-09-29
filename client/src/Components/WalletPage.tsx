import React, { useState } from "react";
import { Box, Typography, Container, List, ListItem, ListItemText, Button } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { QRCodeCanvas } from "qrcode.react";
import Header from "./Header";

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

const WalletPage: React.FC<WalletDetailsProps> = ({
    walletAddress,
    balance,
    transactions,
    walletLabel,
    fee,
}) => {
    // Pagination state for showing limited transactions
    const [itemsToShow, setItemsToShow] = useState(5);

    const sortedTransactions = [...transactions].sort(
        (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    );

    // Chart Data for Network Fee
    const feeChartData = {
        labels: ['Network Fee'],
        datasets: [
            {
                label: 'Fee',
                data: [fee || 0],
                backgroundColor: 'rgba(75,192,192,0.6)',
                borderColor: 'rgba(75,192,192,1)',
                borderWidth: 1,
            },
        ],
    };

    return (
        <>
            <Header />
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
                        <img src="/images/walletBalance.png" alt="Wallet Icon" width="40" />
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
                        {/* QR Code generate */}
                        <QRCodeCanvas value={walletAddress} size={50} style={{ marginRight: '10px' }} />
                        <Box sx={{ ml: 2, textAlign: 'left' }}>
                            <Typography variant="h6">Wallet Address</Typography>
                            <Typography variant="body1">{walletAddress}</Typography>
                        </Box>
                    </Box>
                </Box>


                {/* Transaction History */}
                <Box sx={{ mb: 4 }}>
                    <Typography variant="h6">Transaction History</Typography>
                    <List>
                        {sortedTransactions.slice(0, itemsToShow).map((tx) => (
                            <ListItem key={tx.id} sx={{ borderBottom: '1px solid #ccc', padding: 2, flexDirection: 'column', alignItems: 'flex-start' }}>
                                
                                {/* Transaction ID as Header */}
                                <Typography variant="h6" sx={{ mb: 1 }}>
                                    Transaction ID: {tx.id}
                                </Typography>
                                
                                {/* Transaction Details in Two Columns */}
                                <Box
                                    sx={{
                                        display: 'flex',
                                        flexDirection: 'row',
                                        justifyContent: 'flex-start',
                                        gap: 4, // Gap between the two columns
                                        width: '100%',
                                    }}
                                >
                                    {/* Left Column: Sender, Receiver, Amount */}
                                    <Box sx={{ flex: 1 }}>
                                        <Typography variant="body2">
                                            <strong>Sender:</strong> {tx.sender}
                                        </Typography>
                                        <Typography variant="body2">
                                            <strong>Receiver:</strong> {tx.receiver}
                                        </Typography>
                                        <Typography variant="body2">
                                            <strong>Amount:</strong> {tx.amount} Coins
                                        </Typography>
                                    </Box>

                                    {/* Right Column: Time, Status */}
                                    <Box sx={{ flex: 1 }}>
                                        <Typography variant="body2">
                                            <strong>Time:</strong> {tx.timestamp}
                                        </Typography>
                                        <Typography variant="body2">
                                            <strong>Status:</strong> {tx.status}
                                        </Typography>
                                    </Box>
                                </Box>
                            </ListItem>
                        ))}
                    </List>

                    {/* Pagination Buttons */}
                    <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
                        <Button variant="contained" onClick={() => setItemsToShow(itemsToShow + 5)}>
                            Show More
                        </Button>
                    </Box>
                </Box>
            </Container >
        </>
    );
};


export default WalletPage;