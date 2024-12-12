export interface Transaction {
    CreatedAt?: string;
    FileName: string | "";
    FileHash: string;
    RequesterID: string | "";
    RequesterWallet?: string | "";
    TargetID: string;
    TargetWallet?: string;
    Status?: 'pending' | 'accepted' | 'declined' | 'completed';
    Fee:  number;
    Size: number;
    TransactionID: string;
}
