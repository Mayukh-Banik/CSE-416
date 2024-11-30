import React, { useEffect, useState } from "react";
import Sidebar from "./Sidebar";
import PublishIcon from '@mui/icons-material/Publish';
import DeleteIcon from '@mui/icons-material/Delete';
import { 
  Box, 
  Button, 
  Typography, 
  List, 
  ListItem, 
  Switch,
  IconButton,
  Snackbar,
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  ListItemText,
  LinearProgress,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tooltip,
} from "@mui/material";
import { useTheme } from '@mui/material/styles';
import { FileMetadata } from "../models/fileMetadata";
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward';
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward';

declare global {
  interface Window {
      electron: {
          ipcRenderer: typeof import('electron').ipcRenderer;
          saveFile: (fileData: { fileName: string, fileData: Buffer }) => Promise<{ success: boolean, message: string }>;
      };
  }
}
// import { saveFileMetadata, getFilesForUser, deleteFileMetadata, updateFileMetadata, FileMetadata } from '../utils/localStorage'

// user id/wallet id is hardcoded for now

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const ipcRenderer = window.electron?.ipcRenderer;

interface FilesProp {
  uploadedFiles: FileMetadata[];
  setUploadedFiles: React.Dispatch<React.SetStateAction<FileMetadata[]>>;
  initialFetch: boolean;
  setInitialFetch: React.Dispatch<React.SetStateAction<boolean>>;
}

const FilesPage: React.FC<FilesProp> = ({uploadedFiles, setUploadedFiles, initialFetch, setInitialFetch}) => {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [descriptions, setDescriptions] = useState<{ [key: string]: string }>({}); // Track descriptions
  const [fileHashes, setFileHashes] = useState<{ [key: string]: string }>({}); // Track hashes
  const [downloadedFiles, setDownloadedFiles] = useState<FileMetadata[]>([]);
  const [notification, setNotification] = useState<{ open: boolean; message: string; severity: "success" | "error" }>({ open: false, message: "", severity: "success" });
  const [publishDialogOpen, setPublishDialogOpen] = useState(false); // Control for the modal
  const [currentFileHash, setCurrentFileHash] = useState<string | null>(null); // Track the file being published
  const [fees, setFees] = useState<{ [key: string]: number }>({}); // Track fees of to be uploaded files
  const [names, setNames] = useState<{ [key: string]: string }>({}); // Track names of to be uploaded files
  const theme = useTheme();
  const [loading, setLoading] = useState(false); // Loading state for file upload

  useEffect(() => {
    const fetchAllFiles = async() => {
      await fetchFiles("uploaded");
      await fetchFiles("downloaded");
      localStorage.setItem('filesFetched','true')
    };
    fetchAllFiles(); 
  }, [])

  const fetchFiles = async (fileType: string) => {
    try {
      console.log(`Getting local user's ${fileType} files`);
      const response = await fetch(`http://localhost:8081/files/fetch?file=${fileType}`, {
        method: "GET",
      });
      if (!response.ok) throw new Error(`Failed to load ${fileType} file data`);

      const data = await response.json();
      console.log("Fetched data", data);

      if (fileType === "uploaded") {
        setUploadedFiles(data); // Set the state with the loaded data
      } else {
        setDownloadedFiles(data)
      }
    } catch (error) {
      console.error("Error fetching uploaded files:", error);
    }
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files) {
      const fileArray = Array.from(files);
      setSelectedFiles(prevSelectedFiles => [...prevSelectedFiles, ...fileArray]);
      // Compute hashes for each file
      fileArray.forEach(file => computeSHA256(file));
    }
  };

  const handleUpload = async () => {
    if (selectedFiles.length === 0) return;
    setLoading(true);
  
    try {
      // Prepare individual upload tasks
      const uploadTasks = selectedFiles.map((file) => uploadFile(file));
  
      // Process all uploads and get results
      const uploadResults = await Promise.allSettled(uploadTasks);
  
      const successfulUploads = uploadResults
        .filter((result) => result.status === "fulfilled")
        .map((result) => (result as PromiseFulfilledResult<FileMetadata>).value);
  
      const failedUploads = uploadResults
        .filter((result) => result.status === "rejected")
        .map((result) => (result as PromiseRejectedResult).reason);
  
      // Update the UI based on results
      if (successfulUploads.length > 0) {
        setUploadedFiles((prev) => [...prev, ...successfulUploads]);
        setNotification({
          open: true,
          message: `${successfulUploads.length} file(s) uploaded successfully!`,
          severity: "success",
        });
      }
  
      if (failedUploads.length > 0) {
        console.error("Some files failed to upload:", failedUploads);
        setNotification({
          open: true,
          message: `${failedUploads.length} file(s) failed to upload.`,
          severity: "error",
        });
      }
    } catch (error) {
      console.error("Error during file upload:", error);
      setNotification({ open: true, message: "Failed to upload files.", severity: "error" });
    } finally {
      setLoading(false);
      setSelectedFiles([]);
      setDescriptions({});
      setFees({});
      setFileHashes({});
    }
  };
  
  // Function to handle individual file upload
  const uploadFile = async (file: File) => {
    try {
      console.log("Processing file:", file.name);
  
      // Extract the file extension and base name
      const fileParts = file.name.split(".");
      const fileExtension = fileParts.pop();
      const baseName = fileParts.join(".");
  
      // Read file data
      const arrayBuffer = await file.arrayBuffer();
      const fileData = Buffer.from(arrayBuffer);
  
      // Create metadata object
      const metadata: FileMetadata = {
        Name: names[file.name] || baseName,
        Type: file.type,
        Size: file.size,
        Description: descriptions[file.name] || "",
        Hash: fileHashes[file.name],
        IsPublished: true,
        Fee: fees[file.name] || 0,
        OriginalUploader: true,
        NameWithExtension: fileExtension,
        Rating: "",
      };
  
      // Save file locally using Electron API
      const saveResponse = await window.electron.saveFile({
        fileName: `${metadata.Name}.${fileExtension}`,
        fileData,
      });

      metadata.NameWithExtension = `${metadata.Name}.${fileExtension}`
  
      if (!saveResponse.success) {
        throw new Error(`Failed to save file locally: ${file.name}`);
      }
  
      // Send metadata to backend
      const response = await fetch("http://localhost:8081/files/upload", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(metadata),
      });
  
      if (!response.ok) {
        throw new Error(`Failed to upload metadata for: ${file.name}`);
      }
  
      console.log("File metadata uploaded:", file.name);
      return metadata; // Return metadata for successful uploads
    } catch (error) {
      console.error("Error processing file:", file.name, error);
      throw error;
    }
  };
  
  const handleDescriptionChange = (fileId: string, description: string) => {
    if (!loading) {
      setDescriptions((prev) => ({ ...prev, [fileId]: description }));
    }
  };

  const handleFeeChange = (fileId: string, fee: number) => {
    setFees((prev) => ({ ...prev, [fileId]: fee }));
  };

  const handleNameChange = (fileId: string, name: string) => {
    setNames((prev) => ({ ...prev, [fileId]: name }));
  };

// have to fix deleting file
  const handleDeleteUploadedFile = async (selectedFile: FileMetadata) => {
    console.log("attempting to delete file ", selectedFile.Name)
    try {
      const response = await fetch(`http://localhost:8081/files/delete?hash=${selectedFile.Hash}&originalUploader=${selectedFile.OriginalUploader}&name=${selectedFile.NameWithExtension}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error("Failed to delete file")
      }
      const data = await response.json();
      console.log('file deleted successfully', data);

      setUploadedFiles((prev) => prev.filter((file) => file.Hash !== selectedFile.Hash));
      setDownloadedFiles((prev) => prev.filter((file) => file.Hash !== selectedFile.Hash));
      setSelectedFiles((prev) => prev.filter((file) => file.name !== selectedFile.Hash));
      setNotification({ open: true, message: "File deleted.", severity: "success" });
    } catch (error) {
      console.error("error: ", error);
      setNotification({ open: true, message: "failed to delete file.", severity: "error" });
    }
  }

  const handleCloseNotification = () => {
    setNotification({ ...notification, open: false });
  };

  const handleUploadClick = () => {
    document.getElementById("file-input")?.click();
  };

  const computeSHA256 = async (file: File) => {
    const arrayBuffer = await file.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest("SHA-256", arrayBuffer);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray
      .map(byte => byte.toString(16).padStart(2, "0"))
      .join("");
    
    setFileHashes(prevHashes => ({
      ...prevHashes,
      [file.name]: hashHex,
    }));
  };

  const handleConfirmPublish = async (hash: string, files: FileMetadata[]) => {
    const fileToPublish = files.find(file => file.Hash === hash);
    
    if (!fileToPublish) {
        setNotification({ open: true, message: "File not found", severity: "error" });
        return;
    }
    console.log("old metadata: ", fileToPublish)

    const updatedMetadata = {
      ...fileToPublish,
      IsPublished: !fileToPublish.IsPublished,
    };

    console.log("updated metadata: ", updatedMetadata)

    try {
        const response = await fetch("http://localhost:8081/files/upload", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(updatedMetadata),
        });

        if (response.ok) {
            // Update the file's published status locally in the UI
            setUploadedFiles(prevFiles => 
              prevFiles.map(currentFile =>
                currentFile.Hash === hash
                  ? { ...currentFile, IsPublished: !currentFile.IsPublished }
                  : currentFile
              )
            );

            setDownloadedFiles(prevFiles => 
              prevFiles.map(currentFile =>
                currentFile.Hash === hash
                  ? { ...currentFile, IsPublished: !currentFile.IsPublished }
                  : currentFile
              )
            );

            let message;
            if (updatedMetadata.IsPublished) {
              message = "File published successfully!"
            } else {
              message = "File unpublished successfully!"
            }

            setNotification({ open: true, message: message, severity: "success" });

            const data = await response.text();
            console.log("Publish response: ", data);
        } else {
            const errorData = await response.text();
            console.error("Publish response error:", errorData);
            setNotification({ open: true, message: "Failed to change publish status", severity: "error" });
        }
    } catch (error) {
        console.error("Error publishing file:", error);
        setNotification({ open: true, message: "An error occurred", severity: "error" });
    } finally {
        console.log("Published file:", updatedMetadata);
        // setPublishDialogOpen(false);
    }
};

  const handleDeleteSelectedFile = (hash: string) => {
    setSelectedFiles((prev) => prev.filter((file) => file.name !== hash));
  }

  const handleVote = async (fileHash: string, voteType: string ) => {
    try {
      const response = await fetch(`http://localhost:8081/files/vote?fileHash=${fileHash}&voteType=${voteType}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });

      if (!response.ok) {
        throw new Error("Failed to rate file")
      }
      const data = await response.json();
      console.log('file rated successfully', data);

      setNotification({ open: true, message: "Successfully rated file.", severity: "success" });
    } catch (error) {
      console.error("error: ", error);
      setNotification({ open: true, message: "Unable to rate file.", severity: "error" });
    }
  }

  /*
  will modularize later 
  its a mess rn ...
  */
 
  return (
      <Box
        sx={{
          padding: 2,
          marginTop: '70px',
          marginLeft: `${drawerWidth}px`, 
          transition: 'margin-left 0.3s ease',
          [theme.breakpoints.down('sm')]: {
            marginLeft: `${collapsedDrawerWidth}px`,
          },
        }}
      >
      <Sidebar />
      <Box sx={{ flexGrow: 1}}>
        <Typography variant="h4" gutterBottom>
          My Files
        </Typography>
        <input
          type="file"
          id="file-input"
          multiple
          onChange={handleFileChange}
          style={{ display: 'none' }}
        />
        <Box 
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            width: 200,
            height: 200,
            border: '2px dashed #3f51b5',
            borderRadius: 2,
            cursor: 'pointer',
            marginBottom: 2,
            background: 'white',
            '&:hover': {
              backgroundColor: '#e3f2fd',
            },
          }}
          onClick={handleUploadClick} 
        >
          <PublishIcon sx={{ fontSize: 50 }} />
          <Typography variant="h6" sx={{ marginTop: 1 }}>
            Select Files
          </Typography>
        </Box>
        <Button 
          variant="contained" 
          onClick={handleUpload} 
          disabled={selectedFiles.length === 0}
        >
          Upload Selected
        </Button>
        
        {loading && <LinearProgress sx={{ width: '100%', marginTop: 2 }} />} {/* Progress bar when loading is true */}

        {selectedFiles.length > 0 && (
          <Box sx={{ marginTop: 2 }}>
            <Typography variant="h6">Selected Files:</Typography>
            <TableContainer component={Paper} sx={{ marginTop: 2 }}>
              <Table>
                {/* Table Header */}
                <TableHead>
                  <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>File Size</TableCell>
                    <TableCell>Description</TableCell>
                    <TableCell>Fee</TableCell>
                    <TableCell>Actions</TableCell>
                  </TableRow>
                </TableHead>
                {/* Table Body */}
                <TableBody>
                  {selectedFiles.map((file, index) => (
                    <TableRow key={index}>
                      {/* Editable File Name */}
                      <TableCell sx={{ width: '25%' }}>
                        <TextField
                          label="File Name"
                          variant="outlined"
                          fullWidth
                          value={names[file.name] || file.name} // Use fileNames state or the original name
                          onChange={(e) => handleNameChange(file.name, e.target.value)}
                        />
                      </TableCell>
                      {/* File Size */}
                      <TableCell>{`${(file.size / 1024).toFixed(2)} KB`}</TableCell>
                      {/* Description */}
                      <TableCell sx={{ width: '40%' }}>
                        <TextField
                          label="Description"
                          variant="outlined"
                          fullWidth
                          multiline
                          rows={2} // Makes the description box larger
                          value={descriptions[file.name] || ""}
                          onChange={(e) => handleDescriptionChange(file.name, e.target.value)}
                        />
                      </TableCell>
                      {/* File Fee */}
                      <TableCell>
                        <TextField
                          type="number"
                          variant="outlined"
                          size="small"
                          value={fees[file.name] || 0}
                          onChange={(e) => handleFeeChange(file.name, parseFloat(e.target.value))}
                        />
                      </TableCell>
                      {/* Actions */}
                      <TableCell>
                        <IconButton
                          edge="end"
                          aria-label="delete"
                          onClick={() => handleDeleteSelectedFile(file.name)}
                        >
                          <DeleteIcon />
                        </IconButton>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </Box>
        )}

        {uploadedFiles.length > 0 && (
          <Box sx={{ marginTop: 2 }}>
            <Typography variant="h6">Uploaded Files</Typography>
            <TableContainer component={Paper} sx={{ marginTop: 2 }}>
              <Table>
                {/* Table Header */}
                <TableHead>
                  <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>Hash</TableCell>
                    <TableCell>Size</TableCell>
                    <TableCell>Description</TableCell>
                    <TableCell>Fee</TableCell>
                    <TableCell>Actions</TableCell>
                  </TableRow>
                </TableHead>
                {/* Table Body */}
                <TableBody>
                  {uploadedFiles.map((file, index) => (
                    <TableRow key={index}>
                      <TableCell>{file.Name}</TableCell>
                      <TableCell
                        sx={{
                          maxWidth: '200px', // Set a max width to control the wrapping
                          whiteSpace: 'normal', // Allows the text to break lines
                          wordWrap: 'break-word', // Breaks words if they're too long
                          overflowWrap: 'break-word', // Fallback for compatibility
                        }}
                      >
                        {file.Hash}
                      </TableCell>
                      <TableCell>{`${(file.Size / 1024).toFixed(2)} KB`}</TableCell>
                      <TableCell>{file.Description}</TableCell>
                      <TableCell>{file.Fee}</TableCell>
                      {/* Actions */}
                      <TableCell>
                        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                          <Tooltip title="Publish File" arrow>
                            <Switch
                              edge="end"
                              onChange={() => handleConfirmPublish(file.Hash, uploadedFiles)}
                              checked={file.IsPublished}
                              color="primary"
                              inputProps={{ 'aria-label': `publish ${file.Name}` }}
                            />
                          </Tooltip>
                          <IconButton
                            edge="end"
                            aria-label="delete"
                            onClick={() => handleDeleteUploadedFile(file)}
                          >
                            <DeleteIcon />
                          </IconButton>
                        </Box>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </Box>
        )}

        {downloadedFiles.length > 0 && (
          <Box sx={{ marginTop: 2 }}>
            <Typography variant="h6">Downloaded Files</Typography>
            <TableContainer component={Paper} sx={{ marginTop: 2 }}>
              <Table>
                {/* Table Header */}
                <TableHead>
                  <TableRow>
                    <TableCell></TableCell>
                    <TableCell>Name</TableCell>
                    <TableCell>Hash</TableCell>
                    <TableCell>Size</TableCell>
                    <TableCell>Description</TableCell>
                    <TableCell>Fee</TableCell>
                    <TableCell>Actions</TableCell>
                  </TableRow>
                </TableHead>
                {/* Table Body */}
                <TableBody>
                  {downloadedFiles.map((file, index) => (
                    <TableRow key={index}>
                      <TableCell>
                        <IconButton color="success" onClick={() => handleVote(file.Hash, 'upvote')}>
                          <ArrowUpwardIcon />
                        </IconButton>
                        <IconButton color="error" onClick={() => handleVote(file.Hash, 'downvote')}>
                          <ArrowDownwardIcon />
                        </IconButton>
                      </TableCell>

                      <TableCell>{file.Name}</TableCell>
                      <TableCell
                        sx={{
                          maxWidth: '200px', // Set a max width to control the wrapping
                          whiteSpace: 'normal', // Allows the text to break lines
                          wordWrap: 'break-word', // Breaks words if they're too long
                          overflowWrap: 'break-word', // Fallback for compatibility
                        }}
                      >
                        {file.Hash}
                      </TableCell>
                      <TableCell>{`${(file.Size / 1024).toFixed(2)} KB`}</TableCell>
                      <TableCell>{file.Description}</TableCell>
                      <TableCell>{file.Fee}</TableCell>
                      {/* Actions */}
                      <TableCell>
                        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                          <Tooltip title="Publish File" arrow>
                            <Switch
                              edge="end"
                              onChange={() => handleConfirmPublish(file.Hash, downloadedFiles)}
                              checked={file.IsPublished}
                              color="primary"
                              inputProps={{ 'aria-label': `publish ${file.Name}` }}
                            />
                          </Tooltip>
                          <IconButton
                            edge="end"
                            aria-label="delete"
                            onClick={() => handleDeleteUploadedFile(file)}
                          >
                            <DeleteIcon />
                          </IconButton>
                        </Box>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </Box>
        )}

      </Box>

      

      {/* Publish Modal */}
      {/* <Dialog open={publishDialogOpen} onClose={() => setPublishDialogOpen(false)}>
        <DialogTitle>Set Download Fee</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Enter Fee (in Squid Coins)"
            type="number"
            fullWidth
            variant="outlined"
            value={fee}

            onChange={(e) => setFee(Number(e.target.value))}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setPublishDialogOpen(false)} color="secondary">
            Cancel
          </Button>
          <Button onClick={handleConfirmPublish} color="primary">
            Publish
          </Button>
        </DialogActions>
      </Dialog> */}

      {/* Notification Snackbar */}
      <Snackbar 
        open={notification.open} 
        autoHideDuration={6000} 
        onClose={handleCloseNotification}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={handleCloseNotification} severity={notification.severity} sx={{ width: '100%' }}>
          {notification.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default FilesPage;