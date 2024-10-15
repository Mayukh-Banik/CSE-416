import React, { useState } from "react";
import Sidebar from "./Sidebar";
import PublishIcon from '@mui/icons-material/Publish';
import DeleteIcon from '@mui/icons-material/Delete';
import { 
  Box, 
  Button, 
  Typography, 
  List, 
  ListItem, 
  ListItemText, 
  Switch,
  IconButton,
  Snackbar,
  Alert
} from "@mui/material";

interface UploadedFile {
  id: string; 
  name: string;
  size: number; 
  isPublic: boolean;
}

const FilesPage: React.FC = () => {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([]);
  const [notification, setNotification] = useState<{ open: boolean; message: string; severity: "success" | "error" }>({ open: false, message: "", severity: "success" });
  const [fileHashes, setFileHashes] = useState<{ [key: string]: string }>({});

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files) {
      const fileArray = Array.from(files);
      setSelectedFiles(fileArray);
      // Compute hashes for each file
      fileArray.forEach(file => computeSHA512(file));
    }
  };

  const handleUpload = async () => {
    try {
      // not actual functionality
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // uploaded files
      const newUploadedFiles: UploadedFile[] = selectedFiles.map((file) => ({
        id: `${file.name}-${file.size}-${Date.now()}`, // Simple unique ID
        name: file.name,
        size: file.size,
        isPublic: false,
      }));

      setUploadedFiles((prev) => [...prev, ...newUploadedFiles]);
      setSelectedFiles([]);
      setNotification({ open: true, message: "Files uploaded successfully!", severity: "success" });
    } catch (error) {
      setNotification({ open: true, message: "Failed to upload files.", severity: "error" });
      console.error("Error uploading files:", error);
    }
  };

  const handleTogglePublic = async (id: string) => {
    setUploadedFiles((prev) =>
      prev.map((file) =>
        file.id === id ? { ...file, isPublic: !file.isPublic } : file
      )
    );
  };

  const handleDeleteFile = (id: string) => {
    setUploadedFiles((prev) => prev.filter((file) => file.id !== id));
    setNotification({ open: true, message: "File deleted.", severity: "success" });
  };

  const handleCloseNotification = () => {
    setNotification({ ...notification, open: false });
  };

  const handleUploadClick = () => {
    document.getElementById("file-input")?.click();
  };


  const computeSHA512 = async (file: File) => {
    const arrayBuffer = await file.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest("SHA-512", arrayBuffer);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray
      .map(byte => byte.toString(16).padStart(2, "0"))
      .join("");
    
    setFileHashes(prevHashes => ({
      ...prevHashes,
      [file.name]: hashHex,
    }));
  };

  return (
    <Box sx={{ display: 'flex', padding: 2 }}>
      <Sidebar />
      <Box sx={{ flexGrow: 1, marginLeft: 2, marginTop: 10 }}>
        <Typography variant="h4" gutterBottom>
          Upload Files
        </Typography>
        <input
          type="file"
          id="file-input"
          multiple
          onChange={handleFileChange}
          style={{ display: 'none' }} // Hide the default file input
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
            '&:hover': {
              backgroundColor: '#e3f2fd',
            },
          }}
          onClick={handleUploadClick} // Clicking the box opens the file input
        >
          <PublishIcon sx={{ fontSize: 50 }} />
          <Typography variant="h6" sx={{ marginTop: 1 }}>
            Upload File
          </Typography>
        </Box>
        <Button 
          variant="contained" 
          onClick={handleUpload} 
          disabled={selectedFiles.length === 0}
        >
          Upload Selected
        </Button>

        {/* Show selected files */}
        {selectedFiles.length > 0 && (
          <Box sx={{ marginTop: 2 }}>
            <Typography variant="h6">Selected Files:</Typography>
            <List>
              {selectedFiles.map((file, index) => (
                <ListItem key={index}>
                  <ListItemText 
                    primary={file.name} 
                    secondary={`Size: ${(file.size / 1024).toFixed(2)} KB`} 
                  />
                </ListItem>
              ))}
            </List>
          </Box>
        )}

        {/* uploaded files */}
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
                    secondary={`Size: ${(file.size / 1024).toFixed(2)} KB`} 
                  />
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <Typography variant="body2" component="span" sx={{ marginRight: 1 }}>
                      Public
                    </Typography>
                    <Switch 
                      edge="end" 
                      onChange={() => handleTogglePublic(file.id)} 
                      checked={file.isPublic} 
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
