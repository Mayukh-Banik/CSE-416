import React, { useState } from "react";
import {
  Box,
  Typography,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Button,
  Paper,
} from "@mui/material";
import Sidebar from "./Sidebar"; // Assuming Sidebar is a common component
import { useTheme } from '@mui/material/styles';

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const AccountViewPage: React.FC = () => {
  const theme = useTheme();
  
  // Dummy account data
  const accountDetails = {
    name: "john_doe",
    reputation: 150, // Dummy reputation score
    balance: 1000.50, // Dummy account balance
  };

  // Dummy file data
  const [files, setFiles] = useState([
    { name: "file1.txt", size: 15, date: "2024-10-01", rating: 0 },
    { name: "file2.txt", size: 30, date: "2024-10-02", rating: 0 },
  ]);

  // Functions to handle upvote and downvote
  const handleUpvote = (index: number) => {
    const newFiles = [...files];
    newFiles[index].rating += 1;
    setFiles(newFiles);
  };

  const handleDownvote = (index: number) => {
    const newFiles = [...files];
    newFiles[index].rating -= 1;
    setFiles(newFiles);
  };

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
      <Sidebar />
      <Box sx={{ flexGrow: 1, padding: 3 }}>
        <Typography variant="h4" gutterBottom>
          Account Information
        </Typography>

        {/* Account details */}
        <Box sx={{ mt: 2 }}>
          <Typography variant="h6">Account Name:</Typography>
          <Typography variant="body1">{accountDetails.name}</Typography>
          <Divider sx={{ my: 2 }} />

          <Typography variant="h6">Reputation:</Typography>
          <Typography variant="body1">{accountDetails.reputation}</Typography>
          <Divider sx={{ my: 2 }} />

          <Typography variant="h6">Account Balance:</Typography>
          <Typography variant="body1">${accountDetails.balance.toFixed(2)}</Typography>
          <Divider sx={{ my: 2 }} />
        </Box>

        {/* Uploaded files table */}
        <Typography variant="h6" sx={{ mt: 3 }}>
          Uploaded Files
        </Typography>
        <Table component={Paper}>
          <TableHead>
            <TableRow>
              <TableCell>File Name</TableCell>
              <TableCell>File Size (KB)</TableCell>
              <TableCell>Upload Date</TableCell>
              <TableCell>Rating</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {files.length > 0 ? (
              files.map((file, index) => (
                <TableRow key={index}>
                  <TableCell>{file.name}</TableCell>
                  <TableCell>{file.size}</TableCell>
                  <TableCell>{file.date}</TableCell>
                  <TableCell>{file.rating}</TableCell>
                  <TableCell>
                    <Button onClick={() => handleUpvote(index)}>Upvote</Button>
                    <Button onClick={() => handleDownvote(index)}>Downvote</Button>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={5} align="center">
                  No files uploaded
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>

        {/* Downloaded files table */}
        <Typography variant="h6" sx={{ mt: 3 }}>
          Downloaded Files
        </Typography>
        <Table component={Paper}>
          <TableHead>
            <TableRow>
              <TableCell>File Name</TableCell>
              <TableCell>File Size (KB)</TableCell>
              <TableCell>Download Date</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {/* Empty table for now */}
            <TableRow>
              <TableCell colSpan={3} align="center">
                No files downloaded
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </Box>
    </Box>
  );
};

export default AccountViewPage;
