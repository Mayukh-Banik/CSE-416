import React, { useState } from 'react';
import { Box, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Button, TextField, TablePagination, Paper, Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import Sidebar from './Sidebar';
import { useTheme } from '@mui/material/styles';

export interface IFile {
  fileName: string;
  hash: string;
  reputation: number;
  fileSize: number; // in bytes
  createdAt: Date;
  providers: IProvider[];
}

// adjust model schema to figure out how to connect different hosts and fees to the same file 
export interface IProvider {
  providerId: string;
  providerName: string;
  fee: number;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const MarketplacePage: React.FC = () => {
  const theme = useTheme();
  
  const [files, setFiles] = useState<IFile[]>([
    {
      fileName: 'Vacation_Snapshot.png',
      hash: 'img001',
      reputation: 4,
      fileSize: 2048000,
      createdAt: new Date('2023-09-15'),
      providers: [
        { providerId: '123', providerName: 'John', fee: 0.2 },
        { providerId: '124', providerName: 'Alice', fee: 0.5 },
      ],
    },
    {
      fileName: 'Project_Proposal.pdf',
      hash: 'doc002',
      reputation: 5,
      fileSize: 512000,
      createdAt: new Date('2023-08-10'),
      providers: [
        { providerId: '125', providerName: 'Bob', fee: 1.0 },
      ],
    },
    {
      fileName: 'Family_Photo.jpg',
      hash: 'img003',
      reputation: 3,
      fileSize: 1500000,
      createdAt: new Date('2023-07-22'),
      providers: [{ providerId: '127', providerName: 'Jim', fee: 0.4 }],
    },
  ]);
  
  const [searchTerm, setSearchTerm] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [selectedFile, setSelectedFile] = useState<IFile | null>(null);
  const [open, setOpen] = useState(false);

  const [providers, setProviders] = useState<IProvider[]>([
    {providerId: '123', providerName: 'John', fee: 0.2},
    {providerId: '127', providerName: 'Bob', fee: 1.2},
  ]);

  const filteredFiles = files.filter(file => 
    file.fileName.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleDownloadRequest = (file: IFile) => {
    setSelectedFile(file);
    setOpen(true); // Open the modal for provider selection
  };

  const handleProviderSelect = (provider: string) => {
    console.log(`Selected provider: ${provider} for file: ${selectedFile?.fileName}`);
    // Implement actual download logic here
    setOpen(false); // Close the modal after selecting a provider
  };

  const handleRefresh = () => {
    console.log("HI");
  };

  const handleDownloadByHash = async () =>
  {
    console.log("HI");
    const hash = prompt("Enter the file hash");
    console.log('you entered:', hash)
    if(hash == null || hash.length==0) return;

    const response = await fetch("http://localhost:8081/fetch",{
      method :"POST",
      headers:{
        "Content-Type": "application/json",
      },
      body: JSON.stringify({val: hash}),
    })

    if (!response.ok)
    {
      throw new Error(`HTTP Error: status : ${response.status}`);
    }

    const data = await response.text();
    console.log("File fetching successful:", data);
  }

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
  const currentFiles = filteredFiles.slice(indexOfFirstFile, indexOfLastFile);

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

      <TextField
        label="Search Files"
        variant="outlined"
        fullWidth
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        sx={{ marginBottom: 2, background: "white" }}
      />
      
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
            {currentFiles.map((file) => (
              <TableRow key={file.hash}>
                <TableCell>{file.fileName}</TableCell>
                <TableCell>{(file.fileSize / 1024).toFixed(2)}</TableCell>
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
      />

      <Dialog open={open} onClose={() => setOpen(false)}>
        <DialogTitle>Select a Provider for {selectedFile?.fileName}</DialogTitle>
        <DialogContent>
          {selectedFile?.providers.length ? (
            selectedFile.providers.map((provider) => (
              <Button
                key={provider.providerId}
                variant="outlined"
                onClick={() => handleProviderSelect(provider.providerName)}
                sx={{ margin: 1, display: 'flex', justifyContent: 'space-between', width: '100%' }}
              >
                <span>{provider.providerName}</span>
                <span>{provider.fee} ORC/MB</span>
              </Button>
            ))
          ) : (
            <Typography>No providers available for this file.</Typography>  // Message if no providers
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Cancel</Button>
        </DialogActions>
      </Dialog>


    </Box>
  );
};

export default MarketplacePage;
