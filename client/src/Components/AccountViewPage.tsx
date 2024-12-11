import React, { useEffect, useState } from "react";
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
import { FileMetadata } from "../models/fileMetadata";

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const AccountViewPage: React.FC = () => {
  const [uploadedFiles, setUploadedFiles] = useState<FileMetadata[]>([])
  const [ratings, setRatings] = useState<{ [key: string]: number }>({});
  const theme = useTheme();

  // Initial account data
  const [accountDetails, setAccountDetails] = useState({
    walletId: "gen-public-key-123",
    totalVotes: 10, // Starting with 10 votes
    totalScore: 50,  // Starting with 10 votes, all 5 stars (10 * 5 = 50)
    balance: 100,
  });

  useEffect(() => {
      const fetchFiles = async () => {
        try {
          console.log("Getting local user's uploaded files");
          let fileType = "uploaded"
          const response = await fetch(`http://localhost:8081/files/fetch?file=${fileType}`, {
            method: "GET",
          });
          if (!response.ok) throw new Error(`Failed to load ${fileType} file data`);
    
          const data = await response.json();
          console.log("Fetched data", data);
  
          setUploadedFiles(data); // Set the state with the loaded data
        } catch (error) {
          console.error("Error fetching files:", error);
        }
      };
  
      fetchFiles();
    }, []);


  const handleVote = async (fileHash: string, voteType: 'upvote' | 'downvote') => { 
    try {
      const response = await fetch(`http://localhost:8081/files/vote?fileHash=${fileHash}&voteType=${voteType}`, {
        method: "POST",
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        setRatings(prevRatings => {
          const currentRating = prevRatings[fileHash] || 0;
          const newRating =
            voteType === 'upvote' ? currentRating + 1 : currentRating - 1;
      
          return { ...prevRatings, [fileHash]: newRating };
        });
      } else {
        throw new Error("Failed to update vote");
      }
    } catch (error) {
      console.error("Error updating vote:", error);
    }
  };


  // Function to calculate reputation out of 5 stars
  const calculateReputation = () => {
    return (accountDetails.totalScore / accountDetails.totalVotes).toFixed(2); // Round to 2 decimal points
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

        {/* Account details */}
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
            {uploadedFiles.length > 0 ? (
              uploadedFiles.map((file, index) => (
                <TableRow key={index}>
                  <TableCell>{file.Name}</TableCell>
                  <TableCell>{file.Size}</TableCell>
                  <TableCell>{file.CreatedAt}</TableCell>
                  <TableCell>{file.VoteType}</TableCell>
                  <TableCell>
                    <Button
                      onClick={() => handleVote(file.Hash, "upvote")}
                      disabled={file.HasVoted} // Disable button if already voted
                    >
                      Upvote
                    </Button>
                    <Button
                      onClick={() => handleVote(file.Hash, "downvote")}
                      disabled={file.HasVoted || file.OriginalUploader} // Disable button if already voted
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
        {/* <Typography variant="h6" sx={{ mt: 3 }}>
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
            <TableRow>
              <TableCell colSpan={3} align="center">
                No files downloaded
              </TableCell>
            </TableRow>
          </TableBody>
        </Table> */}

      </Box>
    </Box>
  );
};

export default AccountViewPage;