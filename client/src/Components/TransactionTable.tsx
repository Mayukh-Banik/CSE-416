import React, { useState } from "react";
import {
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TablePagination,
} from "@mui/material";
import { TransactionTableProps } from "../types/interfaces";

const TransactionTable: React.FC<TransactionTableProps> = ({
    transactions,
    search,
    dateFilter,
    statusFilter,
}) => {
    // Helper function to filter transactions based on search, date, and status filters
    const filterTransactions = () => {
        return transactions.filter((transaction) => {
            // Search filter
            const searchMatch =
                transaction.fileName.toLowerCase().includes(search.toLowerCase()) ||
                transaction.transactionId.toLowerCase().includes(search.toLowerCase());

            // Date filter
            const date = new Date(transaction.timestamp);
            const now = new Date();
            let dateMatch = true;

            if (dateFilter === "today") {
                dateMatch = date.toDateString() === now.toDateString();
            } else if (dateFilter === "week") {
                const weekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
                dateMatch = date >= weekAgo;
            } else if (dateFilter === "month") {
                const monthAgo = new Date(now.getFullYear(), now.getMonth() - 1, now.getDate());
                dateMatch = date >= monthAgo;
            } else if (dateFilter === "year") {
                const yearAgo = new Date(now.getFullYear() - 1, now.getMonth(), now.getDate());
                dateMatch = date >= yearAgo;
            }

            // Status filter
            const statusMatch =
                statusFilter === "all" || transaction.status.toLowerCase() === statusFilter.toLowerCase();

            return searchMatch && dateMatch && statusMatch;
        });
    };

    const filteredTransactions = filterTransactions();

    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(10);

    const handleChangePage = (event: unknown, newPage: number) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (
        event: React.ChangeEvent<HTMLInputElement>
    ) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0); // reset page to 0 whenever rows per page changes
    };

    const paginatedTransactions = filteredTransactions.slice(
        page * rowsPerPage,
        page * rowsPerPage + rowsPerPage
    );

    return (
        <TableContainer>
            <Table>
                <TableHead sx={{ backgroundColor: "grey.200" }}>
                    <TableRow>
                        <TableCell>Date & Time</TableCell>
                        <TableCell>Transaction ID</TableCell>
                        <TableCell>File Name</TableCell>
                        <TableCell>File Size</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>Sender</TableCell>
                        <TableCell>Receiver</TableCell>
                        <TableCell>Fee</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {filteredTransactions.length > 0 ? (
                        paginatedTransactions.map((transaction) => (
                            <TableRow key={transaction._id}>
                                <TableCell>{new Date(transaction.timestamp).toLocaleString()}</TableCell>
                                <TableCell>{transaction.transactionId}</TableCell>
                                <TableCell>{transaction.fileName}</TableCell>
                                <TableCell>{transaction.fileSize}</TableCell>
                                <TableCell>{transaction.status}</TableCell>
                                <TableCell>{transaction.sender}</TableCell>
                                <TableCell>{transaction.receiver}</TableCell>
                                <TableCell>{transaction.fee}</TableCell>
                            </TableRow>
                        ))
                    ) : (
                        <TableRow>
                            <TableCell colSpan={8} align="center">
                                No transactions found
                            </TableCell>
                        </TableRow>
                    )}
                </TableBody>
            </Table>
            <TablePagination
                component="div"
                count={filteredTransactions.length}
                page={page}
                onPageChange={handleChangePage}
                rowsPerPage={rowsPerPage}
                onRowsPerPageChange={handleChangeRowsPerPage}
                rowsPerPageOptions={[5, 10, 25]}
            />
        </TableContainer>
    );
};

export {};

export default TransactionTable;
