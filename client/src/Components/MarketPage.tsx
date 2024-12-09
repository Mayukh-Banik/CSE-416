import React, { useEffect, useState } from 'react';
import { Box, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Button, TextField, TablePagination, Paper, Dialog, DialogTitle, DialogContent, DialogActions, Snackbar, Alert, LinearProgress } from '@mui/material';
import Sidebar from './Sidebar';
import { useTheme } from '@mui/material/styles';
import { FileMetadata, Provider, DHTMetadata } from "../models/fileMetadata"
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
  const [loadingRequest, setLoadingRequest] = useState(false);
  const [loadingSearch, setLoadingSearch] = useState(false);
  const [refresh, setRefresh] = useState(false);
  const [ratings, setRatings] = useState<{ [key: string]: number }>({}); // Store ratings by file hash
  const [marketResults, setMarketResults] = useState<FileMetadata[]>([])

  useEffect(() => {
    const initializeMarketplace = async () => {
      localStorage.setItem("marketplaceLoaded", "false");
      const hasLoaded = localStorage.getItem("marketplaceLoaded");
      console.log("Has Loaded from localStorage:", hasLoaded); // Debug log
      if (!hasLoaded || hasLoaded === "false") {
        console.log("First load, calling handleRefresh..."); // Debug log
        await handleRefresh(true);
        localStorage.setItem("marketplaceLoaded", "true");
      }
    };
  
    initializeMarketplace();
  }, []); // Empty dependency array to run only on initial render
  

  useEffect(() => {
    const fetchRatings = async () => {
      setRatings({})
      const updatedRatings: { [key: string]: number } = { ...ratings };
      console.log('search results: ', searchResults)
      for (const file of searchResults) {
        console.log("getting rating of file: ", file)
        if (file.Rating != null) {
          updatedRatings[file.Hash] = file.Rating
        }
      }
      setRatings(updatedRatings); // Batch update once all ratings are fetched
      console.log("all file ratings: ", ratings)
    };
  
    if (searchResults != null) {
    fetchRatings();
  }
  }, [searchResults]); // Runs when search results change
  

  const resetStates = async () => {
    setFileHash("");
    setProviders([]);
    setSelectedFile(null);
    // setSearchTerm("");
    setLoadingRequest(false);
    setLoadingSearch(false);
    setSearchResults(marketResults)
    setRatings({})
  }

  const handleCloseNotification = () => {
    setNotification({ ...notification, open: false });
  };
  

  const handleDownloadRequest = async (file: FileMetadata) => {
    if (refresh) {
      console.log('handling download request for file', file.Hash)
      setFileHash(file.Hash)
      await getFileByHash(file.Hash)
    } else {
      setSelectedFile(file)
    }
    setOpen(true); // Open the modal for provider selection
    setRefresh(false);
  };

  const handleProviderSelect = async (provider: string) => {
    console.log(`requesting file ${fileHash} from provider ${provider}: `)
    setLoadingRequest(true)
    try {
      let request: Transaction = {
        TargetID: provider,
        FileHash: fileHash, 
        RequesterID: "",
        Status: "pending",
        FileName: selectedFile?.Name || "" ,
        Size: selectedFile?.Size || 0,
        TransactionID: "",
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
      } else {
        setNotification({ open: true, message: `Request sent to provider ${provider}`, severity: "success" });
        setOpen(false); // Close the modal after selecting a provider
        setSearchResults(marketResults)
      }
    } catch (error) {
      console.error("Error in handleProviderSelect:", error);
      setNotification({ open: true, message: "Failed to find providers.", severity: "error" });
    } finally {
      setLoadingRequest(false);
    }
  };

  const handleRefresh = async (initial: boolean) => {
    setRefresh(true);
    await resetStates();
    setSearchResults([])
    setLoadingSearch(true);
    console.log("Refreshing marketplace");

    try {
      const response = await fetch(`http://localhost:8081/files/refresh?val=${initial}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Failed to fetch marketplace files.");
      }

      const data = await response.json();
      console.log("Data received:", data);
      setSearchResults(data);
      setMarketResults(data)
      setNotification({ open: true, message: "Marketplace refreshed successfully.", severity: "success" });
    } catch (error) {
      console.error("Error refreshing marketplace:", error);
      setNotification({ open: true, message: "Failed to refresh marketplace.", severity: "error" });
    } finally {
      setLoadingSearch(false);
    }
  };
  
  // only works for complete file hashes
  const handleSearchRequest = async (searchTerm: string) => {
    await resetStates();
    setSearchResults([])
    if (!searchTerm || searchTerm.length === 0) {
      setSearchResults(marketResults);
      return;
    }
    setFileHash(searchTerm);
    await getFileByHash(searchTerm);

    if (selectedFile) {
      setSearchResults([selectedFile])
    }
    // setOpen(true);
  }

  const getFileByHash = async (hash: string) => {
    await resetStates()
    setLoadingSearch(true)
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
        console.log("getFileByHash: providers for file ", hash, data.Providers)
        setSearchResults([data])

        const updatedRatings: { [key: string]: number } = {};
        updatedRatings[hash] = data.Rating
        setRatings(updatedRatings)
    } catch (error) {
        console.error("Error:", error);
        setNotification({ open: true, message: "File not found", severity: "error" });
    } finally {
      setLoadingSearch(false)
    }
  };

  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0); // Reset to the first page
  };

  // pagination
  const indexOfLastFile = (page + 1) * rowsPerPage;
  const indexOfFirstFile = indexOfLastFile - rowsPerPage;
  const currentFiles = searchResults?.slice(indexOfFirstFile, indexOfLastFile);

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

      <Button 
        variant="contained" 
        onClick={() => {handleRefresh(false)}}
        disabled={loadingRequest || loadingSearch}
        sx={{marginBottom:2}}
        >
        Refresh
      </Button>

      {/* {loadingSearch && <LinearProgress sx={{ marginBottom: 2 }} />} */}

      <TextField
        label="Search Files"
        variant="outlined"
        fullWidth
        value={searchTerm}
        onChange={(e) => {
          if (!loadingRequest || !loadingSearch) {
            setSearchTerm(e.target.value)
          }
        }}
        onKeyDown={(e) => {
          if (e.key === "Enter" && (!loadingRequest || !loadingSearch)){
            handleSearchRequest(searchTerm);
          }
        }}
        sx={{ marginBottom: 2, background: "white" }}
        disabled={loadingRequest || loadingSearch}
      />

      {loadingSearch && <LinearProgress sx={{ width: '100%', marginTop: 2 }} />} {/* Progress bar when loading is true */}
      
      {currentFiles != null && currentFiles.length > 0 && <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>File Name</TableCell>
              <TableCell>File Size (KB)</TableCell>
              {/* <TableCell>Rating</TableCell> */}
              {/* <TableCell>Created At</TableCell> */}
              <TableCell></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {currentFiles?.map((file) => (
              <TableRow key={file.Hash}>
                <TableCell>{file.Name}</TableCell>
                <TableCell>{(file.Size / 1024).toFixed(2)}</TableCell>
                {/* <TableCell>{file.Reputation}</TableCell> */}
                {/* <TableCell>{ratings[file.Hash] !== undefined ? ratings[file.Hash] : "0"}</TableCell> Use ratings */}
                {/* <TableCell>{file.createdAt.toLocaleDateString()}</TableCell> */}
                {/* <TableCell>2023-09-15</TableCell> */}

                <TableCell>
                  <Button variant="contained" onClick={() => handleDownloadRequest(file)}>
                    Download
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>}

      {searchResults != null && searchResults.length > 0 && <TablePagination
        component="div"
        count={searchResults.length}
        page={page}
        onPageChange={handleChangePage}
        rowsPerPage={rowsPerPage}
        onRowsPerPageChange={handleChangeRowsPerPage}
        rowsPerPageOptions={[5, 10, 25]}
        sx={{ marginTop: 2 }}
      />}

      <Dialog open={open} onClose={() => setOpen(false)}>
        {loadingRequest && <LinearProgress sx={{ width: '100%', marginTop: 2 }} />} {/* Progress bar when loading is true */}
        <DialogTitle>{selectedFile?.Name}</DialogTitle>
        <DialogContent>
        {selectedFile && (
          <Box sx={{ marginBottom: 2 }}>
            <Typography>Size: {selectedFile.Size} MB</Typography>
            {selectedFile.Description && (
              <Typography>Description: {selectedFile.Description}</Typography>
            )}
            <Typography>Rating: {selectedFile.Rating}</Typography>
          </Box>
        )}
  
          {Object.entries(providers).some(([_, provider]) => provider.IsActive) ? (
          Object.entries(providers)
            .filter(([_, provider]) => provider.IsActive)
            .map(([peerID, provider]) => (
              <Button
                key={peerID} // Ensure this is unique for the key
                variant="outlined"
                onClick={() => handleProviderSelect(peerID)}
                sx={{
                  margin: 1,
                  display: 'flex',
                  justifyContent: 'space-between',
                  width: '100%',
                }}
              >
                <span style={{ wordBreak: 'break-all', whiteSpace: 'normal' }}>
                  {peerID}
                </span>
                <span>{provider.Fee} SQD/MB</span>
              </Button>
            ))
        ) : (
          <Typography>No providers available for this file.</Typography>
        )}

        </DialogContent>
        <DialogActions>
          <Button onClick={() => {setOpen(false), setSearchResults(marketResults)}}>Cancel</Button>
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
