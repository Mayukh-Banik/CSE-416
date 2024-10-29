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

// File model
interface UploadedFile {
  isPublished: boolean | undefined;
  id: string; 
  name: string;
  size: number; 
  description: string; // Add description field
  hash: string; // Add hash field
  fee?: number; // Add a fee property
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const FilesPage: React.FC = () => {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [descriptions, setDescriptions] = useState<{ [key: string]: string }>({}); // Track descriptions
  const [fileHashes, setFileHashes] = useState<{ [key: string]: string }>({}); // Track hashes
  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([]);
  const [notification, setNotification] = useState<{ open: boolean; message: string; severity: "success" | "error" }>({ open: false, message: "", severity: "success" });
  const [publishDialogOpen, setPublishDialogOpen] = useState(false); // Control for the modal
  const [currentFileId, setCurrentFileId] = useState<string | null>(null); // Track the file being published
  const [fee, setFee] = useState<number | undefined>(undefined); // Fee value for publishing
  
  const theme = useTheme();

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files) {
      const fileArray = Array.from(files);
      setSelectedFiles(prevSelectedFiles => [...prevSelectedFiles, ...fileArray]);
      // Compute hashes for each file
      fileArray.forEach(file => computeSHA512(file));
    }
  };

  // const handleUpload = async () => {
  //   try {
  //     // Simulate file upload delay
  //     await new Promise((resolve) => setTimeout(resolve, 1000));

  //     // Create uploaded file objects with descriptions and hashes
  //     const newUploadedFiles: UploadedFile[] = selectedFiles.map((file) => ({
  //       id: `${file.name}-${file.size}-${Date.now()}`, // Unique ID
  //       name: file.name,
  //       size: file.size,
  //       description: descriptions[file.name] || "", // Attach description
  //       hash: fileHashes[file.name], // Attach the computed hash
  //       isPublished: false,
  //     }));

  //     setUploadedFiles((prev) => [...prev, ...newUploadedFiles]);
  //     setSelectedFiles([]);
  //     setDescriptions({}); // Reset descriptions after upload
  //     setFileHashes({}); // Reset hashes after upload
  //     setNotification({ open: true, message: "Files uploaded successfully!", severity: "success" });
  //   } catch (error) {
  //     setNotification({ open: true, message: "Failed to upload files.", severity: "error" });
  //     console.error("Error uploading files:", error);
  //   }
  // };

  // // jerry's upload function
  // const handleUpload = async () => {
  //   if (selectedFiles.length === 0) return;

  //   for (const file of selectedFiles) {
  //     const filedata = await toString(file);
  //     const metadata = {
  //       name: file.name,
  //       type: file.type,
  //       size: file.size,
  //       file_data: filedata,
  //     };
  //     console.log("Uploading file:", file.name);

  //     await fetch("http://localhost:8000/upload", {
  //       method: "POST",
  //       headers: {
  //         "Content-Type": "application/json",
  //       },
  //       body: JSON.stringify(metadata),
  //     })
  //       .then((response) => response.json())
  //       .then((data) => {
  //         console.log("Success", data);
  //       })
  //       .catch((error) => {
  //         console.log("Error", error);
  //       });
  //   }
  // };


  const handleUpload = async () => {
    if (selectedFiles.length === 0) return;
  
    try {
      // Create uploaded file objects with descriptions, hashes, and metadata
      const newUploadedFiles = await Promise.all(
        selectedFiles.map(async (file) => {
          const fileData = file.toString(); // Convert file to string data
          const metadata = {
            name: file.name,
            type: file.type,
            size: file.size,
            file_data: fileData,
            description: descriptions[file.name] || "",
            hash: fileHashes[file.name],
          };
  
          // Send the metadata to the server
          const response = await fetch("http://localhost:8000/upload", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(metadata),
          });
          
          const data = await response.json();
          console.log("File upload successful:", data);
  
          // Return file info for the uploadedFiles state
          return {
            id: `${file.name}-${file.size}-${Date.now()}`,
            name: file.name,
            size: file.size,
            description: descriptions[file.name] || "",
            hash: fileHashes[file.name],
            isPublished: false,
          };
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
    setCurrentFileId(id); 
    if (uploadedFiles.find(file => file.id === id)?.isPublished) {
      setPublishDialogOpen(false);
      setUploadedFiles((prev) =>
        prev.map((file) =>
          file.id === currentFileId ? { ...file, isPublished: false } : file
        )
      );
    } else {
      setPublishDialogOpen(true); 
    }
  };

  const handleDeleteFile = (id: string) => {
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

  const handleConfirmPublish = () => {
    setUploadedFiles((prev) =>
      prev.map((file) =>
        file.id === currentFileId ? { ...file, isPublished: true, fee } : file
      )
    );
    setPublishDialogOpen(false);
    setNotification({ open: true, message: "File published successfully!", severity: "success" });
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
                          {`SHA-512 Hash: ${fileHashes[file.name] || "Computing..."}`}                       
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
                        {`SHA-512: ${file.hash.slice(0, 10)}...${file.hash.slice(-10)}`}
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
