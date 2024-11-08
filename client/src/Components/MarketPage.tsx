import React, { useState } from 'react';
import { Box, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Button, TextField, TablePagination, Paper, Dialog, DialogTitle, DialogContent, DialogActions, Snackbar, Alert } from '@mui/material';
import Sidebar from './Sidebar';
import { useTheme } from '@mui/material/styles';
import { FileMetadata, saveFileMetadata } from "../utils/localStorage";
import { render } from '@testing-library/react';


// export interface IFile {
//   fileName: string;
//   hash: string;
//   reputation: number;
//   fileSize: number; // in bytes
//   createdAt: Date;
//   providers: IProvider[];
// }

// adjust model schema to figure out how to connect different hosts and fees to the same file 
export interface IProvider {
  peerID: string;
  address: string;
  fee: number;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const MarketplacePage: React.FC = () => {
  const theme = useTheme();
  
  // const [files, setFiles] = useState<FileMetadata[]>([
  //   {
  //     fileName: 'Vacation_Snapshot.png',
  //     hash: 'img001',
  //     reputation: 4,
  //     fileSize: 2048000,
  //     createdAt: new Date('2023-09-15'),
  //     providers: [
  //       { providerId: '123', peerID: 'John', fee: 0.2 },
  //       { providerId: '124', peerID: 'Alice', fee: 0.5 },
  //     ],
  //   },
  //   {
  //     fileName: 'Project_Proposal.pdf',
  //     hash: 'doc002',
  //     reputation: 5,
  //     fileSize: 512000,
  //     createdAt: new Date('2023-08-10'),
  //     providers: [
  //       { providerId: '125', peerID: 'Bob', fee: 1.0 },
  //     ],
  //   },
  //   {
  //     fileName: 'Family_Photo.jpg',
  //     hash: 'img003',
  //     reputation: 3,
  //     fileSize: 1500000,
  //     createdAt: new Date('2023-07-22'),
  //     providers: [{ providerId: '127', peerID: 'Jim', fee: 0.4 }],
  //   },
  // ]);
  
  const [searchTerm, setSearchTerm] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [selectedFile, setSelectedFile] = useState<FileMetadata | null>(null);
  const [open, setOpen] = useState(false);
  const [fileHash, setFileHash] = useState('');
  const [notification, setNotification] = useState<{ open: boolean; message: string; severity: "success" | "error" }>({ open: false, message: "", severity: "success" });
  // const [providers, setProviders] = useState<IProvider[]>([
  //   {providerId: '123', peerID: 'John', fee: 0.2},
  //   {providerId: '127', peerID: 'Bob', fee: 1.2},
  // ]);
  const [providers, setProviders] = useState<IProvider[]>([]);

  const handleCloseNotification = () => {
    setNotification({ ...notification, open: false });
  };
  
  const handleDownloadRequest = (file: FileMetadata) => {
    setSelectedFile(file);
    setOpen(true); // Open the modal for provider selection
  };

  const handleProviderSelect = (provider: string) => {
    console.log(`Selected provider: ${provider} for file: ${selectedFile?.name}`);

    // Implement actual download logic here
    setNotification({ open: true, message: `Request sent to provider ${provider}`, severity: "success" });
    setOpen(false); // Close the modal after selecting a provider
  };

  const handleRefresh = () => {
    console.log("HI");
  };

  const handleDownloadByHash = async () => {
    console.log("HI");
    const hash = prompt("Enter the file hash");
    console.log('you entered:', hash)
    if(hash == null || hash.length==0) return;
    setFileHash(hash);
    handleFindProvidersByFileHash(hash);
    // .. other download stuff
    setOpen(true);
  }

  const handleFindProvidersByFileHash = async (hash:string) => {
    try {
        const response = await fetch("http://localhost:8081/fetch", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ val: hash }),
        });

        if (!response.ok) {
            throw new Error(`HTTP Error: status : ${response.status}`);
        }

        const data = await response.json();
        console.log("Providers:", data); // Log providers to check
        setProviders(data); // Assuming `setProviders` is a state setter in React
    } catch (error) {
        console.error("Error:", error);
        setNotification({ open: true, message: "Failed to find providers.", severity: "error" });

    }
};

// const fileHash = "ef5f60e42a6359d8e2d7e03d42b3de27644508d83a3813f67ef4e9b087948fe5"
const getFile = async (fileHash:string) => {
  console.log('in get file button')
  try {
      const response = await fetch(`http://localhost:8081/file/${fileHash}`);
      if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
      }
      console.log("got the file:", response)

  } catch (error) {
      console.error("Error fetching file:", error);
  }
};
const handleDownload = (fileHash:string) => {
  getFile(fileHash);
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

      <Button variant="contained" onClick={() => {handleRefresh}}>
        Refresh
      </Button>
      <Button variant="contained" onClick={() => handleDownloadByHash()}>
                    Download by Hash
                  </Button>
                  <button onClick={() => handleDownload(fileHash)}>Download File</button>
      <TextField
        label="Search Files"
        variant="outlined"
        fullWidth
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        sx={{ marginBottom: 2, background: "white" }}
      />
      
      {/* <TableContainer component={Paper}>
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
            {currentFiles.map((file) => (
              <TableRow key={file.hash}>
                <TableCell>{file.name}</TableCell>
                <TableCell>{(file.size / 1024).toFixed(2)}</TableCell>
                <TableCell>{file.reputation}</TableCell>
                <TableCell>{file.createdAt.toLocaleDateString()}</TableCell>
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
      <TablePagination
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
        <DialogTitle>Select a Provider for {selectedFile?.name}</DialogTitle>
        <DialogContent>
          {providers.length ? (
            // Filter out duplicate providers based on the peerID or address
            providers
              .filter((value, index, self) =>
                index === self.findIndex((t) => (
                  t.peerID === value.peerID // or use t.address if filtering by address
                ))
              )
              .map((provider) => (
                <Button
                  key={provider.peerID} // Ensure this is unique for the key
                  variant="outlined"
                  onClick={() => handleProviderSelect(provider.peerID)}
                  sx={{
                    margin: 1,
                    display: 'flex',
                    justifyContent: 'space-between',
                    width: '100%',
                  }}
                >
                  <span>{provider.peerID.substring(0, 8)}...{provider.peerID.substring(provider.peerID.length - 6)}</span> 
                  <span>{provider.fee} ORC/MB</span>
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
