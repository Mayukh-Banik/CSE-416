// TransactionTable.tsx
import React from "react";
import {
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
} from "@mui/material";

// not final - just here to display transaction page details
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
];

const TransactionTable: React.FC = () => (
    <TableContainer>
        <Table>
        <TableHead>
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
            {transactions.map((transaction) => (
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
            ))}
        </TableBody>
        </Table>
    </TableContainer>
);

export default TransactionTable;
