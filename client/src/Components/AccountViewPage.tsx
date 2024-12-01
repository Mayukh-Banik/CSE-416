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
          const response = await fetch("http://localhost:8081/files/fetchAll");
          if (!response.ok) throw new Error("Failed to load file data");
  
          const data = await response.json();
          console.log("Fetched data", data);
  
          setUploadedFiles(data); // Set the state with the loaded data
        } catch (error) {
          console.error("Error fetching files:", error);
        }
      };
  
      fetchFiles();
    }, []);

  // Function to calculate reputation out of 5 stars
  const calculateReputation = () => {
    return (accountDetails.totalScore / accountDetails.totalVotes).toFixed(2); // Round to 2 decimal points
  };

  // // Function to handle upvote (equivalent to a 5-star vote)
  // const handleUpvote = (index: number) => {
  //   const newFiles = [...files];
  //   if (!newFiles[index].hasVoted) {
  //     newFiles[index].rating += 5;
  //     newFiles[index].hasVoted = true;
  //     setFiles(newFiles);

  //     // Add 5 stars to the total score and increase vote count by 1
  //     setAccountDetails((prevAccountDetails) => ({
  //       ...prevAccountDetails,
  //       totalVotes: prevAccountDetails.totalVotes + 1,
  //       totalScore: prevAccountDetails.totalScore + 5,
  //     }));
  //   }
  // };

  // // Function to handle downvote (equivalent to a 0-star vote)
  // const handleDownvote = (index: number) => {
  //   const newFiles = [...files];
  //   if (!newFiles[index].hasVoted) {
  //     newFiles[index].rating += 0; // No change in rating for a downvote
  //     newFiles[index].hasVoted = true;
  //     setFiles(newFiles);

  //     // Add 0 stars to the total score and increase vote count by 1
  //     setAccountDetails((prevAccountDetails) => ({
  //       ...prevAccountDetails,
  //       totalVotes: prevAccountDetails.totalVotes + 1,
  //     }));
  //   }
  // };

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
                  <TableCell>{file.Reputation}</TableCell>
                  <TableCell>
                    <Button
                      // onClick={() => handleUpvote(index)}
                      // disabled={file.hasVoted} // Disable button if already voted
                    >
                      Upvote
                    </Button>
                    <Button
                      // onClick={() => handleDownvote(index)}
                      // disabled={file.hasVoted} // Disable button if already voted
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
