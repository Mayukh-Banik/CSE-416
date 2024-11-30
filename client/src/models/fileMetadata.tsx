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
  rating?: number;
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
  IsPublished?: boolean;
  Fee: number;
  OriginalUploader: boolean; // true is user acquire file by downloading, false if user themselves uploaded file
  NameWithExtension?: string;
  Rating?: string;
}