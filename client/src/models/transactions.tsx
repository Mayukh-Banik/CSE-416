export interface Transaction {
    CreatedAt?: string;
    FileName: string;
    FileHash: string;
    RequesterID: string | "";
    TargetID: string;
    Status?: 'pending' | 'accepted' | 'declined' | 'completed';
    Fee?:  number;
}

export interface PendingRequest {
    CreatedAt: string;
    FileName: string;
    FileHash: string;
    RequesterID: string;
    Fee:  number;
}