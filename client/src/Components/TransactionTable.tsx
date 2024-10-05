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

// not sure how this would fit with schema
interface Transaction {
    dateTime: string;
    transactionId: string;
    sender: string;
    receiver: string;
    fileName: string;
    fileSize: string;
    status: "Complete" | "Pending" | "Failed";
    fee: string;
}

// filler data for displaying
const transactions: Transaction[] = [
    {
        dateTime: "2024-09-25 10:00",
        transactionId: "TX123456",
        sender: "0xfcfc...2a3a",
        receiver: "0xacda...c452",
        fileName: "file1.zip",
        fileSize: "10MB",
        status: "Complete",
        fee: "0.50",
    },
    {
        dateTime: "2024-09-22 12:30",
        transactionId: "TX654321",
        sender: "0xabcd...1234",
        receiver: "0xdeef...5678",
        fileName: "file2.jpg",
        fileSize: "20MB",
        status: "Pending",
        fee: "0.30",
    },
    {
        dateTime: "2024-09-20 14:45",
        transactionId: "TX789012",
        sender: "0xefgh...3456",
        receiver: "0ijkl...7890",
        fileName: "file3.pdf",
        fileSize: "15MB",
        status: "Failed",
        fee: "0.40",
    },
    {
        dateTime: "2024-09-18 08:10",
        transactionId: "TX456789",
        sender: "0xmnop...9012",
        receiver: "0qrst...3456",
        fileName: "file4.png",
        fileSize: "5MB",
        status: "Complete",
        fee: "0.25",
    },
    {
        dateTime: "2024-09-15 16:20",
        transactionId: "TX234567",
        sender: "0xuvwx...7890",
        receiver: "0xabcd...1234",
        fileName: "file5.docx",
        fileSize: "30MB",
        status: "Complete",
        fee: "0.60",
    },
    {
        dateTime: "2024-09-12 11:00",
        transactionId: "TX098765",
        sender: "0xdefg...3456",
        receiver: "0xhijk...7890",
        fileName: "file6.mp4",
        fileSize: "50MB",
        status: "Pending",
        fee: "0.75",
    },
    {
        dateTime: "2024-09-10 19:35",
        transactionId: "TX876543",
        sender: "0xlmno...5678",
        receiver: "0xpqrs...9012",
        fileName: "file7.mp3",
        fileSize: "12MB",
        status: "Failed",
        fee: "0.35",
    },
    {
        dateTime: "2024-09-07 09:55",
        transactionId: "TX987654",
        sender: "0xtuvw...1234",
        receiver: "0xyz...5678",
        fileName: "file8.avi",
        fileSize: "60MB",
        status: "Complete",
        fee: "1.00",
    },
    {
        dateTime: "2024-09-05 14:25",
        transactionId: "TX543210",
        sender: "0xqrst...9012",
        receiver: "0xuvwx...7890",
        fileName: "file9.csv",
        fileSize: "25MB",
        status: "Complete",
        fee: "0.45",
    },
    {
        dateTime: "2024-09-02 21:50",
        transactionId: "TX102938",
        sender: "0xabcd...5678",
        receiver: "0xmnop...1234",
        fileName: "file10.zip",
        fileSize: "8MB",
        status: "Pending",
        fee: "0.20",
    },
];


// props interface for the TransactionTable
interface TransactionTableProps {
    search: string;
    dateFilter: string;
    statusFilter: string;
}

// helper function to filter transactions 
const filterTransactions = (
    transactions: Transaction[],
    search: string,
    dateFilter: string,
    statusFilter: string
): Transaction[] => {
    return transactions.filter((transaction) => {
        const searchMatch = transaction.fileName.toLowerCase().includes(search.toLowerCase()) ||
                            transaction.transactionId.toLowerCase().includes(search.toLowerCase());
        
        // date filter
        const date = new Date(transaction.dateTime); // convert dateTime to Date object so we can "filter"/match
        const now = new Date();
        let dateMatch = true;

        if (dateFilter === "today") {
            dateMatch = date.toDateString() === now.toDateString();
        } else if (dateFilter === "week") {
            const weekAgo = new Date(now.setDate(now.getDate() - 7));
            dateMatch = date >= weekAgo;
        } else if (dateFilter === "month") {
            const monthAgo = new Date(now.setMonth(now.getMonth() - 1));
            dateMatch = date >= monthAgo;
        } else if (dateFilter === "year") {
            const yearAgo = new Date(now.setFullYear(now.getFullYear() - 1));
            dateMatch = date >= yearAgo;
        }

        // status filter
        const statusMatch = statusFilter === "all" || transaction.status.toLowerCase() === statusFilter.toLowerCase();

        return searchMatch && dateMatch && statusMatch;
    });
};

const TransactionTable: React.FC<TransactionTableProps> = ({ search, dateFilter, statusFilter }) => {
    // apply filters if necessary 
    const filteredTransactions = filterTransactions(transactions, search, dateFilter, statusFilter);

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
                <TableHead sx={{ backgroundColor: "grey.200" }}> {/* Set the background color here */}
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
                            <TableRow key={transaction.transactionId}>
                                <TableCell>{transaction.dateTime}</TableCell>
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
                count={transactions.length}
                page={page}
                onPageChange={handleChangePage}
                rowsPerPage={rowsPerPage}
                onRowsPerPageChange={handleChangeRowsPerPage}
                rowsPerPageOptions={[5, 10, 25]}
            />
        </TableContainer>
    );
};

export default TransactionTable;
