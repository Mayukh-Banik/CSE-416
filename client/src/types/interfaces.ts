export interface TransactionProps {
  _id: string;
  transactionId: string;
  sender: string;
  receiver: string;
  amount: number;
  timestamp: string;
  fileName: string;
  fileId: string;
  fileSize: string;
  fee: number;
  status: "pending" | "completed" | "failed";
}
export interface WalletDetailsProps {
  userId: string;
  walletAddress: string;
  balance: number;
  transactions: TransactionProps[];
  walletLabel?: string;
  fee?: number;
}

export interface TransactionTableProps {
  transactions: TransactionProps[];
  search: string;
  dateFilter: string;
  statusFilter: string;
}