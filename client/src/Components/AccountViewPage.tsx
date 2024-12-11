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
  Card,
  CardContent,
} from "@mui/material";
import Sidebar from "./Sidebar";
import { useTheme } from "@mui/material/styles";
import { FileMetadata } from "../models/fileMetadata";

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const AccountViewPage: React.FC = () => {
  const theme = useTheme();

  // State for uploaded files
  const [uploadedFiles, setUploadedFiles] = useState<FileMetadata[]>([]);

  // State for wallet address and balance
  const [walletDetails, setWalletDetails] = useState({
    walletAddress: "",
    balance: "",
  });

  // Fetch wallet address and balance
  useEffect(() => {
    const fetchWalletDetails = async () => {
      try {
        const response = await fetch(
          "http://localhost:8080/api/auth/getMiningAddressAndBalance"
        );
        if (!response.ok) {
          throw new Error("Failed to fetch wallet details");
        }
        const data = await response.json();
        console.log("Fetched wallet details:", data);

        setWalletDetails({
          walletAddress: data.mining_address,
          balance: data.balance,
        });
      } catch (error) {
        console.error("Error fetching wallet details:", error);
      }
    };

    fetchWalletDetails();

    // Periodically update wallet balance
    const interval = setInterval(() => {
      fetchWalletDetails();
    }, 10000); // Update every 10 seconds

    return () => clearInterval(interval);
  }, []);

  // Fetch uploaded files
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

        {/* Wallet Address and Balance */}
        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Typography variant="h6">Wallet Address</Typography>
            <Typography
              variant="body1"
              sx={{
                wordWrap: "break-word",
                mt: 1,
                mb: 2,
                fontFamily: "monospace",
                backgroundColor: "#f5f5f5",
                padding: "8px",
                borderRadius: "4px",
              }}
            >
              {walletDetails.walletAddress || "Loading..."}
            </Typography>

            <Divider sx={{ my: 2 }} />

            <Typography variant="h6">Account Balance</Typography>
            <Typography
              variant="body1"
              sx={{
                fontWeight: "bold",
                color: theme.palette.primary.main,
                mt: 1,
              }}
            >
              {walletDetails.balance ? `${walletDetails.balance} BTC` : "Loading..."}
            </Typography>
          </CardContent>
        </Card>

        {/* Uploaded Files */}
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
                  <TableCell>{file.Rating}</TableCell>
                  <TableCell>
                    <Button>Upvote</Button>
                    <Button>Downvote</Button>
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

        {/* Downloaded Files */}
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
