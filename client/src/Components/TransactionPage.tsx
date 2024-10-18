import React, { useState } from "react";
import {
    Box,
    Container,
    CssBaseline,
    TextField,
    Select,
    MenuItem,
    InputLabel,
    FormControl,
    SelectChangeEvent
} from "@mui/material";
import TransactionTable from "./TransactionTable";
import Sidebar from "./Sidebar";

// interface TransactionPageProps {
//     open: boolean;
// }

const TransactionPage: React.FC = () => {
    const [search, setSearch] = useState("");
    const [dateFilter, setDateFilter] = useState<string>("all");
    const [statusFilter, setStatusFilter] = useState<string>("all");
    const [sidebarOpen, setSidebarOpen] = useState<boolean>(true); // State to toggle sidebar

    const toggleSidebar = () => {
      setSidebarOpen((prev) => !prev);
    };

    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
        setSearch(event.target.value);
    };

    const handleDateFilterChange = (event: SelectChangeEvent<string>) => {
        setDateFilter(event.target.value as string);
    };

    const handleStatusFilterChange = (event: SelectChangeEvent<string>) => {
        setStatusFilter(event.target.value as string);
    };

    return (
        
        <Box sx={{ display: "flex" }}>
            <CssBaseline />
            <Sidebar />
            <Container component="main" sx={{ flexGrow: 1, p: 3 }}>
                <TextField
                    label="Search Transactions"
                    variant="outlined"
                    onChange={handleSearch}
                    value={search}
                    sx={{ marginBottom: 2 }}
                />

                {/* Date Filter */}
                <FormControl variant="outlined" sx={{ marginBottom: 2, minWidth: 120 }}>
                    <InputLabel>Date</InputLabel>
                    <Select
                        value={dateFilter}
                        onChange={handleDateFilterChange}
                        label="Date"
                    >
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
                    <Select
                        value={statusFilter}
                        onChange={handleStatusFilterChange}
                        label="Status"
                    >
                        <MenuItem value="all">All</MenuItem>
                        <MenuItem value="Complete">Complete</MenuItem>
                        <MenuItem value="Pending">Pending</MenuItem>
                        <MenuItem value="Failed">Failed</MenuItem>
                    </Select>
                </FormControl>

                {/* Transaction Table */}
                {/* <TransactionTable search={search} dateFilter={dateFilter} statusFilter={statusFilter} /> */}
            </Container>
        </Box>
    );
};

export default TransactionPage;
