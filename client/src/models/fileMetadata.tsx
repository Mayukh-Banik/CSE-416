export interface Provider {
    Fee: number;
    PeerID: string;
    IsActive: boolean;
  }
  
export interface dhtFile{
  name: string;
  type: string;
  size: number;
  description: string;
  hash: string;
  reputation?: number;
  isPublished?: boolean;
  providers: Provider[];
}

export interface FileMetadata {
  //id: string;
  Name: string;
  Type: string;
  Size: number;
  Description: string;
  Hash: string;
  CreatedAt?: string;
  Reputation?: number;
  IsPublished?: boolean;
  Fee: number;
  Path: string;
}