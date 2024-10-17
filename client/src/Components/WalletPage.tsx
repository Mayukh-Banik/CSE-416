import React from "react";
import { Box, Typography, Container } from "@mui/material";
import { QRCodeCanvas } from "qrcode.react";
import Sidebar from "./Sidebar"; // Ensure this is imported correctly
import TransactionPage from "./TransactionPage"; // Make sure this component is correctly defined
import { useTheme } from '@mui/material/styles';

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
    // isCollapsed: boolean; // Add this prop to manage sidebar state
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const WalletPage: React.FC<WalletDetailsProps> = ({
    walletAddress,
    balance,
    transactions,
    walletLabel,
    fee,
    // isCollapsed, // Accept the prop
}) => {
    // const marginLeft = isCollapsed ? `${collapsedDrawerWidth}px` : `${drawerWidth}px`; // Determine margin based on sidebar state
    const theme = useTheme();
    return (
        <Box
            sx={{
                padding: 2,
                marginTop: '70px',
                marginLeft: `${drawerWidth}px`, // Default expanded margin
                transition: 'margin-left 0.3s ease', // Smooth transition
                [theme.breakpoints.down('sm')]: {
                    marginLeft: `${collapsedDrawerWidth}px`, // Adjust left margin for small screens
                },
            }}
        >
            <Sidebar />
            <Container maxWidth="md">
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
                        <img src="/squidcoin.png" alt="Squid Icon" width="40" />
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
                        {/* QR Code generation */}
                        <QRCodeCanvas value={walletAddress} size={50} style={{ marginRight: '10px' }} />
                        <Box sx={{ ml: 2, textAlign: 'left' }}>
                            <Typography variant="h6">Wallet Address</Typography>
                            <Typography variant="body1">{walletAddress}</Typography>
                        </Box>
                    </Box>
                </Box>

                <Typography variant="h6">Transaction History</Typography>
                <TransactionPage /> {/* Ensure TransactionPage accepts transactions as a prop */}
            </Container>
        </Box>
    );
};

export default WalletPage;
