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
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField
} from "@mui/material";
import { useTheme } from '@mui/material/styles';

// File model
interface UploadedFile {
  id: string; 
  name: string;
  size: number; 
  isPublished: boolean;
  fee?: number; // Add a fee property
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const FilesPage: React.FC = () => {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([]);
  const [notification, setNotification] = useState<{ open: boolean; message: string; severity: "success" | "error" }>({ open: false, message: "", severity: "success" });
  const [fileHashes, setFileHashes] = useState<{ [key: string]: string }>({});
  const [publishDialogOpen, setPublishDialogOpen] = useState(false); // Control for the modal
  const [currentFileId, setCurrentFileId] = useState<string | null>(null); // Track the file being published
  const [fee, setFee] = useState<number | undefined>(undefined); // Fee value for publishing
  
  const theme = useTheme();

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {    
    const files = event.target.files;
    try {
      if (files) {
        const fileArray = Array.from(files);
        // setSelectedFiles(fileArray);
        // Compute hashes for each file
        fileArray.forEach(file => computeSHA512(file));
        // uploaded files
        const newUploadedFiles: UploadedFile[] = fileArray.map((file) => ({
          id: `${file.name}-${file.size}-${Date.now()}`, // Simple unique ID
          name: file.name,
          size: file.size,
          isPublished: false,
        }));
  
        setUploadedFiles((prev) => [...prev, ...newUploadedFiles]);
        // setSelectedFiles([]);
        setNotification({ open: true, message: "Files uploaded successfully!", severity: "success" });
      }
    } catch (error) {
      setNotification({ open: true, message: "Failed to upload files.", severity: "error" });
      console.error("Error uploading files:", error);
    }
    
  };

  const handleTogglePublished = (id: string) => {
    // fix logic 
    setCurrentFileId(id); // Set the current file to publish
    if (uploadedFiles.find(file => file.id === id)?.isPublished) {
      setPublishDialogOpen(false);
      setUploadedFiles((prev) =>
        prev.map((file) =>
          file.id === currentFileId ? { ...file, isPublished: false } : file
        )
      );
    } else {
      setPublishDialogOpen(true); // Open the modal
    }
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

  // comfirm publishing file to dht?
  const handleConfirmPublish = () => {
    //implement actual logic later
    setUploadedFiles((prev) =>
      prev.map((file) =>
        file.id === currentFileId ? { ...file, isPublished: true, fee } : file
      )
    );
    setPublishDialogOpen(false); // close modal
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
        {/* <Button 
          variant="contained" 
          onClick={handleUpload} 
          disabled={selectedFiles.length === 0}
        >
          Upload Selected
        </Button> */}

        {/* {selectedFiles.length > 0 && (
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
        )} */}

        {uploadedFiles.length > 0 && (
          <Box sx={{ marginTop: 4 }}>
            <Typography variant="h5" gutterBottom>
              Selected Files
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
                      Publish
                    </Typography>
                    <Switch 
                      edge="end" 
                      onChange={() => handleTogglePublished(file.id)} 
                      checked={file.isPublished} 
                      color="primary"
                      inputProps={{ 'aria-label': `publish ${file.name}` }}
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
