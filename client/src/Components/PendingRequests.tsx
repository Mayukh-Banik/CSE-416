import { useEffect, useState } from 'react';
import { Transaction } from '../models/transactions';

function PendingRequests() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8081/ws");

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.status === "pending") {
        setTransactions((prevTransactions) => [...prevTransactions, data]);
      }
    };

    return () => socket.close();
  }, []);

  const handleResponse = (fileHash: string, response: string) => {
    fetch(`/api/handleRequest`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ fileHash, response }),
    }).then(() => {
      // Update state to remove the transaction after responding
      setTransactions((prevTransactions) =>
        prevTransactions.filter((transaction) => transaction.FileHash !== fileHash)
      );
    }).catch((error) => {
      console.error("Error handling request:", error);
    });
  };

  return (
    <div>
      {transactions.map((transaction) => (
        <div key={transaction.FileHash}>
          <p>Download request from {transaction.RequesterID} for file {transaction.FileHash}</p>
          <button onClick={() => handleResponse(transaction.FileHash, "accepted")}>Accept</button>
          <button onClick={() => handleResponse(transaction.FileHash, "declined")}>Decline</button>
        </div>
      ))}
    </div>
  );
}

export default PendingRequests;
