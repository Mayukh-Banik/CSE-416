import React, { useState } from 'react';
import { Box, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Button, TextField, TablePagination, Paper } from '@mui/material';
import Sidebar from './Sidebar';
import { useTheme } from '@mui/material/styles';

export interface IFile {
  fileName: string;
  hash: string;
  reputation: string;
  fileSize: number; // in bytes
  createdAt: Date;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const MarketplacePage: React.FC = () => {
  const theme = useTheme();
  //dummy data from chat
  const [files, setFiles] = useState<IFile[]>([
    { fileName: 'Vacation_Snapshot.png', hash: 'img001', reputation: 'High', fileSize: 2048000, createdAt: new Date('2023-09-15') },
    { fileName: 'Project_Proposal.pdf', hash: 'doc002', reputation: 'Medium', fileSize: 512000, createdAt: new Date('2023-08-10') },
    { fileName: 'Family_Photo.jpg', hash: 'img003', reputation: 'High', fileSize: 1500000, createdAt: new Date('2023-07-22') },
    { fileName: 'Recipe_Book.pdf', hash: 'doc004', reputation: 'Low', fileSize: 1000000, createdAt: new Date('2023-10-01') },
    { fileName: 'Nature_Wallpaper.jpg', hash: 'img005', reputation: 'High', fileSize: 2500000, createdAt: new Date('2023-06-05') },
    { fileName: 'User_Manual.png', hash: 'img006', reputation: 'Medium', fileSize: 350000, createdAt: new Date('2023-05-14') },
    { fileName: 'Presentation_Slides.pptx', hash: 'doc007', reputation: 'Medium', fileSize: 700000, createdAt: new Date('2023-09-20') },
    { fileName: 'Financial_Report.pdf', hash: 'doc008', reputation: 'High', fileSize: 800000, createdAt: new Date('2023-10-05') },
    { fileName: 'Tech_Article.jpg', hash: 'img009', reputation: 'Low', fileSize: 600000, createdAt: new Date('2023-04-30') },
    { fileName: 'Game_Screenshots.png', hash: 'img010', reputation: 'Medium', fileSize: 1200000, createdAt: new Date('2023-09-30') },
  ]);
  
  const [searchTerm, setSearchTerm] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(5);

  const filteredFiles = files.filter(file => 
    file.fileName.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleDownloadRequest = (file: IFile) => {
    //add implementation later
    console.log(`Requesting download for: ${file.fileName}`);
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
    </Box>
  );
};

export default MarketplacePage;
