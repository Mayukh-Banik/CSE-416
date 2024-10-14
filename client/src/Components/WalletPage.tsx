import React, { useEffect, useState } from "react";
import axios, { AxiosResponse } from "axios";
import {
    Box,
    Typography,
    Container,
    Button,
    CssBaseline,
    TextField,
    Select,
    MenuItem,
    InputLabel,
    FormControl,
    SelectChangeEvent,
} from "@mui/material";
import { QRCodeCanvas } from "qrcode.react";
import Sidebar from "./Sidebar";
import { TransactionProps } from "../types/interfaces";

interface User {
    userId: string;
    walletAddress: string;
    balance: number;
    walletLabel: string;
    // 추가적인 필드가 필요하면 여기에 추가
}

const WalletPage: React.FC = () => {
    // State for user details and transactions
    const [user, setUser] = useState<User | null>(null);
    const [transactions, setTransactions] = useState<TransactionProps[]>([]);
    const [itemsToShow, setItemsToShow] = useState(5);
    const [search, setSearch] = useState("");
    const [dateFilter, setDateFilter] = useState<string>("all");
    const [statusFilter, setStatusFilter] = useState<string>("all");

    useEffect(() => {
        // Fetch user information from the server using the token in HttpOnly cookie
        const fetchUserInfo = async () => {
            try {
                console.log("씨발");
                const response = await axios.get("http://localhost:5000/api/users/info", {
                    withCredentials: true,
                });
                console.log("쭈발",response);
                setUser(response.data.user);

                // Fetch transaction history when user information is available
                const fetchTransactions = async () => {
                    try {
                        const transactionsResponse: AxiosResponse<TransactionProps[]> = await axios.get(
                            `http://localhost:5000/api/transaction/${response.data.user.userId}`,
                            { withCredentials: true }
                        );
                        setTransactions(transactionsResponse.data);
                    } catch (error) {
                        console.error("Error fetching transactions:", error);
                    }
                };

                fetchTransactions();
            } catch (error) {
                console.error("Error fetching user info:", error);
            }
        };

        fetchUserInfo();
    }, []);

    // Function to load more transactions
    const loadMoreTransactions = () => {
        setItemsToShow((prev) => prev + 5);
    };

    // Handle search and filter changes
    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
        setSearch(event.target.value);
    };

    const handleDateFilterChange = (event: SelectChangeEvent<string>) => {
        setDateFilter(event.target.value as string);
    };

    const handleStatusFilterChange = (event: SelectChangeEvent<string>) => {
        setStatusFilter(event.target.value as string);
    };

    // Sort transactions by timestamp in descending order
    const sortedTransactions = [...transactions].sort(
        (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    );

    if (!user) {
        return <div>Loading user information...</div>;
    }

    return (
        <>
            <CssBaseline />
            <Sidebar />
            <Container maxWidth="md" sx={{ mt: 4 }}>
                <Typography variant="h4" sx={{ mb: 2 }}>
                    {user.walletLabel || "Wallet Details"}
                </Typography>
                {/* Balance and Wallet Address */}
                <Box
                    sx={{
                        display: "flex",
                        flexDirection: "row",
                        justifyContent: "space-between",
                        alignItems: "center",
                        padding: 2,
                        gap: 2,
                    }}
                >
                    {/* Balance Box */}
                    <Box
                        sx={{
                            display: "flex",
                            flexDirection: "row",
                            alignItems: "center",
                            padding: 2,
                            borderRadius: 2,
                            boxShadow: 2,
                            width: "45%",
                        }}
                    >
                        <img
                            src="/images/walletBalance.png"
                            alt="Wallet Balance Icon"
                            width="40"
                        />
                        <Box sx={{ ml: 2, textAlign: "left" }}>
                            <Typography variant="h6">Balance</Typography>
                            <Typography variant="body1">{user.balance} Coins</Typography>
                        </Box>
                    </Box>

                    {/* Wallet Address Box */}
                    <Box
                        sx={{
                            display: "flex",
                            flexDirection: "row",
                            alignItems: "center",
                            padding: 2,
                            borderRadius: 2,
                            boxShadow: 2,
                            width: "45%",
                            ml: 2,
                        }}
                    >
                        {/* QR Code generate */}
                        <QRCodeCanvas
                            value={user.walletAddress}
                            size={50}
                            style={{ marginRight: "10px" }}
                        />
                        <Box sx={{ ml: 2, textAlign: "left" }}>
                            <Typography variant="h6">Wallet Address</Typography>
                            <Typography variant="body1">{user.walletAddress}</Typography>
                        </Box>
                    </Box>
                </Box>

                {/* Transaction Filters */}
                <Box sx={{ mt: 4, mb: 2 }}>
                    <TextField
                        label="Search Transactions"
                        variant="outlined"
                        onChange={handleSearch}
                        value={search}
                        sx={{ marginBottom: 2, mr: 2 }}
                    />

                    {/* Date Filter */}
                    <FormControl variant="outlined" sx={{ marginBottom: 2, minWidth: 120, mr: 2 }}>
                        <InputLabel>Date</InputLabel>
                        <Select value={dateFilter} onChange={handleDateFilterChange} label="Date">
                            <MenuItem value="all">All</MenuItem>
                            <MenuItem value="today">Today</MenuItem>
                            <MenuItem value="week">This Week</MenuItem>
                            <MenuItem value="month">This Month</MenuItem>
                            <MenuItem value="year">This Year</MenuItem>
                        </Select>
                    </FormControl>

                    {/* Status Filter */}
                    <FormControl variant="outlined" sx={{ marginBottom: 2, minWidth: 120 }}>
                        <InputLabel>Status</InputLabel>
                        <Select value={statusFilter} onChange={handleStatusFilterChange} label="Status">
                            <MenuItem value="all">All</MenuItem>
                            <MenuItem value="pending">Pending</MenuItem>
                            <MenuItem value="completed">Completed</MenuItem>
                            <MenuItem value="failed">Failed</MenuItem>
                        </Select>
                    </FormControl>
                </Box>

                {/* Transaction History */}
                <Typography variant="h6" sx={{ mt: 4 }}>
                    Transaction History
                </Typography>
                {/* <TransactionTable
                    search={search}
                    dateFilter={dateFilter}
                    statusFilter={statusFilter}
                    transactions={sortedTransactions.slice(0, itemsToShow)}
                /> */}

                {itemsToShow < sortedTransactions.length && (
                    <Button variant="contained" sx={{ mt: 2 }} onClick={loadMoreTransactions}>
                        Load More
                    </Button>
                )}
            </Container>
        </>
    );
};

export default WalletPage;