import React, { useState } from "react";
import Sidebar from "./Sidebar";
import { Box, Button, Typography } from "@mui/material";
import Header from "./Header";
const FilesPage: React.FC = () => {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files) {
      setSelectedFiles(Array.from(files)); 
    }
  };

  const handleUpload = () => {
    // uploading logic
    console.log("Uploading files:", selectedFiles);
  };

  return (
    <Box sx={{ padding: 2 }}>
        {/* <Header/> */}
        <Sidebar />
        <Box sx={{ marginLeft: 2 , marginTop: 10}}>
            <Typography variant="h4" gutterBottom>
            Upload Files
            </Typography>
            <input
            type="file"
            multiple
            onChange={handleFileChange}
            style={{ marginBottom: '16px' }}
            />
            <Button variant="contained" onClick={handleUpload}>
            Upload
            </Button>
            {/* show selected files */}
            {selectedFiles.length > 0 && (
            <Box sx={{ marginTop: 2 }}>
                <Typography variant="h6">Selected Files:</Typography>
                <ul>
                {selectedFiles.map((file, index) => (
                    <li key={index}>{file.name}</li>
                ))}
                </ul>
            </Box>
            )}
        </Box>
        </Box>
  );
};

export default FilesPage;
