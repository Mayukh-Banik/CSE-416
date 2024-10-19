import axios from 'axios';
import User from './models/user';
import Transaction from './models/transaction';
import dotenv from 'dotenv';

dotenv.config();

const PORT = process.env.PORT || 3000; 
console.log(process.env.PORT);

// monk data
const monkUsers = [
  {
    name: 'TempUser',
    email: 'tempUser@example.com',
    password: 'password123'
  },
  {
    name: 'TempReceiver',
    email: 'tempReceiver@example.com',
    password: 'password123'
  },
  {
    name: 'TempSender1',
    email: 'sender1@example.com',
    password: 'password123'
  },
  {
    name: 'TempSender2',
    email: 'sender2@example.com',
    password: 'password123'
  },
  {
    name: 'TempSender3',
    email: 'sender3@example.com',
    password: 'password123'
  }
];

const monkTransactions = [
  {
    transactionId: "tx001",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 10,
    timestamp: "2023-10-01T10:00:00",
    fileName: "file1.txt",
    fileId: "file001",
    fileSize: "2MB",
    fee: 1.2,
    status: "completed",
  },
  {
    transactionId: "tx009",
    sender: "TempSender3",
    receiver: "TempUser",
    amount: 90,
    timestamp: "2023-10-09T10:00:00",
    fileName: "file9.zip",
    fileId: "file009",
    fileSize: "5MB",
    fee: 3.5,
    status: "completed",
  },
  {
    transactionId: "tx002",
    sender: "TempSender2",
    receiver: "TempUser",
    amount: 20,
    timestamp: "2023-10-02T10:00:00",
    fileName: "file2.png",
    fileId: "file002",
    fileSize: "1MB",
    fee: 0.8,
    status: "pending",
  },
  {
    transactionId: "tx005",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 50,
    timestamp: "2023-10-05T10:00:00",
    fileName: "file5.pdf",
    fileId: "file005",
    fileSize: "4MB",
    fee: 2.0,
    status: "completed",
  },
  {
    transactionId: "tx011",
    sender: "TempSender2",
    receiver: "TempUser",
    amount: 110,
    timestamp: "2023-10-11T10:00:00",
    fileName: "file11.mp3",
    fileId: "file011",
    fileSize: "7MB",
    fee: 5.0,
    status: "completed",
  },
  {
    transactionId: "tx003",
    sender: "TempSender3",
    receiver: "TempReceiver",
    amount: 30,
    timestamp: "2023-10-03T10:00:00",
    fileName: "file3.docx",
    fileId: "file003",
    fileSize: "3MB",
    fee: 1.5,
    status: "completed",
  },
  {
    transactionId: "tx004",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 40,
    timestamp: "2023-10-04T10:00:00",
    fileName: "file4.jpg",
    fileId: "file004",
    fileSize: "2MB",
    fee: 1.0,
    status: "failed",
  },
  {
    transactionId: "tx006",
    sender: "TempSender2",
    receiver: "TempReceiver",
    amount: 60,
    timestamp: "2023-10-06T10:00:00",
    fileName: "file6.mov",
    fileId: "file006",
    fileSize: "10MB",
    fee: 6.5,
    status: "pending",
  },
  {
    transactionId: "tx007",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 70,
    timestamp: "2023-10-07T10:00:00",
    fileName: "file7.avi",
    fileId: "file007",
    fileSize: "15MB",
    fee: 4.2,
    status: "completed",
  },
  {
    transactionId: "tx013",
    sender: "TempSender2",
    receiver: "TempUser",
    amount: 130,
    timestamp: "2023-10-13T10:00:00",
    fileName: "file13.zip",
    fileId: "file013",
    fileSize: "8MB",
    fee: 7.5,
    status: "failed",
  },
  {
    transactionId: "tx008",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 80,
    timestamp: "2023-10-08T10:00:00",
    fileName: "file8.mkv",
    fileId: "file008",
    fileSize: "12MB",
    fee: 5.1,
    status: "failed",
  },
  {
    transactionId: "tx010",
    sender: "TempSender2",
    receiver: "TempReceiver",
    amount: 100,
    timestamp: "2023-10-10T10:00:00",
    fileName: "file10.wav",
    fileId: "file010",
    fileSize: "9MB",
    fee: 6.0,
    status: "pending",
  },
  {
    transactionId: "tx012",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 120,
    timestamp: "2023-10-12T10:00:00",
    fileName: "file12.pdf",
    fileId: "file012",
    fileSize: "3MB",
    fee: 2.5,
    status: "pending",
  },
  {
    transactionId: "tx018",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 120,
    timestamp: "2023-10-18T10:00:00",
    fileName: "file18.pdf",
    fileId: "file018",
    fileSize: "3MB",
    fee: 2.5,
    status: "pending",
  },
  {
    transactionId: "tx017",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 120,
    timestamp: "2023-10-17T10:00:00",
    fileName: "file17.pdf",
    fileId: "file017",
    fileSize: "3MB",
    fee: 2.5,
    status: "pending",
  },
  {
    transactionId: "tx016",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 120,
    timestamp: "2023-10-16T10:00:00",
    fileName: "file16.pdf",
    fileId: "file016",
    fileSize: "3MB",
    fee: 2.5,
    status: "pending",
  },
  {
    transactionId: "tx015",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 120,
    timestamp: "2023-10-15T10:00:00",
    fileName: "file15.pdf",
    fileId: "file015",
    fileSize: "3MB",
    fee: 2.5,
    status: "pending",
  },
  {
    transactionId: "tx014",
    sender: "TempSender1",
    receiver: "TempUser",
    amount: 120,
    timestamp: "2023-10-14T10:00:00",
    fileName: "file14.pdf",
    fileId: "file014",
    fileSize: "3MB",
    fee: 2.5,
    status: "pending",
  }
];

export const createMonkUsers = async (): Promise<void> => {
  try {
    for (const userData of monkUsers) {
      await axios.post(`http://localhost:${PORT}/api/users/signup`, userData);
      console.log(`User ${userData.name} created successfully`);
    }
  } catch (error) {
    console.error('Error creating monk users:', error);
  }
};

export const createMonkTransactions = async () => {
  try {
    const existingTransactionIds = await Transaction.find(
      { transactionId: { $in: monkTransactions.map((t) => t.transactionId) } },
      { transactionId: 1 }
    ).lean();

    const existingTransactionIdSet = new Set(existingTransactionIds.map((t) => t.transactionId));

    const users = await User.find({
      name: { $in: monkTransactions.flatMap((t) => [t.sender, t.receiver]) },
    });

    const userMap = new Map(users.map((user) => [user.name, user._id]));

    const transactionsToInsert = monkTransactions
      .filter((transaction) => !existingTransactionIdSet.has(transaction.transactionId)) // transactionId로 필터링
      .map((transaction) => {
        const senderId = userMap.get(transaction.sender); 
        const receiverId = userMap.get(transaction.receiver); 

        if (!senderId || !receiverId) {
          console.warn(`Sender or Receiver not found for transaction ${transaction.transactionId}. Skipping...`);
          return null;
        }

        return {
          transactionId: transaction.transactionId,
          sender: senderId, 
          receiver: receiverId, 
          amount: transaction.amount,
          timestamp: transaction.timestamp,
          fileName: transaction.fileName,
          fileId: transaction.fileId, 
          fileSize: transaction.fileSize, 
          fee: transaction.fee,
          status: transaction.status,
        };
      })
      .filter((transaction) => transaction !== null); 

    if (transactionsToInsert.length > 0) {
      await Transaction.insertMany(transactionsToInsert);
      console.log('Temporary transactions inserted successfully');
    } else {
      console.log('No new transactions to insert');
    }
  } catch (error) {
    console.error('Error inserting temporary transactions:', error);
  }
};