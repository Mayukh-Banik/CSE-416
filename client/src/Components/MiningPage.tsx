import React, { useState, useEffect } from "react";
import {
    Box,
    Typography,
    Button,
    CircularProgress,
    Snackbar,
    Alert,
    Paper,
    Grid,
    Card,
    CardContent,
    TextField,
} from "@mui/material";
import { styled, useTheme } from "@mui/material/styles";
import Sidebar from "./Sidebar";
import { useNavigate } from "react-router-dom";

const refreshInterval = 5000; // 5 seconds

// Styled components
const MiningContainer = styled(Box)(({ theme }) => ({
    padding: theme.spacing(4),
    transition: "margin-left 0.3s ease",
    [theme.breakpoints.down("sm")]: {
        marginLeft: 0,
    },
}));

// Define MiningStatus interface
interface MiningStatus {
    mining: boolean;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const MiningPage: React.FC = () => {
    const theme = useTheme();
    const navigate = useNavigate();
    const [isMining, setIsMining] = useState<boolean>(false);

    const [balance, setBalance] = useState<string>("0 SQC");
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [numBlocks, setNumBlocks] = useState<number>(1);

    const API_BASE_URL = "http://localhost:8080/api/btc"; 


    useEffect(() => {
        console.log("Mining status changed:", isMining);
        if (!isMining) {
            setIsLoading(false); // Unloading when mining stops
        }
    }, [isMining]);


    // Fetch mining status periodically
    useEffect(() => {
        const fetchMiningStatus = async () => {
            try {
                const response = await fetch(`${API_BASE_URL}/getminingstatus`, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                    },
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.message || "Failed to fetch mining status");
                }

                const data = await response.json();

                
                setIsMining((prevStatus) => {
                    if (prevStatus !== data.data.mining) {
                        console.log("Mining status changed:", data.data.mining);
                        setIsLoading(false);
                    }
                    return data.data.mining; 
                });
            } catch (err) {
                console.error("Error fetching mining status:", err);
                setError(err instanceof Error ? err.message : "An unknown error occurred.");
                setIsLoading(false); 
            }
        };

        fetchMiningStatus(); 
        const interval = setInterval(fetchMiningStatus, refreshInterval + 1000);

        return () => clearInterval(interval); 
    }, []);

    // Fetch balance periodically
    useEffect(() => {
        const fetchBalance = async () => {
            try {
                const response = await fetch(`${API_BASE_URL}/balance`, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                    },
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.message || "Failed to fetch balance");
                }

                const data = await response.json();
                setBalance(data.data.balance || "0 SQC");
            } catch (err) {
                console.error("Error fetching balance:", err);
                setError(err instanceof Error ? err.message : "An unknown error occurred.");
            }
        };

        fetchBalance(); 

        const interval = setInterval(fetchBalance, refreshInterval); 
        return () => clearInterval(interval); 
    }, []);

    // Start mining
    const handleStartMining = async () => {
        if (numBlocks <= 0) {
            setError("Number of blocks must be greater than 0");
            return;
        }

        console.log("Start Mining: Button clicked");
        setIsLoading(true);

        try {

            const response = await fetch(`${API_BASE_URL}/startmining`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ numBlock: numBlocks }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                console.error("Server Error:", errorData);
                throw new Error(errorData.message || "Failed to start mining");
            }

            const data = await response.json();
            console.log("Mining started successfully:", data);
            setSuccess(data.message || "Mining started successfully");
        } catch (error) {
            console.error("Error starting mining:", error); 
            setIsMining(false); 
            setError(error instanceof Error ? error.message : "An unknown error occurred.");
        } finally {
            console.log("Start Mining: Finished"); 
            setIsLoading(false); 
        }
    };

    // Stop mining
    const handleStopMining = async () => {
        console.log("Stop Mining: Button clicked"); 
        setIsLoading(true);

        try {

            const response = await fetch(`${API_BASE_URL}/stopmining`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                const errorData = await response.json();
                console.error("Server Error:", errorData); 
                throw new Error(errorData.message || "Failed to stop mining");
            }

            const data = await response.json();
            console.log("Mining stopped successfully:", data);
            setSuccess(data.message || "Mining stopped successfully");
        } catch (error) {
            console.error("Error stopping mining:", error);
            setError(error instanceof Error ? error.message : "An unknown error occurred.");
        } finally {
            console.log("Stop Mining: Finished"); 
            setIsLoading(false); 
        }
    };

    const handleCloseSnackbar = () => {
        setError(null);
        setSuccess(null);
    };

    return (
        <Box
            sx={{
                padding: 2,
                marginTop: '70px',
                marginLeft: `${drawerWidth}px`,
                transition: 'margin-left 0.3s ease',
                [theme.breakpoints.down('sm')]: {
                    marginLeft: `${collapsedDrawerWidth}px`,
                },
            }}
        >
            <Sidebar />
            <Typography variant="h4" gutterBottom>
                Mining
            </Typography>
            <MiningContainer>
                <Grid container spacing={4}>
                    <Grid item xs={12}>
                        <Paper elevation={3} sx={{ padding: 4 }}>
                            <Box
                                sx={{
                                    display: "flex",
                                    flexDirection: "column",
                                    alignItems: "center",
                                }}
                            >
                                <Typography variant="h6" gutterBottom>
                                    ~ SquidCoins ~
                                </Typography>
                                <Box
                                    sx={{
                                        display: "flex",
                                        alignItems: "center",
                                        gap: 2,
                                        marginTop: 2,
                                    }}
                                >
                                    <TextField
                                        type="number"
                                        label="Blocks to Mine"
                                        value={numBlocks}
                                        onChange={(e) => setNumBlocks(Number(e.target.value))}
                                        sx={{ width: 150 }}
                                        disabled={isMining || isLoading} 
                                        inputProps={{ min: 1 }}
                                    />
                                    <Button
                                        variant="contained"
                                        color="primary"
                                        onClick={isMining ? handleStopMining : handleStartMining}
                                        disabled={isLoading}
                                        sx={{ height: '56px' }}
                                    >
                                        {isLoading ? (<CircularProgress size={24} />) : isMining ? ("Stop Mining") : ("Start Mining")}
                                    </Button>
                                </Box>
                            </Box>
                        </Paper>
                    </Grid>

                    {/* Mining Status */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6">Current Mining Status</Typography>
                                <Typography variant="body1">{isMining ? "Mining in Progress" : "Idle"}</Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Mined Blocks - Placeholder */}
                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Mined Blocks
                                </Typography>
                                <Typography variant="body1" color="textSecondary">
                                    Coming Soon...
                                </Typography>
                            </CardContent>
                        </Card>
                    </Grid>

                    {/* Balance */}
                    <Grid item xs={12}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Balance
                                </Typography>
                                <Typography variant="body1">{balance}</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                </Grid>

                {/* Feedback */}
                <Snackbar
                    open={!!error}
                    autoHideDuration={6000}
                    onClose={handleCloseSnackbar}
                    anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
                >
                    <Alert onClose={handleCloseSnackbar} severity="error" sx={{ width: "100%" }}>
                        {error}
                    </Alert>
                </Snackbar>

                <Snackbar
                    open={!!success}
                    autoHideDuration={6000}
                    onClose={handleCloseSnackbar}
                    anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
                >
                    <Alert onClose={handleCloseSnackbar} severity="success" sx={{ width: "100%" }}>
                        {success}
                    </Alert>
                </Snackbar>
            </MiningContainer>
        </Box>
    );
};

export default MiningPage;
