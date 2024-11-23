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
//import { FileMetadata } from '../../local_server/models/file';

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

  // initial fetch of uploaded data
  useEffect(() => {
    if (!initialFetch) {
      const fetchFiles = async () => {
        try {
          console.log("Getting local user's uploaded files");
          let fileType = "uploaded"
          const response = await fetch(`http://localhost:8081/files/fetch?file=${fileType}`, {
            method: "GET",
          });
          if (!response.ok) throw new Error("Failed to load uploaded file data");
  
          const data = await response.json();
          console.log("Fetched data", data);
  
          setUploadedFiles(data); // Set the state with the loaded data
          setInitialFetch(true); // Set initialFetch to true to prevent further calls
        } catch (error) {
          console.error("Error fetching uploaded files:", error);
        }
      };
  
      fetchFiles();
    }
  }, [initialFetch, setUploadedFiles, setInitialFetch]);
  
  // called every time page loads to refresh downloaded files
  useEffect(()=> {
    const fetchDownloadedFiles = async () => {
      try {
        console.log("Getting local user's downloaded files");
        let fileType = "downloaded"
        const response = await fetch(`http://localhost:8081/files/fetch?file=${fileType}`, {
          method: "GET",
        });
        if (!response.ok) throw new Error("Failed to load downloaded file data");

        const data = await response.json();
        console.log("Fetched data", data);

        setDownloadedFiles(data); // Set the state with the loaded data
      } catch (error) {
        console.error("Error fetching downloaded files:", error);
      }
    }; 
    fetchDownloadedFiles();
  }, [])

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
      // Process files and upload metadata
      const uploadResults = await Promise.allSettled(
          selectedFiles.map(async (file) => {
              try {
                console.log("Processing file:", file.name);

                // Read file data
                const arrayBuffer = await file.arrayBuffer();
                const fileData = Buffer.from(arrayBuffer); // Convert to Buffer for backend
                console.log("File read successfully:", file.name);

                // Save file locally using Electron API
                const saveResponse = await window.electron.saveFile({
                    fileName: file.name,
                    fileData,
                });

                if (!saveResponse.success) {
                    throw new Error(`Failed to save file: ${file.name}`);
                }
                console.log("File saved locally:", file.name);

                // Create metadata object
                const metadata: FileMetadata = {
                    Name: names[file.name] || file.name,
                    Type: file.type,
                    Size: file.size,
                    Description: descriptions[file.name] || "",
                    Hash: fileHashes[file.name],
                    IsPublished: true,
                    Fee: fees[file.name] || 0,
                    OriginalUploader: true,
                };

                // Send metadata to the backend
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
                throw error; // Allow Promise.allSettled to catch the error
              }
          })
      );

      // Handle successful and failed uploads
      const successfulUploads = uploadResults
        .filter((result) => result.status === "fulfilled")
        .map((result) => (result as PromiseFulfilledResult<FileMetadata>).value);

      const failedUploads = uploadResults
        .filter((result) => result.status === "rejected")
        .map((result) => (result as PromiseRejectedResult).reason);

      if (successfulUploads.length > 0) {
        // Update uploaded files state
        setUploadedFiles((prev) => [...prev, ...successfulUploads]);
        setNotification({ open: true, message: "Files uploaded successfully!", severity: "success" });
      }

      if (failedUploads.length > 0) {
        console.warn("Some files failed to upload:", failedUploads);
        setNotification({
          open: true,
          message: `${failedUploads.length} file(s) failed to upload.`,
          severity: "error",
        });
      }

      // Clear form data for selected files
      setSelectedFiles([]);
      setDescriptions({});
      setFees({});
      setFileHashes({});
    } catch (error) {
        console.error("Error during file upload:", error);
        setNotification({ open: true, message: "Failed to upload files.", severity: "error" });
    } finally {
        setLoading(false);
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
  const handleDeleteUploadedFile = async (hash: string, originalUploader: boolean) => {
    try {
      const response = await fetch(`http://localhost:8081/files/delete?hash=${hash}&originalUploader=${originalUploader}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error("Failed to delete file")
      }
      const data = await response.json();
      console.log('file deleted successfully', data);

      setUploadedFiles((prev) => prev.filter((file) => file.Hash !== hash));
      setSelectedFiles((prev) => prev.filter((file) => file.name !== hash));
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

  const handleConfirmPublish = async (hash: string) => {
    const fileToPublish = uploadedFiles.find(file => file.Hash === hash);
    
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
                              onChange={() => handleConfirmPublish(file.Hash)}
                              checked={file.IsPublished}
                              color="primary"
                              inputProps={{ 'aria-label': `publish ${file.Name}` }}
                            />
                          </Tooltip>
                          <IconButton
                            edge="end"
                            aria-label="delete"
                            onClick={() => handleDeleteUploadedFile(file.Hash, file.OriginalUploader)}
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