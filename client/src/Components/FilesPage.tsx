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
import { FileMetadata, saveFileMetadata } from "../utils/localStorage";
//import { FileMetadata } from '../../local_server/models/file';


// import { saveFileMetadata, getFilesForUser, deleteFileMetadata, updateFileMetadata, FileMetadata } from '../utils/localStorage'

// user id/wallet id is hardcoded for now

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

import fs from 'fs';
import fsPromises from 'fs/promises';
import path from 'path';

const FilesPage: React.FC = () => {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [descriptions, setDescriptions] = useState<{ [key: string]: string }>({}); // Track descriptions
  const [fileHashes, setFileHashes] = useState<{ [key: string]: string }>({}); // Track hashes
  const [uploadedFiles, setUploadedFiles] = useState<FileMetadata[]>([]);
  const [notification, setNotification] = useState<{ open: boolean; message: string; severity: "success" | "error" }>({ open: false, message: "", severity: "success" });
  const [publishDialogOpen, setPublishDialogOpen] = useState(false); // Control for the modal
  const [currentFileHash, setCurrentFileHash] = useState<string | null>(null); // Track the file being published
  // const [fee, setFee] = useState<number | undefined>(undefined); // Fee value for publishing
  const [fees, setFees] = useState<{ [key: string]: number }>({}); // Track fees of to be uploaded files
  
  const theme = useTheme();
  
  // useEffect(() => {
  //   const fetchUploadedFiles = () => {
  //     const files = getFilesForUser('123'); // Adjust user ID as necessary
  //     setUploadedFiles(files);
  //   };
  //   fetchUploadedFiles();
  // }, []);

   useEffect(()=>{
    const fetchUploadedFiles = async () => {
      try{
        const res = await fetch('/fetchUploadedFiles');
        if (!res.ok) {
          throw new Error('errordfhsf');
        }
        const data = await res.json();
        setUploadedFiles(data);
      } catch (error) {
        console.error('failed to get all')
      }
    }
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
                // setPublishDialogOpen(true);
                let metadata:FileMetadata = {
                    //id: `${file.name}-${file.size}-${Date.now()}`, // Unique ID for the uploaded file
                    name: file.name,
                    type: file.type,
                    size: file.size,
                    // file_data: base64FileData, // Encode file data as Base64 if required
                    description: descriptions[file.name] || "",
                    hash: fileHashes[file.name], // not needed - computed on backend
                    isPublished: true, // Initially published
                    fee: fees[file.name],
                };

                // add into the json file in a certain predetermined path
                /** 
                const dir_path = path.join(__dirname, '..','utils','files.json');
                try{
                  if (!fs.existsSync(dir_path)){
                    await fsPromises.mkdir(dir_path);
                    console.log("should've created utils json file");
                  }
                  const file_data = await fsPromises.readFile(dir_path,'utf-8');
                  const files_data: FileMetadata[] = JSON.parse(file_data);

                 let dup = false;
                 for (const file of files_data){
                  if(file.name===metadata.name){
                    dup = true;
                    break;
                  }
                 }

                 if(!dup){
                  files_data.push(metadata);
                  await fsPromises.writeFile(dir_path,JSON.stringify(files_data,null,2),'utf8');
                  console.log("added file metadata");
                }
                else{
                  console.log("duplicate found");
                }
                } catch (error){
                  console.error("error checking/updating files.json",error);
                }
                  */
                
                // Send the metadata to the server
                const response = await fetch("http://localhost:8081/upload", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(metadata),
                });
                console.log("body is ",JSON.stringify(metadata));
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const data = await response.text();
                console.log("File upload successful:", data, metadata);
                
                // update local server/database
                //saveFileMetadata('123', metadata);
                
                return metadata;
            })
        );

        // Update uploadedFiles state
        setUploadedFiles((prev) => [...prev, ...newUploadedFiles]);

        // Clear selected files, descriptions, and hashes after successful upload
        setSelectedFiles([]);
        setDescriptions({});
        setFees({});
        setFileHashes({});

        // Show success notification
        setNotification({ open: true, message: "Files uploaded successfully!", severity: "success" });
    } catch (error) {
        console.error("Error uploading files:", error);

        // Show error notification
        setNotification({ open: true, message: "Failed to upload files.", severity: "error" });
    }
};

// redo
  const updateFile = async (hash: string, updatedData: FileMetadata) => {
    try {
      const response = await fetch(`http://localhost:8082/update/${hash}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updatedData),
      });

      if (!response.ok) {
        throw new Error('Failed to update file');
      }

      const result = await response.json();
      console.log('File updated:', result);
    } catch (error) {
      console.error('Error updating file:', error);
    }
  };

  const handleDescriptionChange = (fileId: string, description: string) => {
    setDescriptions((prev) => ({ ...prev, [fileId]: description }));
  };

  const handleFeeChange = (fileId: string, fee: number) => {
    setFees((prev: any) => ({ ...prev, [fileId]: fee }));
  };

  const handleTogglePublished = (hash: string) => {
    const file = uploadedFiles.find(file => file.hash === hash); // Get the current file

    if (file) {
      updateFile(hash, {
        ...file,
        isPublished: !file.isPublished,
 })
    }
    if (file?.isPublished) {
        // If the file is already published, set the unpublished state
        setUploadedFiles((prev) =>
            prev.map((f) =>
                f.hash === hash ? { ...f, isPublished: false } : f
            )
        );
    } else {
        // If the file is not published, open the publish dialog
        setCurrentFileHash(hash); // Set the current file ID to the one being published
        setPublishDialogOpen(true); 
    }
    
  };


  // deletes file from list of files to be uploaded - these have not been added to the database yet
  const handleDeleteSelectedFile = (hash: string) => {
    // deleteFileMetadata('123', id);
    setUploadedFiles((prev) => prev.filter((file) => file.hash !== hash));
    
    setSelectedFiles((prev) => prev.filter((file) => file.name !== hash));
    setNotification({ open: true, message: "File deleted.", severity: "success" });
  };

  const handleDeleteUploadedFile = async (hash: string) => {
    try {
      const response = await fetch(('http://localhost:8082/delete/${hash}'), {
        method: "DELETE" 
      });

      if (!response.ok) {
        throw new Error('Failed to update file');
      }
  
      const result = await response.json();
      console.log('File deleted:', result);
    } catch (error) {
    console.error('Error deleting file:', error);
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

  const handleConfirmPublish = async () => {
    const fileToPublish = uploadedFiles.find(file => file.hash === currentFileHash);
    
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
        const response = await fetch("http://localhost:8081/publish", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              key: fileToPublish.hash,           // The hash ID generated on upload
              value: JSON.stringify(metadata),      // should not have metadata in value field
          }),
        });

        if (response.ok) {
            setUploadedFiles((prev) =>
                prev.map((file) =>
                    file.hash === currentFileHash ? { ...file, isPublished: true, fees } : file
                )
            );

            setNotification({ open: true, message: "File published successfully!", severity: "success" });
            const data = await response.text();
            updateFile(fileToPublish.hash, {
                ...fileToPublish,
                isPublished: true, 
           })
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
                          {/* {`Description: ${(file.description)}`} <br /> */}

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

                    <Box sx={{ display: 'flex', flexDirection: 'column', marginTop: 1 }}>
                      <TextField
                        label="Fee"
                        type="number"
                        variant="outlined"
                        fullWidth
                        margin="normal"
                        value={fees[file.name] || 0}
                        onChange={(e) => handleFeeChange(file.name, parseFloat(e.target.value))}
                      />
                    </Box>

                  </Box>
                  <IconButton 
                      edge="end" 
                      aria-label="delete" 
                      onClick={() => handleDeleteSelectedFile(file.name)}
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
                <ListItem key={file.hash} divider>
                  <ListItemText 
                    primary={file.name} 
                    secondary={
                      <>
                        {`Size: ${(file.size / 1024).toFixed(2)} KB`} <br />
                        {`Description: ${file.description}`} <br />
                        {`SHA-256: ${file.hash.slice(0, 10)}...${file.hash.slice(-10)}`} <br />
                        {`Fee: ${file.fee}`} <br />
                        </>
                    }                  
                    />
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <Typography variant="body2" component="span" sx={{ marginRight: 1 }}>
                      Publish
                    </Typography>
                    <Switch 
                      edge="end" 
                      onChange={() => handleTogglePublished(file.hash)} 
                      checked={file.isPublished} 
                      color="primary"
                      inputProps={{ 'aria-label': `make ${file.name} public` }}
                    />
                    <IconButton edge="end" aria-label="delete" onClick={() => handleDeleteUploadedFile(file.hash)}>
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
