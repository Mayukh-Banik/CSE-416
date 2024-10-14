import React, { useState } from "react";
import Sidebar from "./Sidebar";
import { Box, Button, Typography } from "@mui/material";
import Header from "./Header";

const FilesPage: React.FC = () => {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
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

  const handleUpload = () => {
    // uploading logic
    console.log("Uploading files:", selectedFiles);
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
    <Box sx={{ padding: 2 }}>
      {/* <Header/> */}
      <Sidebar />
      <Box sx={{ marginLeft: 2, marginTop: 10 }}>
        <Typography variant="h4" gutterBottom>
          Upload Files
        </Typography>
        <input
          type="file"
          multiple
          onChange={handleFileChange}
          style={{ marginBottom: "16px" }}
        />
        <Button variant="contained" onClick={handleUpload}>
          Upload
        </Button>

        {/* Show selected files and their hashes */}
        {selectedFiles.length > 0 && (
          <Box sx={{ marginTop: 2 }}>
            <Typography variant="h6">Selected Files:</Typography>
            <ul>
              {selectedFiles.map((file, index) => (
                <li key={index}>
                  {file.name} - SHA-512:{" "}
                  {fileHashes[file.name] ? fileHashes[file.name] : "Calculating..."}
                </li>
              ))}
            </ul>
          </Box>
        )}
      </Box>
    </Box>
  );
};

export default FilesPage;
