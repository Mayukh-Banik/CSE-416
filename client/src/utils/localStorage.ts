export interface FileMetadata {
  //id: string;
  name: string;
  type: string;
  size: number;
  description: string;
  hash: string;
  //createdAt?: string;
  //reputation?: number;
  isPublished?: boolean;
  //fee?: number;
}

// userID is the same as public wallet ID
// Function to save file metadata to local storage
export const saveFileMetadata = (userId: string, file: FileMetadata) => {
  const existingFiles = localStorage.getItem(userId);
  const filesArray = existingFiles ? JSON.parse(existingFiles) : [];
  
  filesArray.push(file);
  localStorage.setItem(userId, JSON.stringify(filesArray));
  console.log('Successfully saved file to local database');
};

// Function to update file metadata
export const updateFileMetadata = (userId: string, updatedFile: FileMetadata) => {
  const existingFiles = localStorage.getItem(userId);
  if (!existingFiles) {
    console.log('No files found for the user');
    return;
  }

  const filesArray: FileMetadata[] = JSON.parse(existingFiles);
  const fileIndex = filesArray.findIndex(file => file.hash === updatedFile.hash);

  if (fileIndex !== -1) {
    filesArray[fileIndex] = { ...filesArray[fileIndex], ...updatedFile }; // Update the existing file metadata
    localStorage.setItem(userId, JSON.stringify(filesArray));
    console.log('Successfully updated file in local database');
  } else {
    console.log('File not found for update');
  }
};

// Function to delete file metadata
export const deleteFileMetadata = (userId: string, fileHash: string) => {
  const existingFiles = localStorage.getItem(userId);
  if (!existingFiles) {
    console.log('No files found for the user');
    return;
  }

  const filesArray: FileMetadata[] = JSON.parse(existingFiles);
  const newFilesArray = filesArray.filter(file => file.hash !== fileHash);

  if (newFilesArray.length !== filesArray.length) {
    localStorage.setItem(userId, JSON.stringify(newFilesArray));
    console.log('Successfully deleted file from local database');
  } else {
    console.log('File not found for deletion');
  }
};

// Function to retrieve all files for a user from local storage
export const getFilesForUser = (userId: string): FileMetadata[] => {
  const existingFiles = localStorage.getItem(userId);
  return existingFiles ? JSON.parse(existingFiles) : [];
};
