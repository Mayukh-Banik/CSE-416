import React, { useState } from 'react';
import { Box, Typography, List, ListItem, ListItemText, Button, TextField, Pagination } from '@mui/material';
import Header from './Header';
import Sidebar from './Sidebar';

export interface IFile {
  fileName: string;
  hash: string;
  reputation: string;
  fileSize: number; // in bytes
  createdAt: Date;
}

const MarketPage: React.FC = () => {
  const [files, setFiles] = useState<IFile[]>([
    { fileName: 'Example File 1', hash: 'abc123', reputation: 'High', fileSize: 2048, createdAt: new Date() },
    { fileName: 'Example File 2', hash: 'def456', reputation: 'Medium', fileSize: 1024, createdAt: new Date() },
    { fileName: 'Example File 3', hash: 'ghi789', reputation: 'Low', fileSize: 512, createdAt: new Date() },
    { fileName: 'Example File 4', hash: 'jkl012', reputation: 'High', fileSize: 2048, createdAt: new Date() },
    { fileName: 'Example File 5', hash: 'mno345', reputation: 'Medium', fileSize: 1024, createdAt: new Date() },
    { fileName: 'Example File 6', hash: 'pqr678', reputation: 'Low', fileSize: 512, createdAt: new Date() },
    // Add more files as needed
  ]);
  
  const [searchTerm, setSearchTerm] = useState('');
  const [page, setPage] = useState(1);
  const itemsPerPage = 10;

  const filteredFiles = files.filter(file => 
    file.fileName.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleDownloadRequest = (file: IFile) => {
    console.log(`Requesting download for: ${file.fileName}`);
  };

  const handlePageChange = (event: React.ChangeEvent<unknown>, value: number) => {
    setPage(value);
  };

  const indexOfLastFile = page * itemsPerPage;
  const indexOfFirstFile = indexOfLastFile - itemsPerPage;
  const currentFiles = filteredFiles.slice(indexOfFirstFile, indexOfLastFile);

  return (
    <Box sx={{ padding: 2, marginTop:'100px' }}>
      <Sidebar/>
      <Typography variant="h4" gutterBottom>
        Marketplace
      </Typography>
      <TextField
        label="Search Files"
        variant="outlined"
        fullWidth
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        sx={{ marginBottom: 2 }}
      />
      <List>
        {currentFiles.map((file) => (
          <ListItem key={file.hash} divider>
            <ListItemText
              primary={file.fileName}
              secondary={`Size: ${(file.fileSize / 1024).toFixed(2)} KB | Reputation: ${file.reputation} | Created At: ${file.createdAt.toLocaleDateString()}`}
            />
            <Button variant="contained" onClick={() => handleDownloadRequest(file)}>
              Download
            </Button>
          </ListItem>
        ))}
      </List>
      <Pagination
        count={Math.ceil(filteredFiles.length / itemsPerPage)}
        page={page}
        onChange={handlePageChange}
        color="primary"
        sx={{ marginTop: 2 }}
      />
    </Box>
  );
};

export default MarketPage;
