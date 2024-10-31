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
} from "@mui/material";
import { useTheme } from '@mui/material/styles';
import { saveFileMetadata, getFilesForUser, deleteFileMetadata, FileMetadata } from '../utils/localStorage'

// user id/wallet id is hardcoded for now

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const FilesPage: React.FC = () => {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [descriptions, setDescriptions] = useState<{ [key: string]: string }>({}); // Track descriptions
  const [fileHashes, setFileHashes] = useState<{ [key: string]: string }>({}); // Track hashes
  const [uploadedFiles, setUploadedFiles] = useState<FileMetadata[]>([]);
  const [notification, setNotification] = useState<{ open: boolean; message: string; severity: "success" | "error" }>({ open: false, message: "", severity: "success" });
  const [publishDialogOpen, setPublishDialogOpen] = useState(false); // Control for the modal
  const [currentFileId, setCurrentFileId] = useState<string | null>(null); // Track the file being published
  const [fee, setFee] = useState<number | undefined>(undefined); // Fee value for publishing
  
  const theme = useTheme();
  
  useEffect(() => {
    const fetchUploadedFiles = () => {
      const files = getFilesForUser('123'); // Adjust user ID as necessary
      setUploadedFiles(files);
    };
    fetchUploadedFiles();
  }, []);

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

    try {
        // Create uploaded file objects with descriptions, hashes, and metadata
        const newUploadedFiles = await Promise.all(
            selectedFiles.map(async (file) => {
                const fileData = await file.arrayBuffer(); // Read file as ArrayBuffer
                // const base64FileData = btoa(String.fromCharCode(...new Uint8Array(fileData))); // Convert to Base64

                let metadata : FileMetadata = {
                    id: `${file.name}-${file.size}-${Date.now()}`, // Unique ID for the uploaded file
                    name: file.name,
                    type: file.type,
                    size: file.size,
                    // file_data: base64FileData, // Encode file data as Base64 if required
                    description: descriptions[file.name] || "",
                    hash: fileHashes[file.name], // Store the computed hash
                    isPublished: false, // Initially not published
                };

                // Send the metadata to the server
                const response = await fetch("http://localhost:8080/upload", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(metadata),
                });

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const data = await response.json();
                console.log("File upload successful:", data, metadata);
                
                // update local server/database
                saveFileMetadata('123', metadata);
                
                return metadata;
            })
        );

        // Update uploadedFiles state
        setUploadedFiles((prev) => [...prev, ...newUploadedFiles]);

        // Clear selected files, descriptions, and hashes after successful upload
        setSelectedFiles([]);
        setDescriptions({});
        setFileHashes({});

        // Show success notification
        setNotification({ open: true, message: "Files uploaded successfully!", severity: "success" });
    } catch (error) {
        console.error("Error uploading files:", error);

        // Show error notification
        setNotification({ open: true, message: "Failed to upload files.", severity: "error" });
    }
};

  const handleDescriptionChange = (fileId: string, description: string) => {
    setDescriptions((prev) => ({ ...prev, [fileId]: description }));
  };

  const handleTogglePublished = (id: string) => {
    const file = uploadedFiles.find(file => file.id === id); // Get the current file

    if (file?.isPublished) {
        // If the file is already published, set the unpublished state
        setUploadedFiles((prev) =>
            prev.map((f) =>
                f.id === id ? { ...f, isPublished: false } : f
            )
        );
    } else {
        // If the file is not published, open the publish dialog
        setCurrentFileId(id); // Set the current file ID to the one being published
        setPublishDialogOpen(true); 
    }
  };


  const handleDeleteFile = (id: string) => {
    deleteFileMetadata('123', id);
    setUploadedFiles((prev) => prev.filter((file) => file.id !== id));
    setSelectedFiles((prev) => prev.filter((file) => file.name !== id));
    setNotification({ open: true, message: "File deleted.", severity: "success" });
  };

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

  const handleConfirmPublish = async () => {
    const fileToPublish = uploadedFiles.find(file => file.id === currentFileId);
    
    if (!fileToPublish) {
        setNotification({ open: true, message: "File not found", severity: "error" });
        return;
    }

    const metadata = {
        name: fileToPublish.name,
        type: fileToPublish.type,
        size: fileToPublish.size,
        description: fileToPublish.description,
        hash: fileToPublish.hash,
    };

    console.log("Publishing file metadata:", metadata);

    try {
        const response = await fetch("http://localhost:8080/publish", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(metadata),
        });

        if (response.ok) {
            setUploadedFiles((prev) =>
                prev.map((file) =>
                    file.id === currentFileId ? { ...file, isPublished: true, fee } : file
                )
            );
            setNotification({ open: true, message: "File published successfully!", severity: "success" });
            const data = await response.text();
            console.log("Publish response: ", data);
        } else {
            const errorData = await response.text();
            console.error("Publish response error:", errorData);
            setNotification({ open: true, message: "Failed to publish file", severity: "error" });
        }
    } catch (error) {
        console.error("Error publishing file:", error);
        setNotification({ open: true, message: "An error occurred", severity: "error" });
    } finally {
        setPublishDialogOpen(false);
    }
};

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
          Import Files
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

        {selectedFiles.length > 0 && (
          <Box sx={{ marginTop: 2 }}>
            <Typography variant="h6">Selected Files:</Typography>
            <List>
              {selectedFiles.map((file, index) => (
                <ListItem key={index} divider>
                  <Box 
                    sx={{ 
                      display: 'flex', 
                      flexDirection: 'column', // Align text and input fields vertically
                      width: '100%', 
                    }}
                  >
                    {/* selected file details */}
                    <ListItemText
                      sx={{
                        width:'100%',
                        whiteSpace: 'normal',  // wrap text
                        wordBreak: 'break-word',
                        overflowWrap: 'break-word',
                      }}
                      primary={file.name}
                      secondary={
                        <>
                          {`Size: ${(file.size / 1024).toFixed(2)} KB`} <br />
                          {`SHA-256 Hash: ${fileHashes[file.name] || "Computing..."}`}                       
                        </>
                      }  
                    />
                    
                    <Box sx={{ display: 'flex', flexDirection: 'column', marginTop: 1 }}>
                      <TextField
                        label="Description"
                        variant="outlined"
                        fullWidth
                        margin="normal"
                        value={descriptions[file.name] || ""}
                        onChange={(e) => handleDescriptionChange(file.name, e.target.value)}
                      />
                    </Box>
                  </Box>
                  <IconButton 
                      edge="end" 
                      aria-label="delete" 
                      onClick={() => handleDeleteFile(file.name)}
                      sx={{marginTop:15}}
                    >
                      <DeleteIcon />
                    </IconButton>
                </ListItem>
              ))}
            </List>
          </Box>
        )}
        
        {uploadedFiles.length > 0 && (
          <Box sx={{ marginTop: 4 }}>
            <Typography variant="h5" gutterBottom>
              Uploaded Files
            </Typography>
            <List>
              {uploadedFiles.map((file) => (
                <ListItem key={file.id} divider>
                  <ListItemText 
                    primary={file.name} 
                    secondary={
                      <>
                        {`Size: ${(file.size / 1024).toFixed(2)} KB`} <br />
                        {`Description: ${file.description}`} <br />
                        {`SHA-256: ${file.hash.slice(0, 10)}...${file.hash.slice(-10)}`}
                        </>
                    }                  
                    />
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <Typography variant="body2" component="span" sx={{ marginRight: 1 }}>
                      Publish
                    </Typography>
                    <Switch 
                      edge="end" 
                      onChange={() => handleTogglePublished(file.id)} 
                      checked={file.isPublished} 
                      color="primary"
                      inputProps={{ 'aria-label': `make ${file.name} public` }}
                    />
                    <IconButton edge="end" aria-label="delete" onClick={() => handleDeleteFile(file.id)}>
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                </ListItem>
              ))}
            </List>
          </Box>
        )}
      </Box>

      {/* Publish Modal */}
      <Dialog open={publishDialogOpen} onClose={() => setPublishDialogOpen(false)}>
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
      </Dialog>

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
