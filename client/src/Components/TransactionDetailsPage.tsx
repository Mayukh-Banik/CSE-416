import React from "react";
import { useParams, Link } from "react-router-dom";
import { Box, Typography, Button } from "@mui/material";
import { Container } from "@mui/material";

interface Transaction {
  id: string;
  sender: string;
  receiver: string;
  amount: number;
  timestamp: string;
  status: string;
  file: string;
}

const transactions: Transaction[] = [
  {
    id: "tx001",
    sender: "0xsender001",
    receiver: "0xreceiver001",
    amount: 10,
    timestamp: "2023-10-01T10:00:00",
    status: "completed",
    file: "file001.pdf"
  },
  {
    id: "tx009",
    sender: "0xsender009",
    receiver: "0xreceiver009",
    amount: 90,
    timestamp: "2023-10-09T17:20:00",
    status: "completed",
    file: "file009.docx"
  }
  // Add other transactions...
];

const TransactionDetailsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const transaction = transactions.find((tx) => tx.id === id);

  if (!transaction) {
    return (
      <Container>
        <Typography variant="h5">Transaction not found</Typography>
      </Container>
    );
  }

  return (
    <Container>
      <Typography variant="h4" gutterBottom>
        Transaction Details
      </Typography>
      <Box>
        <Typography variant="body1">Transaction ID: {transaction.id}</Typography>
        <Typography variant="body1">Sender: {transaction.sender}</Typography>
        <Typography variant="body1">Receiver: {transaction.receiver}</Typography>
        <Typography variant="body1">Amount: {transaction.amount}</Typography>
        <Typography variant="body1">Timestamp: {transaction.timestamp}</Typography>
        <Typography variant="body1">Status: {transaction.status}</Typography>
        <Typography variant="body1">File: {transaction.file}</Typography>
        <Button
          variant="contained"
          component={Link}
          to={`/fileview/${transaction.file}`}
        >
          View File
        </Button>
      </Box>
    </Container>
  );
};

export default TransactionDetailsPage;
