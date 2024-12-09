export interface Provider {
    Fee: number;
    PeerAddr: string;
    IsActive: boolean;
    Rating: number;
  }
  
export interface DHTMetadata{
  Name:             string
	NameWithExtension: string
	Type:              string
	Size:              number
	Description:       string
	CreatedAt:         string
	Rating:            number               //
  Providers: { [key: string]: Provider }; // Use PeerID as the key
	NumRaters:         number
	Upvote:            number
	Downvote:          number
	Hash:              string
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
  Rating?: number;
  HasVoted: boolean
}