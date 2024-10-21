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
import Sidebar from "./Sidebar";
import { useTheme } from "@mui/material/styles";

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const AccountViewPage: React.FC = () => {
  const theme = useTheme();

  // Initial account data
  const [accountDetails, setAccountDetails] = useState({
    walletId: "gen-public-key-123",
    totalVotes: 10, // Starting with 10 votes
    totalScore: 50,  // Starting with 10 votes, all 5 stars
    balance: 100,
  });

  // Dummy file data
  const [files, setFiles] = useState([
    { name: "file1.txt", size: 15, date: "2024-10-01", rating: 0, hasVoted: false },
    { name: "file2.txt", size: 30, date: "2024-10-02", rating: 0, hasVoted: false },
  ]);

  // Function to calculate reputation out of 5 stars
  const calculateReputation = () => {
    return (accountDetails.totalScore / accountDetails.totalVotes); 
  };

  // New 5 star vote
  const handleUpvote = (index: number) => {
    const newFiles = [...files];
    if (!newFiles[index].hasVoted) {
      newFiles[index].rating += 5;
      newFiles[index].hasVoted = true;
      setFiles(newFiles);

      // temp to be deleted
      setAccountDetails((prevAccountDetails) => ({
        ...prevAccountDetails,
        totalVotes: prevAccountDetails.totalVotes + 1,
        totalScore: prevAccountDetails.totalScore + 5,
      }));
    }
  };

  // delete later
  const handleDownvote = (index: number) => {
    const newFiles = [...files];
    if (!newFiles[index].hasVoted) {
      newFiles[index].rating += 0; 
      newFiles[index].hasVoted = true;
      setFiles(newFiles);

      // delete
      setAccountDetails((prevAccountDetails) => ({
        ...prevAccountDetails,
        totalVotes: prevAccountDetails.totalVotes + 1,
      }));
    }
  };

  return (
    <Box
      sx={{
        padding: 2,
        marginTop: "70px",
        marginLeft: `${drawerWidth}px`,
        transition: "margin-left 0.3s ease",
        [theme.breakpoints.down("sm")]: {
          marginLeft: `${collapsedDrawerWidth}px`,
        },
      }}
    >
      <Sidebar />
      <Box sx={{ flexGrow: 1, padding: 3 }}>
        <Typography variant="h4" gutterBottom>
          Account Information
        </Typography>

        <Box sx={{ mt: 2 }}>
          <Typography variant="h6">Wallet Id:</Typography>
          <Typography variant="body1">{accountDetails.walletId}</Typography>
          <Divider sx={{ my: 2 }} />

          <Typography variant="h6">Reputation (out of 5 stars):</Typography>
          <Typography variant="body1">{calculateReputation()} / 5</Typography>
          <Divider sx={{ my: 2 }} />

          <Typography variant="h6">Account Balance:</Typography>
          <Typography variant="body1">
            {accountDetails.balance.toFixed(2)} coins
          </Typography>
          <Divider sx={{ my: 2 }} />
        </Box>

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
                    <Button // rework buttons later
                      onClick={() => handleUpvote(index)}
                      disabled={file.hasVoted} 
                    >
                      Upvote
                    </Button>
                    <Button
                      onClick={() => handleDownvote(index)}
                      disabled={file.hasVoted} 
                    >
                      Downvote
                    </Button>
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
