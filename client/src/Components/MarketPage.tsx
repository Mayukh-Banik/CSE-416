import React, { useState } from 'react';
import { Box, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Button, TextField, TablePagination, Paper, Dialog, DialogTitle, DialogContent, DialogActions, Snackbar, Alert, LinearProgress } from '@mui/material';
import Sidebar from './Sidebar';
import { useTheme } from '@mui/material/styles';
import { FileMetadata, Provider } from "../models/fileMetadata"
import { Transaction } from '../models/transactions';

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const MarketplacePage: React.FC = () => {
  const theme = useTheme();
  const [searchTerm, setSearchTerm] = useState('');
  const [searchResults, setSearchResults] = useState<FileMetadata[]>([])
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [selectedFile, setSelectedFile] = useState<FileMetadata | null>(null);
  const [open, setOpen] = useState(false);
  const [fileHash, setFileHash] = useState('');
  const [notification, setNotification] = useState<{ open: boolean; message: string; severity: "success" | "error" }>({ open: false, message: "", severity: "success" });
  const [providers, setProviders] = useState<Provider[]>([]);
  const [loading, setLoading] = useState(false)

  const resetStates = async () => {
    setFileHash("");
    setProviders([]);
    setSelectedFile(null);
    setSearchTerm("");
    setLoading(false);
  }

  const handleCloseNotification = () => {
    setNotification({ ...notification, open: false });
  };
  
  const handleDownloadRequest = async (file: FileMetadata) => {
    setSelectedFile(file);
    setOpen(true); // Open the modal for provider selection
  };

  const handleProviderSelect = async (provider: string) => {
    try {
      let request: Transaction = {
        TargetID: provider,
        FileHash: fileHash, 
        RequesterID: "",
        Status: "pending",
        FileName: "",
        // CreatedAt: Date.now().toLocaleString(),
      }
      console.log("Request data being sent:", request);

      const response = await fetch(`http://localhost:8081/download/request`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        const responseText = await response.text(); // Get the error message
        setNotification({ open: true, message: responseText, severity: "error" });
        return
      }
      setNotification({ open: true, message: `Request sent to provider ${provider}`, severity: "success" });
      setOpen(false); // Close the modal after selecting a provider
    } catch (error) {
      console.error("Error in handleProviderSelect:", error);
      setNotification({ open: true, message: "Failed to find providers.", severity: "error" });
    } 
  };

  const handleRefresh = async () => {
    console.log("HI");
  };

  const handleDownloadByHash = async () => {
    console.log("HI");
    const hash = prompt("Enter the file hash");
    console.log('you entered:', hash)
    if (hash == null || hash.length==0) return;
    setFileHash(hash);
    getFileByHash(hash);
    setOpen(true);
  }

  
  // only works for complete file hashes
  const handleSearchRequest = async (searchTerm: string) => {
    // resetStates();
    if (!searchTerm || searchTerm.length === 0) return;
    setFileHash(searchTerm);
    await getFileByHash(searchTerm);

    if (selectedFile) {
      setSearchResults([selectedFile])
    }
    // setOpen(true);
  }

  const getFileByHash = async (hash: string) => {
    setLoading(true)
    try {
        const encodedHash = encodeURIComponent(hash);  // Ensure hash is URL-safe
        const url = `http://localhost:8081/files/getFile?val=${encodedHash}`;
        
        console.log("Request URL:", url); // Log the request URL for debugging

        const response = await fetch(url, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP Error: status : ${response.status}`);
        }

        const data = await response.json();
        console.log("file metadata:", data); 

        setSelectedFile(data);
        setProviders(data.Providers);
        setSearchResults([data])
    } catch (error) {
        console.error("Error:", error);
        setNotification({ open: true, message: "File not found", severity: "error" });
    } finally {
      setLoading(false)
    }
  };


  // const handleDownload = (fileHash:string) => {
  //   getFile(fileHash);
  // };

  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0); // Reset to the first page
  };

  // pagination
  // const indexOfLastFile = (page + 1) * rowsPerPage;
  // const indexOfFirstFile = indexOfLastFile - rowsPerPage;
  // const currentFiles = filteredFiles.slice(indexOfFirstFile, indexOfLastFile);

  return (
    <Box
      sx={{
        padding: 2,
        marginTop: '70px',
        marginLeft: `${drawerWidth}px`, // Default expanded margin
        transition: 'margin-left 0.3s ease', // Smooth transition
        [theme.breakpoints.down('sm')]: {
          marginLeft: `${collapsedDrawerWidth}px`, // Adjust left margin for small screens
        },
      }}
    >
      <Sidebar/>
      <Typography variant="h4" gutterBottom>
        Marketplace
      </Typography>

      <Button variant="contained" onClick={() => {handleRefresh()}}>
        Refresh
      </Button>
      <Button variant="contained" onClick={() => handleDownloadByHash()}>
        Download by Hash
      </Button>

      <TextField
        label="Search Files"
        variant="outlined"
        fullWidth
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === "Enter"){
            handleSearchRequest(searchTerm);
          }
        }}
        sx={{ marginBottom: 2, background: "white" }}
      />

      {loading && <LinearProgress sx={{ width: '100%', marginTop: 2 }} />} {/* Progress bar when loading is true */}
      
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>File Name</TableCell>
              <TableCell>File Size (KB)</TableCell>
              <TableCell>Reputation</TableCell>
              <TableCell>Created At</TableCell>
              <TableCell></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {searchResults?.map((file) => (
              <TableRow key={file.Hash}>
                <TableCell>{file.Name}</TableCell>
                <TableCell>{(file.Size / 1024).toFixed(2)}</TableCell>
                {/* <TableCell>{file.Reputation}</TableCell> */}
                <TableCell>3/5</TableCell>
                {/* <TableCell>{file.createdAt.toLocaleDateString()}</TableCell> */}
                <TableCell>2023-09-15</TableCell>

                <TableCell>
                  <Button variant="contained" onClick={() => handleDownloadRequest(file)}>
                    Download
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      {/* <TablePagination
        component="div"
        count={filteredFiles.length}
        page={page}
        onPageChange={handleChangePage}
        rowsPerPage={rowsPerPage}
        onRowsPerPageChange={handleChangeRowsPerPage}
        rowsPerPageOptions={[5, 10, 25]}
        sx={{ marginTop: 2 }}
      /> */}

      <Dialog open={open} onClose={() => setOpen(false)}>
        <DialogTitle>{selectedFile?.Name}</DialogTitle>
        <DialogContent>
        {selectedFile && (
          <Box sx={{ marginBottom: 2 }}>
            <Typography>Size: {selectedFile.Size} MB</Typography>
            {selectedFile.Description && (
              <Typography>Description: {selectedFile.Description}</Typography>
            )}
          </Box>
        )}
  
          {providers.length ? (
            providers
              .filter((value, index, self) =>
                index === self.findIndex((t) => (
                  t.PeerID === value.PeerID // or use t.address if filtering by address
                ))
              )
              .map((provider) => (
                <Button
                  key={provider.PeerID} // Ensure this is unique for the key
                  variant="outlined"
                  onClick={() => handleProviderSelect(provider.PeerID)}
                  sx={{
                    margin: 1,
                    display: 'flex',
                    justifyContent: 'space-between',
                    width: '100%',
                  }}
                >
                  <span>{provider.PeerID.substring(0,7)}...{provider.PeerID.substring(provider.PeerID.length - 7)}</span> 
                  <span>{provider.Fee} SQD/MB</span>
                </Button>
              ))
          ) : (
            <Typography>No providers available for this file.</Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Cancel</Button>
        </DialogActions>
      </Dialog>


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

export default MarketplacePage;
