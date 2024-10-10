export interface Transaction {
    dateTime: string;
    transactionId: string;
    sender: string;
    receiver: string;
    fileName: string;
    fileSize: string;
    status: "Complete" | "Pending" | "Failed";
    fee: string;
}

export const transactions: Transaction[] = [
    {
        dateTime: '2023-10-01T10:30:00Z',
        transactionId: 'TXN001',
        sender: 'Alice',
        receiver: 'Bob',
        fileName: 'file1.pdf',
        fileSize: '2MB',
        status: 'Complete' as Transaction['status'],  // Explicitly casting status
        fee: '10.50',
    },
    {
        dateTime: '2023-10-01T12:45:00Z',
        transactionId: 'TXN002',
        sender: 'Charlie',
        receiver: 'Dave',
        fileName: 'file2.pdf',
        fileSize: '3MB',
        status: 'Complete' as Transaction['status'],
        fee: '12.00',
    },
    {
        dateTime: '2023-10-02T09:15:00Z',
        transactionId: 'TXN003',
        sender: 'Eve',
        receiver: 'Frank',
        fileName: 'file3.pdf',
        fileSize: '1.5MB',
        status: 'Pending' as Transaction['status'],
        fee: '8.75',
    },
    {
        dateTime: '2023-10-02T14:10:00Z',
        transactionId: 'TXN004',
        sender: 'Grace',
        receiver: 'Heidi',
        fileName: 'file4.pdf',
        fileSize: '5MB',
        status: 'Failed' as Transaction['status'],
        fee: '7.30',
    },
    {
        dateTime: '2023-10-03T11:25:00Z',
        transactionId: 'TXN005',
        sender: 'Ivan',
        receiver: 'Judy',
        fileName: 'file5.pdf',
        fileSize: '6MB',
        status: 'Complete' as Transaction['status'],
        fee: '15.00',
    }
];

export default transactions;
